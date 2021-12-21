package chained_test

import (
	"context"
	"github.com/altiscope/geo-golang"
	"github.com/altiscope/geo-golang/chained"
	"github.com/altiscope/geo-golang/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// geocoder is chained with one data geocoder with address -> location data
// the other has location -> address data
// this will exercise the chained fallback handling
var (
	addressFixture = geo.Address{
		FormattedAddress: "64 Elizabeth Street, Melbourne, Victoria 3000, Australia",
	}
	locationFixture = geo.Location{
		Lat: -37.814107,
		Lng: 144.96328,
	}
	geocoder = chained.Geocoder(
		data.Geocoder(
			data.AddressToLocation{
				addressFixture: locationFixture,
			},
			data.LocationToAddress{},
		),

		data.Geocoder(
			data.AddressToLocation{},
			data.LocationToAddress{
				locationFixture: addressFixture,
			},
		),
	)
)

func TestGeocode(t *testing.T) {
	ctx := context.TODO()
	location, err := geocoder.Geocode(ctx, addressFixture.FormattedAddress)
	assert.NoError(t, err)
	assert.Equal(t, geo.Location{locationFixture.Lat, locationFixture.Lng}, *location)
}

func TestReverseGeocode(t *testing.T) {
	ctx := context.TODO()
	address, err := geocoder.ReverseGeocode(ctx, locationFixture.Lat, locationFixture.Lng)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(address.FormattedAddress, "Melbourne, Victoria 3000, Australia"))
}

func TestReverseGeocodeWithNoResult(t *testing.T) {
	ctx := context.TODO()
	addr, err := geocoder.ReverseGeocode(ctx, 0, 0)
	assert.Nil(t, err)
	assert.Nil(t, addr)
}

func TestChainedGeocode(t *testing.T) {
	ctx := context.TODO()
	mock1 := data.Geocoder(
		data.AddressToLocation{
			geo.Address{FormattedAddress: "Austin,TX"}: geo.Location{Lat: 1, Lng: 2},
		},
		data.LocationToAddress{},
	)

	mock2 := data.Geocoder(
		data.AddressToLocation{
			geo.Address{FormattedAddress: "Dallas,TX"}: geo.Location{Lat: 3, Lng: 4},
		},
		data.LocationToAddress{},
	)

	c := chained.Geocoder(mock1, mock2)

	l, err := c.Geocode(ctx, "Austin,TX")
	assert.NoError(t, err)
	assert.Equal(t, geo.Location{Lat: 1, Lng: 2}, *l)

	l, err = c.Geocode(ctx, "Dallas,TX")
	assert.NoError(t, err)
	assert.Equal(t, geo.Location{Lat: 3, Lng: 4}, *l)

	addr, err := c.Geocode(ctx, "NOWHERE,TX")
	assert.Nil(t, err)
	assert.Nil(t, addr)
}
