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

func (tenantsMonitor *TenantsMonitor) Save() {
	for sheetName := range tenantsMonitor.sheetIdxes {
		tenantsMonitor.saveNameRow(sheetName)
		tenantsMonitor.createGraph(sheetName)
	}
	if err := tenantsMonitor.excelFile.SaveAs("tenantsResourceRecord.xlsx"); err != nil {
		log.Log(log.Custom).Info("tenantsResourceRecord excel file saved err occur")
		log.Log(log.Custom).Info(fmt.Sprintf("%v", err))
	} else {
		dir, _ := os.Getwd()
		log.Log(log.Custom).Info(fmt.Sprintf("excel file saved to %v", dir))
	}
}
