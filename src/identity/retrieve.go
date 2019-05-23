package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

// Retrieve get an identity
func (i Identity) Retrieve() (Identity, error) {
	if i.ID != "" {
		return i.RetrieveEntry()
	}

	ident, err := i.ScanEntry()
	if err != nil {
		return Identity{}, err
	}
	return ident.RetrievePlate(i)
}

// RetrievePlate for an identity identified by email/phone
func (i Identity) RetrievePlate(request Identity) (Identity, error) {
	reqPlate := request.Registrations[0]

	for _, plate := range i.Registrations {
		if plate.Plate == reqPlate.Plate {
			return i, nil
		}
	}

	return Identity{}, errors.New("no plate match")
}

// RetrieveAll get all the identities
func RetrieveAll() ([]Identity, error) {
	i := Identity{}
	return i.ScanEntries()
}

// Retrieve http get the identity
func Retrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "identityID")
	i := Identity{
		ID: id,
	}

	ident, err := i.Retrieve()
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(
		Response{
			Identity: ident,
		})
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve response encoder err: %v", err))
	}
}

// RetrieveAllIdentities http get all the identities
func RetrieveAllIdentities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	identities, err := RetrieveAll()
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{
		Identities: identities,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieveall response encoder error: %v", err))
		return
	}
}
