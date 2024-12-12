package monitor

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance/nodes"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/xuri/excelize/v2"
)

type NodesMonitor struct {
	nodesList      map[string]int // key: node; val: excel col
	numbersOfNodes int
	excelFile      *excelize.File
	currRow        int
	sheetIdxes     map[string]int // key: sheet name; val: sheet index
	startTime  	   time.Time
	first 		   bool
}

func NewNodesMonitor() *NodesMonitor {
	tmp := &NodesMonitor{
		nodesList:      make(map[string]int, 0),
		numbersOfNodes: 0,
		excelFile:      excelize.NewFile(),
		currRow:        2, // rowIdx is 1-indexed, username belongs to row 1, so we started from 2
		sheetIdxes:     make(map[string]int, 0),
		startTime: time.Now(),
		first : false,
	}
	tmp.sheetIdxes["CPU"], _ = tmp.excelFile.NewSheet("CPU")
	tmp.sheetIdxes["Memory"], _ = tmp.excelFile.NewSheet("Memory")
	tmp.sheetIdxes["DL"], _ = tmp.excelFile.NewSheet("DL")
	tmp.sheetIdxes["Gap"], _ = tmp.excelFile.NewSheet("Gap")
	tmp.excelFile.SetActiveSheet(tmp.sheetIdxes["CPU"])
	// 删除預設的 Sheet1
	if err := tmp.excelFile.DeleteSheet("Sheet1"); err != nil {
		log.Log(log.Custom).Info("delete default sheet error")
	}
	return tmp
}

