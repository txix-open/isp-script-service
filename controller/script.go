package controller

import (
	"isp-script-service/domain"
	"isp-script-service/service"
)

var Script scriptController

type scriptController struct{}

// @Tags script
// @Summary Выполнить скрипт под конкретным ID
// @Description Возвращает результат выполнения скрипта и ошибку (если есть)
// @Accept  json
// @Produce  json
// @Param body body domain.ExecuteByIdRequest true "идентификатор скрипта"
// @Success 200 {object} domain.ScriptResp
// @Failure 404 {object} structure.GrpcError
// @Router /script/execute_by_id [POST].
func (c scriptController) ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error) {
	return service.Script.ExecuteById(req)
}

// @Tags script
// @Summary Выполнить скрипт без учёта идентификатора
// @Description Возвращает результат выполнения скрипта и ошибку (если есть)
// @Accept  json
// @Produce  json
// @Param body body domain.ExecuteRequest true "Скрипт необходимый к выполнению"
// @Success 200 {object} domain.ScriptResp
// @Router /script/execute [POST].
func (c scriptController) Execute(req domain.ExecuteRequest) *domain.ScriptResp {
	return service.Script.Execute(req)
}

// @Tags script
// @Summary Выполнить набор скриптов под конкретными ID
// @Description Возвращает результат выполнения скриптов и ошибок (если есть)
// @Accept  json
// @Produce  json
// @Param body body []domain.ExecuteByIdRequest true "Набор идентификаторов и аргументов"
// @Success 200 {array} domain.ScriptResp
// @Router /script/batch_execute [POST].
func (c scriptController) BatchExecute(req []domain.ExecuteByIdRequest) []domain.ScriptResp {
	return service.Script.BatchExecute(req)
}

// @Tags script
// @Summary Выполнить набор скриптов под конкретными ID с идентичным аргументом для всех
// @Description Возвращает результат выполнения скриптов и ошибок (если есть)
// @Accept  json
// @Produce  json
// @Param body body domain.BatchExecuteByIdsRequest true "Набор идентификаторов и единый аргумент"
// @Success 200 {array} domain.ScriptResp
// @Router /script/batch_execute_by_ids [POST].
func (c scriptController) BatchExecuteById(req domain.BatchExecuteByIdsRequest) []domain.ScriptResp {
	return service.Script.BatchExecuteById(req)
}
