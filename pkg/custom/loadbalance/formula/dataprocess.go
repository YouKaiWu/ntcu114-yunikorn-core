package formula

import (
	"math"
)


func GetMean(objects []float64) float64{
	sum := 0.0
	n := len(objects)
	for i := 0; i < n; i++ {
		sum += objects[i]
	}

	return sum / float64(n)
}

func GetStdDev(objects []float64, meanOfObjects float64) float64{
	n := len(objects)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += math.Pow(math.Abs(objects[i] - meanOfObjects), 2)
	}

	return math.Sqrt(sum / float64(n))
}

func GetStandardizationScore(object, meanOfObjects, stdDevOfObjects float64) float64{
	return (object - meanOfObjects) / stdDevOfObjects
}