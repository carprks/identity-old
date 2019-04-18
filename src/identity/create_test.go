package identity_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"main/src/identity"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIdentity_Create(t *testing.T) {
	godotenv.Load()

	tests := []struct{
		request identity.Identity
		expect identity.Identity
		err error
	}{
		{
			request: identity.Identity{
				Email: "test@test.test",
				Phone: "123",
				Company: false,
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized: false,
					},
				},
			},
			expect: identity.Identity{
				ID: "fde8ffe8-75c6-5448-b44e-b4c81526a1eb",
				Email: "test@test.test",
				Phone: "123",
				Company: false,
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized: false,
					},
				},
			},
			err: nil,
		},
		{
			request: identity.Identity{
				Email: "test@test.test",
				Phone: "123",
				Company: false,
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized: false,
					},
				},
			},
			expect: identity.Identity{},
			err: errors.New("identity already exists"),
		},
		{
			request: identity.Identity{
				Email: "test@test.test",
				Phone: "123",
				Company: true,
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized: false,
					},
				},
			},
			expect: identity.Identity{
				ID: "780d270b-bb70-5ae9-96af-b7803c3c7b62",
				Email: "test@test.test",
				Phone: "123",
				Company: true,
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized: false,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		response, err := test.request.Create()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}

	for _, test := range tests {
		test.expect.DeleteEntry()
	}
}

func TestCreate(t *testing.T) {
	tests := []struct{
		request identity.Identity
		expect identity.IdentityResponse
		err error
	}{
		{
			request: identity.Identity{
				Email: "test@test.test",
				Phone: "123",
				Company: false,
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized: false,
					},
				},
			},
			expect: identity.IdentityResponse{
				Identity: identity.Identity{
					ID: "fde8ffe8-75c6-5448-b44e-b4c81526a1eb",
					Email: "test@test.test",
					Phone: "123",
					Company: false,
					Registrations: []identity.Registration{
						{
							Plate: "test123",
							VehicleType: identity.VehicleTypeCar,
							Oversized: false,
						},
					},
				},
			},
			err: nil,
		},
		{
			request: identity.Identity{
				Email: "test@test.test",
				Phone: "123",
				Company: true,
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized: false,
					},
				},
			},
			expect: identity.IdentityResponse{
				Identity: identity.Identity{
					ID: "780d270b-bb70-5ae9-96af-b7803c3c7b62",
					Email: "test@test.test",
					Phone: "123",
					Company: true,
					Registrations: []identity.Registration{
						{
							Plate: "test123",
							VehicleType: identity.VehicleTypeCar,
							Oversized: false,
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		jpr, _ := json.Marshal(test.request)
		request, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jpr))
		response := httptest.NewRecorder()
		identity.Create(response, request)
		assert.Equal(t, 201, response.Code)

		body, _ := ioutil.ReadAll(response.Body)
		i := identity.IdentityResponse{}
		json.Unmarshal(body, &i)
		assert.Equal(t, test.expect, i)
	}

	for _, test := range tests {
		test.expect.Identity.DeleteEntry()
	}
}