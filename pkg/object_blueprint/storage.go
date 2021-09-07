package object_blueprint

import (
	"encoding/hex"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

/*
 * Takes in the objectId and the url, attempting to create or retrieve a record in object storage
 * and attach it to that object
 */
func StoreObjectFile(objectId int64, data []byte) (*types.ObjectBlueprintStorage, error) {
	// we have our zip file in "bytes", now we need to hash these bytes and check if it already exists or not
	hasher := sha3.New256()
	hasher.Write(data)

	hash := hex.EncodeToString(hasher.Sum(nil))

	item, err := db.BlueprintObjectStorageGetByHash(hash)

	if item == nil {
		// does not exist we can create
		item, err = db.BlueprintObjectStorageCreateFromBytes(hash, data)

		if err != nil {
			log.Error(err)

			return nil, err
		}
	}

	relationExists := false
	for _, op := range item.ObjectBlueprints {
		if op.Id == objectId {
			relationExists = true
			break
		}
	}

	if relationExists == false {
		// item now exists we can attach
		_, err = db.Db.Query("INSERT INTO object_blueprint_storage_relations (object_blueprint_id, object_blueprint_storage_id) VALUES (?, ?)", objectId, item.Id)

		if err != nil {
			return nil, err
		}

		// now refresh with relations
		item, err = db.BlueprintObjectStorageGetByHash(hash)

		if err != nil {
			return nil, err
		}
	}

	return item, nil
}
