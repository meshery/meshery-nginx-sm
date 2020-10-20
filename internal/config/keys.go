package config

const (
	local = "local" // local is the key for local config

	// Operation keys
	InstallNginxLatest = "install-nginx-060" // InstallNginxLatest is the key to install nginx

	InstallSampleBookInfo = "install-sample-bookinfo" // InstallSampleBookInfo is the key to install sample bookinfo application

	ValidateSmiConformance = "validate-smi-conformance" // ValidateSmiConformance is the key to run and validate smi conformance test

	RunningMeshVersion = "running_mesh_version" // RunningMeshVersion is the key to store the current running version of the mesh
)
