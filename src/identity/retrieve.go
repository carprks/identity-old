package identity

// Retrieve get an identity
func (i Identity)Retrieve() (Identity, error) {
	if i.ID != ""{
		return i.RetrieveEntry()
	}

	return i.ScanEntry()
}