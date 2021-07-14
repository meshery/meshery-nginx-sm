package nginx

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"

	internalconfig "github.com/layer5io/meshery-nginx/internal/config"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (nginx *Nginx) installNginx(del bool, version string) (string, error) {
	st := status.Installing
	if del {
		st = status.Removing
	}

	err := nginx.Config.GetObject(adapter.MeshSpecKey, nginx)
	if err != nil {
		return st, ErrMeshConfig(err)
	}

	if del {
		err = nginx.runUninstallCmd()
	} else {
		err = nginx.runInstallCmd(version)
	}
	if err != nil {
		nginx.Log.Error(ErrInstallNginx(err))
		return st, ErrInstallNginx(err)
	}

	if del {
		return status.Removed, nil
	}
	return status.Installed, nil
}

func (nginx *Nginx) runInstallCmd(version string) error {
	var er, out bytes.Buffer

	cmd := exec.Command(
		internalconfig.NginxExecutable,
		"deploy",
		"--nginx-mesh-api-image",
		fmt.Sprintf("nginx/nginx-mesh-api:%s", version),
		"--nginx-mesh-sidecar-image",
		fmt.Sprintf("nginx/nginx-mesh-sidecar:%s", version),
		"--nginx-mesh-init-image",
		fmt.Sprintf("nginx/nginx-mesh-init:%s", version),
		"--nginx-mesh-metrics-image",
		fmt.Sprintf("nginx/nginx-mesh-metrics:%s", version),
	)
	cmd.Stderr = &er
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		nginx.Log.Debug(out.String())
		return ErrExecDeploy(err, er.String())
	}

	return nil
}

func (nginx *Nginx) runUninstallCmd() error {
	var er bytes.Buffer

	// Remove the service mesh from kubernetes
	cmd := exec.Command(internalconfig.NginxExecutable, "remove", "-y")
	cmd.Stderr = &er

	if err := cmd.Run(); err != nil {
		return ErrExecRemove(err, er.String())
	}

	// TODO: Remove sidecar proxy from deployments
	return nil
}

func (nginx *Nginx) applyManifest(manifest []byte, isDel bool, namespace string) error {
	err := nginx.MesheryKubeclient.ApplyManifest(manifest, mesherykube.ApplyOptions{
		Namespace: namespace,
		Update:    true,
		Delete:    isDel,
	})
	if err != nil {
		return err
	}

	return nil
}
