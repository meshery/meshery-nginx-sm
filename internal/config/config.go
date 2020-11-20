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

	"github.com/layer5io/meshery-adapter-library/config"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
	"github.com/layer5io/meshkit/utils"
)

const (
	NginxOperation = "nginx"
	Development    = "development"
	Production     = "production"
)

var (
	configRootPath = path.Join(utils.GetHome(), ".meshery")

	kubeconfigFilename = "kubeconfig"
	kubeconfigFiletype = "yaml"
	KubeconfigPath     = path.Join(configRootPath, fmt.Sprintf("%s.%s", kubeconfigFilename, kubeconfigFiletype))
)

// New creates a new config instance
func New(provider string) (config.Handler, error) {

	// Default config
	opts := configprovider.Options{}
	environment := os.Getenv("MESHERY_ENV")
	if len(environment) < 1 {
		environment = Development
	}

	// Config environment
	switch environment {
	case Production:
		opts = ProductionConfig
	case Development:
		opts = DevelopmentConfig
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}

	return nil, ErrEmptyConfig
}

func NewKubeconfigBuilder(provider string) (config.Handler, error) {

	opts := configprovider.Options{}
	environment := os.Getenv("MESHERY_ENV")
	if len(environment) < 1 {
		environment = Development
	}

	// Config environment
	switch environment {
	case Production:
		opts.ProviderConfig = productionKubeConfig
	case Development:
		opts.ProviderConfig = developmentKubeConfig
	}

	// Config provider
	switch provider {
	case configprovider.ViperKey:
		return configprovider.NewViper(opts)
	case configprovider.InMemKey:
		return configprovider.NewInMem(opts)
	}
	return nil, ErrEmptyConfig
}

// InitialiseNSMCtl looks for the "nginx-meshctl" in the $PATH
// if it doesn't finds it then it downloads and installs it in
// "configRootPath" and returns the path for the executable
func InitialiseNSMCtl() (string, error) {
	// Look for the executable in the path
	fmt.Println("Looking for nginx-meshctl in the path...")
	executable, err := exec.LookPath("nginx-meshctl")

	fmt.Println("Looking for nginx-meshctl in", configRootPath, "...")
	executable = path.Join(configRootPath, "nginx-meshctl")
	if _, err := os.Stat(executable); err == nil {
		return executable, nil
	}

	// Proceed to download the binary in the config root path
	fmt.Println("nginx-meshctl not found in the path, downloading...")
	res, err := downloadBinary(runtime.GOOS)
	if err != nil {
		return "", err
	}

	// Install the binary
	fmt.Println("Installing...")
	if err = installBinary(path.Join(configRootPath, "nginx-meshctl"), runtime.GOOS, res); err != nil {
		return "", err
	}

	fmt.Println("Done")
	return path.Join(configRootPath, "nginx-meshctl"), nil
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
	defer out.Close()

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

		r.Close()

		out.Chmod(0755)
	case "windows":
	}

	// Close the response body
	res.Body.Close()
	return nil
}
