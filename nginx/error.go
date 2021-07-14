package nginx

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	// ErrCustomOperationCode should really have an error code defined by now.
	ErrCustomOperationCode = "nginx_test_code"
	// ErrInstallNginxCode provisioning failure
	ErrInstallNginxCode = "nginx_test_code"
	// ErrMeshConfigCode   service mesh configuration failure
	ErrMeshConfigCode = "nginx_test_code"
	// ErrClientConfigCode adapter configuration failure
	ErrClientConfigCode = "nginx_test_code"
	// ErrStreamEventCode  failure
	ErrStreamEventCode = "nginx_test_code"
	// ErrExecDeployCode   failure
	ErrExecDeployCode = "nginx_test_code"
	// ErrExecRemoveCode   failure
	ErrExecRemoveCode = "nginx_test_code"
	// ErrSampleAppCode    failure
	ErrSampleAppCode = "nginx_test_code"
	// ErrOpInvalidCode failure
	ErrOpInvalidCode = "1014"

	// ErrOpInvalid is an error when an invalid operation is requested
	ErrOpInvalid = errors.New(ErrOpInvalidCode, errors.Alert, []string{"Invalid operation"}, []string{}, []string{}, []string{})
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

// ErrExecDeploy is the error for deploying nginx service mesh
func ErrExecDeploy(err error, des string) error {
	return errors.New(ErrExecDeployCode, errors.Alert, []string{"Error executing deploy command"}, []string{des, err.Error()}, []string{}, []string{})
}

// ErrExecRemove is the error for removing nginx service mesh
func ErrExecRemove(err error, des string) error {
	return errors.New(ErrExecRemoveCode, errors.Alert, []string{"Error executing remove command", des}, []string{err.Error()}, []string{}, []string{})
}

// ErrSampleApp is the error for operations on the sample apps
func ErrSampleApp(err error) error {
	return errors.New(ErrSampleAppCode, errors.Alert, []string{"Error with sample app operation"}, []string{err.Error()}, []string{}, []string{})
}

// ErrCustomOperation is the error for custom operations
func ErrCustomOperation(err error) error {
	return errors.New(ErrCustomOperationCode, errors.Alert, []string{"Error with applying custom operation"}, []string{err.Error()}, []string{}, []string{})
}
