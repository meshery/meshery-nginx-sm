package nginx

import (
	"fmt"
	"os"
	"os/exec"
)

// MeshInstance holds the information of the instance of the mesh
type MeshInstance struct {
	InstallRegistry string `json:"installregistry,omitempty"`
	InstallVersion  string `json:"installversion,omitempty"`
	MgmtAddr        string `json:"mgmtaddr,omitempty"`
	Nginxaddr       string `json:"nginxaddr,omitempty"`
}

// CreateInstance installs and creates a mesh environment up and running
func (h *handler) installNginx(del bool, version string, registry string) (string, error) {
	status := "installing"

	if del {
		status = "removing"
	}

	meshinstance := &MeshInstance{
		InstallVersion:  version,
		InstallRegistry: registry,
	}
	err := h.config.Mesh(meshinstance)
	if err != nil {
		return status, ErrMeshConfig(err)
	}

	h.log.Info("Installing Nginx")
	err = meshinstance.installUsingNginxctl(del)
	if err != nil {
		h.log.Err("Nginx installation failed", ErrInstallMesh(err).Error())
		return status, ErrInstallMesh(err)
	}
	if del {
		return "removed", nil
	}

	return "deployed", nil
}

// installSampleApp installs and creates a sample bookinfo application up and running
func (h *handler) installSampleApp(name string) (string, error) {
	// Needs implementation
	return "deployed", nil
}

// installMesh installs the mesh in the cluster or the target location
func (m *MeshInstance) installUsingNginxctl(del bool) error {

	Executable, err := exec.LookPath("./scripts/deploy.sh")
	if err != nil {
		return err
	}

	if del {
		Executable, err = exec.LookPath("./scripts/delete.sh")
		if err != nil {
			return err
		}
	}

	cmd := &exec.Cmd{
		Path:   Executable,
		Args:   []string{Executable},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("NGINX_DOCKER_REGISTRY=%s", m.InstallRegistry),
		fmt.Sprintf("NGINX_MESH_VER=%s", m.InstallVersion),
	)

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (m *MeshInstance) portForward() error {
	// Needs implementation
	return nil
}
