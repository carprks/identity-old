package identity

// Update the identity
func (i Identity)Update(n Identity) (Identity, error) {
	_, err := i.Retrieve()
	if err != nil {
		return Identity{}, err
	}

	return i.UpdateEntry(n)
}