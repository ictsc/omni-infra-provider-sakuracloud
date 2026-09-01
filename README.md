# Omni Infrastructure Provider for Sakura Cloud

Dynamic Sakura Cloud machine lifecycle provider for Sidero Omni.

## Status

Initial, mutation-safe scaffold. It registers with Omni, publishes the Sakura
Cloud machine-class schema, persists per-request state in Omni, and defines the
provision/deprovision boundary. Sakura Cloud API mutations are deliberately
blocked until the resource implementation and ownership checks are complete.

The first supported target is an elastic worker pool. Existing control-plane and
Longhorn storage nodes remain fixed.

## Architecture

```text
MachineClass / desired worker count
                 |
                 v
        Omni MachineRequest
                 |
                 v
 omni-infra-provider-sakuracloud
                 |
                 v
 Server + boot disk + cidata + 2 NICs + packet filter
                 |
                 v
           Talos joins Omni
```

Terraform continues to own shared networks, router + switch, VIP pools, and the
Omni VM. This provider owns only resources created for a MachineRequest.

## Development

```bash
make test
make build
```

Run against Omni with mutations disabled:

```bash
export OMNI_ENDPOINT=https://omni.example.com
export OMNI_SERVICE_ACCOUNT_KEY='...'

./bin/omni-infra-provider-sakuracloud
```

Create the provider credential with:

```bash
omnictl infraprovider create sakuracloud
```

`--allow-mutations` currently only passes the safety gate; the Sakura API client
still returns `not implemented`. This guarantees the scaffold cannot create paid
resources.

See [`examples/machine-class.yaml`](examples/machine-class.yaml) for the planned
provider data shape.

## MVP milestones

1. Implement authenticated Sakura API client and read-only validation.
2. Cache Talos NoCloud images by schematic and Talos version.
3. Allocate external/internal addresses with collision-safe leases.
4. Create disk, cidata, server, NICs, and packet-filter attachment idempotently.
5. Verify create, Omni join, replacement, and guarded deprovision for one worker.
6. Add a separate autoscaler adapter after desired-size scaling is reliable.

## License

MPL-2.0.
