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
		//{
		//	ServiceName: "service1",
		//	Endpoint: api.ServiceEndpoint{
		//		Type: "tcp",
		//		Value: "10.0.0.2:80",
		//	},
		//	Tags: []string{"tag1"},
		//},
		//{
		//	ServiceName: "service1",
		//	Endpoint: api.ServiceEndpoint{
		//		Type: "tcp",
		//		Value: "10.0.0.3:80",
		//	},
		//	Tags: []string{"tag2"},
		//},
		//{
		//	ServiceName: "service1",
		//	Endpoint: api.ServiceEndpoint{
		//		Type: "tcp",
		//		Value: "10.0.0.4:80",
		//	},
		//	Tags: []string{"tag1", "tag2"},
		//},
		{
			ServiceName: "service2",
			Endpoint: api.ServiceEndpoint{
				Type:  "https",
				Value: "10.0.0.5:80",
			},
		},
	}

	rules := []rules.Rule{
	//{
	//	ID: "abcdef",
	//	Destination: "x",
	//	Route: &RouteRule{
	//		Backends: []Backend{
	//			{
	//				Name: "service1",
	//				Tags: []string{"tag1"},
	//			},
	//		},
	//	},
	//},
	//{
	//	ID: "abcdef",
	//	Destination: "x",
	//	Route: &RouteRule{
	//		Backends: []Backend{
	//			{
	//				Name: "service1",
	//				Tags: []string{"tag1", "tag2"},
	//			},
	//		},
	//	},
	//},
	//{
	//	ID: "abcdef",
	//	Destination: "x",
	//	Actions: json.RawMessage([]byte{}),
	//},
	}

	configRoot, err := generateConfig(rules, instances)
	if err != nil {
		panic("OH JESUS!")
	}

	data, err := json.MarshalIndent(configRoot, "", "  ")
	if err != nil {
		panic("OH SWEET JESUS")
	}

	fmt.Println(string(data))
}
