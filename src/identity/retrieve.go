package identity

import "errors"

// Retrieve get an identity
func (i Identity)Retrieve() (Identity, error) {
	if i.ID != ""{
		return i.RetrieveEntry()
	}

	ident, err := i.ScanEntry()
	if err != nil {
		return Identity{}, err
	}
	return ident.RetrievePlate(i)
}

// RetrievePlate for an identity identified by email/phone
func (i Identity)RetrievePlate(request Identity) (Identity, error) {
	reqPlate := request.Registrations[0]

	for _, plate := range i.Registrations {
		if plate.Plate == reqPlate.Plate {
			return i, nil
		}
	}

	return Identity{}, errors.New("no plate match")
}

func RetrieveAll() ([]Identity, error) {
	i := Identity{}
	return i.ScanEntries()
}