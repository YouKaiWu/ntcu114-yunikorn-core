package formula

import (
	"fmt"
	"math"

	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	"github.com/apache/yunikorn-core/pkg/log"
)

func TOPSIS(reqestResource *resources.Resource, nodes nodes.Nodes) string {
	resourceTypes := []string{sicommon.CPU, sicommon.Memory}
	dominantLoads := make([]float64, 0)
	nodesIdx := make([]string, 0) // map idx to nodeId
	for _, node := range nodes {
		nodesIdx = append(nodesIdx, node.NodeID)
		availableResource, capacityResource := node.GetAvailableResource(), node.GetCapacity()
		dominantLoads = append(dominantLoads, GetDominantLoad(availableResource, capacityResource, resourceTypes))
	}

	means := make([]float64, 0)
	stdDevs := make([]float64, 0)

	for idx, node := range nodes {
		originalDominantLoad := dominantLoads[idx]
		availableResource, capacityResource := node.GetAvailableResource(), node.GetCapacity()
		remainResource := resources.Sub(availableResource, reqestResource)
		newDominantLoad := GetDominantLoad(remainResource, capacityResource, resourceTypes)
		dominantLoads[idx] = newDominantLoad
		means = append(means, GetMean(dominantLoads))
		stdDevs = append(stdDevs, GetStdDev(dominantLoads, means[idx]))
		dominantLoads[idx] = originalDominantLoad
	}

	// standardlizationScore
	scoresOfMeans := make([]float64, 0)
	scoresOfStdDevs := make([]float64, 0)

	meanOfMeans := GetMean(means)
	stdDevOfMeans := GetStdDev(means, meanOfMeans)
	meanOfStdDevs := GetMean(stdDevs)
	stdDevOfStdDevs := GetStdDev(stdDevs, meanOfStdDevs)
	
	for idx := range nodes {
		log.Log(log.Custom).Info(fmt.Sprintf("cur node : %v", nodesIdx[idx]))
		// standardlizationScore of mean
		scoreOfMean := 0.0
		if stdDevOfMeans != 0{
			scoreOfMean = GetStandardizationScore(means[idx], meanOfMeans, stdDevOfMeans)
		}
		scoresOfMeans = append(scoresOfMeans, scoreOfMean)
		// standardlizationScore of stdDev
		scoreOfStdDev := 0.0
		if stdDevOfStdDevs != 0{
			scoreOfStdDev = GetStandardizationScore(stdDevs[idx], meanOfStdDevs, stdDevOfStdDevs)
		}
		log.Log(log.Custom).Info(fmt.Sprintf("scoreOfMean: %v, scoreOfStdDev: %v", scoreOfMean, scoreOfStdDev))
		scoresOfStdDevs = append(scoresOfStdDevs, scoreOfStdDev)
	}

	// select solution
	minMean, minStdDev := GetBestSol(scoresOfMeans, scoresOfStdDevs)
	bestSol := []float64{minMean, minStdDev}

	maxMean, maxStdDev := GetWorstSol(scoresOfMeans, scoresOfStdDevs)
	worstSol := []float64{maxMean, maxStdDev}

	RCVals := make([]float64, 0)
	
	for idx := range nodes{
		curSol := []float64{scoresOfMeans[idx], scoresOfStdDevs[idx]} 
		distanceOfBest := GetEuclideanDistance(curSol, bestSol)
		distanceOfWorst := GetEuclideanDistance(curSol, worstSol)
		RCVal := GetRCVal(distanceOfBest, distanceOfWorst)
		RCVals = append(RCVals, RCVal)
	}

	selectedNode := GetNodeId(RCVals, nodesIdx);
	log.Log(log.Custom).Info(fmt.Sprintf("selectedNode: %v", selectedNode))
	return selectedNode
}

func GetBestSol(meansOfObjects, stdDevsOfObjects []float64) (float64, float64) { //(minMean, minStdDev)
	var minMean, minStdDev float64
	minMean = meansOfObjects[0]
	minStdDev = stdDevsOfObjects[0]

	for i := 1; i < len(meansOfObjects); i++ {
		if meansOfObjects[i] < minMean {
			minMean = meansOfObjects[i]
		}
	}

	for i := 1; i < len(stdDevsOfObjects); i++ {
		if stdDevsOfObjects[i] < minStdDev {
			minStdDev = stdDevsOfObjects[i]
		}
	}
	return minMean, minStdDev
}

func GetWorstSol(meansOfObjects, stdDevsOfObjects []float64) (float64, float64) { //(maxMean, maxStdDev)
	var maxMean, maxStdDev float64
	maxMean = meansOfObjects[0]
	maxStdDev = stdDevsOfObjects[0]

	for i := 1; i < len(meansOfObjects); i++ {
		if meansOfObjects[i] > maxMean {
			maxMean = meansOfObjects[i]
		}
	}

	for i := 1; i < len(stdDevsOfObjects); i++ {
		if stdDevsOfObjects[i] > maxStdDev {
			maxStdDev = stdDevsOfObjects[i]
		}
	}
	return maxMean, maxStdDev
}

func GetEuclideanDistance(x, y []float64) float64 {
	var distance float64
	for i := 0; i < len(x); i++ {
		distance += math.Pow((x[i] - y[i]), 2)
	}
	distance = math.Sqrt(distance)
	return distance
}

func GetRCVal(DistanceOfBest, DistanceOfWorst float64) float64 {
	rc := DistanceOfWorst / (DistanceOfBest + DistanceOfWorst)
	return rc
}

func GetNodeId(RCVals []float64, nodesIdx []string) string{
	maxRCVal := 0.0 // the worst sol RC val is 0
	selectedNode := ""
	for idx, RCVal := range RCVals{
		if RCVal >= maxRCVal{
			maxRCVal = RCVal
			selectedNode = nodesIdx[idx]
		}
	}
	return selectedNode
}