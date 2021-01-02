package nginx

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	// ErrInstallNginxCode provisioning failure
	ErrInstallNginxCode = "nginx_test_code"
	// ErrMeshConfigCode   service mesh configuration failure
	ErrMeshConfigCode   = "nginx_test_code"
	// ErrClientConfigCode adapter configuration failure
	ErrClientConfigCode = "nginx_test_code"
	// ErrStreamEventCode  failure
	ErrStreamEventCode  = "nginx_test_code"
	// ErrExecDeployCode   failure
	ErrExecDeployCode   = "nginx_test_code"
	// ErrExecRemoveCode   failure
	ErrExecRemoveCode   = "nginx_test_code"
	// ErrSampleAppCode    failure
	ErrSampleAppCode    = "nginx_test_code"
	// ErrOpInvalid failure
	ErrOpInvalid = errors.NewDefault(errors.ErrOpInvalid, "Invalid operation")
)

// ErrInstallNginx is the error for install mesh
func ErrInstallNginx(err error) error {
	return errors.NewDefault(ErrInstallNginxCode, fmt.Sprintf("Error installing nginx: %s", err.Error()))
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.NewDefault(ErrMeshConfigCode, fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.NewDefault(ErrClientConfigCode, fmt.Sprintf("Error setting client config: %s", err.Error()))
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.NewDefault(ErrStreamEventCode, fmt.Sprintf("Error streaming event: %s", err.Error()))
}

// ErrExecDeploy is the error for deploying nginx service mesh
func ErrExecDeploy(err error, des string) error {
	return errors.NewDefault(ErrExecDeployCode, fmt.Sprintf("Error executing deploy command: %s", des))
}

// ErrExecRemove is the error for removing nginx service mesh
func ErrExecRemove(err error, des string) error {
	return errors.NewDefault(ErrExecRemoveCode, fmt.Sprintf("Error executing remove command: %s", des))
}

// ErrSampleApp is the error for operations on the sample apps
func ErrSampleApp(err error) error {
	return errors.NewDefault(ErrSampleAppCode, fmt.Sprintf("Error with sample app operation: %s", err.Error()))
}
