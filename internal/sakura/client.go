// Package sakura defines the Sakura Cloud mutation boundary.
package sakura

import (
	"context"
	"errors"
)

// ErrMutationsDisabled prevents accidental creation of paid resources.
var ErrMutationsDisabled = errors.New("Sakura Cloud mutations are disabled; pass --allow-mutations after the API implementation is reviewed")

// ErrNotImplemented marks the current scaffold boundary.
var ErrNotImplemented = errors.New("Sakura Cloud machine lifecycle is not implemented")

// MachineRequest contains all material needed to create one Talos machine.
type MachineRequest struct {
	RequestID         string
	Zone              string
	ExternalNetworkID string
	InternalSwitchID  string
	PacketFilterID    string
	ExternalIPPool    string
	InternalIPPool    string
	Gateway           string
	Schematic         string
	TalosVersion      string
	JoinConfig        string
	CPU               int
	MemoryGiB         int
	BootDiskGiB       int
}

// Machine identifies all Sakura resources created for a request.
type Machine struct {
	UUID          string
	ServerID      string
	BootDiskID    string
	CIDATACDROMID string
	ExternalIP    string
	InternalIP    string
}

// Client is the idempotent Sakura Cloud lifecycle API used by the provisioner.
type Client interface {
	Provision(context.Context, MachineRequest) (Machine, error)
	Deprovision(context.Context, string) error
}

// GuardedClient keeps the initial scaffold mutation-safe.
type GuardedClient struct {
	allowMutations bool
}

// NewGuardedClient creates a client with an explicit mutation gate.
func NewGuardedClient(allowMutations bool) *GuardedClient {
	return &GuardedClient{allowMutations: allowMutations}
}

// Provision refuses mutations until the Sakura API implementation exists.
func (client *GuardedClient) Provision(context.Context, MachineRequest) (Machine, error) {
	if !client.allowMutations {
		return Machine{}, ErrMutationsDisabled
	}

	return Machine{}, ErrNotImplemented
}

// Deprovision refuses mutations until ownership checks are implemented.
func (client *GuardedClient) Deprovision(context.Context, string) error {
	if !client.allowMutations {
		return ErrMutationsDisabled
	}

	return ErrNotImplemented
}
