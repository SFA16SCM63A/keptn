package helm

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/helm/pkg/proto/hapi/chart"
	"sigs.k8s.io/yaml"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/apis/networking/istio/v1alpha3"
	"github.com/keptn/keptn/helm-service/pkg/jsonutils"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"

	kyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/helm/pkg/chartutil"
)

const projectName = "sockshop"
const serviceName = "carts"

type GeneratedResource struct {
	URI         string
	FileContent []string
}

var cartsCanaryIstioDestinationRuleGen = GeneratedResource{
	URI: "templates/carts-canary-istio-destinationrule.yaml",
	FileContent: []string{`
{
  "apiVersion": "networking.istio.io/v1alpha3",
  "kind": "DestinationRule",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-canary"
  },
  "spec": {
    "host": "carts-canary.sockshop-production.svc.cluster.local"
  }
}
`},
}

var cartsIstioVirtualserviceGen = GeneratedResource{
	URI: "templates/carts-istio-virtualservice.yaml",
	FileContent: []string{`
{
  "apiVersion": "networking.istio.io/v1alpha3",
  "kind": "VirtualService",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts"
  },
  "spec": {
    "gateways": [
      "sockshop-production-gateway.sockshop-production",
      "mesh"
    ],
    "hosts": [
      "carts.sockshop-production.mydomain.sh",
      "carts",
      "carts.sockshop-production"
    ],
    "http": [
      {
        "route": [
          {
            "destination": {
              "host": "carts-canary.sockshop-production.svc.cluster.local"
            }
          },
          {
            "destination": {
              "host": "carts-primary.sockshop-production.svc.cluster.local"
            },
            "weight": 100
          }
        ]
      }
    ]
  }
}
`},
}

var cartsPrimaryIstioDestinationRuleGen = GeneratedResource{
	URI: "templates/carts-primary-istio-destinationrule.yaml",
	FileContent: []string{`
{
  "apiVersion": "networking.istio.io/v1alpha3",
  "kind": "DestinationRule",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-primary"
  },
  "spec": {
    "host": "carts-primary.sockshop-production.svc.cluster.local"
  }
}
`},
}

var deploymentGen = GeneratedResource{
	URI: "templates/deployment.yml",
	FileContent: []string{`
{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-primary"
  },
  "spec": {
    "replicas" : 1,
    "selector": {
      "matchLabels": {
        "app": "carts-primary"
      }
    },
    "strategy": {
      "rollingUpdate": {
        "maxUnavailable": 0
      },
      "type": "RollingUpdate"
    },
    "template": {
      "metadata": {
        "creationTimestamp": null,
        "labels": {
          "app": "carts-primary"
        }
      },
      "spec": {
        "containers": [
          {            
            "image": "docker.io/keptnexamples/carts:0.8.1",
            "imagePullPolicy": "IfNotPresent",
            "livenessProbe": {
              "httpGet": {
                "path": "/health",
                "port": 8080
              },
              "initialDelaySeconds": 60,
              "periodSeconds": 10,
              "timeoutSeconds": 15
            },
            "name": "carts",
            "ports": [
              {
                "containerPort": 8080,
                "name": "http",
                "protocol": "TCP"
              }
            ],
            "readinessProbe": {
              "httpGet": {
                "path": "/health",
                "port": 8080
              },
              "initialDelaySeconds": 60,
              "periodSeconds": 10,
              "timeoutSeconds": 15
            },
            "resources": {
              "limits": {
                "cpu": "500m",
                "memory": "2Gi"
              },
              "requests": {
                "cpu": "250m",
                "memory": "1Gi"
              }
            }
          }
        ]
      }
    }
  },
  "status": {
  }
}`},
}

var serviceGen = GeneratedResource{
	URI: "templates/service.yaml",
	FileContent: []string{`
{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-canary"
  },
  "spec": {
    "ports": [
      {
        "name": "http",
        "port": 80,
        "protocol": "TCP",
        "targetPort": 8080
      }
    ],
    "selector": {
      "app": "carts"
    },
    "type": "LoadBalancer"
  },
  "status": {
    "loadBalancer": {
    }
  }
}`, `
{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "creationTimestamp": null,
    "name": "carts-primary"
  },
  "spec": {
    "ports": [
      {
        "name": "http",
        "port": 80,
        "protocol": "TCP",
        "targetPort": 8080
      }
    ],
    "selector": {
      "app": "carts-primary"
    },
    "type": "LoadBalancer"
  },
  "status": {
    "loadBalancer": {
    }
  }
}`}}

