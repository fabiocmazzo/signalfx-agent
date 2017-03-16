package config

import (
	"fmt"

	"github.com/signalfx/neo-agent/services"
)

// Plugin describes a collectd plugin
type Plugin struct {
	Plugin   string
	Name     string
	Template string
	Dims     string
	Host     string
	Port     uint16
	Config   map[string]interface{}
}

// PLUGINS is a mapping to create plugin instances with defaults
var PLUGINS = map[services.ServiceType]func(string) *Plugin{
	services.ApacheService: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:   "apache",
			Template: "apache.default.conf.tmpl",
			Name:      instanceName,
			Host:      "localhost",
			Port:      80}
	},
	services.DockerService: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:   "docker",
			Template: "docker.conf.tmpl",
			Name:      instanceName,
			Config: map[string]interface{}{
				"url": "unix:///var/run/docker.sock",
			},
		}
	},
	services.MongoDBService: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:   "mongodb",
			Template: "mongodb.default.conf.tmpl",
			Name:      instanceName,
			Host:      "localhost",
			Port:      27017}
	},
	services.RedisService: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:   "redis",
			Template: "redis.default.conf.tmpl",
			Name:      instanceName,
			Host:      "localhost",
			Port:      6379,
		}
	},
	services.SignalfxService: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:   "signalfx",
			Template: "signalfx.conf.tmpl",
			Config: map[string]interface{}{
				"url": "https://ingest.signalfx.com",
			},
			Name: instanceName,
		}
	},
	services.ZookeeperService: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:   "zookeeper",
			Template: "zookeeper.default.conf.tmpl",
			Name:      instanceName,
			Host:      "localhost",
			Port:      2181}
	},
	services.WriteHTTPService: func(instanceName string) *Plugin {
		return &Plugin{
			Plugin:   "write-http",
			Template: "write-http.conf.tmpl",
			Config: map[string]interface{}{
				"url": "https://ingest.signalfx.com",
			},
			Name: instanceName,
		}
	},
}

// NewPlugin constructs a plugin with default values depending on the service type
func NewPlugin(pluginType services.ServiceType, pluginName string) (*Plugin, error) {
	if create, ok := PLUGINS[pluginType]; ok {
		return create(pluginName), nil
	}
	return nil, fmt.Errorf("plugin %s is unsupported", pluginType)
}

// GroupByPlugin creates a map of instances by plugin
func GroupByPlugin(instances []*Plugin) map[string][]*Plugin {
	pluginMap := map[string][]*Plugin{}
	for _, instance := range instances {
		if val, ok := pluginMap[instance.Plugin]; ok {
			pluginMap[instance.Plugin] = append(val, instance)
		} else {
			pluginMap[instance.Plugin] = []*Plugin{instance}
		}
	}
	return pluginMap
}

// CollectdConfig are global collectd settings
type CollectdConfig struct {
	Interval             uint
	Timeout              uint
	ReadThreads          uint
	WriteQueueLimitHigh  uint `yaml:"writeQueueLimitHigh"`
	WriteQueueLimitLow   uint `yaml:"writeQueueLimitLow"`
	CollectInternalStats bool
	Hostname             string
	Plugins              []map[string]interface{}
}

// AppConfig is the top-level configuration object consumed by templates.
type AppConfig struct {
	AgentConfig *CollectdConfig
	Plugins     map[string][]*Plugin
}

// NewCollectdConfig creates a default collectd config instance
func NewCollectdConfig() *CollectdConfig {
	return &CollectdConfig{
		Interval:             15,
		Timeout:              2,
		ReadThreads:          5,
		WriteQueueLimitHigh:  500000,
		WriteQueueLimitLow:   400000,
		CollectInternalStats: true,
	}
}