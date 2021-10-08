package wordpress

import (
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func newTestExecutor(mockOutput string) execution.CommandExecutor {
	return &execution.DebugCommandExecutor{
		Site:       &types.Site{Key: "test-site"},
		MockOutput: mockOutput,
	}
}
