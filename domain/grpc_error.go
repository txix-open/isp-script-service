package domain

type GrpcError struct {
	ErrorMessage string
	ErrorCode    string
	Details      []interface{}
}
