package monitor

import (
	"fmt"
	"github.com/apache/yunikorn-core/pkg/common/resources"
	"github.com/apache/yunikorn-core/pkg/custom/fairness/users"
	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/xuri/excelize/v2"
	"os"
	"time"
)

type TenantsMonitor struct {
	tenantsList      map[string]int // key: user; val: excel col
	numbersOfTenants int
	excelFile        *excelize.File
	allocateCnt          int
	scheduleCnt      int
	sheetIdxes       map[string]int // key: sheet name; val: sheet index
}

func NewTenantsMonitor() *TenantsMonitor {
	tmp := &TenantsMonitor{
		tenantsList:      make(map[string]int, 0),
		numbersOfTenants: 0,
		excelFile:        excelize.NewFile(),
		allocateCnt:          1, // rowIdx is 1-indexed, username belongs to row 1, so we started from 2
		scheduleCnt:      1,
		sheetIdxes:       make(map[string]int, 0),
	}
	tmp.sheetIdxes["CPU"], _ = tmp.excelFile.NewSheet("CPU")
	tmp.sheetIdxes["Memory"], _ = tmp.excelFile.NewSheet("Memory")
	tmp.sheetIdxes["DR"], _ = tmp.excelFile.NewSheet("DR")
	tmp.sheetIdxes["WaitTime"], _ = tmp.excelFile.NewSheet("WaitTime")
	tmp.sheetIdxes["ScheduleInterval"], _ = tmp.excelFile.NewSheet("ScheduleInterval")
	tmp.sheetIdxes["AllocatedCnt"], _ = tmp.excelFile.NewSheet("AllocatedCnt")
	tmp.excelFile.SetActiveSheet(tmp.sheetIdxes["DR"])
	// 删除預設的 Sheet1
	if err := tmp.excelFile.DeleteSheet("Sheet1"); err != nil {
		log.Log(log.Custom).Info("delete default sheet error")
	}
	return tmp
}

func (tenantsMonitor *TenantsMonitor) saveNameRow(sheetName string) {
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue(sheetName, cell, "time")
	for username, colIdx := range tenantsMonitor.tenantsList {
		cell, err := excelize.CoordinatesToCellName(colIdx, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue(sheetName, cell, username)
	}
}

func (tenantsMonitor *TenantsMonitor) AddUser(username string) {
	if _, exist := tenantsMonitor.tenantsList[username]; !exist {
		tenantsMonitor.tenantsList[username] = tenantsMonitor.numbersOfTenants + 2 // colIdx is 1-indexed and time belongs to col 1, so we started from 2
		tenantsMonitor.numbersOfTenants++
	}
}

func (tenantsMonitor *TenantsMonitor) Record(tenants *users.Users, clusterResources *resources.Resource) { // record one row
	tenantsMonitor.recordCPU(tenants, clusterResources)
	tenantsMonitor.recordMemory(tenants, clusterResources)
	tenantsMonitor.recordDR(tenants, clusterResources)
	tenantsMonitor.recordAllocatedRequest(tenants)
	tenantsMonitor.allocateCnt++
}

func (tenantsMonitor *TenantsMonitor) recordCPU(tenants *users.Users, clusterResources *resources.Resource) { // record one row about CPU usage
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.allocateCnt+1)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("CPU", cell, tenantsMonitor.allocateCnt)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		CPUUsage := user.GetCPUUsage(clusterResources)
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.allocateCnt+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("CPU", cell, CPUUsage)
	}
}

func (tenantsMonitor *TenantsMonitor) recordMemory(tenants *users.Users, clusterResources *resources.Resource) { // record one row about Memory usage
	
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.allocateCnt+1)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("Memory", cell, tenantsMonitor.allocateCnt)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		memoryUsage := user.GetMemoryUsage(clusterResources)
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.allocateCnt+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("Memory", cell, memoryUsage)
	}
}

func (tenantsMonitor *TenantsMonitor) recordDR(tenants *users.Users, clusterResources *resources.Resource) { // record one row about Dominant Resource usage
	
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.allocateCnt+1)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("DR", cell, tenantsMonitor.allocateCnt)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		dominantResource, _ := user.GetDRS(clusterResources)
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.allocateCnt+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("DR", cell, dominantResource)
	}
}

