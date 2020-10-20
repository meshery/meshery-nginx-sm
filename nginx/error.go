package nginx

import (
	"fmt"

	"github.com/layer5io/gokit/errors"
)

var (
	ErrOpInvalid = errors.New(errors.ErrOpInvalid, "Invalid operation")
)

// ErrInstallMesh is the error for install mesh
func ErrInstallMesh(err error) error {
	return errors.New(errors.ErrInstallMesh, fmt.Sprintf("Error installing mesh: %s", err.Error()))
}

// ErrMeshConfig is the error for mesh config
func ErrMeshConfig(err error) error {
	return errors.New(errors.ErrMeshConfig, fmt.Sprintf("Error configuration mesh: %s", err.Error()))
}

// ErrPortForward is the error for mesh port forward
func ErrPortForward(err error) error {
	return errors.New(errors.ErrPortForward, fmt.Sprintf("Error portforwarding mesh gui: %s", err.Error()))
}

// ErrClientConfig is the error for setting client config
func ErrClientConfig(err error) error {
	return errors.New(errors.ErrClientConfig, fmt.Sprintf("Error setting client config: %s", err.Error()))
}

// ErrPortForward is the error for setting clientset
func ErrClientSet(err error) error {
	return errors.New(errors.ErrClientSet, fmt.Sprintf("Error setting clientset: %s", err.Error()))
}

// ErrStreamEvent is the error for streaming event
func ErrStreamEvent(err error) error {
	return errors.New(errors.ErrStreamEvent, fmt.Sprintf("Error streaming event: %s", err.Error()))
}

// ErrApplyOperation is the error for applying operation event
func ErrApplyOperation(err error) error {
	return errors.New(errors.ErrApplyOperation, fmt.Sprintf("Error applying operation: %s", err.Error()))
}

// ErrListOperations is the error for listing operations event
func ErrListOperations(err error) error {
	return errors.New(errors.ErrListOperations, fmt.Sprintf("Error listing operations: %s", err.Error()))
}
