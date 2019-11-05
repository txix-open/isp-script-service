package helper

import (
	"isp-script-service/controller"
	"isp-script-service/domain"
)

type ScriptHandlers struct {
	ExecuteById func(request domain.ExecuteByIdRequest) (*domain.ScriptResp, error) `method:"execute_by_id" group:"script"`
	Execute     func(request domain.ExecuteRequest) *domain.ScriptResp              `method:"execute" group:"script"`
}

func GetTaskHandler() *ScriptHandlers {
	return &ScriptHandlers{
		ExecuteById: controller.Script.ExecuteById,
		Execute:     controller.Script.Execute,
	}
}

func GetAllHandlers() []interface{} {
	return []interface{}{
		GetTaskHandler(),
	}
}
