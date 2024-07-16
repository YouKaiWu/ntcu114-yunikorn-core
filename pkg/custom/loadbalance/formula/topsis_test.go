package formula

import (
	// "github.com/apache/yunikorn-core/pkg/common/resources"
	// "github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"

	// sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	
	"testing"
	"gotest.tools/v3/assert"
)


// func TestTOPSIS(t *testing.T) {
	
// 	// 创建模拟资源和节点数据
// 	requestResource := &resources.Resource{
// 		Resources: map[string]resources.Quantity{
// 			sicommon.CPU:    2,
// 			sicommon.Memory: 10,
// 		},
// 	}

// 	// 创建两个模拟节点
// 	node1 := &nodes.Node{
// 		NodeID: "node1",
// 		AvailableResource: &resources.Resource{
// 			Resources: map[string]resources.Quantity{
// 				sicommon.CPU:    16,
// 				sicommon.Memory: 20,
// 			},
// 		},
// 		CapacityResource: &resources.Resource{
// 			Resources: map[string]resources.Quantity{
// 				sicommon.CPU:    20,
// 				sicommon.Memory: 30,
// 			},
// 		},
// 	}

// // 4/20;10/30  1/3
// // 10/30;4/20  1/3

// // 2 ; 10
// // 6/20;20/30  2/3
// // 12/30;14/20  7/10

	
// 	node2 := &nodes.Node{
// 		NodeID: "node2",
// 		AvailableResource: &resources.Resource{
// 			Resources: map[string]resources.Quantity{
// 				sicommon.CPU:    20,
// 				sicommon.Memory: 16,
// 			},
// 		},
// 		CapacityResource: &resources.Resource{
// 			Resources: map[string]resources.Quantity{
// 				sicommon.CPU:    30,
// 				sicommon.Memory: 20,
// 			},
// 		},
// 	}

// 	// 创建节点集合
// 	nodeList := nodes.NewNodes()
// 	nodeList.AddNode(node1)
// 	nodeList.AddNode(node2)

// 	// 调用 TOPSIS 函数
// 	selectedNode := TOPSIS(requestResource, *nodeList)

// 	// 断言选择的节点
// 	expectedNode := "node1" // 根据你的实际情况修改期望的节点 ID
// 	assert.Equal(t, expectedNode, selectedNode)
// }


func TestGetBestSol(t *testing.T) {
	meansOfObjects := [...]float64{1, 2, 3}
	stdDevsOfObjectsOfObjects := [...]float64{1, 2, 3}
	minMean, minStdDev := getBestSol(meansOfObjects[:], stdDevsOfObjectsOfObjects[:])
	expectMinMean := 1.0
	expectMinStdDev := 1.0
	assert.Equal(t, minMean, expectMinMean)
	assert.Equal(t, minStdDev, expectMinStdDev)
}

func TestGetWorstSol(t *testing.T) {
	meansOfObjects := [...]float64{1, 2, 3}
	stdDevsOfObjects := [...]float64{1, 2, 3}
	maxMean, maxStdDev := getWorstSol(meansOfObjects[:], stdDevsOfObjects[:])
	expectMaxMean := 3.0
	expectMaxStdDev := 3.0
	assert.Equal(t, maxMean, expectMaxMean)
	assert.Equal(t, maxStdDev, expectMaxStdDev)
}

func TestGetEuclideanDistance(t *testing.T) {
	x := [...]float64{4, 5, 6, 7}
	y := [...]float64{1, 2, 3, 4}
	distance := getEuclideanDistance(x[:], y[:])
	expect := 6
	assert.Equal(t, distance, expect)
}

func TestGetRCVal(t *testing.T) {
	best := 1.0
	worst := 4.0
	rc := getRCVal(best, worst)
	expect := 0.8
	assert.Equal(t, rc, expect)
}