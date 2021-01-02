package config

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/layer5io/meshery-adapter-library/common"
	"github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshery-adapter-library/status"
	"github.com/layer5io/meshkit/utils"
)

const (
	NginxOperation = "nginx"
	LabelNamespace = "label-namespace"
)

var (
	configRootPath  = path.Join(utils.GetHome(), ".meshery")
	NginxExecutable = "nginx-meshctl"

	Config = configprovider.Options{
		ServerConfig:   ServerConfig,
		MeshSpec:       MeshSpec,
		ProviderConfig: ProviderConfig,
		Operations:     Operations,
	}

	ServerConfig = map[string]string{
		"name":    "nginx-adapter",
		"port":    "10010",
		"version": "v1.0.0",
	}

	MeshSpec = map[string]string{
		"name":     "nginx",
		"status":   status.None,
		"traceurl": status.None,
		"version":  status.None,
	}

	ProviderConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "nginx",
	}

	// KubeConfig - Controlling the kubeconfig lifecycle with viper
	KubeConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kubeconfig",
	}

	Operations = getOperations(common.Operations)
)

// New creates a new config instance
func New(provider string) (config.Handler, error) {
	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(Config)
	case configprovider.InMemKey:
		return configprovider.NewInMem(Config)
	}

	return nil, ErrEmptyConfig
}

func NewKubeconfigBuilder(provider string) (config.Handler, error) {
	opts := configprovider.Options{}
	opts.ProviderConfig = KubeConfig

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, ErrEmptyConfig
}

// RootPath returns the config root path for the adapter
func RootPath() string {
	return configRootPath
}

// InitialiseNSMCtl looks for the "nginx-meshctl" in the $PATH
// if it doesn't finds it then it downloads and installs it in
// "configRootPath" and returns the path for the executable
func InitialiseNSMCtl() error {
	// Look for the executable in the path
	_, err := exec.LookPath(NginxExecutable)
	if err == nil {
		return nil
	}

	executable := path.Join(configRootPath, "nginx-meshctl")
	if _, err := os.Stat(executable); err == nil {
		return nil
	}

	// Proceed to download the binary in the config root path
	res, err := downloadBinary(runtime.GOOS)
	if err != nil {
		return err
	}

	// Install the binary
	if err = installBinary(NginxExecutable, runtime.GOOS, res); err != nil {
		return err
	}

	return nil
}

func downloadBinary(platform string) (*http.Response, error) {
	var url = "http://downloads08.f5.com/esd/download.sv?loc=downloads08.f5.com/downloads/"
	switch platform {
	case "darwin":
		url += "60a64c24-957a-40cc-9ca2-daf9f078a409/nginx-meshctl_darwin.gz"
	case "windows":
		url += "c88ac035-189e-4ee5-8f26-13580f42e492/nginx-meshctl_windows.exe"
	case "linux":
		url += "4c110d24-86be-4452-97b3-f802394a82cd/nginx-meshctl_linux.gz"
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status: %s", resp.Status)
	}

	return resp, nil
}

func installBinary(location, platform string, res *http.Response) error {
	out, err := os.Create(location)
	if err != nil {
		return err
	}

	switch platform {
	case "darwin":
		fallthrough
	case "linux":
		r, err := gzip.NewReader(res.Body)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, r)
		if err != nil {
			return err
		}

		err = r.Close()
		if err != nil {
			return err
		}

		if err = out.Chmod(0755); err != nil {
			return ErrInstallBinary(err)
		}
	case "windows":
	}

	// Close the response body
	err = out.Close()
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
