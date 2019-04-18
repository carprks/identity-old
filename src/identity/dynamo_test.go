package identity_test

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"main/src/identity"
	"testing"
)

func TestIdentity_CreateEntry(t *testing.T) {
	godotenv.Load()

	ident := identity.Identity{
		ID: "test",
		Company: false,
		Phone: "123",
		Email: "test@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1233",
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
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}
}

func TestIdentity_UpdateEntry(t *testing.T) {
	godotenv.Load()

	identOrig := identity.Identity{
		ID: "test",
		Company: false,
		Phone: "123",
		Email: "test@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1233",
				VehicleType: identity.VehicleTypeCar,
				Oversized: false,
			},
		},
	}
	identUpdate := identity.Identity{
		ID: "test",
		Company: true,
		Phone: "123",
		Email: "test@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1233",
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

func TestIdentity_ScanEntries(t *testing.T) {
	godotenv.Load()

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
	godotenv.Load()

	ident := identity.Identity{
		ID: "test",
		Company: true,
		Phone: "123",
		Email: "test@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1233",
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
	godotenv.Load()

	ident := identity.Identity{
		ID: "test",
		Company: false,
		Phone: "123",
		Email: "test@test.test",
		Registrations: []identity.Registration{
			{
				Plate: "test1233",
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