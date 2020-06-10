package domain

type ExecuteByIdRequest struct {
	Id  string `valid:"required~Required"`
	Arg interface{}
}

type ExecuteRequest struct {
	Script string `valid:"required~Required"`
	Arg    interface{}
}

type BatchExecuteByIdsRequest struct {
	Ids []string `valid:"required~Required"`
	Arg interface{}
}
