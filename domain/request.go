package domain

type ExecuteByIdRequest struct {
	Id  string `valid:"required"`
	Arg interface{}
}

type ExecuteRequest struct {
	Script string `valid:"required"`
	Arg    interface{}
}

type BatchExecuteByIdsRequest struct {
	Ids []string `valid:"required"`
	Arg interface{}
}
