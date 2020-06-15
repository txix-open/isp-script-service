package controller

import (
	"isp-script-service/domain"
	"isp-script-service/service"
)

var Script scriptController

type scriptController struct{}

func (c scriptController) ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error) {
	return service.Script.ExecuteById(req)
}

func (c scriptController) Execute(req domain.ExecuteRequest) *domain.ScriptResp {
	return service.Script.Execute(req)
}

func (c scriptController) BatchExecuteById(req domain.BatchExecuteByIdsRequest) []domain.ScriptResp {
	return service.Script.BatchExecuteById(req)
}
