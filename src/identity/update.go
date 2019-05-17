package identity

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
)

// Update the identity
func (i Identity) Update(n Identity) (Identity, error) {
	_, err := i.Retrieve()
	if err != nil {
		return Identity{}, err
	}

	return i.UpdateEntry(n)
}

// Update http update the identity
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "identityID")
	i := Identity{
		ID: id,
	}

	req := Identity{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(fmt.Sprintf("update body err: %v", err))
		eErr := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if eErr != nil {
			fmt.Println(fmt.Sprintf("update encode body err: %v", eErr))
		}
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(fmt.Sprintf("update unmarshall err: %v", err))
		eErr := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if eErr != nil {
			fmt.Println(fmt.Sprintf("update unmarshall encode err: %v", eErr))
		}
		return
	}
	resp, err := i.Update(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(fmt.Sprintf("update err: %v", err))
		eErr := json.NewEncoder(w).Encode(Response{
			Error: err,
		})
		if eErr != nil {
			fmt.Println(fmt.Sprintf("update encode err: %v", eErr))
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(Response{
		Identity: resp,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("update worked err: %v", err))
	}
}
