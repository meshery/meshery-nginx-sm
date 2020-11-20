package config

import (
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
)

var (
	ServiceName = "service_name"
)

func getOperations(dev adapter.Operations) adapter.Operations {

	dev[NginxOperation] = &adapter.Operation{
		Type:        int32(meshes.OpCategory_INSTALL),
		Description: "Nginx Service Mesh",
		Versions: []adapter.Version{
			"0.6.0",
		},
		Templates: []adapter.Template{
			"templates/nginx.yaml",
		},
		AdditionalProperties: map[string]string{
			ServiceName: NginxOperation,
		},
	}

	return dev
}
