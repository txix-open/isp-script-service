package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/client"
	md "google.golang.org/grpc/metadata"
)

type Router struct {
	cli *client.Client
}

func NewRouter(cli *client.Client) Router {
	return Router{
		cli: cli,
	}
}

func (i Router) Invoke(methodPath string, request interface{}, metadata map[string]interface{}) (interface{}, error) {
	ctx := context.Background()
	newMd := md.New(nil)
	for key, val := range metadata {
		newMd.Set(key, fmt.Sprint(val))
	}
	ctx = md.NewOutgoingContext(ctx, newMd)

	var response interface{}

	err := i.cli.Invoke(methodPath).JsonRequestBody(request).JsonResponseBody(&response).Do(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "invoke")
	}

	return response, nil
}
