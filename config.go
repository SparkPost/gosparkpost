package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

//go:generate stringer -type=InstanceType
//go:generate jsonenums -type=InstanceType
type InstanceType int

const (
	Elite InstanceType = iota
	Momentum
	SparkPost
)

const (
	EliteBase     string = "https://%s.msyscloud.com/api/%s/%s"
	SparkPostBase string = "https://api.sparkpost.com/api/%s/%s"
)

//go:generate jsonenums -type=InjectionProtocol
//go:generate stringer -type=InjectionProtocol
type InjectionProtocol int

const (
	REST InjectionProtocol = iota
	SMTP
)

type Binding struct {
	Name    string
	Domains []string
	IPs     []string
}

// Complete config details for one Momentum / SparkPost (Elite?) instance
type Config struct {
	File         string
	Name         string
	Type         InstanceType
	ApiKey       string
	ApiBase      string
	LinkDomain   string
	BaseDomain   string
	Protocols    []InjectionProtocol
	Bindings     map[string]Binding
	BindingNames []string
	TestBinding  string
	TestDomain   string
	rand         *rand.Rand
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Summarize(cfg Config) string {
	if cfg.Type == SparkPost {
		return fmt.Sprintf("Testing %s instance %s (%s)", cfg.Type, cfg.Name, cfg.Protocols)
	} else {
		return fmt.Sprintf("Testing %s instance %s (%s) > binding (%s) > domain (%s)",
			cfg.Type, cfg.Name, cfg.TestBinding, cfg.TestDomain)
	}
}

// Read config file, return populated config object for the specified name
func Load(filename string) (*Config, error) {
	// instance to test comes from environment variable
	instanceName := os.Getenv("MSYS_SMOKE_CONFIG")
	if instanceName == "" {
		return nil, fmt.Errorf("No instance to test - set MSYS_SMOKE_CONFIG")
	}

	// open and parse json config file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := map[string]Config{}
	err = decoder.Decode(&config)
	if err != nil {
		// FIXME: this error message can be super non-intuitive
		return nil, err
	}

	// get config for the specified instance
	test, ok := config[instanceName]
	if !ok {
		return nil, fmt.Errorf("No config for MSYS_SMOKE_CONFIG = [%s]", instanceName)
	}

	test.File = filename
	test.Name = instanceName
	// auto-set name member of each Binding
	for name, binding := range test.Bindings {
		binding.Name = name
		test.Bindings[name] = binding
	}

	// if no protocol is specified for this instance, default to REST
	if len(test.Protocols) == 0 {
		test.Protocols = make([]InjectionProtocol, 1)
		test.Protocols[0] = REST
	}

	// TODO: set default values based on instance type

	// binding to test comes from an environment variable
	bindingName := os.Getenv("MSYS_SMOKE_BINDING")
	if test.Type == SparkPost {
		// ignore test binding - irrelevant for SparkPost
	} else if bindingName != "" {
		binding, ok := test.Bindings[bindingName]
		if !ok {
			return nil, fmt.Errorf("No such binding [%s] for [%s]", bindingName, instanceName)
		}
		// remove all but specified binding from config
		test.Bindings = map[string]Binding{bindingName: binding}
		test.TestBinding = bindingName
	}

	test.BindingNames = make([]string, len(test.Bindings))
	idx := 0
	// auto-populate binding names array
	for key := range test.Bindings {
		test.BindingNames[idx] = key
		idx++
	}

	// initialize pseudo-random number generator
	test.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	/*
		NB: A "test domain" may be specified without a "test binding", however that
		    requires the domain to be present in each binding in the instance.
				If that is not the case, a runtime error will result.
	*/

	// domain to test comes from an environment variable
	domain := os.Getenv("MSYS_SMOKE_DOMAIN")
	if test.Type == SparkPost {
		// ignore test domain - irrelevant for SparkPost
	} else if domain != "" {
		// iterate over remaining bindings, removing all but specified domain
		for bindingName, binding := range test.Bindings {
			if !stringInSlice(domain, binding.Domains) {
				return nil, fmt.Errorf("No such domain [%s] for [%s]/[%s]", domain, instanceName, bindingName)
			} else {
				binding.Domains = []string{domain}
			}
		}
		test.TestDomain = domain
	}

	// protocol to test comes from environment variable
	protocol := os.Getenv("MSYS_SMOKE_PROTOCOL")
	if protocol != "" {
		var p InjectionProtocol
		var ok bool
		if p, ok = _InjectionProtocolNameToValue[protocol]; !ok {
			return nil, fmt.Errorf("Unrecognized injection protocol [%s]", protocol)
		}
		test.Protocols = make([]InjectionProtocol, 1)
		test.Protocols[0] = p
	}

	return &test, nil
}

// Name of the current config object
// (unnecessary: this is part of the config struct)

// A randomly-chosen binding
func (c Config) RandomBinding() Binding {
	keyidx := c.rand.Intn(len(c.Bindings))
	bname := c.BindingNames[keyidx]
	return c.Bindings[bname]
}

// List of all binding names
// (unnecessary: this is part of the config struct)

// Base URL to use for all API requests

// List of domains that must accept FBL messages

// List of domains for the specified binding
func (c Config) BindingDomains(binding string) []string {
	return c.Bindings[binding].Domains
}

// List of configured protocols
// (unnecessary: this is part of the config struct)

// List of ips for the specified binding
func (c Config) BindingIPs(binding string) (ips []net.IP, err error) {
	nIPs := len(c.Bindings[binding].IPs)
	if nIPs == 0 {
		return ips, fmt.Errorf("No IPs for binding [%s]", binding)
	}
	ips = make([]net.IP, nIPs)
	i := 0
	for _, ipStr := range c.Bindings[binding].IPs {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return ips, fmt.Errorf("Invalid IP [%s] for binding [%s]", ipStr, binding)
		}
		ips[i] = ip
		i++
	}
	return ips, nil
}
