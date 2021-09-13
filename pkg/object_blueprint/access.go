package object_blueprint

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func GenerateStorageAccessUrl(object *types.ObjectBlueprint) (string, error) {
	var file *types.ObjectBlueprintStorage

	if object == nil {
		return "", nil
	}

	for _, item := range object.ObjectBlueprintStorage {
		file = item
		break
	}

	if file == nil {
		return "", errors.New("could not locate the file for this object")
	}

	return fmt.Sprintf("%s/storage/%s", config.Config.StorageHost, file.Hash), nil
}
