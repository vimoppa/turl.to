package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vimoppa/turl.to/internal/app"
	"github.com/vimoppa/turl.to/internal/storage"
)

type createURLPayload struct {
	URL string `json:"url"`
}

// AnyURLs retreives all the URL resources avaliable.
func AnyURLs(s storage.Accessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := app.GetAllRecords(s)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := SuccessResponse{
			Code:   http.StatusOK,
			Status: "Success",
			Data:   result,
		}

		RespondWithJSON(w, http.StatusOK, response)
	}
}

// CreateURL creates a new URL resource and returns the short URL.
func CreateURL(s storage.Accessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload createURLPayload

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&payload); err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		hash := app.GenerateHash(payload.URL)
		if exists := s.LookUp(hash); !exists {
			if err := s.WriteOnce(hash, payload.URL); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Failed to create url resource")
				return
			}
		}

		response := SuccessResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Successfully Created URL Resource",
			Data: app.RecordsItem{
				Hash:    hash,
				LongURL: payload.URL,
			},
		}

		RespondWithJSON(w, http.StatusOK, response)
	}
}

// FindOneURL find a long URL matching the short URL.
func FindOneURL(s storage.Accessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		result, err := s.ReadOne(vars["hash"])
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := SuccessResponse{
			Code:    http.StatusOK,
			Status:  "Success",
			Message: "Successfully found matching resource",
			Data:    result,
		}

		if result == "" {
			response.Code = http.StatusNotFound
			response.Status = "Not Found"
			response.Message = "Resource does not exist"
		}

		RespondWithJSON(w, http.StatusOK, response)
	}
}
