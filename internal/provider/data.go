package provider

import (
	"fmt"
	"net/netip"
)

// Data is the provider-specific MachineClass configuration.
type Data struct {
	Zone              string `yaml:"zone"`
	ExternalNetworkID string `yaml:"external_network_id"`
	InternalSwitchID  string `yaml:"internal_switch_id"`
	PacketFilterID    string `yaml:"packet_filter_id"`
	ExternalIPPool    string `yaml:"external_ip_pool"`
	InternalIPPool    string `yaml:"internal_ip_pool"`
	Gateway           string `yaml:"gateway"`
	CPU               int    `yaml:"cpu"`
	MemoryGiB         int    `yaml:"memory_gib"`
	BootDiskGiB       int    `yaml:"boot_disk_gib"`
}

// Validate rejects incomplete or unsafe machine requests before cloud access.
func (data Data) Validate() error {
	switch {
	case data.Zone == "":
		return fmt.Errorf("zone is required")
	case data.ExternalNetworkID == "":
		return fmt.Errorf("external_network_id is required")
	case data.InternalSwitchID == "":
		return fmt.Errorf("internal_switch_id is required")
	case data.PacketFilterID == "":
		return fmt.Errorf("packet_filter_id is required")
	case data.ExternalIPPool == "":
		return fmt.Errorf("external_ip_pool is required")
	case data.InternalIPPool == "":
		return fmt.Errorf("internal_ip_pool is required")
	case data.CPU < 1:
		return fmt.Errorf("cpu must be at least 1")
	case data.MemoryGiB < 2:
		return fmt.Errorf("memory_gib must be at least 2")
	case data.BootDiskGiB < 20:
		return fmt.Errorf("boot_disk_gib must be at least 20")
	}

	if address, err := netip.ParseAddr(data.Gateway); err != nil || !address.Is4() {
		return fmt.Errorf("gateway must be a valid IPv4 address")
	}

	return nil
}
