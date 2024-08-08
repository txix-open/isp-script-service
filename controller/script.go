package controller

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-script-service/domain"
)

type scriptService interface {
	Execute(req domain.ExecuteRequest) *domain.ScriptResp
	ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error)
	BatchExecute(req []domain.ExecuteByIdRequest) []domain.ScriptResp
	BatchExecuteById(req domain.BatchExecuteByIdsRequest) []domain.ScriptResp
}

type Script struct {
	service scriptService
}

func NewScript(service scriptService) Script {
	return Script{service: service}
}

// ExecuteById
// @Tags script
// @Summary Выполнить скрипт под конкретным ID
// @Description Возвращает результат выполнения скрипта и ошибку (если есть)
// @Accept  json
// @Produce  json
// @Param body body domain.ExecuteByIdRequest true "идентификатор скрипта"
// @Success 200 {object} domain.ScriptResp
// @Failure 404 {object} domain.GrpcError
// @Router /script/execute_by_id [POST].
func (c Script) ExecuteById(req domain.ExecuteByIdRequest) (*domain.ScriptResp, error) {
	data, err := c.service.ExecuteById(req)

	switch {
	case errors.Is(err, domain.ErrScriptNotFound):
		return nil, status.Errorf(codes.NotFound, "not defined script for id %s", req.Id)
	case err != nil:
		return nil, err
	default:
		return data, nil
	}
}

// Execute
// @Tags script
// @Summary Выполнить скрипт без учёта идентификатора
// @Description Возвращает результат выполнения скрипта и ошибку (если есть)
// @Accept  json
// @Produce  json
// @Param body body domain.ExecuteRequest true "Скрипт необходимый к выполнению"
// @Success 200 {object} domain.ScriptResp
// @Router /script/execute [POST].
func (c Script) Execute(req domain.ExecuteRequest) *domain.ScriptResp {
	return c.service.Execute(req)
}

// BatchExecute
// @Tags script
// @Summary Выполнить набор скриптов под конкретными ID
// @Description Возвращает результат выполнения скриптов и ошибок (если есть)
// @Accept  json
// @Produce  json
// @Param body body []domain.ExecuteByIdRequest true "Набор идентификаторов и аргументов"
// @Success 200 {array} domain.ScriptResp
// @Router /script/batch_execute [POST].
func (c Script) BatchExecute(req []domain.ExecuteByIdRequest) []domain.ScriptResp {
	return c.service.BatchExecute(req)
}

// BatchExecuteById
// @Tags script
// @Summary Выполнить набор скриптов под конкретными ID с идентичным аргументом для всех
// @Description Возвращает результат выполнения скриптов и ошибок (если есть)
// @Accept  json
// @Produce  json
// @Param body body domain.BatchExecuteByIdsRequest true "Набор идентификаторов и единый аргумент"
// @Success 200 {array} domain.ScriptResp
// @Router /script/batch_execute_by_ids [POST].
func (c Script) BatchExecuteById(req domain.BatchExecuteByIdsRequest) []domain.ScriptResp {
	return c.service.BatchExecuteById(req)
}
