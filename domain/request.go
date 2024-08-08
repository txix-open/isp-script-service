package domain

type ExecuteByIdRequest struct {
	Id  string `validate:"required"`
	Arg interface{}
}

type ExecuteRequest struct {
	Script string `validate:"required"`
	Arg    interface{}
}

type BatchExecuteByIdsRequest struct {
	Ids []string `validate:"required"`
	Arg interface{}
}
