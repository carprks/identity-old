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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(fmt.Errorf("body err: %v", err))
		eErr := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if eErr != nil {
			fmt.Println(fmt.Errorf("body encode err: %v", eErr))
		}
		return
	}

	err = json.Unmarshal(body, &i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(fmt.Errorf("json err: %v", err))
		eErr := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if eErr != nil {
			fmt.Println(fmt.Errorf("json encode err: %v", eErr))
		}
		return
	}

	resp, err := i.Create()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(fmt.Errorf("create err: %v", err))
		eErr := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if eErr != nil {
			fmt.Println(fmt.Errorf("create encode err: %v", eErr))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(Response{
		Identity: resp,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("created err: %v", err))
	}
}
