// Package nginx - Common operations for the adapter
package nginx

import (
	"github.com/layer5io/meshery-adapter-library/status"
)

func (nginx *Nginx) applyCustomOperation(namespace string, manifest string, isDel bool, kubeconfigs []string) (string, error) {
	st := status.Starting

	err := nginx.applyManifest([]byte(manifest), isDel, namespace, kubeconfigs)
	if err != nil {
		return st, ErrCustomOperation(err)
	}

	return status.Completed, nil
}
