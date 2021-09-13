package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/flow"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/samvaughton/wpcommand/v2/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func loadBlueprintObjectRevisionsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["bpUuid"]

	bp, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	obj, err := db.BlueprintObjectGetAllRevisionsSafe(vars["objUuid"], bp.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(resp).Encode(obj)
}

func loadBlueprintObjectHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["bpUuid"]

	bp, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	revId, err := strconv.Atoi(vars["revId"])

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	obj, err := db.BlueprintObjectGetByUuidAndRevisionSafe(vars["objUuid"], revId, bp.Id)

	// now we can allow the file to be downloaded
	var sItem *types.ObjectBlueprintStorage

	// pluck the version one
	for _, storage := range obj.ObjectBlueprintStorage {
		sItem = storage
		break
	}

	if sItem == nil {
		// storage dont exist, attempt a pull
		go flow.VerifyAndStoreObjectFile(obj)
	}

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(resp).Encode(obj)
}

func createBlueprintObjectHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	payload, err := types.NewCreateObjectBlueprintPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	validationErrors := payload.Validate()

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	account := req.Context().Value("account").(*types.Account)

	vars := mux.Vars(req)

	blueprintSet, err := db.BlueprintGetByUuidSafe(vars["uuid"], account.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		json.NewEncoder(resp).Encode(map[string]string{
			"Status":  "NOT_FOUND",
			"Message": "Could not locate the blueprint",
		})
		return
	}

	object, err := flow.CreateObjectBlueprintFromCreatePayload(payload, blueprintSet.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]string{
			"Status":  "COULD_NOT_CREATE",
			"Message": "Something went wrong creating the object.",
		})
		return
	}

	json.NewEncoder(resp).Encode(object)
}

func loadBlueprintObjectFileHandler(resp http.ResponseWriter, req *http.Request) {
	//userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["bpUuid"]

	bp, err := db.BlueprintGetByUuid(key)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	revId, err := strconv.Atoi(vars["revId"])

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	obj, err := db.BlueprintObjectGetByUuidAndRevisionSafe(vars["objUuid"], revId, bp.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	// now we can allow the file to be downloaded
	var sItem *types.ObjectBlueprintStorage

	// pluck the version one
	for _, storage := range obj.ObjectBlueprintStorage {
		sItem = storage
		break
	}

	if sItem == nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	contentType := http.DetectContentType(sItem.File)
	resp.Header().Set("Content-Type", contentType)
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", fmt.Sprintf("%s-%s", obj.ExactName, obj.Version)))

	_, err = resp.Write(sItem.File)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func loadFileFromHashHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	file, err := db.BlueprintObjectStorageGetByHash(vars["hash"])

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	contentType := http.DetectContentType(file.File)
	resp.Header().Set("Content-Type", contentType)
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", file.Hash))

	_, err = resp.Write(file.File)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func createBlueprintObjectRevisionHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["bpUuid"]

	bp, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	obj, err := db.BlueprintObjectGetLatestRevisionByUuidSafe(vars["objUuid"], bp.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	payload, err := types.NewUpdatedVersionObjectBlueprintPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusContentMalformed, "Could not decode the payload", util.HttpEmptyErrors())
		return
	}

	validationErrors := payload.Validate()

	if validationErrors != nil {
		util.HttpErrorEncode(resp, util.HttpStatusValidationErrors, "Validation errors", validationErrors)
		return
	}

	object, err := flow.CreateObjectBlueprintRevisionFromNewVersionPayload(obj, payload)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInvalidPayload, err.Error(), util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(object)
}

func updateBlueprintObjectHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["bpUuid"]

	bp, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	obj, err := db.BlueprintObjectGetLatestRevisionByUuidSafe(vars["objUuid"], bp.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	payload, err := types.NewUpdateObjectBlueprintPayloadFromHttpRequest(req)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	validationErrors := payload.Validate()

	if validationErrors != nil {
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(map[string]interface{}{"Status": "VALIDATION_ERRORS", "Errors": validationErrors})

		return
	}

	obj.Name = payload.Name // assign

	// we want to update all objects with the same uuid for this
	res, err := db.Db.NewUpdate().Model(obj).Where("uuid = ?", obj.Uuid).Exec(context.Background())

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	ra, err := res.RowsAffected()

	if err != nil || ra != 1 {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(obj)
}

func deleteBlueprintObjectHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["bpUuid"]

	bp, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	obj, err := db.BlueprintObjectGetLatestRevisionByUuidSafe(vars["objUuid"], bp.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	// uuid is the same for all types of the same object
	_, err = db.Db.Query("DELETE FROM object_blueprints WHERE uuid = ?", obj.Uuid)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	json.NewEncoder(resp).Encode(obj)
}

func deleteBlueprintObjectRevisionHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	userAccount := req.Context().Value("userAccount").(*types.UserAccount)

	vars := mux.Vars(req)
	key := vars["bpUuid"]

	bp, err := db.BlueprintGetByUuidSafe(key, userAccount.AccountId)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	revId, err := strconv.Atoi(vars["revId"])

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	obj, err := db.BlueprintObjectGetByUuidAndRevisionSafe(vars["objUuid"], revId, bp.Id)

	if err != nil {
		log.Error(err)
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = db.Db.Query("DELETE FROM object_blueprints WHERE id = ?", obj.Id)

	if err != nil {
		log.Error(err)
		util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
		return
	}

	if obj.Active {
		// this object is deleted we need to make the latest one active
		latest, err := db.BlueprintObjectGetLatestRevisionByUuid(obj.Uuid)

		if err != nil {
			log.Error(err)
			util.HttpErrorEncode(resp, util.HttpStatusInternalServerError, "Something went wrong.", util.HttpEmptyErrors())
			return
		}

		// set all inactive
		_, err = db.Db.Query("UPDATE object_blueprints SET active = false WHERE uuid = ?", latest.Uuid)
		// make latest active
		_, err = db.Db.Query("UPDATE object_blueprints SET active = true WHERE id = ?", latest.Id)
	}

	json.NewEncoder(resp).Encode(obj)
}
