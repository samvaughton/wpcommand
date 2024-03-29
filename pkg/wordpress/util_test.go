package wordpress

import (
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func newTestExecutor(mockOutput string) execution.CommandExecutor {
	return &execution.DebugCommandExecutor{
		Site:       &types.Site{types.ApiSiteCore{Key: "test-site"}, types.ApiSiteCredentials{}},
		MockOutput: mockOutput,
	}
}
