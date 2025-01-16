package tests_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/isp-kit/grpc/client"
	"github.com/txix-open/isp-kit/test"
	"github.com/txix-open/isp-kit/test/grpct"
	"isp-script-service/assembly"
	"isp-script-service/conf"
)

func TestMergeScriptTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &MergeScript{})
}

type MergeScript struct {
	suite.Suite
	test    *test.Test
	require *require.Assertions
}

func (s *MergeScript) TestMergeScriptsDuplicateId() {
	s.test, s.require = test.New(s.T())

	cli := s.initMockServer()

	cfg := conf.Remote{
		ScriptExecutionTimeoutMs: 10000,
		Scripts: conf.AllScripts{
			CommonScripts: []conf.ScriptDefinition{
				{
					Id:     "1",
					Script: "return true",
				},
			},
			CustomScripts: []conf.ScriptDefinition{
				{
					Id:     "1",
					Script: "return false",
				},
			},
		},
	}

	locator := assembly.NewLocator(s.test.Logger(), cli)
	_, err := locator.Handler(cfg)

	s.Require().Error(err)
}

func (s *MergeScript) TestMergeScriptsHappyPath() {
	s.test, s.require = test.New(s.T())

	cli := s.initMockServer()

	cfg := conf.Remote{
		ScriptExecutionTimeoutMs: 10000,
		Scripts: conf.AllScripts{
			CommonScripts: []conf.ScriptDefinition{
				{
					Id:     "1",
					Script: "return true",
				},
				{
					Id:     "2",
					Script: "return false",
				},
			},
			CustomScripts: []conf.ScriptDefinition{
				{
					Id:     "3",
					Script: "return true",
				},
			},
		},
	}

	locator := assembly.NewLocator(s.test.Logger(), cli)
	_, err := locator.Handler(cfg)

	s.Require().NoError(err)
}

func (s *MergeScript) initMockServer() *client.Client {
	server, cli := grpct.NewMock(s.test)
	server.Mock("mock", func(req bool) bool {
		return req
	})

	return cli
}
