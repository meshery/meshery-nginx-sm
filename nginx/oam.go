package nginx

import (
	"context"
	"fmt"
	"strings"

	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-nginx/nginx/oam"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
)

// ProcessOAM will handles the grpc invocation for handling OAM objects
func (h *Nginx) ProcessOAM(ctx context.Context, oamReq adapter.OAMRequest) (string, error) {
	var comps []v1alpha1.Component
	for _, acomp := range oamReq.OamComps {
		comp, err := oam.ParseApplicationComponent(acomp)
		if err != nil {
			h.Log.Error(ErrParseOAMComponent)
			continue
		}

		comps = append(comps, comp)
	}

	config, err := oam.ParseApplicationConfiguration(oamReq.OamConfig)
	if err != nil {
		h.Log.Error(ErrParseOAMConfig)
	}
	// If operation is delete then first HandleConfiguration and then handle the deployment
	if oamReq.DeleteOp {
		// Process configuration
		msg2, err := h.HandleApplicationConfiguration(config, oamReq.DeleteOp)
		if err != nil {
			return msg2, ErrProcessOAM(err)
		}

		// Process components
		msg1, err := h.HandleComponents(comps, oamReq.DeleteOp)
		if err != nil {
			return msg1 + "\n" + msg2, ErrProcessOAM(err)
		}

		return msg1 + "\n" + msg2, nil
	}
	// Process components
	msg1, err := h.HandleComponents(comps, oamReq.DeleteOp)
	if err != nil {
		return msg1, ErrProcessOAM(err)
	}

	// Process configuration
	msg2, err := h.HandleApplicationConfiguration(config, oamReq.DeleteOp)
	if err != nil {
		return msg1 + "\n" + msg2, ErrProcessOAM(err)
	}

	return msg1 + "\n" + msg2, nil
}

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Nginx, v1alpha1.Component, bool) (string, error)

func (h *Nginx) HandleComponents(comps []v1alpha1.Component, isDel bool) (string, error) {
	var errs []error
	var msgs []string

	compFuncMap := map[string]CompHandler{
		"NginxMesh": handleComponentNginxMesh,
	}
	for _, comp := range comps {
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			msg, err := handleNginxCoreComponents(h, comp, isDel, "", "")
			if err != nil {
				errs = append(errs, err)
				continue
			}

			msgs = append(msgs, msg)
			continue
		}

		msg, err := fnc(h, comp, isDel)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		msgs = append(msgs, msg)
	}
	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}
func (h *Nginx) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			msgs = append(msgs, fmt.Sprintf("applied trait \"%s\" on service \"%s\"", trait.Name, comp.ComponentName))
		}
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

func handleComponentNginxMesh(c *Nginx, comp v1alpha1.Component, isDelete bool) (string, error) {
	// Get the Nginx version from the settings
	// we are sure that the version of Nginx would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Settings["version"].(string)

	msg, err := c.installNginx(isDelete, version, comp.Namespace)
	if err != nil {
		return fmt.Sprintf("%s: %s", comp.Name, msg), err
	}

	return fmt.Sprintf("%s: %s", comp.Name, msg), nil
}

func handleNginxCoreComponents(
	c *Nginx,
	comp v1alpha1.Component,
	isDel bool,
	apiVersion,
	kind string) (string, error) {
	if apiVersion == "" {
		apiVersion = getAPIVersionFromComponent(comp)
		if apiVersion == "" {
			return "", ErrNginxCoreComponentFail(fmt.Errorf("failed to get API Version for: %s", comp.Name))
		}
	}

	if kind == "" {
		kind = getKindFromComponent(comp)
		if kind == "" {
			return "", ErrNginxCoreComponentFail(fmt.Errorf("failed to get kind for: %s", comp.Name))
		}
	}
	component := map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":        comp.Name,
			"annotations": comp.Annotations,
			"labels":      comp.Labels,
		},
		"spec": comp.Spec.Settings,
	}

	// Convert to yaml
	yamlByt, err := yaml.Marshal(component)
	if err != nil {
		err = ErrParseNginxCoreComponent(err)
		c.Log.Error(err)
		return "", err
	}

	msg := fmt.Sprintf("created %s \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	if isDel {
		msg = fmt.Sprintf("deleted %s config \"%s\" in namespace \"%s\"", kind, comp.Name, comp.Namespace)
	}

	return msg, c.MesheryKubeclient.ApplyManifest(yamlByt, mesherykube.ApplyOptions{
		Namespace: comp.Namespace,
		Update:    true,
		Delete:    isDel,
	})
}
func getAPIVersionFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sAPIVersion"]
}
func getKindFromComponent(comp v1alpha1.Component) string {
	return comp.Annotations["pattern.meshery.io.mesh.workload.k8sKind"]
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []string

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return fmt.Errorf(strings.Join(errMsgs, "\n"))
}

func mergeMsgs(strs []string) string {
	return strings.Join(strs, "\n")
}
