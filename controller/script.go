package controller

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-script-service/compile"
	"isp-script-service/domain"
	"isp-script-service/script"
)

var Script scriptController

type scriptController struct{}

func (c scriptController) ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error) {
	scr := compile.Script.GetById(req.Id)
	if scr == nil {
		return nil, status.Errorf(codes.NotFound, "not defined script for id %s", req.Id)
	}
	return c.executeScript(scr, req.Arg), nil
}

func (c scriptController) Execute(req domain.ExecuteRequest) *domain.ScriptResp {
	scr, err := compile.Script.Create(req.Script)
	if err != nil {
		return c.respError(err, domain.ErrorCompile)
	}
	return c.executeScript(scr, req.Arg)
}

func (c scriptController) executeScript(scr script.Script, arg interface{}) *domain.ScriptResp {
	response, err := script.Default().Execute(scr, arg)
	if err != nil {
		return c.respError(err, domain.ErrorRunTime)
	}
	if response == nil {
		return c.respError(errors.New("empty answer, maybe lost return"), domain.ErrorRunTime)
	}
	return &domain.ScriptResp{
		Result: response,
	}
}

func (c scriptController) respError(err error, errorType string) *domain.ScriptResp {
	respError := domain.Error{
		Type:        errorType,
		Description: err.Error(),
	}
	return &domain.ScriptResp{
		Error: &respError,
	}
}
