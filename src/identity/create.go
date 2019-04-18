package identity

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (i Identity)Create() (Identity, error) {
	i.ID = i.createIdentifier()

	return i.CreateEntry()
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	i := Identity{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(IdentityResponse{
			Error: err,
		})
		return
	}

	jsonErr := json.Unmarshal(body, &i)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(IdentityResponse{
			Error: jsonErr,
		})
		return
	}

	resp, err := i.Create()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(IdentityResponse{
			Error: err,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(IdentityResponse{
		Identity: resp,
	})
	return
}