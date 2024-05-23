package formula

import (
	"math"
)

func GetBestSol(meansOfObjects, stdDevsOfObjects []float64) (float64, float64){ //(minMean, minStdDev)
	var minMean ,minStdDev float64
	minMean = meansOfObjects[0]
	minStdDev = stdDevsOfObjects[0]
 
	for i := 1; i < len(meansOfObjects); i++ {
	   if meansOfObjects[i] < minMean{
		  minMean = meansOfObjects[i]
	   }
	}
 
	for i := 1; i < len(stdDevsOfObjects); i++ {
	   if stdDevsOfObjects[i] < minStdDev{
		  minStdDev = stdDevsOfObjects[i]
	   }
	}
	return minMean,minStdDev
 }


 func GetWorstSol(meansOfObjects, stdDevsOfObjects []float64) (float64, float64){ //(maxMean, maxStdDev)
	var maxMean ,maxStdDev float64
  maxMean = meansOfObjects[0]
  maxStdDev = stdDevsOfObjects[0]

  for i := 1; i < len(meansOfObjects); i++ {
	 if meansOfObjects[i] > maxMean{
	 maxMean = meansOfObjects[i]
	 }
  }

  for i := 1; i < len(stdDevsOfObjects); i++ {
	 if stdDevsOfObjects[i] > maxStdDev{
	 maxStdDev = stdDevsOfObjects[i]
	 }   
  }
  return maxMean,maxStdDev
}


func GetEuclideanDistance(x, y []float64) float64{  
	var distance float64
	for i := 0; i < len(x); i++ {
	   distance = distance + (x[i]-y[i]) * (x[i]-y[i])
	}
	distance = math.Sqrt(distance)
	return distance
}

func GetRCVal(DistanceOfBest, DistanceOfWorst float64) float64{
	var rc float64
	rc = DistanceOfWorst / (DistanceOfBest+DistanceOfWorst)
	return rc
 }