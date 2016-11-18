package envoy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/amalgam8/amalgam8/controller/rules"
	"github.com/amalgam8/amalgam8/registry/api"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	instances := []api.ServiceInstance{
		{
			ServiceName: "freeflow1",
			Endpoint: api.ServiceEndpoint{
				Type:  "tcp",
				Value: "10.11.0.4:80",
			},
		},
	}

	clusterManager := registryInstancesToClusterManager(instances)

	assert.Len(t, clusterManager.Clusters, 1)
	assert.Equal(t, "freeflow1", clusterManager.Clusters[0].Name)
	assert.Len(t, clusterManager.Clusters[0].Hosts, 1)
	assert.Equal(t, "tcp://10.11.0.4:80", clusterManager.Clusters[0].Hosts[0].URL)
}

func TestConvertFancy(t *testing.T) {
	instances := []api.ServiceInstance{
		{
			ServiceName: "service1",
			Endpoint: api.ServiceEndpoint{
				Type:  "tcp",
				Value: "10.0.0.1:80",
			},
		},
		{
			ServiceName: "service1",
			Endpoint: api.ServiceEndpoint{
				Type:  "tcp",
				Value: "10.0.0.2:80",
			},
		},
		{
			ServiceName: "service2",
			Endpoint: api.ServiceEndpoint{
				Type:  "https",
				Value: "10.0.0.3:80",
			},
		},
	}

	clusterManager := registryInstancesToClusterManager(instances)

	assert.Len(t, clusterManager.Clusters, 2)
	assert.Equal(t, "service1", clusterManager.Clusters[0].Name)
	assert.Len(t, clusterManager.Clusters[0].Hosts, 2)
	assert.Equal(t, "tcp://10.0.0.1:80", clusterManager.Clusters[0].Hosts[0].URL)
	assert.Equal(t, "tcp://10.0.0.2:80", clusterManager.Clusters[0].Hosts[1].URL)
	assert.Equal(t, "service2", clusterManager.Clusters[1].Name)
	assert.Len(t, clusterManager.Clusters[1].Hosts, 1)
	assert.Equal(t, "https://10.0.0.3:80", clusterManager.Clusters[1].Hosts[0].URL)

	writer := bytes.NewBufferString("")
	err := writeConfig(writer, clusterManager)
	assert.NoError(t, err)

	//	conf := `{"cluster_manager": {
	//  "clusters": [
	//  {
	//"name": "freeflow1",
	//"connect_timeout_ms": 10000,
	//"type": "static",
	//"lb_type": "round_robin",
	//"max_requests_per_connection" : 10000,
	//"hosts": [
	//{
	//"url": "tcp://10.11.0.4:80"
	//}
	//]
	//}
	//]
	//}`
	//
	//	assert.Equal(t, conf, writer.String())
}

func TestEndpointConvert(t *testing.T) {
	endpoint := api.ServiceEndpoint{
		Type:  "http",
		Value: "10.1.2.45",
	}
	host := endpointToHost(endpoint)

	assert.Equal(t, host.URL, "http://10.1.2.45")
}

func TestConvert2(t *testing.T) {
	instances := []api.ServiceInstance{
		{
			ServiceName: "service1",
			Endpoint: api.ServiceEndpoint{
				Type:  "tcp",
				Value: "10.0.0.1:80",
			},
			Tags: []string{},
		},
		{
			ServiceName: "service1",
			Endpoint: api.ServiceEndpoint{
				Type:  "tcp",
				Value: "10.0.0.2:80",
			},
			Tags: []string{"tag1"},
		},
		{
			ServiceName: "service1",
			Endpoint: api.ServiceEndpoint{
				Type:  "tcp",
				Value: "10.0.0.3:80",
			},
			Tags: []string{"tag2"},
		},
		{
			ServiceName: "service1",
			Endpoint: api.ServiceEndpoint{
				Type:  "tcp",
				Value: "10.0.0.4:80",
			},
			Tags: []string{"tag1", "tag2"},
		},
		{
			ServiceName: "service2",
			Endpoint: api.ServiceEndpoint{
				Type:  "https",
				Value: "10.0.0.5:80",
			},
		},
	}

	rules := []rules.Rule{
		{
			ID:          "abcdef",
			Destination: "service1",
			Route: &rules.Route{
				Backends: []rules.Backend{
					{
						Name: "service1",
						Tags: []string{"tag1"},
					},
				},
			},
		},
		{
			ID:          "abcdef",
			Destination: "service1",
			Route: &rules.Route{
				Backends: []rules.Backend{
					{
						Name: "service1",
						Tags: []string{"tag1", "tag2"},
					},
				},
			},
		},
		{
			ID:          "abcdef",
			Destination: "service2",
			Actions:     json.RawMessage([]byte{}),
		},
	}

	configRoot, err := generateConfig(rules, instances)
	assert.NoError(t, err)

	data, err := json.MarshalIndent(configRoot, "", "  ")
	assert.NoError(t, err)

	fmt.Println(string(data))
}

