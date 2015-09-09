package config

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
	EliteBase     string = "https://%s.msyscloud.com"
	SparkPostBase string = "https://api.sparkpost.com"
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
	ApiVersion   int
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

func (c Config) Summarize() string {
	if c.Type == SparkPost {
		return fmt.Sprintf("Testing %s instance %s %s", c.Type, c.Name, c.Protocols)
	} else {
		testBinding := c.TestBinding
		if testBinding == "" {
			testBinding = "all"
		}
		testDomain := c.TestDomain
		if testDomain == "" {
			testDomain = "all"
		}
		return fmt.Sprintf("Testing %s instance %s %s > binding (%s) > domain (%s)",
			c.Type, c.Name, c.Protocols, testBinding, testDomain)
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

	// set the API version if not specified
	// (v1 is the only version as of 2015-09-09)
	if test.ApiVersion == 0 {
		test.ApiVersion = 1
	}

	// set up the base url to be used for api calls
	// NB: ApiBase must not include "/api/v1", that will be added
	if test.ApiBase == "" {
		if test.Type == Elite {
			test.ApiBase = fmt.Sprintf(EliteBase, test.Name)
		} else if test.Type == SparkPost {
			test.ApiBase = SparkPostBase
		} else if test.Type == Momentum {
			// ApiBase must be manually configured for Momentum AKA onprem
			return nil, fmt.Errorf("Missing ApiBase for Momentum instance [%s]", test.Name)
		}
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
		var found bool
		for _, p = range test.Protocols {
			if protocol == p.String() {
				found = true
				break
			}
		}
		// may only choose from configured protocols
		if !found {
			return nil, fmt.Errorf("Injection protocol [%s] unavailable for instance [%s]", protocol, test.Name)
		}
		test.Protocols = make([]InjectionProtocol, 1)
		test.Protocols[0] = p
	}

	return &test, nil
}

// API URL builder
func (c Config) ApiUrl(suffix string) string {
	return fmt.Sprintf("%s/api/v%d%s", c.ApiBase, c.ApiVersion, suffix)
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
