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

	id := chi.URLParam(r, "identityId")
	i := Identity{
		ID: id,
	}

	req := Identity{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	resp, err := i.Update(req)
	if err != nil {
		ErrorResponse(w, err)
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
