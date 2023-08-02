package tests_test

import (
	"context"
	"testing"

	"github.com/integration-system/isp-kit/grpc/client"
	"github.com/integration-system/isp-kit/test"
	"github.com/integration-system/isp-kit/test/grpct"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-script-service/assembly"
	"isp-script-service/conf"
	"isp-script-service/domain"
)

func TestApiScriptTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ApiScript{})
}

type ApiScript struct {
	suite.Suite
	test    *test.Test
	require *require.Assertions
	apiCli  *client.Client
}

func (s *ApiScript) SetupSuite() {
	s.test, s.require = test.New(s.T())

	cli := s.initMockServer()

	cfg := conf.Remote{
		ScriptExecutionTimeoutMs: 10000,
		Scripts: []conf.ScriptDefinition{
			{
				Id:     "1",
				Script: "return true",
			},
			{
				Id:     "2",
				Script: "return false",
			},
		},
	}

	locator := assembly.NewLocator(s.test.Logger(), cli)
	handler := locator.Handler(cfg)

	_, s.apiCli = grpct.TestServer(s.test, handler)
}

func (s *ApiScript) TestExternalInvoke() {
	req := domain.ExecuteRequest{
		Script: "return external.invoke('mock', true, {})",
		Arg:    struct{}{},
	}
	res := domain.ScriptResp{}
	err := s.apiCli.Invoke("script/script/execute").
		JsonRequestBody(req).
		JsonResponseBody(&res).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(res.Result, true)
}

func (s *ApiScript) TestExecuteHappyPath() {
	req := domain.ExecuteRequest{
		Script: "return arg.a",
		Arg: map[string]bool{
			"a": true,
		},
	}
	res := domain.ScriptResp{}
	err := s.apiCli.Invoke("script/script/execute").
		JsonRequestBody(req).
		JsonResponseBody(&res).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(res.Result, true)
}

func (s *ApiScript) TestExecuteByIdHappyPath() {
	req := domain.ExecuteByIdRequest{
		Id:  "1",
		Arg: struct{}{},
	}
	res := domain.ScriptResp{}
	err := s.apiCli.Invoke("script/script/execute_by_id").
		JsonRequestBody(req).
		JsonResponseBody(&res).
		Do(context.Background())
	s.Require().NoError(err)
	s.Require().Equal(res.Result, true)
}

func (s *ApiScript) TestExecuteByIdNotFound() {
	req := domain.ExecuteByIdRequest{
		Id:  "unknown",
		Arg: struct{}{},
	}
	res := domain.ScriptResp{}
	err := s.apiCli.Invoke("script/script/execute_by_id").
		JsonRequestBody(req).
		JsonResponseBody(&res).
		Do(context.Background())
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ApiScript) TestBatchExecuteHappyPath() {
	req := []domain.ExecuteByIdRequest{
		{
			Id:  "1",
			Arg: struct{}{},
		},
		{
			Id:  "2",
			Arg: struct{}{},
		},
	}
	res := make([]domain.ScriptResp, 0)
	err := s.apiCli.Invoke("script/script/batch_execute").
		JsonRequestBody(req).
		JsonResponseBody(&res).
		Do(context.Background())
	s.Require().NoError(err)
	expected := []domain.ScriptResp{
		{Result: true, Error: nil},
		{Result: false, Error: nil},
	}
	s.Require().ElementsMatch(expected, res)
}

func (s *ApiScript) TestBatchExecuteMixedNotFound() {
	req := []domain.ExecuteByIdRequest{
		{
			Id:  "1",
			Arg: struct{}{},
		},
		{
			Id:  "unknown",
			Arg: struct{}{},
		},
	}
	res := make([]domain.ScriptResp, 0)
	err := s.apiCli.Invoke("script/script/batch_execute").
		JsonRequestBody(req).
		JsonResponseBody(&res).
		Do(context.Background())
	s.Require().NoError(err)
	expected := []domain.ScriptResp{
		{Result: true, Error: nil},
		{Result: nil, Error: &domain.Error{
			Type:        domain.ErrorCompile,
			Description: "not defined script for id unknown",
		}},
	}
	s.Require().ElementsMatch(expected, res)
}

func (s *ApiScript) TestBatchExecuteByIdHappyPath() {
	req := domain.BatchExecuteByIdsRequest{
		Ids: []string{"1", "2"},
		Arg: struct{}{},
	}
	res := make([]domain.ScriptResp, 0)
	err := s.apiCli.Invoke("script/script/batch_execute_by_ids").
		JsonRequestBody(req).
		JsonResponseBody(&res).
		Do(context.Background())
	s.Require().NoError(err)
	expected := []domain.ScriptResp{
		{Result: true, Error: nil},
		{Result: false, Error: nil},
	}
	s.Require().ElementsMatch(expected, res)
}

func (s *ApiScript) initMockServer() *client.Client {
	server, cli := grpct.NewMock(s.test)
	server.Mock("mock", func(req bool) bool {
		return req
	})

	return cli
}
