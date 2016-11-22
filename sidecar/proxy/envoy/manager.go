package envoy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/amalgam8/amalgam8/controller/rules"
	"github.com/amalgam8/amalgam8/registry/api"
)

type Manager interface {
	Update(instances []api.ServiceInstance, rules []rules.Rule) error
}

func NewManager(serviceName string) Manager {
	return &manager{
		serviceName: serviceName,
		service:     NewService("/etc/envoy.json"),
	}
}

type manager struct {
	serviceName string
	service     Service
}

func (m *manager) Update(instances []api.ServiceInstance, rules []rules.Rule) error {
	conf, err := generateConfig(rules, instances, m.serviceName)
	if err != nil {
		return err
	}

	file, err := os.Create("/etc/envoy.json")
	if err != nil {
		return err
	}

	if err := writeConfig(file, conf); err != nil {
		file.Close()
		return err
	}
	file.Close()

	if err := m.service.Reload(); err != nil {
		return err
	}

	return nil
}

func writeConfig(writer io.Writer, conf interface{}) error {
	out, err := json.MarshalIndent(&conf, "", "  ")
	if err != nil {
		return err
	}

	_, err = writer.Write(out)
	return err
}

func registryInstancesToClusterManager(instances []api.ServiceInstance) ClusterManager {
	clusterManager := ClusterManager{}

	clusterMap := make(map[string][]Host)
	for _, instance := range instances {
		key := instance.ServiceName
		host := endpointToHost(instance.Endpoint)
		clusterMap[key] = append(clusterMap[key], host)
	}

	for service, hosts := range clusterMap {
		sort.Sort(ByHost(hosts))

		cluster := Cluster{ // TODO: fill in the rest of the values?
			Name:   service,
			Type:   "static",
			LbType: "round_robin",
			Hosts:  hosts,
		}

		clusterManager.Clusters = append(clusterManager.Clusters, cluster)
	}

	sort.Sort(ByName(clusterManager.Clusters))

	return clusterManager
}

func endpointToHost(endpoint api.ServiceEndpoint) Host {
	return Host{
		URL: fmt.Sprintf(endpoint.Type + "://" + endpoint.Value),
	}
}

func generateConfig(rules []rules.Rule, instances []api.ServiceInstance, serviceName string) (Root, error) {
	clusters, err := convert(rules, instances)
	if err != nil {
		return Root{}, err
	}

	clusterNames := make([]string, len(clusters))
	for i, cluster := range clusters {
		clusterNames[i] = cluster.Name
	}

	routes, err := buildRoutes(clusterNames)
	if err != nil {
		return Root{}, err
	}

	filters, err := buildFaults(rules, serviceName)
	if err != nil {
		return Root{}, err
	}

	return Root{
		Listeners: []Listener{
			{
				Port: 6379,
				Filters: []NetworkFilter{
					{
						Type: "read",
						Name: "http_connection_manager",
						Config: NetworkFilterConfig{
							CodecType:  "auto",
							StatPrefix: "ingress_http",
							RouteConfig: RouteConfig{
								VirtualHosts: []VirtualHost{
									{
										Name:    "backend",
										Domains: []string{"*"},
										Routes:  routes,
									},
								},
							},
							Filters: filters,
						},
					},
				},
			},
		},
		Admin: Admin{
			AccessLogPath: "/var/log/envoy_access.log",
			Port:          8001,
		},
		ClusterManager: ClusterManager{
			Clusters: clusters,
		},
	}, nil
}

type uniqueRoute struct {
	Service string
	Tags    []string
}

