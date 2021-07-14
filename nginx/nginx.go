package nginx

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/common"
	adapterconfig "github.com/layer5io/meshery-adapter-library/config"
	"github.com/layer5io/meshery-adapter-library/status"
	internalconfig "github.com/layer5io/meshery-nginx/internal/config"
	"github.com/layer5io/meshkit/logger"
)

// Nginx defines a model for this adapter
type Nginx struct {
	adapter.Adapter // Type Embedded
}

// New initializes nginx handler.
func New(c adapterconfig.Handler, l logger.Handler, kc adapterconfig.Handler) adapter.Handler {
	return &Nginx{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
		},
	}
}

// ApplyOperation applies the operation on nginx
func (nginx *Nginx) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest) error {

	operations := make(adapter.Operations)
	err := nginx.Config.GetObject(adapter.OperationsKey, &operations)
	if err != nil {
		return err
	}

	stat := status.Deploying

	e := &adapter.Event{
		Operationid: opReq.OperationID,
		Summary:     status.Deploying,
		Details:     status.None,
	}

	switch opReq.OperationName {
	case internalconfig.NginxOperation:
		go func(hh *Nginx, ee *adapter.Event) {
			version := string(operations[opReq.OperationName].Versions[0])
			if stat, err = hh.installNginx(opReq.IsDeleteOperation, version); err != nil {
				e.Summary = fmt.Sprintf("Error while %s Nginx service mesh", stat)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Nginx service mesh %s successfully", stat)
			ee.Details = fmt.Sprintf("The Nginx service mesh is now %s.", stat)
			hh.StreamInfo(e)
		}(nginx, e)
	case common.SmiConformanceOperation:
		go func(hh *Nginx, ee *adapter.Event) {
			name := operations[opReq.OperationName].Description
			_, err := hh.RunSMITest(adapter.SMITestOptions{
				Ctx:         context.TODO(),
				OperationID: ee.Operationid,
				Manifest:    string(operations[opReq.OperationName].Templates[0]),
				Namespace:   "meshery",
				Labels:      make(map[string]string),
				Annotations: make(map[string]string),
			})
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s %s test", status.Running, name)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s test %s successfully", name, status.Completed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(nginx, e)
	case common.BookInfoOperation, common.HTTPBinOperation, common.ImageHubOperation, common.EmojiVotoOperation:
		go func(hh *Nginx, ee *adapter.Event) {
			appName := operations[opReq.OperationName].AdditionalProperties[common.ServiceName]
			stat, err := hh.installSampleApp(opReq.Namespace, opReq.IsDeleteOperation, operations[opReq.OperationName].Templates)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s %s application", stat, appName)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("%s application %s successfully", appName, stat)
			ee.Details = fmt.Sprintf("The %s application is now %s.", appName, stat)
			hh.StreamInfo(e)
		}(nginx, e)
	case common.CustomOperation:
		go func(hh *Nginx, ee *adapter.Event) {
			stat, err := hh.applyCustomOperation(opReq.Namespace, opReq.CustomBody, opReq.IsDeleteOperation)
			if err != nil {
				e.Summary = fmt.Sprintf("Error while %s custom operation", stat)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Manifest %s successfully", status.Deployed)
			ee.Details = ""
			hh.StreamInfo(e)
		}(nginx, e)
	default:
		nginx.StreamErr(e, ErrOpInvalid)
	}

	return nil
}
