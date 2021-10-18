package registry

import (
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/wordpress"
	"strconv"
	"strings"
)

func GetHousecleaningCommand(site *types.Site, config map[string]interface{}) pipeline.SiteCommand {
	return &pipeline.DynamicArgsCommand{
		Name: CmdWpHousecleaning,
		GetArgs: func(pipeline *pipeline.SiteCommandPipeline) ([]string, error) {
			// Delete any default posts
			posts, err := wordpress.GetSitePostAndPages(pipeline.Executor)

			if err != nil {
				return []string{}, err
			}

			// transform into ID set for matching ones
			var deleteIds []string
			postsToBeDeleted := map[string]bool{
				"hello-world":    true,
				"sample-page":    true,
				"privacy-policy": true,
			}

			for _, post := range posts {
				// postId check to make sure we arent just deleting something like a privacy policy that has been
				// actually created
				if postsToBeDeleted[post.Name] && post.Id <= 3 {
					deleteIds = append(deleteIds, strconv.Itoa(post.Id))
				}
			}

			if len(deleteIds) > 0 {
				return []string{fmt.Sprintf("wp post delete %s --force", strings.Join(deleteIds, " "))}, nil
			}

			return []string{}, nil
		},
	}
}
