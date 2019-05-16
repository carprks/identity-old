package identity

import "errors"

// Delete remove the identity
func (i Identity) Delete() (Identity, error) {
	if i.ID != "" && i.Email != "" && i.Phone != "" {
		getEntry, err := i.Retrieve()
		if err != nil {
			return Identity{}, err
		}
		resp, err := getEntry.DeleteEntry()
		if err != nil {
			return Identity{}, err
		}

		return resp, err
	}

	return Identity{}, errors.New("you need to have all the details to delete")
}
