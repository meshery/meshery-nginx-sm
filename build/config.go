package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/manifests"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultGenerationMethod string
var DefaultGenerationURL string
var LatestVersion string
var MeshModelPath string
var AllVersions []string

const Component = "NGINX Service Mesh"

var Meshmodelmetadata = make(map[string]interface{})

var MeshModelConfig = adapter.MeshModelConfig{ //Move to build/config.go
	Category: "Cloud Native Network",
	Metadata: map[string]interface{}{},
}

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_NGINX_SERVICE_MESH)],
		Type:        Component,
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
		},
	}
}

func init() {
	f, _ := os.Open("./build/meshmodel_metadata.json")
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()
	byt, _ := io.ReadAll(f)

	_ = json.Unmarshal(byt, &MeshModelConfig.Metadata)
	wd, _ := os.Getwd()
	MeshModelPath = filepath.Join(wd, "templates", "meshmodel", "components")
	AllVersions, _ = utils.GetLatestReleaseTagsSorted("nginxinc", "nginx-service-mesh")
	if len(AllVersions) == 0 {
		return
	}
	// @TODO need to update this, because NGINX changes how they relate the version of the NGINX
	// Service Mesh to the official release/branch which is v1.7.0 -> v0.7.0
	// LatestVersion = AllVersions[len(AllVersions)-1]
	LatestVersion = "0.7.0"
	DefaultGenerationMethod = adapter.HelmCHARTS
	DefaultGenerationURL = "https://github.com/nginxinc/helm-charts/blob/master/stable/nginx-service-mesh-" + LatestVersion + ".tgz?raw=true"
}
