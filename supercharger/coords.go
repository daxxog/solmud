package supercharger
import (
	"fmt"
)

type IDMS interface {
	Degrees() int16 // Positive latitudes are north of the equator, negative latitudes are south of the equator; Positive longitudes are east of the Prime Meridian; negative longitudes are west of the Prime Meridian.
	Minutes() uint8
	Seconds() uint8
	IsLat() bool // !IsLon()
	IsLon() bool // !IsLat()
	fmt.Stringer // 89°35'09.2"W
}

type ICoords interface {
	fmt.Stringer // 44°31'24.3"N 89°35'09.2"W
	Lon() IDMS // Longitude
	Lat() IDMS // Latitude
	GoogleMaps() // https://www.google.com/maps/place/44%C2%B031'24.3%22N+89%C2%B035'09.2%22W
}

func NewCoords(loc string) ICoords {
}
