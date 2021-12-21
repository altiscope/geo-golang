package main

import (
	"context"
	"fmt"
	"github.com/altiscope/geo-golang"
	"github.com/altiscope/geo-golang/chained"
	"github.com/altiscope/geo-golang/mapbox"
	"os"

)

const (
	addr         = "Melbourne VIC"
	lat, lng     = -37.813611, 144.963056
	radius       = 50
	zoom         = 18
	addrFR       = "Champs de Mars Paris"
	latFR, lngFR = 48.854395, 2.304770
)

func main() {
	ExampleGeocoder()
}

// ExampleGeocoder demonstrates the different geocoding services
func ExampleGeocoder() {

	fmt.Println("Mapbox API")
	try(mapbox.Geocoder(os.Getenv("MAPBOX_API_KEY")))

	// Chained geocoder will fallback to subsequent geocoders
	fmt.Println("ChainedAPI[OpenStreetmap -> Google]")
	try(chained.Geocoder(
		mapbox.Geocoder(os.Getenv("MAPBOX_API_KEY")),
		mapbox.Geocoder(os.Getenv("MAPBOX_API_KEY")),
	))

	// Mapbox API
	// Melbourne VIC location is (-37.814200, 144.963200)
	// Address of (-37.813611,144.963056) is Elwood Park Playground, Melbourne, Victoria 3000, Australia
	// Detailed address: &geo.Address{FormattedAddress:"Elwood Park Playground, Melbourne, Victoria 3000, Australia",
	// 	Street:"Elwood Park Playground", HouseNumber:"", Suburb:"", Postcode:"3000", State:"Victoria", StateDistrict:"",
	// 	County:"", Country:"Australia", CountryCode:"AU", City:"Melbourne"}

}

func try(geocoder geo.Geocoder) {
	ctx := context.TODO()
	location, _ := geocoder.Geocode(ctx, addr)
	if location != nil {
		fmt.Printf("%s location is (%.6f, %.6f)\n", addr, location.Lat, location.Lng)
	} else {
		fmt.Println("got <nil> location")
	}
	address, _ := geocoder.ReverseGeocode(ctx, lat, lng)
	if address != nil {
		fmt.Printf("Address of (%.6f,%.6f) is %s\n", lat, lng, address.FormattedAddress)
		fmt.Printf("Detailed address: %#v\n", address)
	} else {
		fmt.Println("got <nil> address")
	}
	fmt.Print("\n")
}

func tryOnlyFRData(geocoder geo.Geocoder) {
	ctx := context.TODO()
	location, _ := geocoder.Geocode(ctx, addrFR)
	if location != nil {
		fmt.Printf("%s location is (%.6f, %.6f)\n", addrFR, location.Lat, location.Lng)
	} else {
		fmt.Println("got <nil> location")
	}
	address, _ := geocoder.ReverseGeocode(ctx, latFR, lngFR)
	if address != nil {
		fmt.Printf("Address of (%.6f,%.6f) is %s\n", latFR, lngFR, address.FormattedAddress)
		fmt.Printf("Detailed address: %#v\n", address)
	} else {
		fmt.Println("got <nil> address")
	}
	fmt.Print("\n")
}
