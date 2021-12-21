package cached_test

import (
	"context"
	"github.com/altiscope/geo-golang"
	"github.com/altiscope/geo-golang/cached"
	"github.com/altiscope/geo-golang/data"
	"strings"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

var geoCache = cache.New(5*time.Minute, 30*time.Second)

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
	geocoder = cached.Geocoder(
		data.Geocoder(
			data.AddressToLocation{
				addressFixture: locationFixture,
			},
			data.LocationToAddress{
				locationFixture: addressFixture,
			},
		),
		geoCache,
	)
)

func TestGeocode(t *testing.T) {
	ctx := context.TODO()
	location, err := geocoder.Geocode(ctx, "64 Elizabeth Street, Melbourne, Victoria 3000, Australia")
	assert.NoError(t, err)
	assert.Equal(t, locationFixture, *location)
}

func TestReverseGeocode(t *testing.T) {
	ctx := context.TODO()
	address, err := geocoder.ReverseGeocode(ctx, locationFixture.Lat, locationFixture.Lng)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(address.FormattedAddress, "Melbourne, Victoria 3000, Australia"))
}

func TestReverseGeocodeWithNoResult(t *testing.T) {
	ctx := context.TODO()
	addr, err := geocoder.ReverseGeocode(ctx,1, 2)
	assert.Nil(t, err)
	assert.Nil(t, addr)
}

func TestCachedGeocode(t *testing.T) {
	ctx := context.TODO()
	mockAddr := geo.Address{
		FormattedAddress: "42, Some Street, Austin, Texas",
	}
	mock1 := data.Geocoder(
		data.AddressToLocation{
			mockAddr: geo.Location{Lat: 1, Lng: 2},
		},
		data.LocationToAddress{},
	)

	c := cached.Geocoder(mock1, geoCache)

	l, err := c.Geocode(ctx, "42, Some Street, Austin, Texas")
	assert.NoError(t, err)
	assert.Equal(t, geo.Location{Lat: 1, Lng: 2}, *l)

	// Should be cached
	// TODO: write a mock Cache impl to test cache is being used
	l, err = c.Geocode(ctx, "42, Some Street, Austin, Texas")
	assert.NoError(t, err)
	assert.Equal(t, geo.Location{Lat: 1, Lng: 2}, *l)

	addr, err := c.Geocode(ctx, "NOWHERE,TX")
	assert.Nil(t, err)
	assert.Nil(t, addr)
}
