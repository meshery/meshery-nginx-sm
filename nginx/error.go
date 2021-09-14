package nginx

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// ErrCustomOperationCode should really have an error code defined by now.
	ErrCustomOperationCode = "1005"
	// ErrInstallNginxCode provisioning failure
	ErrInstallNginxCode = "1006"
	// ErrMeshConfigCode   service mesh configuration failure
	ErrMeshConfigCode = "1007"
	// ErrClientConfigCode adapter configuration failure
	ErrClientConfigCode = "1008"
	// ErrStreamEventCode  failure
	ErrStreamEventCode = "1009"
	// ErrSampleAppCode    failure
	ErrSampleAppCode = "1010"
	// ErrOpInvalidCode failure
	ErrOpInvalidCode = "1011"
	// ErrNilClientCode represents the error code which is
	// generated when kubernetes client is nil
	ErrNilClientCode = "1012"
	// ErrApplyHelmChartCode represents the error generated
	// during the process of applying helm chart
	ErrApplyHelmChartCode = "1013"

	// ErrOpInvalid is an error when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})

	// ErrNilClient represents the error generated when kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"kubernetes client not initialized"}, []string{"Reconnect the adaptor to Meshery server"})
)

// ErrInstallNginx is the error for install mesh
func ErrInstallNginx(err error) error {
	return errors.New(ErrInstallNginxCode, errors.Alert, []string{"Error with Nginx installation"}, []string{err.Error()}, []string{}, []string{})

}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(ErrMeshConfigCode, errors.Alert, []string{"Error configuration mesh"}, []string{err.Error()}, []string{}, []string{})
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.New(ErrClientConfigCode, errors.Alert, []string{"Error setting client config"}, []string{err.Error()}, []string{}, []string{})
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.New(ErrStreamEventCode, errors.Alert, []string{"Error streaming events"}, []string{err.Error()}, []string{}, []string{})
}

// ErrSampleApp is the error for operations on the sample apps
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrCustomOperation is the error for custom operations
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with applying custom operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrApplyHelmChart is the occurend while applying helm chart
func ErrApplyHelmChart(err error) error {
	return errors.New(ErrApplyHelmChartCode, errors.Alert, []string{"Error occured while applying Helm Chart"}, []string{err.Error()}, []string{}, []string{})
}
