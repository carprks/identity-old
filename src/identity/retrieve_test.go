package identity_test

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"main/src/identity"
	"testing"
)

func TestIdentity_Retrieve(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(fmt.Errorf("godotenv err: %v", err))
	}

	tests := []struct{
		create identity.Identity
		request identity.Identity
		expect identity.Identity
		err error
	}{
		{
			create: identity.Identity{
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
			request: identity.Identity{
				Email: "test@test.test",
				Registrations: []identity.Registration{
					{
						Plate: "test123",
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
	}

	for _, test := range tests {
		_, err := test.create.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("retrieve create err: %v", err))
		}

		response, err := test.request.Retrieve()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}
}