func convert(rules []rules.Rule, instances []api.ServiceInstance) ([]Cluster, error) {
	// Find unique routes
	uniqueRoutes := make(map[string]uniqueRoute)
	for _, rule := range rules {
		if rule.Route != nil {
			for _, backend := range rule.Route.Backends {
				tags := make([]string, len(backend.Tags))
				copy(tags, backend.Tags)
				sort.Strings(tags)

				buf := bytes.NewBufferString("")

				// Backend name is optional: we default to the rule destination
				var service string
				if backend.Name != "" {
					service = backend.Name
				} else {
					service = rule.Destination
				}
				buf.WriteString(service)

				for _, tag := range tags {
					buf.WriteString("_")
					buf.WriteString(tag)
				}

				key := buf.String()
				uniqueRoutes[key] = uniqueRoute{
					Service: service,
					Tags:    tags,
				}
			}
		}
	}

	// Find endpoints for routes
	hostMap := make(map[string][]Host)
	for _, instance := range instances {
		host := Host{
			URL: "tcp://" + instance.Endpoint.Value,
		} // endpointToHost(instance.Endpoint)
		hostMap[instance.ServiceName] = append(hostMap[instance.ServiceName], host)

		instanceTags := make(map[string]struct{})
		for _, tag := range instance.Tags {
			instanceTags[tag] = struct{}{}
		}

		for key, uniqueRoute := range uniqueRoutes {
			if uniqueRoute.Service == instance.ServiceName {
				isSubset := true
				for _, tag := range uniqueRoute.Tags {
					if _, exists := instanceTags[tag]; !exists {
						isSubset = false
						break
					}
				}

				if isSubset {
					hostMap[key] = append(hostMap[key], host)
				}
			}
		}
	}

	clusters := make([]Cluster, 0, len(hostMap))
	for clusterName, hosts := range hostMap {
		sort.Sort(ByHost(hosts))

		cluster := Cluster{ // TODO: fill in the rest of the values?
			Name:             clusterName,
			Type:             "static",
			LbType:           "round_robin",
			ConnectTimeoutMs: 1000,
			Hosts:            hosts,
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func buildRoutes(clusters []string) ([]Route, error) {
	routes := make([]Route, len(clusters))
	for i, cluster := range clusters {
		routes[i] = Route{
			Prefix:        "/" + cluster + "/",
			PrefixRewrite: "/",
			Cluster:       cluster,
		}
	}

	return routes, nil
}

type Source struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type Match struct {
	Headers map[string]string `json:"headers"`
	Source  Source            `json:"source"`
}

func buildFaults(ctlrRules []rules.Rule, serviceName string) ([]HTTPFilter, error) {

	filters := []HTTPFilter{}

	for _, rule := range ctlrRules {

		headers := make([]HTTPHeader, 0)
		if rule.Match != nil {
			matchBytes, err := rule.Match.MarshalJSON()
			if err != nil {
				logrus.WithError(err).Error("Could not marshall match field")
				return filters, err
			}
			match := Match{}
			err = json.Unmarshal(matchBytes, &match)
			if err != nil {
				logrus.WithError(err).Error("could not unmarshall match field")
				return filters, err
			}

			headers = make([]HTTPHeader, 0, len(match.Headers))
			for key, val := range match.Headers {
				headers = append(headers, HTTPHeader{
					Name:  key,
					Value: val,
				})
			}
		}

		if rule.Destination == serviceName {
			for _, action := range rule.Actions {
				switch action.GetType() {
				case "delay":
					delay := action.Internal().(rules.DelayAction)
					filter := HTTPFilter{
						Type: "decoder",
						Name: "fault",
						Config: &HTTPFilterFaultConfig{
							Delay: &HTTPDelayFilter{
								Type:     "fixed",
								Percent:  delay.Probability,
								Duration: delay.Duration,
							},
							Headers: headers,
						},
					}
					filters = append(filters, filter)
				case "abort":
					abort := action.Internal().(rules.AbortAction)
					filter := HTTPFilter{
						Type: "decoder",
						Name: "fault",
						Config: &HTTPFilterFaultConfig{
							Abort: &HTTPAbortFilter{
								Percent:    abort.Probability,
								HTTPStatus: abort.ReturnCode,
							},
							Headers: headers,
						},
					}
					filters = append(filters, filter)
				}
			}
		}
	}

	filters = append(filters, HTTPFilter{
		Type:   "decoder",
		Name:   "router",
		Config: HTTPFilterRouterConfig{},
	})
	return filters, nil
}
