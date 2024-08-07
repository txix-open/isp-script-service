// nolint:wrapcheck
package service

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/log"
	"github.com/txix-open/isp-script"
	"isp-script-service/conf"
	"isp-script-service/domain"
)

type scriptServiceCfg struct {
	Scripts          map[string]CompiledScript
	ExecutionTimeout time.Duration
	SharedScript     string
}

type CompiledScript struct {
	scripts.Script
	Compiled bool
}

type Router interface {
	Invoke(methodPath string, request interface{}, metadata map[string]interface{}) (interface{}, error)
}

type Script struct {
	store        scriptServiceCfg
	scriptEngine *scripts.Engine
	logger       log.Logger
	router       Router
}

func NewScript(invoker Router, logger log.Logger, sd []conf.ScriptDefinition, shared string, timeout int) *Script {
	scr := scriptsFromConfig(context.Background(), logger, sd, shared, timeout)
	scriptInstance := Script{
		store:        scr,
		scriptEngine: scripts.NewEngine(),
		logger:       logger,
		router:       invoker,
	}

	return &scriptInstance
}

func (s *Script) Execute(req domain.ExecuteRequest) *domain.ScriptResp {
	cfg := s.store
	scr, err := createScript(req.Script, cfg.SharedScript)
	if err != nil {
		return s.respError(err, domain.ErrorCompile)
	}

	return s.executeScript(CompiledScript{scr, true}, req.Arg, cfg.ExecutionTimeout)
}

func (s *Script) ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error) {
	cfg := s.store
	scr, ok := cfg.Scripts[req.Id]
	if !ok {
		return nil, domain.ErrScriptNotFound
	}

	return s.executeScript(scr, req.Arg, cfg.ExecutionTimeout), nil
}

func (s *Script) BatchExecute(req []domain.ExecuteByIdRequest) []domain.ScriptResp {
	wg := sync.WaitGroup{}
	cfg := s.store
	response := make([]domain.ScriptResp, len(req))
	for i := range req {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if scr, ok := cfg.Scripts[req[i].Id]; !ok {
				response[i] = *s.respError(errors.Errorf("not defined script for id %s", req[i].Id), domain.ErrorCompile)
			} else {
				response[i] = *s.executeScript(scr, req[i].Arg, cfg.ExecutionTimeout)
			}
		}(i)
	}
	wg.Wait()

	return response
}

func (s *Script) BatchExecuteById(req domain.BatchExecuteByIdsRequest) []domain.ScriptResp {
	wg := sync.WaitGroup{}
	cfg := s.store
	response := make([]domain.ScriptResp, len(req.Ids))
	for i := range req.Ids {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if scr, ok := cfg.Scripts[req.Ids[i]]; !ok {
				response[i] = *s.respError(errors.Errorf("not defined script for id %s", req.Ids[i]), domain.ErrorCompile)
			} else {
				response[i] = *s.executeScript(scr, req.Arg, cfg.ExecutionTimeout)
			}
		}(i)
	}
	wg.Wait()

	return response
}

var errEmpty = errors.New("empty answer, maybe lost return")

func (s *Script) executeScript(scr CompiledScript, arg interface{}, timeout time.Duration) *domain.ScriptResp {
	if !scr.Compiled {
		return s.respError(errors.New("invalid script configuration"), domain.ErrorCompile)
	}

	response, err := s.scriptEngine.Execute(scr.Script, arg,
		scripts.WithTimeout(timeout),
		scripts.WithDefaultToolkit(),
		scripts.WithSet("external", map[string]interface{}{
			"invoke":         s.router.Invoke,
			"hashSha256":     Sha256,
			"hashSha512":     Sha512,
			"generateUUIDv4": UUIDv4,
		}),
	)
	if err != nil {
		return s.respError(err, domain.ErrorRunTime)
	}
	if response == nil {
		return s.respError(errEmpty, domain.ErrorRunTime)
	}

	return &domain.ScriptResp{
		Result: response,
	}
}

func (*Script) respError(err error, errorType string) *domain.ScriptResp {
	respError := domain.Error{
		Type:        errorType,
		Description: err.Error(),
	}

	return &domain.ScriptResp{
		Error: &respError,
	}
}

func Sha256(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func Sha512(value string) string {
	sum := sha512.Sum512([]byte(value))
	return hex.EncodeToString(sum[:])
}

func UUIDv4() string {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		// will be handled as exception in js
		panic(err)
	}
	return randomUUID.String()
}

func scriptsFromConfig(
	ctx context.Context,
	logger log.Logger,
	scriptDef []conf.ScriptDefinition,
	sharedScript string,
	executionTimeoutMs int,
) scriptServiceCfg {
	var compiled bool
	cfg := scriptServiceCfg{
		Scripts:          make(map[string]CompiledScript),
		ExecutionTimeout: time.Duration(executionTimeoutMs) * time.Millisecond,
		SharedScript:     sharedScript,
	}
	for _, value := range scriptDef {
		compiled = true
		scr, err := createScript(value.Script, sharedScript)
		if err != nil {
			logger.Error(ctx, errors.WithMessage(err, "create script from config"),
				log.String("scriptId", value.Id))
			compiled = false
		}
		cfg.Scripts[value.Id] = CompiledScript{scr, compiled}
	}

	return cfg
}

func createScript(script, sharedScript string) (scripts.Script, error) {
	return scripts.NewScript([]byte(sharedScript),
		[]byte("(function() {\n"), []byte(script), []byte("\n})();"))
}
