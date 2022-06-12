package nginx

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// ErrCustomOperationCode should really have an error code defined by now.
	ErrCustomOperationCode = "1006"
	// ErrInstallNginxCode provisioning failure
	ErrInstallNginxCode = "1007"
	// ErrMeshConfigCode   service mesh configuration failure
	ErrMeshConfigCode = "1008"
	// ErrClientConfigCode adapter configuration failure
	ErrClientConfigCode = "1009"
	// ErrStreamEventCode  failure
	ErrStreamEventCode = "1010"
	// ErrSampleAppCode    failure
	ErrSampleAppCode = "1011"
	// ErrOpInvalidCode failure
	ErrOpInvalidCode = "1012"
	// ErrNilClientCode represents the error code which is
	// generated when kubernetes client is nil
	ErrNilClientCode = "1013"
	// ErrApplyHelmChartCode represents the error generated
	// during the process of applying helm chart
	ErrApplyHelmChartCode = "1014"

	//ErrParseOAMComponentCode represents error in parsing oam components
	ErrParseOAMComponentCode = "1015"
	//ErrParseOAMConfigCode represents error in parsing oam config
	ErrParseOAMConfigCode = "1016"
	//ErrProcessOAMCode represents error which is thrown when an OAM operations fails
	ErrProcessOAMCode = "1017"
	//ErrNginxCoreComponentFailCode when core Nginx component processing fails
	ErrNginxCoreComponentFailCode = "1018"
	//ErrParseNginxCoreComponentCode when Nginx core component manifest parsing fails
	ErrParseNginxCoreComponentCode = "1019"

	//ErrLoadNamespaceCode occur during the process of applying namespace
	ErrLoadNamespaceCode = "1020"

	// ErrOpInvalid is an error when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})

	// ErrNilClient represents the error generated when kubernetes client is nil
	ErrNilClient = errors.New(ErrNilClientCode, errors.Alert, []string{"kubernetes client not initialized"}, []string{"Kubernetes client is nil"}, []string{"kubernetes client not initialized"}, []string{"Reconnect the adaptor to Meshery server"})

	// ErrParseOAMComponent represents the error which is
	// generated during the OAM component parsing
	ErrParseOAMComponent = errors.New(ErrParseOAMComponentCode, errors.Alert, []string{"error parsing the component"}, []string{"Error occurred while parsing application component in the OAM request made by Meshery Server"}, []string{"Could not unmarshall configuration component received via ProcessOAM gRPC call into a valid Component struct"}, []string{"Check if Meshery Server is creating valid component for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also, check if Meshery Server and adapters are referring to same component struct provided in MeshKit."})

	// ErrParseOAMConfig represents the error which is
	// generated during the OAM configuration parsing
	ErrParseOAMConfig = errors.New(ErrParseOAMConfigCode, errors.Alert, []string{"error parsing the configuration"}, []string{"Error occured while parsing configuration in the request made by Meshery Server"}, []string{"Could not unmarshall OAM config recieved via ProcessOAM gRPC call into a valid Config struct"}, []string{"Check if Meshery Server is creating valid config for ProcessOAM gRPC call. This error should never happen and can be reported as a bug in Meshery Server. Also, confirm that Meshery Server and Adapters are referring to same config struct provided in MeshKit"})
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
func ErrLoadNamespace(err error) error {
	return errors.New(ErrLoadNamespaceCode, errors.Alert, []string{"Error occured while applying namespace "}, []string{err.Error()}, []string{"Trying to access a namespace which is not available"}, []string{"Verify presence of namespace. Confirm Meshery ServiceAccount permissions"})

}
