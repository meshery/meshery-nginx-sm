package build

import (
	"errors"
	"fmt"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-nginx/nginx/oam"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"
	"github.com/layer5io/meshkit/utils/walker"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var DefaultVersion string
var DefaultGenerationURL string
var DefaultGenerationMethod string
var WorkloadPath string

const (
	repo  = "https://helm.nginx.com/stable"
	chart = "nginx-service-mesh"
)

// NewConfig creates the configuration for creating components
func NewConfig(version string) manifests.Config {
	return manifests.Config{
		Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_NGINX_SERVICE_MESH)],
		MeshVersion: version,
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			// trim the spaces
			for _, crd := range crds {
				crd = strings.TrimSpace(crd)
			}
			return crds
		},
	}
}

func getLatestVersion() (string, error) {
	filename := []string{}
	if err := walker.NewGit().
		Owner("nginxinc").
		Repo("helm-charts").
		Branch("master").
		Root("stable/").
		RegisterFileInterceptor(func(f walker.File) error {
			if strings.HasSuffix(f.Name, ".tgz") && strings.HasPrefix(f.Name, "nginx-service-mesh") {
				filename = append(filename, strings.TrimSuffix(strings.TrimPrefix(f.Name, "nginx-service-mesh-"), ".tgz"))
			}
			return nil
		}).Walk(); err != nil {
		return "", err
	}
	filename = utils.SortDottedStringsByDigits(filename)
	if len(filename) == 0 {
		return "", errors.New("no files found")
	}
	return filename[len(filename)-1], nil
}
func init() {
	version, err := getLatestVersion()
	if err != nil {
		fmt.Println("could not get chart version ", err.Error())
		return
	}
	DefaultVersion, err = kubernetes.HelmChartVersionToAppVersion(repo, chart, version)
	if err != nil {
		fmt.Println("could not get version ", err.Error())
		return
	}
	DefaultGenerationURL = "https://github.com/nginxinc/helm-charts/blob/master/stable/nginx-service-mesh-" + version + ".tgz?raw=true"
	DefaultGenerationMethod = adapter.HelmCHARTS
	WorkloadPath = oam.WorkloadPath
}
