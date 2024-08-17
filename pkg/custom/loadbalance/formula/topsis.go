package formula

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
	"github.com/apache/yunikorn-core/pkg/log"
	// "go.uber.org/zap"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"

	"fmt"
	"math"
)

func TOPSIS(reqestResource *resources.Resource, nodes nodes.Nodes) string {
	resourceTypes := []string{sicommon.CPU, sicommon.Memory}
	dominantLoads := make([]float64, 0)
	nodesIdx := make([]string, 0) // map idx to nodeId
	for _, node := range nodes {
		nodesIdx = append(nodesIdx, node.NodeID)
		availableResource, capacityResource := node.GetAvailableResource(), node.GetCapacity()
		dominantLoads = append(dominantLoads, getDominantLoad(availableResource, capacityResource, resourceTypes))
	}

	means := make([]float64, 0)
	stdDevs := make([]float64, 0)

	for idx, node := range nodes {
		originalDominantLoad := dominantLoads[idx]
		availableResource, capacityResource := node.GetAvailableResource(), node.GetCapacity()
		remainResource := resources.Sub(availableResource, reqestResource)
		newDominantLoad := getDominantLoad(remainResource, capacityResource, resourceTypes)
		dominantLoads[idx] = newDominantLoad
		means = append(means, getMean(dominantLoads))
		stdDevs = append(stdDevs, getStdDev(dominantLoads, means[idx]))
		dominantLoads[idx] = originalDominantLoad

		// log.Log(log.Custom).Info(node.String(), zap.Bool("isSchedulable", node.IsSchedulable()))
	}

	// standardlizationScore
	scoresOfMeans := make([]float64, 0)
	scoresOfStdDevs := make([]float64, 0)

	meanOfMeans := getMean(means)
	stdDevOfMeans := getStdDev(means, meanOfMeans)
	meanOfStdDevs := getMean(stdDevs)
	stdDevOfStdDevs := getStdDev(stdDevs, meanOfStdDevs)
	
	for idx := range nodes {
		// log.Log(log.Custom).Info(fmt.Sprintf("cur node : %v", nodesIdx[idx]))
		// standardlizationScore of mean
		scoreOfMean := 0.0
		if stdDevOfMeans != 0{
			scoreOfMean = getStandardizationScore(means[idx], meanOfMeans, stdDevOfMeans)
		}
		scoresOfMeans = append(scoresOfMeans, scoreOfMean)
		// standardlizationScore of stdDev
		scoreOfStdDev := 0.0
		if stdDevOfStdDevs != 0{
			scoreOfStdDev = getStandardizationScore(stdDevs[idx], meanOfStdDevs, stdDevOfStdDevs)
		}
		// log.Log(log.Custom).Info(fmt.Sprintf("scoreOfMean: %v, scoreOfStdDev: %v", scoreOfMean, scoreOfStdDev))
		scoresOfStdDevs = append(scoresOfStdDevs, scoreOfStdDev)
	}

	// select solution
	minMean, minStdDev := getBestSol(scoresOfMeans, scoresOfStdDevs)
	bestSol := []float64{minMean, minStdDev}

	maxMean, maxStdDev := getWorstSol(scoresOfMeans, scoresOfStdDevs)
	worstSol := []float64{maxMean, maxStdDev}

	RCVals := make([]float64, 0)
	
	for idx := range nodes{
		curSol := []float64{scoresOfMeans[idx], scoresOfStdDevs[idx]} 
		distanceOfBest := getEuclideanDistance(curSol, bestSol)
		distanceOfWorst := getEuclideanDistance(curSol, worstSol)
		RCVal := getRCVal(distanceOfBest, distanceOfWorst)
		RCVals = append(RCVals, RCVal)
	}

	selectedNode := getNodeId(RCVals, nodesIdx);
	log.Log(log.Custom).Info(fmt.Sprintf("selectedNode: %v", selectedNode))
	return selectedNode
}

func getBestSol(meansOfObjects, stdDevsOfObjects []float64) (float64, float64) { //(minMean, minStdDev)
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

func getWorstSol(meansOfObjects, stdDevsOfObjects []float64) (float64, float64) { //(maxMean, maxStdDev)
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

func getEuclideanDistance(x, y []float64) float64 {
	var distance float64
	for i := 0; i < len(x); i++ {
		distance += math.Pow((x[i] - y[i]), 2)
	}
	distance = math.Sqrt(distance)
	return distance
}

func getRCVal(DistanceOfBest, DistanceOfWorst float64) float64 {
	rc := DistanceOfWorst / (DistanceOfBest + DistanceOfWorst)
	return rc
}

func getNodeId(RCVals []float64, nodesIdx []string) string{
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