package conf_test

import (
	"testing"

	"github.com/txix-open/isp-kit/test/rct"
	"isp-script-service/conf"
)

func TestDefaultRemoteConfig(t *testing.T) {
	t.Parallel()
	rct.Test(t, "default_remote_config.json", conf.Remote{})
}