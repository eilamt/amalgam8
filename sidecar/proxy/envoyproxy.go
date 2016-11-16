// Copyright 2016 IBM Corporation
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package proxy

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/amalgam8/amalgam8/controller/rules"
	"github.com/amalgam8/amalgam8/registry/client"
	"github.com/amalgam8/amalgam8/sidecar/proxy/envoy"
	"github.com/amalgam8/amalgam8/sidecar/proxy/monitor"
	"github.com/amalgam8/amalgam8/sidecar/proxy/nginx"
)

// NGINXProxy updates NGINX to reflect changes in the A8 controller and A8 registry
type EnvoyProxy interface {
	monitor.ControllerListener
	monitor.RegistryListener
	GetState() ([]client.ServiceInstance, []rules.Rule)
}

type envoyProxy struct {
	instances []client.ServiceInstance
	rules     []rules.Rule
	envoy     envoy.Manager
	mutex     sync.Mutex
}

// NewEnvoyProxy instantiates a new instance
func NewEnvoyProxy(envoyManager nginx.Manager) EnvoyProxy {
	return &envoyProxy{
		rules:     []rules.Rule{},
		instances: []client.ServiceInstance{},
		envoy:     envoyManager,
	}
}

// CatalogChange updates NGINX on a change in the catalog
func (n *envoyProxy) CatalogChange(instances []client.ServiceInstance) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.instances = instances
	return n.updateNGINX()
}

// RuleChange updates NGINX on a change in the proxy configuration
func (n *envoyProxy) RuleChange(rules []rules.Rule) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.rules = rules
	return n.updateNGINX()
}

func (n *envoyProxy) updateNGINX() error {
	logrus.Debug("Updating NGINX")
	return n.envoy.Update(n.instances, n.rules)
}

func (n *envoyProxy) GetState() ([]client.ServiceInstance, []rules.Rule) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	return n.instances, n.rules
}
