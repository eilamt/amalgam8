package envoy

type HTTPAbortFilter struct {
	Percent    int `json:"abort_percent,omitempty"`
	HTTPStatus int `json:"http_status,omitempty"`
}

type HTTPDelayFilter struct {
	Type     string `json:"type,omitempty"`
	Percent  int    `json:"fixed_delay_percent,omitempty"`
	Duration int    `json:"fixed_duration_ms,omitempty"`
}

type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HTTPFilterFaultConfig struct {
	Abort   *HTTPAbortFilter `json:"abort,omitempty"`
	Delay   *HTTPDelayFilter `json:"delay,omitempty"`
	Headers []HTTPHeader     `json:"headers,omitempty"`
}

type HTTPFilterRouterConfig struct {
	DynamicStats bool `json:"dynamic_stats"`
}

type HTTPFilter struct {
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Config interface{} `json:"config"`
}

type Route struct {
	Prefix        string `json:"prefix"`
	PrefixRewrite string `json:"prefix_rewrite"`
	Cluster       string `json:"cluster"`
}

type VirtualHost struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	Routes  []Route  `json:"routes"`
}

type RouteConfig struct {
	VirtualHosts []VirtualHost `json:"virtual_hosts"`
}

type NetworkFilterConfig struct {
	CodecType         string       `json:"codec_type"`
	StatPrefix        string       `json:"stat_prefix"`
	GenerateRequestID bool         `json:"generate_request_id"`
	RouteConfig       RouteConfig  `json:"route_config"`
	Filters           []HTTPFilter `json:"filters"`
}

type NetworkFilter struct {
	Type   string              `json:"type"`
	Name   string              `json:"name"`
	Config NetworkFilterConfig `json:"config"`
}

type Listener struct {
	Port    int             `json:"port"`
	Filters []NetworkFilter `json:"filters"`
}

type Admin struct {
	AccessLogPath string `json:"access_log_path"`
	Port          int    `json:"port"`
}

type Host struct {
	URL string `json:"url"`
}

type Cluster struct {
	Name                     string `json:"name"`
	ConnectTimeoutMs         int    `json:"connect_timeout_ms"`
	Type                     string `json:"type"`
	LbType                   string `json:"lb_type"`
	MaxRequestsPerConnection int    `json:"max_requests_per_connection,omitempty"`
	Hosts                    []Host `json:"hosts"`
}

type ClusterManager struct {
	Clusters []Cluster `json:"clusters"`
}

type Root struct {
	Listeners      []Listener `json:"listeners"`
	Admin          Admin      `json:"admin"`
	ClusterManager `json:"cluster_manager"`
}

type ByName []Cluster

func (a ByName) Len() int {
	return len(a)
}

func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

type ByHost []Host

func (a ByHost) Len() int {
	return len(a)
}

func (a ByHost) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByHost) Less(i, j int) bool {
	return a[i].URL < a[j].URL
}
