package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Create crete identity
func (i Identity) Create() (Identity, error) {
	i.ID = i.createIdentifier()

	if len(i.Registrations) >= 1 {
		return i.CreateEntry()
	}

	return Identity{}, errors.New("need at least 1 registration")
}

// Create http
func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	i := Identity{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	err = json.Unmarshal(body, &i)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	resp, err := i.Create()
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(Response{
		Identity: resp,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("created err: %v", err))
	}
}
