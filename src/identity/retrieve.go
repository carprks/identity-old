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

// RetrieveUnknown gets the identity from an email and a plate
func (i Identity) RetrieveUnknown(request Identity) (Identity, error) {
	return Identity{}, nil
}

// RetrieveUnknown gets the identify from an unknown id
func RetrieveUnknown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := chi.URLParam(r, "email")
	plateReq := chi.URLParam(r, "plate")

	ident := Identity{
		Email: email,
		Registrations: []Registration{
			{
				Plate: plateReq,
			},
		},
	}
	resp, err := ident.ScanEntry()
	if err != nil {
		ErrorResponse(w, err)
	}

	found := false
	if len(resp.Registrations) >= 1 {
		for _, plate := range resp.Registrations {
			if plate.Plate == plateReq {
				found = true
			}
		}
	}

	if !found {
		ErrorResponse(w, errors.New("no plate match to email"))
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{
		Identity: resp,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("retrieve login response encoder err: %v", err))
	}
}

// Retrieve http get the identity
func Retrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "identityId")
	i := Identity{
		ID: id,
	}

	ident, err := i.Retrieve()
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{
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
