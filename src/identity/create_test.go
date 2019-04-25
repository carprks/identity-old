package identity_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"main/src/identity"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIdentity_Create(t *testing.T) {
	fmt.Println(fmt.Printf("Environ: %v", os.Environ()))

	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

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
		correct := assert.IsType(t, test.err, err)
		if !correct {
			fmt.Println(fmt.Errorf("create entry: %v", err))
		}
		assert.Equal(t, test.expect, response)
	}

	for _, test := range tests {
		if test.expect.ID != "" {
			_, err := test.expect.DeleteEntry()
			if err != nil {
				fmt.Println(fmt.Errorf("create delete entry: %v", err))
			}
		}
	}
}

func TestCreate(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	tests := []struct{
		request identity.Identity
		expect identity.Response
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
			expect: identity.Response{
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
			expect: identity.Response{
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
		i := identity.Response{}
		err := json.Unmarshal(body, &i)
		if err != nil {
			fmt.Println(fmt.Errorf("create json error: %v", err))
		}
		assert.Equal(t, test.expect, i)
	}

	for _, test := range tests {
		_, err := test.expect.Identity.DeleteEntry()
		if err != nil {
			fmt.Println(fmt.Errorf("create delete http: %v", err))
		}
	}
}