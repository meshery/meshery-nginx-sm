package nginx

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"

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
		nginx.Executable,
		"deploy",
		"--nginx-mesh-api-image",
		fmt.Sprintf("%s/nginx-mesh-api:%s", nginx.DockerRegistry, version),
		"--nginx-mesh-sidecar-image",
		fmt.Sprintf("%s/nginx-mesh-sidecar:%s", nginx.DockerRegistry, version),
		"--nginx-mesh-init-image",
		fmt.Sprintf("%s/nginx-mesh-init:%s", nginx.DockerRegistry, version),
		"--nginx-mesh-metrics-image",
		fmt.Sprintf("%s/nginx-mesh-metrics:%s", nginx.DockerRegistry, version),
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
	cmd := exec.Command(nginx.Executable, "remove", "-y")
	cmd.Stderr = &er

	if err := cmd.Run(); err != nil {
		return ErrExecRemove(err, er.String())
	}

	// TODO: Remove sidecar proxy from deployments
	return nil
}

func (nginx *Nginx) applyManifest(manifest []byte) error {
	kclient, err := mesherykube.New(nginx.KubeClient, nginx.RestConfig)
	if err != nil {
		return err
	}

	err = kclient.ApplyManifest(manifest, mesherykube.ApplyOptions{})
	if err != nil {
		return err
	}

	return nil
}
