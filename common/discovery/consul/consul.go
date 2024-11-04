package consul

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr string, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}


func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid hostPort")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	host := parts[0]
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceID,
		Address: host,
		Port:    port,
		Name:    serviceName,
		Check:   &consul.AgentServiceCheck{
			CheckID:  instanceID,
			TLSSkipVerify: true,
			TTL: "5s",
			Timeout: "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}


func (r *Registry) DeRegister(ctx context.Context, instanceID, servicename string) error {
	log.Printf("de-registering %s from %s", instanceID, servicename)

	return r.client.Agent().ServiceDeregister(instanceID)
}


func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var instances []string
	for _, entry := range entries {
		instances = append(instances, 
							fmt.Sprintf("%s:%d", 
										entry.Service.Address, 
										entry.Service.Port))
	}
	return instances, nil
}


func (r *Registry) HealthCheck(instanceID, servicename string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", consul.HealthPassing)
}