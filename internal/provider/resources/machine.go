// Package resources defines state persisted in Omni for this provider.
package resources

import (
	"encoding/json"
	"fmt"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"
	"github.com/siderolabs/omni/client/pkg/infra"
	"google.golang.org/protobuf/types/known/wrapperspb"

	providermeta "github.com/ictsc/omni-infra-provider-sakuracloud/internal/provider/meta"
)

// MachineState is the durable record of Sakura resources owned by one request.
type MachineState struct {
	RequestID     string `json:"request_id,omitempty"`
	Schematic     string `json:"schematic,omitempty"`
	TalosVersion  string `json:"talos_version,omitempty"`
	MachineUUID   string `json:"machine_uuid,omitempty"`
	ServerID      string `json:"server_id,omitempty"`
	BootDiskID    string `json:"boot_disk_id,omitempty"`
	CIDATACDROMID string `json:"cidata_cdrom_id,omitempty"`
	ExternalIP    string `json:"external_ip,omitempty"`
	InternalIP    string `json:"internal_ip,omitempty"`
}

// NewMachine creates an empty provider-owned machine state resource.
func NewMachine(namespace, id string) *Machine {
	return typed.NewResource[MachineSpec, MachineExtension](
		resource.NewMetadata(namespace, infra.ResourceType("Machine", providermeta.ProviderID), id, resource.VersionUndefined),
		protobuf.NewResourceSpec(wrapperspb.String("{}")),
	)
}

// Machine is persisted by the Omni infra provider SDK.
type Machine = typed.Resource[MachineSpec, MachineExtension]

// MachineSpec stores the state as JSON in a protobuf string wrapper.
type MachineSpec = protobuf.ResourceSpec[wrapperspb.StringValue, *wrapperspb.StringValue]

// MachineExtension describes the provider state resource to COSI.
type MachineExtension struct{}

// ResourceDefinition implements typed.Extension.
func (MachineExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             infra.ResourceType("Machine", providermeta.ProviderID),
		DefaultNamespace: infra.ResourceNamespace(providermeta.ProviderID),
	}
}

// Decode returns the structured state.
func Decode(machine *Machine) (MachineState, error) {
	var state MachineState

	if err := json.Unmarshal([]byte(machine.TypedSpec().Value.Value), &state); err != nil {
		return MachineState{}, fmt.Errorf("decode machine state: %w", err)
	}

	return state, nil
}

// Encode replaces the resource's structured state.
func Encode(machine *Machine, state MachineState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("encode machine state: %w", err)
	}

	machine.TypedSpec().Value.Value = string(data)

	return nil
}
