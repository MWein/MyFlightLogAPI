package utils

import "math"

func LatLong2Meters(lat float64, long float64) [2]float64 {
	var meters [2]float64
	meters[0] = long * 20037508.34 / 180.0
	meters[1] = math.Log(math.Tan((90.0+lat)*math.Pi/360)) / (math.Pi / 180)
	meters[1] = meters[1] * 20037508.34 / 180

	return meters
}
