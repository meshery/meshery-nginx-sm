package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/layer5io/meshery-nginx/nginx"
	"github.com/layer5io/meshery-nginx/nginx/oam"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/manifests"

	// "github.com/layer5io/meshkit/tracing"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/api/grpc"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-nginx/internal/config"
	smp "github.com/layer5io/service-mesh-performance/spec"
)

var (
	serviceName = "nginx-adaptor"
)

type Data struct {
	Name         string //`json:"name"`
	Path         string //`json:"path"`
	Sha          string //`json:"sha"`
	Size         int    //`json:"size"`
	Url          string //`json:"url"`
	Html_url     string //`json:"html_url"`
	Git_url      string //`json:"git_url"`
	Download_url string //`json:"download_url"`
	Types        string //`json:"type"`
	Link         string //`json:"link"`
}

// creates the ~/.meshery directory
func init() {
	err := os.MkdirAll(path.Join(config.RootPath(), "bin"), 0750)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// main is the entrypoint of the adaptor
func main() {

	// Initialize Logger instance
	log, err := logger.New(serviceName, logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Setenv("KUBECONFIG", path.Join(
		config.KubeConfig[configprovider.FilePath],
		fmt.Sprintf("%s.%s", config.KubeConfig[configprovider.FileName], config.KubeConfig[configprovider.FileType])),
	)
	if err != nil {
		// Fail silently
		log.Warn(err)
	}

	// Initialize application specific configs and dependencies
	// App and request config
	cfg, err := config.New(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	service := &grpc.Service{}
	err = cfg.GetObject(adapter.ServerKey, service)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	kubeconfigHandler, err := config.NewKubeconfigBuilder(configprovider.ViperKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// // Initialize Tracing instance
	// tracer, err := tracing.New(service.Name, service.TraceURL)
	// if err != nil {
	// 	log.Err("Tracing Init Failed", err.Error())
	// 	os.Exit(1)
	// }

	// Initialize Handler intance
	handler := nginx.New(cfg, log, kubeconfigHandler)
	handler = adapter.AddLogger(log, handler)

	service.Handler = handler
	service.Channel = make(chan interface{}, 10)
	service.StartedAt = time.Now()

	go registerCapabilities(service.Port, log)        //Registering static capabilities
	go registerDynamicCapabilities(service.Port, log) //Registering latest capabilities periodically

	// Server Initialization
	log.Info("Adapter Listening at port: ", service.Port)
	err = grpc.Start(service, nil)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func mesheryServerAddress() string {
	meshReg := os.Getenv("MESHERY_SERVER")

	if meshReg != "" {
		if strings.HasPrefix(meshReg, "http") {
			return meshReg
		}

		return "http://" + meshReg
	}

	return "http://localhost:9081"
}

func serviceAddress() string {
	svcAddr := os.Getenv("SERVICE_ADDR")

	if svcAddr != "" {
		return svcAddr
	}

	return "mesherylocal.layer5.io"
}

func registerCapabilities(port string, log logger.Handler) {
	// Register workloads
	if err := oam.RegisterWorkloads(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Info(err.Error())
	}

	// Register traits
	if err := oam.RegisterTraits(mesheryServerAddress(), serviceAddress()+":"+port); err != nil {
		log.Info(err.Error())
	}
}

func registerDynamicCapabilities(port string, log logger.Handler) {
	registerWorkloads(port, log)
	//Start the ticker
	const reRegisterAfter = 24
	ticker := time.NewTicker(reRegisterAfter * time.Hour)
	for {
		<-ticker.C
		registerWorkloads(port, log)
	}

}
func registerWorkloads(port string, log logger.Handler) {
	release, err := config.GetLatestReleases(1)
	if err != nil {
		log.Info("Could not get latest version")
		return
	}
	version := release[0].TagName
	log.Info("Registering latest workload components for version ", version)

	str, err := ChangeReleaseString()
	if err != nil {
		log.Info("Could not change the version string")
	}

	// Register workloads
	if err := adapter.RegisterWorkLoadsDynamically(mesheryServerAddress(), serviceAddress()+":"+port, &adapter.DynamicComponentsConfig{
		TimeoutInMinutes: 60,
		URL:              "https://github.com/nginxinc/helm-charts/blob/master/stable/nginx-service-mesh-" + str + ".tgz?raw=true",
		GenerationMethod: adapter.HelmCHARTS,
		Config: manifests.Config{
			Name:        smp.ServiceMesh_Type_name[int32(smp.ServiceMesh_NGINX_SERVICE_MESH)],
			MeshVersion: version,
			Filter: manifests.CrdFilter{
				RootFilter:    []string{"$[?(@.kind==\"CustomResourceDefinition\")]"},
				NameFilter:    []string{"$..[\"spec\"][\"names\"][\"kind\"]"},
				VersionFilter: []string{"$..spec.versions[0]", " --o-filter", "$[0]"},
				GroupFilter:   []string{"$..spec", " --o-filter", "$[]"},
				SpecFilter:    []string{"$..openAPIV3Schema.properties.spec", " --o-filter", "$[]"},
			},
		},
		Operation: config.NginxOperation,
	}); err != nil {
		log.Info(nginx.ErrRegisteringWorkload(err))
		return
	}
	log.Info("Latest workload components successfully registered.")
}

func ChangeReleaseString() (string, error) {
	url := "https://api.github.com/repos/nginxinc/helm-charts/contents/stable?raw=true"

	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var p []Data
	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err.Error())
	}
	length := len(p)
	res1 := strings.Replace(p[length-1].Name, "nginx-service-mesh-", "", 1)
	version := strings.Replace(res1, ".tgz", "", 1)
	return version, nil
}
