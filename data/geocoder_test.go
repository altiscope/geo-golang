package data_test

import (
	"context"
	"github.com/altiscope/geo-golang"
	"github.com/altiscope/geo-golang/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	addressFixture = geo.Address{
		FormattedAddress: "64 Elizabeth Street, Melbourne, Victoria 3000, Australia",
	}
	locationFixture = geo.Location{
		Lat: -37.814107,
		Lng: 144.96328,
	}
	geocoder = data.Geocoder(
		data.AddressToLocation{
			addressFixture: locationFixture,
		},
		data.LocationToAddress{
			locationFixture: addressFixture,
		},
	)
)

func TestGeocode(t *testing.T) {
	ctx := context.TODO()
	location, err := geocoder.Geocode(ctx, addressFixture.FormattedAddress)
	assert.NoError(t, err)
	assert.Equal(t, geo.Location{Lat: -37.814107, Lng: 144.96328}, *location)
}

func TestReverseGeocode(t *testing.T) {
	ctx := context.TODO()
	address, err := geocoder.ReverseGeocode(ctx, locationFixture.Lat, locationFixture.Lng)
	assert.Nil(t, err)
	assert.NotNil(t, address)
	assert.True(t, strings.Contains(address.FormattedAddress, "Melbourne, Victoria 3000, Australia"))
}

func TestReverseGeocodeWithNoResult(t *testing.T) {
	ctx := context.TODO()
	addr, err := geocoder.ReverseGeocode(ctx, 1, 2)
	assert.Nil(t, err)
	assert.Nil(t, addr)
}
