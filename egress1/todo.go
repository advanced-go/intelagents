package egress1

import "github.com/advanced-go/stdlib/messaging"

// What should the company do??
// 1. - Provide a listing of agents and their class
// 2. - Manage external requests/response exchanges

// So a CIA Case Officer will interact with a Director of Operations Staff Operations Officer.
//

// How does startup work?
//

// All agents utilize: observations, guidance, and experience as input to inference -> to create actions
// observations + experience + guidance -> inference -> actions

// NewIngressAgent - ingress traffic controller for a host, one agent per host.
func NewIngressAgent(uri string, assignment any, ctrlHandler messaging.Handler) (messaging.Agent, error) {
	return nil, nil
}

// NewEgressAgent - egress traffic controller for a host, one agent per route
func NewEgressAgent(uri string, assignment any, ctrlHandler messaging.Handler) (messaging.Agent, error) {
	return nil, nil
}

// NewEgressRoutingAgent - egress traffic routing for a host, one agent per rout
func NewEgressRoutingAgent(uri string, assignment any, ctrlHandler messaging.Handler) (messaging.Agent, error) {
	return nil, nil
}
