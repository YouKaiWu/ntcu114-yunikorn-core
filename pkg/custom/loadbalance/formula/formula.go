package formula

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	
	"math"
)


func getDominantLoad(availResource, capResource *resources.Resource, resourceTypes []string) float64 {
	dominantLoad := 0.0
	for _, resourceType := range resourceTypes {
		usage := float64(capResource.Resources[resourceType] - availResource.Resources[resourceType]) / float64(capResource.Resources[resourceType])
		if usage > dominantLoad {
			dominantLoad = usage
		}
	}
	return dominantLoad
}

func getMean(objects []float64) float64 {
	sum := 0.0
	n := len(objects)
	for i := 0; i < n; i++ {
		sum += objects[i]
	}

	return sum / float64(n)
}

func getStdDev(objects []float64, meanOfObjects float64) float64 {
	n := len(objects)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += math.Pow(math.Abs(objects[i]-meanOfObjects), 2)
	}

	return math.Sqrt(sum / float64(n))
}

func getStandardizationScore(object, meanOfObjects, stdDevOfObjects float64) float64 {
	return (object - meanOfObjects) / stdDevOfObjects
}
