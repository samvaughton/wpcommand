package flow

import (
	"context"
	"errors"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/object_blueprint"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func DownloadAndVerifyZip(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("status code received is: %v, something went wrong", resp.StatusCode))
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(data)

	if contentType != "application/x-gzip" && contentType != "application/zip" {
		return nil, errors.New("provided file is not a zip file")
	}

	return data, nil
}

func CreateObjectBlueprintFromCreatePayload(payload *types.CreateObjectBlueprintPayload, blueprintSetId int64) (*types.ObjectBlueprint, error) {
	tx, err := db.Db.BeginTx(context.Background(), nil)

	// first we need to download the object
	data, err := DownloadAndVerifyZip(payload.Url)

	if err != nil {
		return nil, err
	}

	// file all looks ok

	object, err := db.BlueprintObjectCreateFromPayload(payload, blueprintSetId)

	if err != nil {
		tx.Rollback()

		return nil, err
	}

	// now we need to store the object
	_, err = object_blueprint.StoreObjectFile(object.Id, data)

	err = tx.Commit()

	if err != nil {
		log.Error(err)

		return nil, err
	}

	return object, nil
}

func CreateObjectBlueprintRevisionFromNewVersionPayload(object *types.ObjectBlueprint, payload *types.UpdatedVersionObjectBlueprintPayload) (*types.ObjectBlueprint, error) {
	tx, err := db.Db.BeginTx(context.Background(), nil)

	// first we need to download the object
	data, err := DownloadAndVerifyZip(object.OriginalObjectUrl)

	if err != nil {
		return nil, err
	}

	// file all looks ok
	// before we create set all revisions as inactive as the new revision will be the active one
	_, err = db.Db.Query("UPDATE object_blueprints SET active = false WHERE uuid = ?", object.Uuid)

	if err != nil {
		log.Error(err)
		tx.Rollback()

		return nil, err
	}

	newObj, err := db.BlueprintObjectCreateNewRevisionFromPayload(object.Uuid, payload)

	if err != nil {
		tx.Rollback()

		return nil, err
	}

	// now we need to store the object
	_, err = object_blueprint.StoreObjectFile(newObj.Id, data)

	if err != nil {
		tx.Rollback()

		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		log.Error(err)

		return nil, err
	}

	return newObj, nil
}

func VerifyAndStoreObjectFile(object *types.ObjectBlueprint) {
	data, err := DownloadAndVerifyZip(object.OriginalObjectUrl)

	if err != nil {
		log.Error(err)
		return
	}

	object_blueprint.StoreObjectFile(object.Id, data)
}
