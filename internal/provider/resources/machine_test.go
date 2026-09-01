package resources

import "testing"

func TestMachineStateRoundTrip(t *testing.T) {
	t.Parallel()

	machine := NewMachine("test", "request-1")
	want := MachineState{RequestID: "request-1", ServerID: "server-1", InternalIP: "192.168.100.128"}

	if err := Encode(machine, want); err != nil {
		t.Fatal(err)
	}

	got, err := Decode(machine)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf("state mismatch: got %#v, want %#v", got, want)
	}
}
