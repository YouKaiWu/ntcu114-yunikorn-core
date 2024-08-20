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
	sheetIdx         int
}

func NewTenantsMonitor() *TenantsMonitor {
	tmp := &TenantsMonitor{
		tenantsList:      make(map[string]int, 0),
		numbersOfTenants: 0,
		excelFile:        excelize.NewFile(),
		currRow:          2,  // rowIdx is 1-indexed, username belongs to row 1, so we started from 2 
		sheetIdx:         -1,
	}
	tmp.sheetIdx, _ = tmp.excelFile.NewSheet("log")
	return tmp
}

func (tenantsMonitor *TenantsMonitor) saveNameRow() {
	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("log", cell, "time")
	for username, colIdx := range tenantsMonitor.tenantsList {
		cell, err := excelize.CoordinatesToCellName(colIdx, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("log", cell, username)
	}
}

func (tenantsMonitor *TenantsMonitor) AddUser(username string) {
	if _, exist := tenantsMonitor.tenantsList[username]; !exist {
		tenantsMonitor.tenantsList[username] = tenantsMonitor.numbersOfTenants + 2  // colIdx is 1-indexed and time belongs to col 1, so we started from 2
		tenantsMonitor.numbersOfTenants++
	}
}

func (tenantsMonitor *TenantsMonitor) Record(timestamp time.Time, tenants *users.Users, clusterResources *resources.Resource) { // record one row
	formattedTime := timestamp.Format("2006-01-02 15:04:05")
	cell, err := excelize.CoordinatesToCellName(1, tenantsMonitor.currRow)
	if err != nil {
		fmt.Println(err)
		return
	}
	tenantsMonitor.excelFile.SetCellValue("log", cell, formattedTime)
	for username, colIdx := range tenantsMonitor.tenantsList {
		user := tenants.GetUser(username)
		dominantResource, _ := user.GetDRS(clusterResources)
		cell, err := excelize.CoordinatesToCellName(colIdx, tenantsMonitor.currRow)
		if err != nil {
			fmt.Println(err)
			return
		}
		tenantsMonitor.excelFile.SetCellValue("log", cell, dominantResource)
	}
	tenantsMonitor.currRow++
}

func (tenantsMonitor *TenantsMonitor) Save() {
	tenantsMonitor.saveNameRow()
	tenantsMonitor.createGraph()
	if err := tenantsMonitor.excelFile.SaveAs("tenantsResourceRecord.xlsx"); err != nil {
		log.Log(log.Custom).Info("excel file saved err occur")
		log.Log(log.Custom).Info(fmt.Sprintf("%v", err))
	} else {
		dir, _ := os.Getwd()
		log.Log(log.Custom).Info(fmt.Sprintf("excel file saved to %v", dir))
	}
}
