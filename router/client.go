package router

import (
	"fmt"
	"github.com/integration-system/isp-lib/backend"
	"github.com/integration-system/isp-lib/modules"
	log "github.com/integration-system/isp-log"
	"google.golang.org/grpc"
	md "google.golang.org/grpc/metadata"
	"isp-script-service/codes"
)

var (
	Client = backend.NewRxGrpcClient(
		backend.WithDialOptions(grpc.WithInsecure(), grpc.WithBlock()),
		backend.WithDialingErrorHandler(func(err error) {
			log.Warnf(codes.InvokeRouterServiceError, "Router client dialing err: %v", err)
		}),
	)
)

func Invoke(methodPath string, request interface{}, metadata map[string]interface{}) (interface{}, error) {
	newMd := md.New(map[string]string{})
	for key, info := range metadata {
		newMd.Set(key, fmt.Sprintf("%v", info))
	} //todo metadata

	var response interface{}
	if err := Client.Visit(func(c *backend.InternalGrpcClient) error {
		resp, err := c.InvokeWithDynamicStruct(
			methodPath,
			modules.RouterModuleId,
			request,
			//todo metadata
		)
		response = resp
		return err
	}); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}
