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

type Nginx struct {
	adapter.Adapter // Type Embedded

	// DockerRegistry is the registry for
	// nginx service mesh related images
	DockerRegistry string

	// Executable is the path where the
	// nginx-meshctl is located
	Executable string
}

// New initializes nginx handler.
func New(c adapterconfig.Handler, l logger.Handler, kc adapterconfig.Handler, executable string) adapter.Handler {
	return &Nginx{
		Adapter: adapter.Adapter{
			Config:            c,
			Log:               l,
			KubeconfigHandler: kc,
		},
		DockerRegistry: "layer5",
		Executable:     executable,
	}
}

// ApplyOperation applies the operation on nginx
func (nginx *Nginx) ApplyOperation(ctx context.Context, opReq adapter.OperationRequest) error {

	operations := make(adapter.Operations, 0)
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
			err := hh.ValidateSMIConformance(&adapter.SmiTestOptions{
				Ctx:  context.TODO(),
				OpID: ee.Operationid,
			})
			if err != nil {
				return
			}
		}(nginx, e)
	default:
		nginx.StreamErr(e, ErrOpInvalid)
	}

	return nil
}
