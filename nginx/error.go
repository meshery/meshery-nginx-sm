package nginx

import (
	"github.com/layer5io/meshkit/errors"
)

var (
	ErrInstallNginxCode    = "nginx_test_code"
	ErrMeshConfigCode      = "nginx_test_code"
	ErrClientConfigCode    = "nginx_test_code"
	ErrStreamEventCode     = "nginx_test_code"
	ErrExecDeployCode      = "nginx_test_code"
	ErrExecRemoveCode      = "nginx_test_code"
	ErrSampleAppCode       = "nginx_test_code"
	ErrCustomOperationCode = "nginx_test_code"

	ErrOpInvalid = errors.NewDefault(errors.ErrOpInvalid, "Invalid operation")
)

// ErrInstallNginx is the error for install mesh
func ErrInstallNginx(err error) error {
	return errors.NewDefault(ErrInstallNginxCode, "Error installing nginx", err.Error())
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.NewDefault(ErrMeshConfigCode, "Error configuration mesh", err.Error())
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.NewDefault(ErrClientConfigCode, "Error setting client config", err.Error())
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.NewDefault(ErrStreamEventCode, "Error streaming event", err.Error())
}

// ErrExecDeploy is the error for deploying nginx service mesh
func ErrExecDeploy(err error, des string) error {
	return errors.NewDefault(ErrExecDeployCode, "Error executing deploy command", des)
}

// ErrExecRemove is the error for removing nginx service mesh
func ErrExecRemove(err error, des string) error {
	return errors.NewDefault(ErrExecRemoveCode, "Error executing remove command", des)
}

// ErrSampleApp is the error for operations on the sample apps
func ErrSampleApp(err error) error {
	return errors.NewDefault(ErrSampleAppCode, "Error with sample app operation", err.Error())
}

// ErrCustomOperation is the error for custom operations
func ErrCustomOperation(err error) error {
	return errors.NewDefault(ErrCustomOperationCode, "Error with applying custom operation", err.Error())
}
