package service

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/integration-system/isp-lib/v2/config"
	log "github.com/integration-system/isp-log"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	log_code "isp-script-service/codes"
	"isp-script-service/conf"
	"isp-script-service/domain"
	"isp-script-service/script"
)

var Script = &scriptService{}

type scriptService struct {
	store atomic.Value
}

func (s *scriptService) ReceiveConfiguration(scriptDef []conf.ScriptDefinition) {
	newStore := make(map[string]script.Script)
	for i, value := range scriptDef {
		scr, err := s.Create(value.Script)
		if err != nil {
			log.Errorf(log_code.CreateScriptFromConfigError, "create script from config (number %d): %v", i, err)
			continue
		}
		newStore[value.Id] = scr
	}
	s.store.Store(newStore)
}

func (s *scriptService) Execute(req domain.ExecuteRequest) *domain.ScriptResp {
	scr, err := s.Create(req.Script)
	if err != nil {
		return s.respError(err, domain.ErrorCompile)
	}
	return s.executeScript(scr, req.Arg)
}

func (s *scriptService) ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error) {
	scr := s.store.Load().(map[string]script.Script)[req.Id]
	if scr == nil {
		return nil, status.Errorf(codes.NotFound, "not defined script for id %s", req.Id)
	}
	return s.executeScript(scr, req.Arg), nil
}

func (s *scriptService) BatchExecute(req []domain.ExecuteByIdRequest) []domain.ScriptResp {
	wg := sync.WaitGroup{}
	store := s.store.Load().(map[string]script.Script)
	response := make([]domain.ScriptResp, len(req))
	for i := range req {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if store[req[i].Id] == nil {
				response[i] = *s.respError(errors.Errorf("not defined script for id %s", req[i].Id), domain.ErrorCompile)
			} else {
				response[i] = *s.executeScript(store[req[i].Id], req[i].Arg)
			}
		}(i)
	}
	wg.Wait()
	return response
}

func (s *scriptService) BatchExecuteById(req domain.BatchExecuteByIdsRequest) []domain.ScriptResp {
	wg := sync.WaitGroup{}
	store := s.store.Load().(map[string]script.Script)
	response := make([]domain.ScriptResp, len(req.Ids))
	for i := range req.Ids {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if store[req.Ids[i]] == nil {
				response[i] = *s.respError(errors.Errorf("not defined script for id %s", req.Ids[i]), domain.ErrorCompile)
			} else {
				response[i] = *s.executeScript(store[req.Ids[i]], req.Arg)
			}
		}(i)
	}
	wg.Wait()
	return response
}

func (*scriptService) Create(scr string) (script.Script, error) {
	cfg := config.GetRemote().(*conf.RemoteConfig)
	scr = fmt.Sprintf("%s%s%s%s", cfg.SharedScript, "(function() {", scr, "})();")
	return script.Default().Compile(scr)
}

func (s *scriptService) executeScript(scr script.Script, arg interface{}) *domain.ScriptResp {
	response, err := script.Default().Execute(scr, arg)
	if err != nil {
		return s.respError(err, domain.ErrorRunTime)
	}
	if response == nil {
		return s.respError(errors.New("empty answer, maybe lost return"), domain.ErrorRunTime)
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
