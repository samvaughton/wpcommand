package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/wordpress"
	"strconv"
	"strings"
)

func GetAcfSyncFieldsCommand(site *types.Site) pipeline.SiteCommand {
	return &pipeline.SimplePipelineCommand{
		Name: CmdWpSyncAcfFields,
		Commands: []pipeline.SiteCommand{
			// delete ACF field groups, ready for re-sync
			&pipeline.DynamicArgsCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpSyncAcfFields, "delete-existing"),
				GetArgs: func(pipeline *pipeline.SiteCommandPipeline) ([]string, error) {
					// now we need to delete all lazy blocks and import the new ones
					// Delete any default posts
					var ids []string

					acfFieldGroups, err := wordpress.GetSiteAcfFieldGroups(pipeline.Executor)

					if err != nil {
						return []string{}, err
					}

					for _, block := range acfFieldGroups {
						ids = append(ids, strconv.Itoa(block.Id))
					}

					if len(ids) > 0 {
						return []string{fmt.Sprintf("wp post delete %s --force", strings.Join(ids, " "))}, nil
					}

					return []string{}, nil // no execution
				},
			},
			&pipeline.SimpleCommand{
				Name: fmt.Sprintf("%s.%s", CmdWpSyncAcfFields, "sync"),
				Args: []string{"wp eval-file /opt/bitnami/eval-index.php acf-field-group-sync"},
			},
		},
	}
}
