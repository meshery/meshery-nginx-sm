package config

import (
	"fmt"

	"github.com/layer5io/gokit/utils"
)

var (

	// server holds server configuration
	server = map[string]string{
		"name":    "nginx-adapter",
		"port":    "10010",
		"version": "v1.0.0",
	}

	// mesh holds mesh configuration
	mesh = map[string]string{
		"name":     "Nginx Service Mesh",
		"status":   "not installed",
		"traceurl": "http://localhost:14268/api/traces",
		"version":  "0.6.0",
	}

	// operations holds the supported operations inside mesh
	operations = map[string]interface{}{
		InstallNginxLatest: map[string]interface{}{
			"type": "0",
			"properties": map[string]string{
				"description": "Install Nginx service mesh (0.6.0)",
				"version":     "0.6.0",
			},
		},
		InstallSampleBookInfo: map[string]interface{}{
			"type": "1",
			"properties": map[string]string{
				"description": "Install BookInfo Application",
				"version":     "latest",
			},
		},
		ValidateSmiConformance: map[string]interface{}{
			"type": "3",
			"properties": map[string]string{
				"description": "Validate SMI conformance",
				"version":     "latest",
			},
		},
	}

	// Viper configuration
	filepath = fmt.Sprintf("%s/.meshery/nginx", utils.GetHome())
	filename = "config"
	filetype = "yaml"
)
