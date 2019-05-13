package identity_test

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"main/src/identity"
	"os"
	"testing"
)

func TestRetrieveAll(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	create := []struct{
		identity.Identity
	}{
		{
			identity.Identity{
				Email: "test@test.test",
				Phone: "123",
				Registrations: []identity.Registration{
					{
						Plate: "test123",
						VehicleType: identity.VehicleTypeCar,
					},
				},
			},
		},
		{
			identity.Identity{
				Email: "test123@test.test",
				Phone: "456",
				Registrations: []identity.Registration{
					{
						Plate: "test456",
						VehicleType: identity.VehicleTypeBike,
					},
				},
			},
		},
	}

	idents := []identity.Identity{}
	for _, test := range create {
		resp, err := test.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("RetreiveAll Create Err: %v", err))
		}

		idents = append(idents, resp)
	}

	tests := []struct{
		expect []identity.Identity
		err error
	}{
		{
			expect: idents,
			err: nil,
		},
	}

	for _, test := range tests {
		resp, err := identity.RetrieveAll()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, resp)
	}

	for _, ident := range idents {
		_, err := ident.DeleteEntry()
		if err != nil {
			fmt.Println(fmt.Errorf("retriveall delete err: %v", err))
		}
	}
}

func TestIdentity_Retrieve(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	create := []struct{
		identity.Identity
	}{
		{
			identity.Identity{
				Email:   "test@test.test",
				Phone:   "123",
				Company: false,
				Registrations: []identity.Registration{
					{
						Plate:       "test123",
						VehicleType: identity.VehicleTypeCar,
						Oversized:   false,
					},
				},
			},
		},
	}

	tests := []struct{
		request identity.Identity
		expect identity.Identity
		err error
	}{
		{
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
		{
			request: identity.Identity{
				Email: "test@test.test",
				Registrations: []identity.Registration{
					{
						Plate: "123test",
					},
				},
			},
			expect: identity.Identity{},
			err: errors.New("no plate match"),
		},
		{
			request: identity.Identity{
				Email: "test@test.test",
			},
			expect: identity.Identity{},
			err: errors.New("need at least 1 plate"),
		},
	}

	// Create
	created := []identity.Identity{}
	for _, test := range create {
		res, err := test.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("retrieve create err: %v", err))
		}

		created = append(created, res)
	}

	// Retrieve
	for _, test := range tests {
		response, err := test.request.Retrieve()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}

	// Delete
	for _, test := range created {
		if test.ID != "" {
			_, err  := test.DeleteEntry()
			if err != nil {
				fmt.Println(fmt.Errorf("retrive delete err: %v", err))
			}
		}
	}
}