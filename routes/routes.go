package routes

import (
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/grpc/endpoint"
	"github.com/integration-system/isp-kit/grpc/isp"
	"isp-script-service/controller"
)

type Controllers struct {
	Script controller.Script
}

func EndpointDescriptors() []cluster.EndpointDescriptor {
	return endpointDescriptors(Controllers{})
}

func Handler(wrapper endpoint.Wrapper, c Controllers) isp.BackendServiceServer {
	muxer := grpc.NewMux()
	for _, descriptor := range endpointDescriptors(c) {
		muxer.Handle(descriptor.Path, wrapper.Endpoint(descriptor.Handler))
	}
	return muxer
}

func endpointDescriptors(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:             "script/execute",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.Execute,
		},
		{
			Path:             "script/batch_execute",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.BatchExecute,
		},
		{
			Path:             "script/execute_by_id",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.ExecuteById,
		},
		{
			Path:             "script/batch_execute_by_ids",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.BatchExecuteById,
		},
	}
}
