package parser

import (
	"github.com/apache/yunikorn-core/pkg/common/resources"
	// "github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"

	// "fmt"
)

var (
	ResourceType        = [2]string{sicommon.CPU, sicommon.Memory}
)

func ParseRequestInfo(request *objects.AllocationAsk) (appID string, requestResource *resources.Resource) {
	appID = request.GetApplicationID()
	requestResource = request.GetAllocatedResource().Clone()
	// log.Log(log.Custom).Info(fmt.Sprintf("appID:%v consume resource", appID))
	// for _, rType := range ResourceType{
	// 	log.Log(log.Custom).Info(fmt.Sprintf("resourceType: %v, quantity: %v", rType, requestResource.Resources[rType]))
	// }
	return
}