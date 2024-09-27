package formula

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	// "fmt"
)

func TOPSIS(req *resources.Resource, nodes nodes.Nodes) string {
	MIGs := make([]float64, 0)
	CPUUtilizations := make([]float64, 0)
	MemoryUtilizations := make([]float64, 0)
	devs := make([]float64, 0)
	nodesID := make([]string, 0)
	for _, targetNode := range nodes {
		nodeID := targetNode.NodeID
		nodesID = append(nodesID, nodeID)
		mig := GetMIG(req, targetNode)
		usageOfResource := GetUsages(req, nodeID, nodes)
		dev := GetDev(req, nodeID, nodes)
		MIGs = append(MIGs, mig)
		CPUUtilizations = append(CPUUtilizations, usageOfResource[0])
		MemoryUtilizations = append(MemoryUtilizations, usageOfResource[1])
		devs = append(devs, dev)
	}

	// Normalize
	NorCPUs := Normalized(CPUUtilizations)
	NorMems := Normalized(MemoryUtilizations)
	NorMIGs := Normalized(MIGs)
	NorDevs := Normalized(devs)

	objectNames := []string{"CPUUtilization", "MemoryUtilization", "MIG", "dev"}
	weightedCPUs := Weight(NorCPUs, objectNames)
	weightedMems := Weight(NorMems, objectNames)
	weightedMIGs := Weight(NorMIGs, objectNames)
	weightedDevs := Weight(NorDevs, objectNames)

	// A+ and A-
	APlustCPU := APlusOfUsages(weightedCPUs)
	APlustMem := APlusOfUsages(weightedMems)
	APlusMIG := APlus(weightedMIGs)
	APlusDevs := APlus(weightedDevs)

	AMinusCPU := AMinusOfUsages(weightedCPUs)
	AMinusMem := AMinusOfUsages(weightedMems)
	AMinusMIG := AMinus(weightedMIGs)
	AMinusDevs := AMinus(weightedDevs)

	// SM+ and SM-
	weighted := [][]float64{weightedCPUs, weightedMems, weightedMIGs, weightedDevs}
	APlusObjective := []float64{APlustCPU, APlustMem, APlusMIG, APlusDevs}
	AMinusObjective := []float64{AMinusCPU, AMinusMem, AMinusMIG, AMinusDevs}
	SMPlusObject := SM(weighted, APlusObjective)
	SMMinusObject := SM(weighted, AMinusObjective)

	nodeIndex, _ := IndexOfMaxRC(SMPlusObject, SMMinusObject)
	return nodesID[nodeIndex]
}

func GetObjectives(req *resources.Resource, n *objects.Node) (float64, float64) {
	// mig float64(resources.GetMIGFromNodeUtilization())
	// usage resources.AverageUsage()
	return GetMIG(req, n), GetNodeUsage(n)
}

func GetNodeUsage(n *objects.Node) float64 {
	return resources.AverageUsage(n.GetUtilizedResource())
}

func GetMIG(req *resources.Resource, n *objects.Node) float64 {
	change := resources.Sub(n.GetCapacity(), resources.Sub(n.GetAvailableResource(), req))
	return float64(resources.GetMIGFromNodeUtilization(change))
}

func GetUsages(req *resources.Resource, assignNode string, nodes nodes.Nodes) []float64 {
	ns := make([]*resources.Resource, 0)
	for _, n := range nodes {
		res := n.GetAvailableResource()
		if n.NodeID == assignNode {
			res = resources.Sub(res, req)
		}
		ns = append(ns, resources.CalculateAbsUsedCapacity(n.GetCapacity(), resources.Sub(n.GetCapacity(), res)))
	}
	ave := resources.Average(ns)
	return resources.GetCPUandMemoryUtilizations(ave)
}

func GetDev(eq *resources.Resource, assignNode string, nodes nodes.Nodes) float64 {
	ns := make([]*resources.Resource, 0)
	for _, n := range nodes {
		res := n.GetAvailableResource()
		if n.NodeID == assignNode {
			res = resources.Sub(res, eq)
		}
		ns = append(ns, resources.CalculateAbsUsedCapacity(n.GetCapacity(), resources.Sub(n.GetCapacity(), res)))
	}
	ave := resources.Average(ns)
	return resources.GetDeviationFromNodes(ns, ave)
}