func (nodesMonitor *NodesMonitor) saveNameRow(sheetName string) {
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	nodesMonitor.excelFile.SetCellValue(sheetName, cell, "time")
	for nodeID, colIdx := range nodesMonitor.nodesList {
		cell, err := excelize.CoordinatesToCellName(colIdx, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		nodesMonitor.excelFile.SetCellValue(sheetName, cell, nodeID)
	}
}

func (nodesMonitor *NodesMonitor) AddNode(nodeID string) {
	if _, exist := nodesMonitor.nodesList[nodeID]; !exist {
		nodesMonitor.nodesList[nodeID] = nodesMonitor.numbersOfNodes + 2 // colIdx is 1-indexed and time belongs to col 1, so we started from 2
		nodesMonitor.numbersOfNodes++
	}
}

func (nodesMonitor *NodesMonitor) Record(recordTime time.Time, nodes *nodes.Nodes) { // record one row of various usage
	if !nodesMonitor.first{
		nodesMonitor.startTime = recordTime
		nodesMonitor.first = true
	}
	nodesMonitor.recordCPU(recordTime, nodes)
	nodesMonitor.recordMemory(recordTime, nodes)
	nodesMonitor.recordDL(recordTime, nodes)
	nodesMonitor.recordGap(recordTime, nodes)
	nodesMonitor.currRow++
}
func (nodesMonitor *NodesMonitor) recordCPU(recordTime time.Time, nodes *nodes.Nodes) { // record one row about cpu usage of nodes
	cell, err := excelize.CoordinatesToCellName(1, nodesMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	scheduleTime := recordTime.Sub(nodesMonitor.startTime).Seconds()
	nodesMonitor.excelFile.SetCellValue("CPU", cell, scheduleTime)
	for _, node := range *nodes {
		colIdx := nodesMonitor.nodesList[node.NodeID]
		cell, err := excelize.CoordinatesToCellName(colIdx, nodesMonitor.currRow)
		if err != nil {
			fmt.Println(err)
			return
		}
		usedResource := resources.Sub(node.GetCapacity(), node.GetAvailableResource())
		nodeUsage := float64(usedResource.Resources["vcore"]) / float64(node.GetCapacity().Resources["vcore"])
		nodesMonitor.excelFile.SetCellValue("CPU", cell, nodeUsage)
	}
}

func (nodesMonitor *NodesMonitor) recordMemory(recordTime time.Time, nodes *nodes.Nodes) { // record one row about memory usage of nodes
	cell, err := excelize.CoordinatesToCellName(1, nodesMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	scheduleTime := recordTime.Sub(nodesMonitor.startTime).Seconds()
	nodesMonitor.excelFile.SetCellValue("Memory", cell, scheduleTime)
	for _, node := range *nodes {
		colIdx := nodesMonitor.nodesList[node.NodeID]
		cell, err := excelize.CoordinatesToCellName(colIdx, nodesMonitor.currRow)
		if err != nil {
			fmt.Println(err)
			return
		}
		usedResource := resources.Sub(node.GetCapacity(), node.GetAvailableResource())
		nodeUsage := float64(usedResource.Resources["memory"]) / float64(node.GetCapacity().Resources["memory"])
		nodesMonitor.excelFile.SetCellValue("Memory", cell, nodeUsage)
	}
}

func (nodesMonitor *NodesMonitor) recordDL(recordTime time.Time, nodes *nodes.Nodes) { // record one row about Dominant Load(most stressed resource) usage of nodes
	cell, err := excelize.CoordinatesToCellName(1, nodesMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	scheduleTime := recordTime.Sub(nodesMonitor.startTime).Seconds()
	nodesMonitor.excelFile.SetCellValue("DL", cell, scheduleTime)
	for _, node := range *nodes {
		colIdx := nodesMonitor.nodesList[node.NodeID]
		cell, err := excelize.CoordinatesToCellName(colIdx, nodesMonitor.currRow)
		if err != nil {
			fmt.Println(err)
			return
		}
		usedResource := resources.Sub(node.GetCapacity(), node.GetAvailableResource())
		memoryUsage := float64(usedResource.Resources["memory"]) / float64(node.GetCapacity().Resources["memory"])
		CPUUsage := float64(usedResource.Resources["vcore"]) / float64(node.GetCapacity().Resources["vcore"])
		DL := max(memoryUsage, CPUUsage)
		nodesMonitor.excelFile.SetCellValue("DL", cell, DL)
	}
}

func (nodesMonitor *NodesMonitor) recordGap(recordTime time.Time, nodes *nodes.Nodes) { // record one row about gap between cpu usage and memory of nodes
	cell, err := excelize.CoordinatesToCellName(1, nodesMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	scheduleTime := recordTime.Sub(nodesMonitor.startTime).Seconds()
	nodesMonitor.excelFile.SetCellValue("Gap", cell, scheduleTime)
	for _, node := range *nodes {
		colIdx := nodesMonitor.nodesList[node.NodeID]
		cell, err := excelize.CoordinatesToCellName(colIdx, nodesMonitor.currRow)
		if err != nil {
			fmt.Println(err)
			return
		}
		usedResource := resources.Sub(node.GetCapacity(), node.GetAvailableResource())
		memoryUsage := float64(usedResource.Resources["memory"]) / float64(node.GetCapacity().Resources["memory"])
		CPUUsage := float64(usedResource.Resources["vcore"]) / float64(node.GetCapacity().Resources["vcore"])
		Gap := math.Abs(memoryUsage - CPUUsage)
		nodesMonitor.excelFile.SetCellValue("Gap", cell, Gap)
	}
}

func (nodesMonitor *NodesMonitor) Save() {
	for sheetName := range nodesMonitor.sheetIdxes {
		nodesMonitor.saveNameRow(sheetName)
		nodesMonitor.createGraph(sheetName)
	}
	if err := nodesMonitor.excelFile.SaveAs("nodesResourceRecord.xlsx"); err != nil {
		log.Log(log.Custom).Info("nodesResourceRecord excel file saved err occur")
		log.Log(log.Custom).Info(fmt.Sprintf("%v", err))
	} else {
		dir, _ := os.Getwd()
		log.Log(log.Custom).Info(fmt.Sprintf("excel file saved to %v", dir))
	}
}
