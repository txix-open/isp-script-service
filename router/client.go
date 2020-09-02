package router

import (
	"fmt"

	"github.com/integration-system/isp-lib/v2/backend"
	"google.golang.org/grpc"
	md "google.golang.org/grpc/metadata"
)

var Client = backend.NewRxGrpcClient(
	backend.WithDialOptions(grpc.WithInsecure()),
)

func Invoke(methodPath string, request interface{}, metadata map[string]interface{}) (interface{}, error) {
	newMd := md.New(nil)
	for key, val := range metadata {
		newMd.Set(key, fmt.Sprint(val))
	}

	var response interface{}
	err := Client.Invoke(methodPath, -3, request, &response, backend.WithMetadata(newMd))
	if err != nil {
		return nil, err
	}

	return response, nil
}
