package config

import (
	"github.com/layer5io/meshery-adapter-library/common"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
)

var (
	ProductionConfig = configprovider.Options{
		ServerConfig:   productionServerConfig,
		MeshSpec:       productionMeshSpec,
		ProviderConfig: productionProviderConfig,
		Operations:     productionOperations,
	}

	productionServerConfig = map[string]string{
		"name":    "nginx-adapter",
		"port":    "10010",
		"version": "v1.0.0",
	}

	productionMeshSpec = map[string]string{
		"name":     "nginx",
		"status":   "none",
		"traceurl": "none",
		"version":  "none",
	}

	productionProviderConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "nginx",
	}

	// Controlling the kubeconfig lifecycle with viper
	productionKubeConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "",
		configprovider.FileName: "kubeconfig",
	}

	productionOperations = getOperations(common.Operations)
)
