package object_blueprint

import (
	"context"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/sha3"
)

/*
 * Takes in the objectId and the url, attempting to create or retrieve a record in object storage
 * and attach it to that object
 */
func StoreObjectFile(tx bun.IDB, objectId int64, data []byte) (*types.ObjectBlueprintStorage, error) {
	// we have our zip file in "bytes", now we need to hash these bytes and check if it already exists or not
	hasher := sha3.New256()
	hasher.Write(data)

	hash := hex.EncodeToString(hasher.Sum(nil))

	item, err := db.BlueprintObjectStorageGetByHash(tx, hash)

	if item == nil {
		// does not exist we can create
		item, err = db.BlueprintObjectStorageCreateFromBytes(tx, hash, data)

		if err != nil {
			return nil, err
		}
	}

	if item == nil {
		// to stop ide checks
		return nil, errors.New("item var within StoreObjectFile is still nil")
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
		ob := &types.ObjectBlueprintStorageRelation{
			ObjectBlueprintId:        objectId,
			ObjectBlueprintStorageId: item.Id,
		}
		_, err := tx.NewInsert().Model(ob).Returning("*").Exec(context.Background())

		if err != nil {
			return nil, err
		}

		// now refresh with relations
		item, err = db.BlueprintObjectStorageGetByHash(tx, hash)

		if err != nil {
			return nil, err
		}
	}

	return item, nil
}
