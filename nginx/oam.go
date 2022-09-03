package nginx

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/layer5io/meshery-adapter-library/adapter"
	"github.com/layer5io/meshery-adapter-library/meshes"
	"github.com/layer5io/meshery-nginx/nginx/oam"
	"github.com/layer5io/meshery-traefik-mesh/internal/config"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"gopkg.in/yaml.v2"
)

// ProcessOAM will handles the grpc invocation for handling OAM objects
func (nginx *Nginx) ProcessOAM(ctx context.Context, oamReq adapter.OAMRequest, hchan *chan interface{}) (string, error) {
	nginx.SetChannel(hchan)
	err := nginx.CreateKubeconfigs(oamReq.K8sConfigs)
	if err != nil {
		return "", err
	}
	kubeConfigs := oamReq.K8sConfigs

	var comps []v1alpha1.Component
	for _, acomp := range oamReq.OamComps {
		comp, err := oam.ParseApplicationComponent(acomp)
		if err != nil {
			nginx.Log.Error(ErrParseOAMComponent)
			continue
		}

		comps = append(comps, comp)
	}

	config, err := oam.ParseApplicationConfiguration(oamReq.OamConfig)
	if err != nil {
		nginx.Log.Error(ErrParseOAMConfig)
	}
	// If operation is delete then first HandleConfiguration and then handle the deployment
	if oamReq.DeleteOp {
		// Process configuration
		msg2, err := nginx.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeConfigs)
		if err != nil {
			return msg2, ErrProcessOAM(err)
		}

		// Process components
		msg1, err := nginx.HandleComponents(comps, oamReq.DeleteOp, kubeConfigs)
		if err != nil {
			return msg1 + "\n" + msg2, ErrProcessOAM(err)
		}

		return msg1 + "\n" + msg2, nil
	}
	// Process components
	msg1, err := nginx.HandleComponents(comps, oamReq.DeleteOp, kubeConfigs)
	if err != nil {
		return msg1, ErrProcessOAM(err)
	}

	// Process configuration
	msg2, err := nginx.HandleApplicationConfiguration(config, oamReq.DeleteOp, kubeConfigs)
	if err != nil {
		return msg1 + "\n" + msg2, ErrProcessOAM(err)
	}

	return msg1 + "\n" + msg2, nil
}

// CompHandler is the type for functions which can handle OAM components
type CompHandler func(*Nginx, v1alpha1.Component, bool, []string) (string, error)

//HandleComponents handles the parsed oam components from pattern file
func (nginx *Nginx) HandleComponents(comps []v1alpha1.Component, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string
	stat1 := "deploying"
	stat2 := "deployed"
	if isDel {
		stat1 = "removing"
		stat2 = "removed"
	}

	compFuncMap := map[string]CompHandler{
		"NginxMesh": handleComponentNginxMesh,
	}
	for _, comp := range comps {
		ee := &meshes.EventsResponse{
			OperationId:   uuid.New().String(),
			Component:     config.ServerConfig["type"],
			ComponentName: config.ServerConfig["name"],
		}
		fnc, ok := compFuncMap[comp.Spec.Type]
		if !ok {
			msg, err := handleNginxCoreComponents(nginx, comp, isDel, "", "", kubeconfigs)
			if err != nil {
				ee.Summary = fmt.Sprint("Error while %s %s", stat1, comp.Spec.Type)
				mesh.streamErr(ee.Summary, ee, err)
				errs = append(errs, err)
				continue
			}
			ee.Summary = fmt.Sprintf("%s %s successfully", comp.Spec.Type, stat2)
			ee.Details = fmt.Sprintf("The %s is now %s.", comp.Spec.Type, stat2)
			mesh.StreamInfo(ee)
			msgs = append(msgs, msg)
			continue
		}

		msg, err := fnc(nginx, comp, isDel, kubeconfigs)
		if err != nil {
			ee.Summary = fmt.Sprintf("Error while %s %s", stat1, comp.Spec.Type)
			mesh.streamErr(ee.Summary, ee, err)
			errs = append(errs, err)
			continue
		}
		ee.Summary = fmt.Sprintf("%s %s %s successfully", comp.Name, comp.Spec.Type, stat2)
		ee.Details = fmt.Sprintf("The %s %s is now %s.", comp.Name, comp.Spec.Type, stat2)
		mesh.StreamInfo(ee)
		msgs = append(msgs, msg)
	}
	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

// HandleApplicationConfiguration handles the processing of OAM application configuration
func (nginx *Nginx) HandleApplicationConfiguration(config v1alpha1.Configuration, isDel bool, kubeconfigs []string) (string, error) {
	var errs []error
	var msgs []string
	for _, comp := range config.Spec.Components {
		for _, trait := range comp.Traits {
			/*
				Handling of different traits will be added here in code
			*/
			msgs = append(msgs, fmt.Sprintf("applied trait \"%s\" on service \"%s\"", trait.Name, comp.ComponentName))
		}
	}

	if err := mergeErrors(errs); err != nil {
		return mergeMsgs(msgs), err
	}

	return mergeMsgs(msgs), nil
}

func handleComponentNginxMesh(c *Nginx, comp v1alpha1.Component, isDelete bool, kubeconfigs []string) (string, error) {
	// Get the Nginx version from the settings
	// we are sure that the version of Nginx would be present
	// because the configuration is already validated against the schema
	version := comp.Spec.Version
	if version == "" {
		return "", fmt.Errorf("pass valid version inside service for Nginx SM installation")
	}
	//TODO: When no version is passed in service, use the latest Nginx SM version

	msg, err := c.installNginx(isDelete, version, comp.Namespace, kubeconfigs)
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
	kind string,
	kubeconfigs []string) (string, error) {
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

	return msg, c.applyManifest(yamlByt, isDel, comp.Namespace, kubeconfigs)
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
