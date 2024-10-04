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
	currRow          int
	sheetIdxes       map[string]int // key: sheet name; val: sheet index
}

func NewTenantsMonitor() *TenantsMonitor {
	tmp := &TenantsMonitor{
		tenantsList:      make(map[string]int, 0),
		numbersOfTenants: 0,
		excelFile:        excelize.NewFile(),
		currRow:          2, // rowIdx is 1-indexed, username belongs to row 1, so we started from 2
		sheetIdxes:       make(map[string]int, 0),
	}
	tmp.sheetIdxes["CPU"], _ = tmp.excelFile.NewSheet("CPU")
	tmp.sheetIdxes["Memory"], _ = tmp.excelFile.NewSheet("Memory")
	tmp.sheetIdxes["DR"], _ = tmp.excelFile.NewSheet("DR")
	tmp.sheetIdxes["WaitTime"], _ = tmp.excelFile.NewSheet("WaitTime")
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

func (tenantsMonitor *TenantsMonitor) Record(timestamp time.Time, tenants *users.Users, clusterResources *resources.Resource) { // record one row
	tenantsMonitor.recordCPU(timestamp, tenants, clusterResources)
	tenantsMonitor.recordMemory(timestamp, tenants, clusterResources)
	tenantsMonitor.recordDR(timestamp, tenants, clusterResources)
	tenantsMonitor.currRow++
}

func (tenantsMonitor *TenantsMonitor) recordCPU(timestamp time.Time, tenants *users.Users, clusterResources *resources.Resource) { // record one row about CPU usage
	formattedTime := timestamp.Format("2006-01-02 15:04:05")
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("CPU", cell, formattedTime)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		CPUUsage := user.GetCPUUsage(clusterResources)
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.currRow)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("CPU", cell, CPUUsage)
	}
}

func (tenantsMonitor *TenantsMonitor) recordMemory(timestamp time.Time, tenants *users.Users, clusterResources *resources.Resource) { // record one row about Memory usage
	formattedTime := timestamp.Format("2006-01-02 15:04:05")
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("Memory", cell, formattedTime)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		memoryUsage := user.GetMemoryUsage(clusterResources)
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.currRow)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("Memory", cell, memoryUsage)
	}
}

func (tenantsMonitor *TenantsMonitor) recordDR(timestamp time.Time, tenants *users.Users, clusterResources *resources.Resource) { // record one row about Dominant Resource usage
	formattedTime := timestamp.Format("2006-01-02 15:04:05")
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("DR", cell, formattedTime)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		dominantResource, _ := user.GetDRS(clusterResources)
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.currRow)
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


func (tenantsMonitor *TenantsMonitor) Save(tenants *users.Users) {
	for sheetName := range tenantsMonitor.sheetIdxes {
		if sheetName != "WaitTime"{
			tenantsMonitor.saveNameRow(sheetName)
			tenantsMonitor.createLineChart(sheetName)
		}
	}
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