func TestBookInfo(t *testing.T) {
	ruleBytes := []byte(`[
    {
      "id": "ad95f5d6-fa7b-448d-8c27-928e40b37202",
      "priority": 2,
      "destination": "reviews",
      "match": {
        "headers": {
          "Cookie": ".*?user=jason"
        }
      },
      "route": {
        "backends": [
          {
            "tags": [
              "v2"
            ]
          }
        ]
      }
    },
    {
      "id": "e31da124-8394-4b12-9abf-ebdc7db679a9",
      "priority": 1,
      "destination": "details",
      "route": {
        "backends": [
          {
            "tags": [
              "v1"
            ]
          }
        ]
      }
    },
    {
      "id": "ab823eb5-e56c-485c-901f-0f29adfa8e4f",
      "priority": 1,
      "destination": "productpage",
      "route": {
        "backends": [
          {
            "tags": [
              "v1"
            ]
          }
        ]
      }
    },
    {
      "id": "03b97f82-40c5-4c51-8bf9-b1057a73019b",
      "priority": 1,
      "destination": "ratings",
      "route": {
        "backends": [
          {
            "tags": [
              "v1"
            ]
          }
        ]
      }
    },
    {
      "id": "c67226e2-8506-4e75-9e47-84d9d24f0326",
      "priority": 1,
      "destination": "reviews",
      "route": {
        "backends": [
          {
            "tags": [
              "v1"
            ]
          }
        ]
      }
    }
  ]
`)

	instanceBytes := []byte(`[
    {
      "id": "74d2a394184327f5",
      "service_name": "productpage",
      "endpoint": {
        "type": "http",
        "value": "172.17.0.6:9080"
      },
      "ttl": 60,
      "status": "UP",
      "last_heartbeat": "2016-11-18T17:02:32.822819186Z",
      "tags": [
        "v1"
      ]
    },
    {
      "id": "26b250bc98d8a74c",
      "service_name": "ratings",
      "endpoint": {
        "type": "http",
        "value": "172.17.0.11:9080"
      },
      "ttl": 60,
      "status": "UP",
      "last_heartbeat": "2016-11-18T17:02:33.784740831Z",
      "tags": [
        "v1"
      ]
    },
    {
      "id": "9f7a75cdbbf492c7",
      "service_name": "details",
      "endpoint": {
        "type": "http",
        "value": "172.17.0.7:9080"
      },
      "ttl": 60,
      "status": "UP",
      "last_heartbeat": "2016-11-18T17:02:32.986290003Z",
      "tags": [
        "v1"
      ]
    },
    {
      "id": "05f853b7b4ab8b37",
      "service_name": "reviews",
      "endpoint": {
        "type": "http",
        "value": "172.17.0.10:9080"
      },
      "ttl": 60,
      "status": "UP",
      "last_heartbeat": "2016-11-18T17:02:33.559542468Z",
      "tags": [
        "v3"
      ]
    },
    {
      "id": "a4a740e9af065016",
      "service_name": "reviews",
      "endpoint": {
        "type": "http",
        "value": "172.17.0.8:9080"
      },
      "ttl": 60,
      "status": "UP",
      "last_heartbeat": "2016-11-18T17:02:33.18906562Z",
      "tags": [
        "v1"
      ]
    },
    {
      "id": "5f940f0ddee732bb",
      "service_name": "reviews",
      "endpoint": {
        "type": "http",
        "value": "172.17.0.9:9080"
      },
      "ttl": 60,
      "status": "UP",
      "last_heartbeat": "2016-11-18T17:02:33.349101984Z",
      "tags": [
        "v2"
      ]
    }
  ]`)
	var ruleList []rules.Rule
	err := json.Unmarshal(ruleBytes, &ruleList)
	assert.NoError(t, err)

	var instances []api.ServiceInstance
	err = json.Unmarshal(instanceBytes, &instances)
	assert.NoError(t, err)

	configRoot, err := generateConfig(ruleList, instances)
	assert.NoError(t, err)

	data, err := json.MarshalIndent(configRoot, "", "  ")
	assert.NoError(t, err)

	fmt.Println(string(data))
}
