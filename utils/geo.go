package utils

import (
	"bytes"
	"math"
	"strings"
)

var (
	bits      = []int{16, 8, 4, 2, 1}
	base32    = []byte("0123456789bcdefghjkmnpqrstuvwxyz")
	neighbors = [][]string{
		[]string{
			"p0r21436x8zb9dcf5h7kjnmqesgutwvy",
			"bc01fg45238967deuvhjyznpkmstqrwx",
		},
		[]string{
			"bc01fg45238967deuvhjyznpkmstqrwx",
			"p0r21436x8zb9dcf5h7kjnmqesgutwvy",
		},
		[]string{
			"14365h7k9dcfesgujnmqp0r2twvyx8zb",
			"238967debc01fg45kmstqrwxuvhjyznp",
		},
		[]string{
			"238967debc01fg45kmstqrwxuvhjyznp",
			"14365h7k9dcfesgujnmqp0r2twvyx8zb",
		},
	}
	borders = [][]string{
		[]string{
			"prxz",
			"bcfguvyz",
		},
		[]string{
			"bcfguvyz",
			"prxz",
		},
		[]string{
			"028b",
			"0145hjnp",
		},
		[]string{
			"0145hjnp",
			"028b",
		},
	}
)

// Calculates adjacent geohashes.
func CalculateAdjacent(s, dir string) string {
	s = strings.ToLower(s)
	lastChr := s[(len(s) - 1):]
	oddEven := (len(s) % 2) // 0=even; 1=odd;
	var dirInt int
	switch dir {
	default:
		dirInt = 0
	case "right":
		dirInt = 1
	case "bottom":
		dirInt = 2
	case "left":
		dirInt = 3
	}
	// base := s[0:]
	base := s[:(len(s) - 1)]
	if strings.Index(borders[dirInt][oddEven], lastChr) != -1 {
		base = CalculateAdjacent(base, dir)
	}
	return base + string(base32[strings.Index(neighbors[dirInt][oddEven], lastChr)])
}

func refineInterval(interval []float64, cd, mask int) []float64 {
	if cd&mask > 0 {
		interval[0] = (interval[0] + interval[1]) / 2
	} else {
		interval[1] = (interval[0] + interval[1]) / 2
	}
	return interval
}

// Get LatLng coordinates from a geohash
func GeoHashDecode(geohash string) (resultLat float64, resultLng float64) {
	isEven := true
	lat := []float64{-90, 90}
	lng := []float64{-180, 180}
	latErr := float64(90)
	lngErr := float64(180)
	var c string
	var cd int
	for i := 0; i < len(geohash); i++ {
		c = geohash[i : i+1]
		cd = bytes.Index(base32, []byte(c))
		for j := 0; j < 5; j++ {
			if isEven {
				lngErr /= 2
				lng = refineInterval(lng, cd, bits[j])
			} else {
				latErr /= 2
				lat = refineInterval(lat, cd, bits[j])
			}
			isEven = !isEven
		}
	}
	resultLat = (lat[0] + lat[1]) / 2
	resultLng = (lng[0] + lng[1]) / 2
	return
}

// Create a geohash with 12 positions based on LatLng coordinates
func GeoHashEncode(latitude, longitude float64, len int) string {
	return EncodeWithPrecision(latitude, longitude, len)
}

// Create a geohash with given precision (number of characters of the resulting
// hash) based on LatLng coordinates
func EncodeWithPrecision(latitude, longitude float64, precision int) string {
	isEven := true
	lat := []float64{-90, 90}
	lng := []float64{-180, 180}
	bit := 0
	ch := 0
	var geohash bytes.Buffer
	var mid float64
	for geohash.Len() < precision {
		if isEven {
			mid = (lng[0] + lng[1]) / 2
			if longitude > mid {
				ch |= bits[bit]
				lng[0] = mid
			} else {
				lng[1] = mid
			}
		} else {
			mid = (lat[0] + lat[1]) / 2
			if latitude > mid {
				ch |= bits[bit]
				lat[0] = mid
			} else {
				lat[1] = mid
			}
		}
		isEven = !isEven
		if bit < 4 {
			bit++
		} else {
			geohash.WriteByte(base32[ch])
			bit = 0
			ch = 0
		}
	}
	return geohash.String()
}

// GeoHashMatch 匹配最近的hash位置
func GeoHashMatch(hashCode string, hashCodes ...string) (result string) {
	lat, lng := GeoHashDecode(hashCode)
	minDistance := math.MaxFloat64
	result = hashCodes[0]

	for _, code := range hashCodes {
		codeLat, codeLng := GeoHashDecode(code)
		distance := GeoDistance(lat, lng, codeLat, codeLng)
		if distance < minDistance {
			minDistance = distance
			result = code
		}
	}
	return
}

// GeoDistance 2点之间距离
func GeoDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6378138.0 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	if math.IsNaN(dist) {
		dist = 0
	}
	return dist * radius
}

// Deg2Rad converts a degree to randian
// Result:
//  - Type: float64
//  - Metric: Radian
func Deg2Rad(degree float64) float64 {
	return degree * math.Pi / 180
}

// Rad2Deg converts a randian to degree
// Result:
//  - Type: float64
//  - Metric: Degree
func Rad2Deg(radian float64) float64 {
	return radian * 180 / math.Pi
}

// SquarePoint 计算某个经纬度周围某段距离的正方形的四个点
func SquarePoint(lon, lat, distance float64) (result struct {
	LeftTop struct {
		Lon float64
		Lat float64
	}
	RightTop struct {
		Lon float64
		Lat float64
	}
	LeftBottom struct {
		Lon float64
		Lat float64
	}
	RightBottom struct {
		Lon float64
		Lat float64
	}
}) {
	radius := 6378138.0 / 1000
	dlng := 2 * math.Asin(math.Sin(distance/(2*radius))/math.Cos(Deg2Rad(lat)))
	dlngnew := Rad2Deg(dlng)

	dlat := distance / radius
	dlatnew := Rad2Deg(dlat)

	result.LeftTop.Lon = lon - dlngnew
	result.LeftTop.Lat = lat + dlatnew

	result.RightTop.Lon = lon + dlngnew
	result.RightTop.Lat = lat + dlatnew

	result.LeftBottom.Lon = lon - dlngnew
	result.LeftBottom.Lat = lat - dlatnew

	result.RightBottom.Lon = lon + dlngnew
	result.RightBottom.Lat = lat - dlatnew

	return
}
