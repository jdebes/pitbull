package api

import (
	"encoding/json"
	"net/http"
)

func UnmarshalRequest(w http.ResponseWriter, r *http.Request, model RequestModel) error {
	err := json.NewDecoder(r.Body).Decode(model)
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return err
	}

	err = model.Valid()
	if err != nil {
		Error(w, err, http.StatusBadRequest)
		return err
	}

	return nil
}
