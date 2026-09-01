// Package provider implements the Omni provisioning state machine.
package provider

import (
	"context"
	"fmt"

	"github.com/siderolabs/omni/client/pkg/infra/provision"
	omniinfra "github.com/siderolabs/omni/client/pkg/omni/resources/infra"
	"go.uber.org/zap"

	"github.com/ictsc/omni-infra-provider-sakuracloud/internal/provider/resources"
	"github.com/ictsc/omni-infra-provider-sakuracloud/internal/sakura"
)

// Provisioner translates Omni MachineRequests into Sakura Cloud machines.
type Provisioner struct {
	client sakura.Client
}

// NewProvisioner creates a provider using the supplied mutation boundary.
func NewProvisioner(client sakura.Client) *Provisioner {
	return &Provisioner{client: client}
}

// ProvisionSteps implements the Omni dynamic infrastructure provider SDK.
func (provider *Provisioner) ProvisionSteps() []provision.Step[*resources.Machine] {
	return []provision.Step[*resources.Machine]{
		provision.NewStep("validateRequest", func(_ context.Context, _ *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			var data Data
			if err := pctx.UnmarshalProviderData(&data); err != nil {
				return fmt.Errorf("decode provider data: %w", err)
			}

			if err := data.Validate(); err != nil {
				return fmt.Errorf("validate provider data: %w", err)
			}

			state, err := resources.Decode(pctx.State)
			if err != nil {
				return err
			}

			state.RequestID = pctx.GetRequestID()
			state.TalosVersion = pctx.GetTalosVersion()

			return resources.Encode(pctx.State, state)
		}),
		provision.NewStep("createSchematic", func(ctx context.Context, logger *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			state, err := resources.Decode(pctx.State)
			if err != nil {
				return err
			}

			if state.Schematic != "" {
				return nil
			}

			state.Schematic, err = pctx.GenerateSchematicID(ctx, logger, provision.WithoutConnectionParams())
			if err != nil {
				return fmt.Errorf("generate Talos schematic: %w", err)
			}

			return resources.Encode(pctx.State, state)
		}),
		provision.NewStep("provisionMachine", func(ctx context.Context, _ *zap.Logger, pctx provision.Context[*resources.Machine]) error {
			state, err := resources.Decode(pctx.State)
			if err != nil {
				return err
			}

			if state.ServerID != "" {
				pctx.SetMachineUUID(state.MachineUUID)

				return nil
			}

			var data Data
			if err = pctx.UnmarshalProviderData(&data); err != nil {
				return fmt.Errorf("decode provider data: %w", err)
			}

			machine, err := provider.client.Provision(ctx, sakura.MachineRequest{
				RequestID:         pctx.GetRequestID(),
				Zone:              data.Zone,
				ExternalNetworkID: data.ExternalNetworkID,
				InternalSwitchID:  data.InternalSwitchID,
				PacketFilterID:    data.PacketFilterID,
				ExternalIPPool:    data.ExternalIPPool,
				InternalIPPool:    data.InternalIPPool,
				Gateway:           data.Gateway,
				Schematic:         state.Schematic,
				TalosVersion:      state.TalosVersion,
				JoinConfig:        pctx.ConnectionParams.JoinConfig,
				CPU:               data.CPU,
				MemoryGiB:         data.MemoryGiB,
				BootDiskGiB:       data.BootDiskGiB,
			})
			if err != nil {
				return fmt.Errorf("provision Sakura Cloud machine: %w", err)
			}

			if machine.UUID == "" || machine.ServerID == "" {
				return fmt.Errorf("Sakura Cloud client returned incomplete machine identity")
			}

			state.MachineUUID = machine.UUID
			state.ServerID = machine.ServerID
			state.BootDiskID = machine.BootDiskID
			state.CIDATACDROMID = machine.CIDATACDROMID
			state.ExternalIP = machine.ExternalIP
			state.InternalIP = machine.InternalIP
			pctx.SetMachineUUID(machine.UUID)

			return resources.Encode(pctx.State, state)
		}),
	}
}

// Deprovision removes only resources recorded as owned by this request.
func (provider *Provisioner) Deprovision(ctx context.Context, _ *zap.Logger, machine *resources.Machine, request *omniinfra.MachineRequest) error {
	state, err := resources.Decode(machine)
	if err != nil {
		return err
	}

	if state.ServerID == "" {
		return nil
	}

	if state.RequestID == "" || state.RequestID != request.Metadata().ID() {
		return fmt.Errorf("refusing to delete server %q: ownership marker does not match request %q", state.ServerID, request.Metadata().ID())
	}

	if err = provider.client.Deprovision(ctx, state.RequestID); err != nil {
		return fmt.Errorf("deprovision Sakura Cloud machine: %w", err)
	}

	return nil
}
