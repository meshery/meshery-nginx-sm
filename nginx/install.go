package nginx

import (
	"fmt"
	"strings"
	"sync"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"

	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

const (
	repo  = "https://helm.nginx.com/stable"
	chart = "nginx-service-mesh"
)

// Installs NGINX service mesh using helm charts.
// Unlike other adapters, doesn't keep CLI as a fallback method
func (nginx *Nginx) installNginx(del bool, version, namespace string, kubeconfigs []string) (string, error) {
	nginx.Log.Debug(fmt.Sprintf("Requested install of version: %s", version))
	nginx.Log.Debug(fmt.Sprintf("Requested action is delete: %v", del))
	nginx.Log.Debug(fmt.Sprintf("Requested action is in namespace: %s", namespace))

	st := status.Installing
	if del {
		st = status.Removing
	}

	err := nginx.Config.GetObject(adapter.MeshSpecKey, nginx)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	err = nginx.applyHelmChart(del, version, namespace, kubeconfigs)
	if err != nil {
		nginx.Log.Error(ErrInstallNginx(err))
		return st, ErrInstallNginx(err)
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}

func (nginx *Nginx) applyHelmChart(del bool, version, namespace string, kubeconfigs []string) error {
	chartVersion, err := mesherykube.HelmAppVersionToChartVersion(repo, chart, version)
	if err != nil {
		version = strings.TrimPrefix(version, "v")
		chartVersion, err = mesherykube.HelmAppVersionToChartVersion(repo, chart, version)
		if err != nil {
			return ErrApplyHelmChart(err)
		}
	}

	var act mesherykube.HelmChartAction
	if del {
		act = mesherykube.UNINSTALL
	} else {
		act = mesherykube.INSTALL
	}

	// Set namespace to "nginx-system" if it is undefined, "default", or "meshery".
	// NGINX SM should be in it's own namespace, so that's why these namespaces are overridden.
	forbiddenNamespaces := []string{"", "default", "meshery"}
	for _, n := range forbiddenNamespaces {
		if strings.ToLower(namespace) == n {
			namespace = "nginx-system"
			break
		}
	}

	// Set deployment override flag to disable automatic sidecar injection in the meshery namespace.
	// This is to prevent Meshery from having connectivy issues with other meshes or non-meshed services.
	// This is equal to using the Helm flag: --set autoInjection.disabledNamespaces={"meshery"}
	overrideVal := map[string]interface{}{
		"autoInjection": map[string]interface{}{
			"disabledNamespaces": []string{"meshery"},
		},
	}

	// Create Helm config used to install charts.
	c := mesherykube.ApplyHelmChartConfig{
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: repo,
			Chart:      chart,
			Version:    chartVersion,
		},
		Namespace:       namespace,
		Action:          act,
		CreateNamespace: true,
		ReleaseName:     chart,
		OverrideValues:  overrideVal,
	}

	nginx.Log.Info(fmt.Sprintf("Installing NGINX Service Mesh %s using Helm chart: %+v\n", version, c))

	var wg sync.WaitGroup
	var errs []error
	var errMx sync.Mutex
	
	for _, config := range kubeconfigs {
		wg.Add(1)
		go func(config string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(config))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

			// Install Helm chart.
			err = kClient.ApplyHelmChart(c)
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}
		}(config)
	}
	wg.Wait()
	if len(errs) == 0 {
		return nil
	}

	mergedErrors := mergeErrors(errs)
	return ErrApplyHelmChart(mergedErrors)
}

func (nginx *Nginx) applyManifest(manifest []byte, isDel bool, namespace string, kubeconfigs []string) error {
	var wg sync.WaitGroup
	var errs []error
	var errMx sync.Mutex

	for _, config := range kubeconfigs {
		wg.Add(1)
		go func(config string) {
			defer wg.Done()
			kClient, err := mesherykube.New([]byte(config))
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

			err = kClient.ApplyManifest(manifest, mesherykube.ApplyOptions{
				Namespace: namespace,
				Update: true,
				Delete: isDel,
			})
			if err != nil {
				errMx.Lock()
				errs = append(errs, err)
				errMx.Unlock()
				return
			}

		}(config)
	}

	wg.Wait()
	if len(errs) == 0 {
		return nil
	}

	return mergeErrors(errs)
}
