package identity_test

import (
	"errors"
	"fmt"
	"github.com/carprks/identity/src/identity"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIdentity_Delete(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Sprintf("godotenv err: %v", err))
		}
	}

	create := []struct {
		request identity.Identity
		expect  identity.Identity
		err     error
	}{
		{
			request: identity.Identity{
				Email: "testdel@test.test",
				Phone: "123",
				Registrations: []identity.Registration{
					{
						Plate:       "test123",
						VehicleType: identity.VehicleTypeCar,
					},
				},
			},
			expect: identity.Identity{
				ID:    "46ff1881-b13b-5d59-8a0f-b653656d3b15",
				Email: "testdel@test.test",
				Phone: "123",
				Registrations: []identity.Registration{
					{
						Plate:       "test123",
						VehicleType: identity.VehicleTypeCar,
					},
				},
			},
			err: nil,
		},
		{
			request: identity.Identity{
				Email: "testdel1@test.test",
				Phone: "123",
				Registrations: []identity.Registration{
					{
						Plate:       "123",
						VehicleType: identity.VehicleTypeBike,
					},
				},
			},
			expect: identity.Identity{
				ID:    "241faf05-2d23-5088-845d-8782d84eee47",
				Email: "testdel1@test.test",
				Phone: "123",
				Registrations: []identity.Registration{
					{
						Plate:       "123",
						VehicleType: identity.VehicleTypeBike,
					},
				},
			},
			err: nil,
		},
	}
	for _, test := range create {
		resp, err := test.request.Create()
		correct := assert.IsType(t, test.err, err)
		if !correct {
			fmt.Println(fmt.Sprintf("delete create err: %v", err))
		}
		assert.Equal(t, test.expect, resp)
	}

	tests := []struct {
		request identity.Identity
		expect  identity.Identity
		err     error
	}{
		{
			request: identity.Identity{
				ID:    "46ff1881-b13b-5d59-8a0f-b653656d3b15",
				Email: "testdel@test.test",
				Phone: "123",
			},
			expect: identity.Identity{},
			err:    nil,
		},
		{
			request: identity.Identity{
				ID: "241faf05-2d23-5088-845d-8782d84eee47",
			},
			expect: identity.Identity{},
			err:    errors.New("you need to have all the details to delete"),
		},
	}
	for _, test := range tests {
		resp, err := test.request.Delete()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, resp)
	}

	_, err := identity.Identity{
		ID: "241faf05-2d23-5088-845d-8782d84eee47",
	}.DeleteEntry()
	if err != nil {
		fmt.Println(fmt.Sprintf("delete fail test: %v", err))
	}
}