func (tenantsMonitor *TenantsMonitor) recordWaitTime(tenants *users.Users) {
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, "username")

	for username, colIdx := range tenantsMonitor.tenantsList {
		cell, err := excelize.CoordinatesToCellName(colIdx, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, username)
	}

	cell, err = excelize.CoordinatesToCellName(1, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, "totalWaitTime")
	maxAvgWaitTime := 0.0
	minAvgWaitTime := float64(1<<31 - 1)
	totalJobWaitTime := 0.0
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		totalWaitTime := user.GetWaitTime().Seconds()
		cell, err := excelize.CoordinatesToCellName(colIdx, 2)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, totalWaitTime)

		avgWaitTime := totalWaitTime / float64(user.GetCompletedRequestCnt())
		cell, err = excelize.CoordinatesToCellName(colIdx, 3)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, avgWaitTime)

		if avgWaitTime > maxAvgWaitTime {
			maxAvgWaitTime = avgWaitTime
		}
		if avgWaitTime < minAvgWaitTime {
			minAvgWaitTime = avgWaitTime
		}

		totalJobWaitTime += totalWaitTime
	}

	cell, err = excelize.CoordinatesToCellName(1, 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, "maxAvgWaitTime")
	tenantsMonitor.excelFile.SetCellValue("WaitTime", "B4", maxAvgWaitTime)

	cell, err = excelize.CoordinatesToCellName(1, 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, "maxAvgWaitTime - minAvgWaitTime")
	tenantsMonitor.excelFile.SetCellValue("WaitTime", "B5", maxAvgWaitTime-minAvgWaitTime)

	cell, err = excelize.CoordinatesToCellName(1, 6)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("WaitTime", cell, "totalJobWaitTime")
	tenantsMonitor.excelFile.SetCellValue("WaitTime", "B6", totalJobWaitTime)
}

func (tenantsMonitor *TenantsMonitor) recordAllocatedRequest(tenants *users.Users) { // record one row about users' requests allocated count
	
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.allocateCnt+1)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("AllocatedCnt", cell, tenantsMonitor.allocateCnt)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		allocatedCnt := user.CurrAllocatedCnt
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.allocateCnt+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("AllocatedCnt", cell, allocatedCnt)
	}
}

func (tenantsMonitor *TenantsMonitor) RecordScheduleInterval(tenants *users.Users, scheduledApp string) { // record pending interval between two schedule
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.scheduleCnt+1)
	if err != nil {
		log.Log(log.Custom).Info(fmt.Sprintf("record data err:%v", err))
		return
	}
	tenantsMonitor.excelFile.SetCellValue("ScheduleInterval", cell, tenantsMonitor.scheduleCnt)
	cell, err = excelize.CoordinatesToCellName(2, tenantsMonitor.scheduleCnt+1)
	if err != nil {
		log.Log(log.Custom).Info(fmt.Sprintf("record data err:%v", err))
		return
	}
	curTime := time.Now()
	timeInterval := curTime.Sub(tenants.LastScheduleTime)
	if tenants.LastScheduleTime.IsZero(){
		timeInterval = 0
	}
	tenants.LastScheduleTime = curTime
	tenantsMonitor.excelFile.SetCellValue("ScheduleInterval", cell, timeInterval.Seconds())
	cell, err = excelize.CoordinatesToCellName(3, tenantsMonitor.scheduleCnt+1)
	if err != nil {
		log.Log(log.Custom).Info(fmt.Sprintf("record data err:%v", err))
		return
	}
	tenantsMonitor.excelFile.SetCellValue("ScheduleInterval", cell, scheduledApp)
	tenantsMonitor.scheduleCnt++
}

func (tenantsMonitor *TenantsMonitor) Save(tenants *users.Users) {
	for sheetName := range tenantsMonitor.sheetIdxes {
		if sheetName != "WaitTime" && sheetName != "ScheduleInterval" {
			tenantsMonitor.saveNameRow(sheetName)
			tenantsMonitor.createLineChart(sheetName)
		}
	}
	tenantsMonitor.createLineChartOfScheduleInterval("ScheduleInterval")
	tenantsMonitor.recordWaitTime(tenants)
	tenantsMonitor.createTotalWaitTimeBarChart("WaitTime")
	tenantsMonitor.createAvgWaitTimeBarChart("WaitTime")
	if err := tenantsMonitor.excelFile.SaveAs("tenantsResourceRecord.xlsx"); err != nil {
		log.Log(log.Custom).Info("tenantsResourceRecord excel file saved err occur")
		log.Log(log.Custom).Info(fmt.Sprintf("%v", err))
	} else {
		dir, _ := os.Getwd()
		log.Log(log.Custom).Info(fmt.Sprintf("excel file saved to %v", dir))
	}
}
