package helper

import (
	"isp-script-service/controller"
	"isp-script-service/domain"
)

type ScriptHandlers struct {
	ExecuteById       func(request domain.ExecuteByIdRequest) (*domain.ScriptResp, error) `method:"execute_by_id" group:"script"`
	Execute           func(request domain.ExecuteRequest) *domain.ScriptResp              `method:"execute" group:"script"`
	BatchExecute      func(request []domain.ExecuteByIdRequest) []domain.ScriptResp       `method:"batch_execute" group:"script"`
	BatchExecuteByIds func(request domain.BatchExecuteByIdsRequest) []domain.ScriptResp   `method:"batch_execute_by_ids" group:"script"`
}

func GetTaskHandler() *ScriptHandlers {
	return &ScriptHandlers{
		ExecuteById:       controller.Script.ExecuteById,
		Execute:           controller.Script.Execute,
		BatchExecute:      controller.Script.BatchExecute,
		BatchExecuteByIds: controller.Script.BatchExecuteById,
	}
}

func GetAllHandlers() []interface{} {
	return []interface{}{
		GetTaskHandler(),
	}
}