var valuesGen = GeneratedResource{
	URI: "values.yaml",
	FileContent: []string{`
{
  "image": "docker.io/keptnexamples/carts:0.8.1",
  "replicas": 1
}`},
}

func TestGenerateManagedChart(t *testing.T) {

	data := CreateHelmChartData(t)

	h := NewGeneratedChartHandler(mesh.NewIstioMesh(), NewCanaryOnDeploymentGenerator(), "mydomain.sh")
	inputChart, err := keptnutils.LoadChart(data)
	if err != nil {
		t.Error(err)
	}
	gen, err := h.GenerateManagedChart(inputChart, projectName, "production")
	assert.Nil(t, err, "Generating the managed Chart should not return any error")

	workingPath, err := ioutil.TempDir("", "helm-test")
	defer os.RemoveAll(workingPath)
	packagedChartFilePath := filepath.Join(workingPath, serviceName)
	err = ioutil.WriteFile(packagedChartFilePath, gen, 0644)
	if err != nil {
		t.Error(err)
	}

	ch, err := chartutil.Load(packagedChartFilePath)

	// Compare values
	yReader := kyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader([]byte(ch.Values.Raw))))
	yamlData, err := yReader.Read()
	if err != nil {
		t.Error(err)
	}
	jsonData, err := jsonutils.ToJSON(yamlData)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(jsonData), valuesGen.FileContent[0])

	// Compare templates
	generatedTemplateResources := []GeneratedResource{cartsCanaryIstioDestinationRuleGen, cartsIstioVirtualserviceGen, cartsPrimaryIstioDestinationRuleGen,
		deploymentGen, serviceGen}

	for _, resource := range generatedTemplateResources {

		reader := ioutil.NopCloser(bytes.NewReader(getTemplateByName(ch, resource.URI).Data))
		decoder := kyaml.NewDocumentDecoder(reader)

		for i := 0; ; i++ {
			b1 := make([]byte, 4096)
			n1, err := decoder.Read(b1)
			if err == io.EOF {
				break
			}
			assert.Nil(t, err, "")

			jsonData, err := jsonutils.ToJSON(b1[:n1])
			if err != nil {
				t.Error(err)
			}

			ja := jsonassert.New(t)
			ja.Assertf(string(jsonData), resource.FileContent[i])
		}
	}
}

func TestUpdateCanaryWeight(t *testing.T) {
	data := CreateHelmChartData(t)

	h := NewGeneratedChartHandler(mesh.NewIstioMesh(), NewCanaryOnDeploymentGenerator(), "mydomain.sh")
	chart, err := keptnutils.LoadChart(data)
	if err != nil {
		t.Error(err)
	}
	genChartData, err := h.GenerateManagedChart(chart, "sockshop", "production")
	if err != nil {
		t.Error(err)
	}
	genChart, err := keptnutils.LoadChart(genChartData)
	if err != nil {
		t.Error(err)
	}

	h.UpdateCanaryWeight(genChart, 48)
	assert.Nil(t, err, "Generating the managed Chart should not return any error")

	template := getTemplateByName(genChart, "templates/carts-istio-virtualservice.yaml")
	assert.NotNil(t, template, "Template must not be null")
	vs := v1alpha3.VirtualService{}
	err = yaml.Unmarshal(template.Data, &vs)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(vs.Spec.Http))
	assert.Equal(t, 2, len(vs.Spec.Http[0].Route))
	assert.Equal(t, int32(48), vs.Spec.Http[0].Route[0].Weight)
	assert.True(t, strings.Contains(vs.Spec.Http[0].Route[0].Destination.Host, "canary"))
	assert.Equal(t, int32(52), vs.Spec.Http[0].Route[1].Weight)
	assert.True(t, strings.Contains(vs.Spec.Http[0].Route[1].Destination.Host, "primary"))
}

func getTemplateByName(chart *chart.Chart, templateName string) *chart.Template {

	for _, template := range chart.Templates {
		if template.Name == templateName {
			return template
		}
	}
	return nil
}
