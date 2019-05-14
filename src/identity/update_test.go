package identity

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIdentity_Update(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	create := []struct {
		ident Identity
		err   error
	}{
		{
			ident: Identity{
				Email:   "test@test.test",
				Phone:   "123",
				Company: false,
				Registrations: []Registration{
					{
						Plate:       "test123",
						VehicleType: VehicleTypeCar,
						Oversized:   false,
					},
				},
			},
			err: nil,
		},
	}
	created := Identity{}
	for _, test := range create {
		created, _ = test.ident.Create()
	}

	tests := []struct {
		request Identity
		expect  Identity
		err     error
	}{
		{
			request: Identity{
				ID:    created.ID,
				Email: "testUpdate@test.test",
			},
			expect: Identity{
				ID:      "fde8ffe8-75c6-5448-b44e-b4c81526a1eb",
				Email:   "testUpdate@test.test",
				Phone:   "123",
				Company: false,
				Registrations: []Registration{
					{
						VehicleType: VehicleTypeCar,
						Oversized:   false,
						Plate:       "test123",
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		resp, err := created.Update(test.request)
		correct := assert.IsType(t, test.err, err)
		if !correct {
			fmt.Println(fmt.Errorf("update err: %v", err))
		}
		assert.Equal(t, test.expect, resp)
	}

	_, err := created.DeleteEntry()
	if err != nil {
		fmt.Println(fmt.Errorf("delete update err: %v", err))
	}
}
