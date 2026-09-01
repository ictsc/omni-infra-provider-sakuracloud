# Design decisions

## Ownership boundary

Terraform owns shared Sakura Cloud infrastructure: router + switch, internal
switches, VIP address ranges, and the Omni VM. This provider owns only resources
created for a single Omni `MachineRequest`: server, boot disk, cidata CD-ROM,
network attachments, packet filter attachment, and the allocated node addresses.

Never manage the same Sakura Cloud resource from both Terraform and this
provider.

## Safety

Cloud mutations are disabled unless `--allow-mutations` is explicitly set. A
future Sakura API implementation must remain idempotent by using the Omni machine
request ID as its stable key. Deprovisioning must refuse to delete resources that
do not carry that ownership marker.

OS disks and generated cidata are ephemeral. Longhorn/data disks are not part of
the initial provider scope and must never be deleted implicitly.

## Network model

Machines have two NICs. `eth0` attaches to the existing router + switch so Talos
API and ingress VIPs can use L2 announcements. `eth1` attaches to the internal
switch. Static addresses and Omni join configuration are injected through
NoCloud cidata because Sakura Cloud disk modification does not support Talos.

The Talos API remains protected by Sakura packet filters. Autoscaled machines are
workers only until control-plane replacement and quorum behavior are separately
implemented and tested.
