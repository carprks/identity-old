package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Create crete identity
func (i Identity)Create() (Identity, error) {
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
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if err != nil {
			fmt.Println(fmt.Errorf("body err: %v", err))
		}
		return
	}

	jsonErr := json.Unmarshal(body, &i)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Response{
			Error: jsonErr,
		})
		if err != nil {
			fmt.Println(fmt.Errorf("json Err: %v", err))
		}
		return
	}

	resp, err := i.Create()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if err != nil {
			fmt.Println(fmt.Errorf("create Err: %v", err))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(Response{
		Identity: resp,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("created Err: %v", err))
	}
	return
}