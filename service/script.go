package service

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/integration-system/isp-lib/v2/scripts"
	log "github.com/integration-system/isp-log"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	log_code "isp-script-service/codes"
	"isp-script-service/conf"
	"isp-script-service/domain"
	"isp-script-service/router"
)

var Script = &scriptService{
	scriptEngine: scripts.NewEngine(),
}

type scriptService struct {
	store        atomic.Value
	scriptEngine *scripts.Engine
}

type scriptServiceCfg struct {
	Scripts          map[string]CompiledScript
	ExecutionTimeout time.Duration
	SharedScript     string
}

type CompiledScript struct {
	scripts.Script
	Compiled bool
}

func (s *scriptService) ReceiveConfiguration(scriptDef []conf.ScriptDefinition, sharedScript string, executionTimeoutMs int) {
	var compiled bool
	cfg := &scriptServiceCfg{
		Scripts:          make(map[string]CompiledScript),
		ExecutionTimeout: time.Duration(executionTimeoutMs) * time.Millisecond,
		SharedScript:     sharedScript,
	}
	for i, value := range scriptDef {
		compiled = true
		scr, err := s.Create(value.Script, sharedScript)
		if err != nil {
			log.Errorf(log_code.CreateScriptFromConfigError, "create script from config (number %d): %v", i, err)
			compiled = false
		}
		cfg.Scripts[value.Id] = CompiledScript{scr, compiled}
	}
	s.store.Store(cfg)
}

func (s *scriptService) Execute(req domain.ExecuteRequest) *domain.ScriptResp {
	cfg := s.loadCfg()
	scr, err := s.Create(req.Script, cfg.SharedScript)
	if err != nil {
		return s.respError(err, domain.ErrorCompile)
	}

	return s.executeScript(CompiledScript{scr, true}, req.Arg, cfg.ExecutionTimeout)
}

func (s *scriptService) ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error) {
	cfg := s.loadCfg()
	scr, ok := cfg.Scripts[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "not defined script for id %s", req.Id)
	}

	return s.executeScript(scr, req.Arg, cfg.ExecutionTimeout), nil
}

func (s *scriptService) BatchExecute(req []domain.ExecuteByIdRequest) []domain.ScriptResp {
	wg := sync.WaitGroup{}
	cfg := s.loadCfg()
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

func (s *scriptService) BatchExecuteById(req domain.BatchExecuteByIdsRequest) []domain.ScriptResp {
	wg := sync.WaitGroup{}
	cfg := s.loadCfg()
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

func (*scriptService) Create(script, sharedScript string) (scripts.Script, error) {
	return scripts.NewScript([]byte(sharedScript),
		[]byte("(function() {\n"), []byte(script), []byte("\n})();"))
}

var errEmpty = errors.New("empty answer, maybe lost return")

func (s *scriptService) executeScript(scr CompiledScript, arg interface{}, timeout time.Duration) *domain.ScriptResp {
	if !scr.Compiled {
		return s.respError(errors.New("invalid script configuration"), domain.ErrorCompile)
	}

	response, err := s.scriptEngine.Execute(scr.Script, arg,
		scripts.WithScriptTimeout(timeout),
		// TODO: remove. invoke is deprecated, all functions should be inside `external` object
		scripts.WithSet("invoke", router.Invoke),
		scripts.WithSet("external", map[string]interface{}{
			"invoke":         router.Invoke,
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

func (*scriptService) respError(err error, errorType string) *domain.ScriptResp {
	respError := domain.Error{
		Type:        errorType,
		Description: err.Error(),
	}

	return &domain.ScriptResp{
		Error: &respError,
	}
}

func (s *scriptService) loadCfg() *scriptServiceCfg {
	return s.store.Load().(*scriptServiceCfg)
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
