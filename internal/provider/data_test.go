package provider

import "testing"

func validData() Data {
	return Data{
		Zone:              "tk1b",
		ExternalNetworkID: "external",
		InternalSwitchID:  "internal",
		PacketFilterID:    "filter",
		ExternalIPPool:    "192.0.2.16-192.0.2.31",
		InternalIPPool:    "192.168.100.128-192.168.100.191",
		Gateway:           "192.0.2.1",
		CPU:               4,
		MemoryGiB:         8,
		BootDiskGiB:       100,
	}
}

func TestDataValidate(t *testing.T) {
	t.Parallel()

	if err := validData().Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDataValidateRejectsSmallDisk(t *testing.T) {
	t.Parallel()

	data := validData()
	data.BootDiskGiB = 5

	if err := data.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}
