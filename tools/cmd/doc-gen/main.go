/*
Copyright 2019 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"

	"kube-aggregation/pkg/apiserver/runtime"
	"kube-aggregation/pkg/constants"
	"kube-aggregation/pkg/informers"
	monitoringv1alpha3 "kube-aggregation/pkg/kapis/monitoring/v1alpha3"
	resourcesv1alpha2 "kube-aggregation/pkg/kapis/resources/v1alpha2"
	resourcesv1alpha3 "kube-aggregation/pkg/kapis/resources/v1alpha3"
	"kube-aggregation/pkg/simple/client/k8s"
	"kube-aggregation/pkg/version"
)

var output string

func init() {
	flag.StringVar(&output, "output", "./api/ks-openapi-spec/swagger.json", "--output=./api.json")
}

func main() {
	flag.Parse()
	swaggerSpec := generateSwaggerJson()

	err := validateSpec(swaggerSpec)
	if err != nil {
		klog.Warningf("Swagger specification has errors")
	}
}

func validateSpec(apiSpec []byte) error {

	swaggerDoc, err := loads.Analyzed(apiSpec, "")
	if err != nil {
		return err
	}

	// Attempts to report about all errors
	validate.SetContinueOnErrors(true)

	v := validate.NewSpecValidator(swaggerDoc.Schema(), strfmt.Default)
	result, _ := v.Validate(swaggerDoc)

	if result.HasWarnings() {
		log.Printf("See warnings below:\n")
		for _, desc := range result.Warnings {
			log.Printf("- WARNING: %s\n", desc.Error())
		}

	}
	if result.HasErrors() {
		str := fmt.Sprintf("The swagger spec is invalid against swagger specification %s.\nSee errors below:\n", swaggerDoc.Version())
		for _, desc := range result.Errors {
			str += fmt.Sprintf("- %s\n", desc.Error())
		}
		log.Println(str)
		return errors.New(str)
	}

	return nil
}

func generateSwaggerJson() []byte {
	container := runtime.Container
	clientsets := k8s.NewNullClient()

	informerFactory := informers.NewNullInformerFactory()
	urlruntime.Must(monitoringv1alpha3.AddToContainer(container, clientsets.Kubernetes(), nil, nil, informerFactory, nil))
	urlruntime.Must(resourcesv1alpha2.AddToContainer(container, clientsets.Kubernetes(), informerFactory, ""))
	urlruntime.Must(resourcesv1alpha3.AddToContainer(container, informerFactory, nil))

	config := restfulspec.Config{
		WebServices:                   container.RegisteredWebServices(),
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}

	swagger := restfulspec.BuildSwagger(config)
	swagger.Info.Extensions = make(spec.Extensions)
	swagger.Info.Extensions.Add("x-tagGroups", []struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}{
		{
			Name: "Authentication",
			Tags: []string{constants.AuthenticationTag},
		},
		{
			Name: "Identity Management",
			Tags: []string{
				constants.UserTag,
			},
		},
		{
			Name: "Resources",
			Tags: []string{
				constants.ClusterResourcesTag,
				constants.NamespaceResourcesTag,
			},
		},
		{
			Name: "Other",
			Tags: []string{
				constants.RegistryTag,
				constants.GitTag,
				constants.ToolboxTag,
				constants.TerminalTag,
			},
		},
		{
			Name: "Monitoring",
			Tags: []string{
				constants.ClusterMetricsTag,
				constants.NodeMetricsTag,
				constants.NamespaceMetricsTag,
				constants.WorkloadMetricsTag,
				constants.PodMetricsTag,
				constants.ContainerMetricsTag,
				constants.WorkspaceMetricsTag,
				constants.ComponentMetricsTag,
				constants.ComponentStatusTag,
			},
		},
		{
			Name: "Logging",
			Tags: []string{constants.LogQueryTag},
		},
	})

	data, _ := json.MarshalIndent(swagger, "", "  ")
	err := ioutil.WriteFile(output, data, 420)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully written to %s", output)

	return data
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "KubeSphere",
			Description: "KubeSphere OpenAPI",
			Version:     version.Get().GitVersion,
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "KubeSphere",
					URL:   "https://kubesphere.io/",
					Email: "kubesphere@yunify.com",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache 2.0",
					URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
				},
			},
		},
	}

	// setup security definitions
	swo.SecurityDefinitions = map[string]*spec.SecurityScheme{
		"jwt": spec.APIKeyAuth("Authorization", "header"),
	}
	swo.Security = []map[string][]string{{"jwt": []string{}}}
}

func apiTree(container *restful.Container) {
	buf := bytes.NewBufferString("\n")
	for _, ws := range container.RegisteredWebServices() {
		for _, route := range ws.Routes() {
			buf.WriteString(fmt.Sprintf("%s %s\n", route.Method, route.Path))
		}
	}
	log.Println(buf.String())
}
