package main

// The configuration (environment) data structure
type Config_s struct {
	CAcert string `json:"cacert"`
	//	Cert        string `json:"cert"`
	//	Key         string `json:"key"`
	ListenerURL string `json:"listenerurl"`
}

// The Client information
type PrometheusTarget_s struct {
	Targets []string          `json:"targets" yaml:"targets"`
	Labels  map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// CommandPayload is the full payload including the command and PrometheusTarget
type CommandPayload_s struct {
	Command          string             `json:"command"`
	PrometheusTarget PrometheusTarget_s `json:"prometheus_target"`
}
