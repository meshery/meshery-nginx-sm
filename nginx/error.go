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

	//ErrParseOAMComponentCode represents error in parsing oam components
	ErrParseOAMComponentCode = "1014"
	//ErrParseOAMConfigCode represents error in parsing oam config
	ErrParseOAMConfigCode = "1015"
	//ErrProcessOAMCode represents error which is thrown when an OAM operations fails
	ErrProcessOAMCode = "1016"
	//ErrNginxCoreComponentFailCode when core Nginx component processing fails
	ErrNginxCoreComponentFailCode = "1017"
	//ErrParseNginxCoreComponentCode when Nginx core component manifest parsing fails
	ErrParseNginxCoreComponentCode = "1018"

	//ErrLoadNamespaceCode occur during the process of applying namespace
	ErrLoadNamespaceCode = "1015"

	// ErrOpInvalid is an error when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})

	// ErrNilClient represents the error generated when kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"kubernetes client not initialized"}, []string{"Reconnect the adaptor to Meshery server"})
	//ErrParseOAMComponent represents the error generated when oam component could not be parsed
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occurred while prasing application component in the OAM request made"}, []string{"Invalid OAM component passed in OAM request"}, []string{"Check if your request has vaild OAM components"})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Could not generate application configuration from given json"}, []string{"Invalid OAM config passed in MeshOps request"}, []string{"Confirm that the request has valid OAM config"})
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

// ErrProcessOAM is a generic error which is thrown when an OAM operations fails
func ErrProcessOAM(err error) error {
	return errors.New(ErrProcessOAMCode, errors.Alert, []string{"error performing OAM operations"}, []string{err.Error()}, []string{}, []string{})
}

// ErrNginxCoreComponentFail is the error when core Nginx component processing fails
func ErrNginxCoreComponentFail(err error) error {
	return errors.New(ErrNginxCoreComponentFailCode, errors.Alert, []string{"error in NGINX Service Mesh core component"}, []string{err.Error()}, []string{}, []string{})
}

// ErrParseNginxCoreComponent is the error when Nginx core component manifest parsing fails
func ErrParseNginxCoreComponent(err error) error {
	return errors.New(ErrParseNginxCoreComponentCode, errors.Alert, []string{"Failure to parse core component manifest for NGINX Service Mesh"}, []string{err.Error()}, []string{}, []string{})
}

// ErrLoadNamespace is the occurend while applying namespace
func ErrLoadNamespace(err error, s string) error {
	return errors.New(ErrLoadNamespaceCode, errors.Alert, []string{"Error occured while applying namespace "}, []string{err.Error()}, []string{"Trying to access a namespace which is not available"}, []string{"Verify presence of namespace. Confirm Meshery ServiceAccount permissions"})

}
