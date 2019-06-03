package identity

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
)

// Response struct
type Response struct {
	Identities []Identity `json:"identities,omitempty"`
	Identity   Identity   `json:"identity,omitempty"`
	Error      error      `json:"error,omitempty"`
}

// Identity struct
type Identity struct {
	ID            string         `json:"id,omitempty"`
	Email         string         `json:"email,omitempty"`
	Phone         string         `json:"phone,omitempty"`
	Company       bool           `json:"company,omitempty"`
	Registrations []Registration `json:"registrations,omitempty"`
}

// Registration struct
type Registration struct {
	VehicleType VehicleType `json:"vehicleType"`
	Oversized   bool        `json:"oversized,omitempty"`
	Plate       string      `json:"plate"`
}

// VehicleType self explanatory
type VehicleType string

const (
	// VehicleTypeBike Motorbike
	VehicleTypeBike VehicleType = "Bike"
	VehicleTypeMotorBike VehicleType = "Motorbike"

	// VehicleTypeCar Car
	VehicleTypeCar VehicleType = "Car"

	// VehicleTypeVan Van
	VehicleTypeVan VehicleType = "Van"

	// VehicleTypeUnknown Unknown
	VehicleTypeUnknown VehicleType = "Unknown"
)

func (v VehicleType) convertToString() string {
	switch v {
	case VehicleTypeBike:
		return string(v)
	case VehicleTypeMotorBike:
		return string(v)
	case VehicleTypeCar:
		return string(v)
	case VehicleTypeVan:
		return string(v)
	}

	return string(VehicleTypeUnknown)
}

// GetVehicleType returns the vehicle type
func GetVehicleType(s string) VehicleType {
	// Title case the type
	s = strings.Title(s)

	if strings.Contains(string(VehicleTypeCar), s) {
		return VehicleTypeCar
	}

	if strings.Contains(string(VehicleTypeVan), s) {
		return VehicleTypeVan
	}

	if strings.Contains(string(VehicleTypeBike), s) {
		return VehicleTypeBike
	} else if strings.Contains(string(VehicleTypeMotorBike), s) {
		return VehicleTypeMotorBike
	}

	return VehicleTypeUnknown
}

func (i Identity) createIdentifier() string {
	if i.Company {
		u := uuid.NewV5(uuid.NamespaceURL, fmt.Sprintf("https://identity.carprk.com/company/%s:%s", i.Email, i.Phone))
		return u.String()
	}

	u := uuid.NewV5(uuid.NamespaceURL, fmt.Sprintf("https://identity.carprk.com/user/%s:%s", i.Email, i.Phone))
	return u.String()
}

// ErrorResponse default response
func ErrorResponse(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Println(fmt.Sprintf("err: %v", e))
	eErr := json.NewEncoder(w).Encode(Response{
		Error: e,
	})
	if eErr != nil {
		fmt.Println(fmt.Sprintf("encode err: %v", eErr))
	}
}