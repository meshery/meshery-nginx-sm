package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

func getOperations(dev adapter.Operations) adapter.Operations {

	dev[NginxOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "Nginx Service Mesh",
		Versions: []adapter.Version{
			"0.6.0",
		},
	}

	return dev
}
