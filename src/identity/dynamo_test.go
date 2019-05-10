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

func TestIdentity_CreateEntry(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	ident := identity.Identity{
		ID: "testDynamo",
		Company: false,
		Phone: "123",
		Email: "testDynamo@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1234",
				VehicleType: identity.VehicleTypeCar,
				Oversized: false,
			},
		},
	}

	tests := []struct{
		request identity.Identity
		expect identity.Identity
		err error
	}{
		{
			request: ident,
			expect: ident,
			err: nil,
		},
	}

	for _, test := range tests {
		response, err := test.request.CreateEntry()
		if err != nil {
			fmt.Println(fmt.Errorf("dynamo create: %v", err))
		}
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}
}

func TestIdentity_UpdateEntry(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	identOrig := identity.Identity{
		ID: "testDynamo",
		Company: false,
		Phone: "123",
		Email: "testDynamo@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1234",
				VehicleType: identity.VehicleTypeCar,
				Oversized: false,
			},
		},
	}
	identUpdate := identity.Identity{
		ID: "testDynamo",
		Company: true,
		Phone: "123",
		Email: "testDynamo@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1234",
				VehicleType: identity.VehicleTypeCar,
				Oversized: false,
			},
		},
	}

	tests := []struct{
		request identity.Identity
		expect identity.Identity
		err error
	}{
		{
			request: identOrig,
			expect: identUpdate,
			err: nil,
		},
	}

	for _, test := range tests {
		response, err := test.request.UpdateEntry(identUpdate)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}
}

func TestIdentity_ScanEntry(t *testing.T) {
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
				Email: "testDynamo@test.test",
			},
			expect: identity.Identity{},
			err: errors.New("need at least 1 plate"),
		},
		{
			request: identity.Identity{
				Email: "testDynamo@test.test",
				Registrations: []identity.Registration{
					{
						Plate: "test1234",
					},
				},
			},
			expect: identity.Identity{
				ID:      "testDynamo",
				Company: true,
				Phone:   "123",
				Email:   "testDynamo@test.test",
				Registrations: []identity.Registration{
					{
						Plate:       "test1234",
						VehicleType: identity.VehicleTypeCar,
						Oversized:   false,
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		response, err := test.request.ScanEntry()
		correct := assert.IsType(t, test.err, err)
		if !correct {
			fmt.Println(fmt.Errorf("scan test err: %v", err))
		}
		assert.Equal(t, test.expect, response)
	}
}

func TestIdentity_ScanEntries(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	tests := []struct{
		request identity.Identity
		expect int
		err error
	}{
		{
			request: identity.Identity{},
			expect: 1,
			err: nil,
		},
	}

	for _, test := range tests {
		response, err := test.request.ScanEntries()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, len(response))
	}
}

func TestIdentity_RetrieveEntry(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	ident := identity.Identity{
		ID: "testDynamo",
		Company: true,
		Phone: "123",
		Email: "testDynamo@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1234",
				VehicleType: identity.VehicleTypeCar,
				Oversized: false,
			},
		},
	}

	tests := []struct{
		request identity.Identity
		expect identity.Identity
		err error
	}{
		{
			request: ident,
			expect: ident,
			err: nil,
		},
	}

	for _, test := range tests {
		response, err := test.request.RetrieveEntry()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}
}

func TestIdentity_DeleteEntry(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godotenv err: %v", err))
		}
	}

	ident := identity.Identity{
		ID: "testDynamo",
		Company: false,
		Phone: "123",
		Email: "testDynamo@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1234",
				VehicleType: identity.VehicleTypeCar,
				Oversized: false,
			},
		},
	}

	tests := []struct{
		request identity.Identity
		expect identity.Identity
		err error
	}{
		{
			request: ident,
			expect: identity.Identity{},
			err: nil,
		},
	}

	for _, test := range tests {
		response, err := test.request.DeleteEntry()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}
}

func TestScanAll(t *testing.T) {
	if os.Getenv("AWS_DB_TABLE") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println(fmt.Errorf("godoterr: %v", err))
		}
	}

	idents := []identity.Identity{}
	for i := 0; i < 10; i++ {
		ident := identity.Identity{
			ID:    fmt.Sprintf("testDyamo%d", i),
			Email: fmt.Sprintf("test%d@test.test", i),
			Phone: fmt.Sprintf("%d", i),
		}
		idents = append(idents, ident)

		_, err := ident.CreateEntry()
		if err != nil {
			fmt.Println(fmt.Errorf("scan create err: %v", err))
		}
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
		resp, err := identity.ScanAll()
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, resp)
	}

	for _, ident := range idents {
		_, err := ident.DeleteEntry()
		if err != nil {
			fmt.Println(fmt.Errorf("delete scan err: %v", err))
		}
	}
}