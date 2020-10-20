package nginx

// Spec holds the specifications for nginx adapter
type Spec struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

// GetName returns the name of the mesh
func (h *handler) GetName() string {
	spec := &Spec{}
	err := h.config.Mesh(&spec)
	if err != nil {
		h.log.Err("1000", err.Error())
		return "Not set"
	}
	return spec.Name
}
