// nolint:ireturn
package routes

import (
	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"isp-script-service/controller"
)

type Controllers struct {
	Script controller.Script
}

func EndpointDescriptors() []cluster.EndpointDescriptor {
	return endpointDescriptors(Controllers{})
}

func Handler(wrapper endpoint.Wrapper, c Controllers) *grpc.Mux {
	muxer := grpc.NewMux()
	for _, descriptor := range endpointDescriptors(c) {
		muxer.Handle(descriptor.Path, wrapper.Endpoint(descriptor.Handler))
	}
	return muxer
}

func endpointDescriptors(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:             "script/script/execute",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.Execute,
		},
		{
			Path:             "script/script/batch_execute",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.BatchExecute,
		},
		{
			Path:             "script/script/execute_by_id",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.ExecuteById,
		},
		{
			Path:             "script/script/batch_execute_by_ids",
			Inner:            true,
			UserAuthRequired: false,
			Handler:          c.Script.BatchExecuteById,
		},
	}
}
