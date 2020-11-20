package nginx

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/status"
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
	var er bytes.Buffer

	cmd := exec.Command(
		nginx.Executable,
		"deploy",
		fmt.Sprintf("--nginx-mesh-api-image \"%s/nginx-mesh-api:%s\"", nginx.DockerRegistry, version),
		fmt.Sprintf("--nginx-mesh-sidecar-image \"%s/nginx-mesh-sidecar:%s\"", nginx.DockerRegistry, version),
		fmt.Sprintf("--nginx-mesh-init-image \"%s/nginx-mesh-init:%s\"", nginx.DockerRegistry, version),
		fmt.Sprintf("--nginx-mesh-metrics-image \"%s/nginx-mesh-metrics:%s\"", nginx.DockerRegistry, version),
	)
	cmd.Stderr = &er

	if err := cmd.Run(); err != nil {
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
