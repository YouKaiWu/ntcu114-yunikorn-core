package custom

import (
	"github.com/apache/yunikorn-core/pkg/custom/fairness"
	"github.com/apache/yunikorn-core/pkg/custom/loadbalance"
)

var fairnessManager *fairness.FairnessManager
var loadbalanceManager *loadbalance.LoadbalanceManager


func Init(){
	fairnessManager = fairness.NewFairnessManager()
	loadbalanceManager = loadbalance.NewLoadBalanceManager()
}

func GetFairnessManager() *fairness.FairnessManager{
	return fairnessManager
}

func GetLoadBalanceManager() *loadbalance.LoadbalanceManager{
	return loadbalanceManager
}