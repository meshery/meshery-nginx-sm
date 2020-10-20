package nginx

import (
	"context"
	"fmt"

	cfg "github.com/layer5io/meshery-nginx/internal/config"
)

// Operation holds the informormation for list of operations
type Operation struct {
	Type       int32             `json:"type,string,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// Operations hold a map of Operation objects
type Operations map[string]*Operation

// ApplyOperation applies the operation on nginx
func (h *handler) ApplyOperation(ctx context.Context, op string, id string, del bool) error {

	operations := make(Operations, 0)
	err := h.config.Operations(&operations)
	if err != nil {
		return ErrApplyOperation(err)
	}

	status := "deploying"
	e := &Event{
		Operationid: id,
		Summary:     "Deploying",
		Details:     "None",
	}

	switch op {
	case cfg.InstallNginxLatest:
		go func(hh *handler, ee *Event) {
			if status, err := hh.installNginx(del, operations[op].Properties["version"]); err != nil {
				e.Summary = fmt.Sprintf("Error while %s Nginx service mesh", status)
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			h.config.SetKey(cfg.RunningMeshVersion, operations[op].Properties["version"])
			ee.Summary = fmt.Sprintf("Nginx service mesh %s successfully", status)
			ee.Details = fmt.Sprintf("The Nginx service mesh is now %s.", status)
			hh.StreamInfo(e)
		}(h, e)
	case cfg.InstallSampleBookInfo:
		go func(hh *handler, ee *Event) {
			if status, err := hh.installSampleApp(operations[op].Properties["description"]); err != nil {
				e.Summary = fmt.Sprintf("Error while %s Sample %s application", status, operations[op].Properties["description"])
				e.Details = err.Error()
				hh.StreamErr(e, err)
				return
			}
			ee.Summary = fmt.Sprintf("Sample %s application %s successfully", operations[op].Properties["description"], status)
			ee.Details = fmt.Sprintf("The Sample %s application is now %s.", operations[op].Properties["description"], status)
			hh.StreamInfo(e)
		}(h, e)
	case cfg.ValidateSmiConformance:
		go func(hh *handler, ee *Event) {
			version, err := h.config.GetKey(cfg.RunningMeshVersion)
			if err != nil {
				return
			}
			err = hh.validateSMIConformance(ee.Operationid, version)
			if err != nil {
				return
			}
		}(h, e)
	default:
		h.StreamErr(e, ErrOpInvalid)
	}

	return nil
}

// ListOperations lists the operations available
func (h *handler) ListOperations() (Operations, error) {
	operations := make(Operations, 0)
	err := h.config.Operations(&operations)
	if err != nil {
		return nil, ErrListOperations(err)
	}
	return operations, nil
}
