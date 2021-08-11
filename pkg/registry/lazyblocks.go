package registry

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/wordpress"
	"strconv"
	"strings"
)

func GetLazyblocksSyncCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdWpSyncLazyblocks,
		Commands: []pipeline.SiteCommand{
			// delete old lazy blocks
			&pipeline.DynamicArgsCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpSyncLazyblocks, "delete-existing"),
				GetArgs: func(pipeline *pipeline.SiteCommandPipeline) ([]string, error) {
					// now we need to delete all lazy blocks and import the new ones
					// Delete any default posts
					var lazyIds []string
					lazyblocks, err := wordpress.GetSiteLazyblocksPosts(pipeline.Executor)

					if err != nil {
						return []string{}, errors.New("failed to get lazy block list, cannot sync")
					}

					for _, block := range lazyblocks {
						lazyIds = append(lazyIds, strconv.Itoa(block.Id))
					}

					if len(lazyIds) > 0 {
						return []string{fmt.Sprintf("wp post delete %s --force", strings.Join(lazyIds, " "))}, nil
					}

					return []string{}, nil
				},
			},
			// import lazy blocks
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpSyncLazyblocks, "import"),
				Args: []string{"wp eval-file /opt/bitnami/eval-index.php lazyblocks-import"},
			},
		},
	}
}
