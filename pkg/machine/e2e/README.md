# Working README for running the machine tests

Note: you must not have any machines defined before running tests
## Linux

### QEMU

`make localmachine`

## Microsoft Windows

### HyperV

1. Open a powershell as admin
1. $env:CONTAINERS_MACHINE_PROVIDER="hyperv"
1. $env:MACHINE_IMAGE="https://fedorapeople.org/groups/podman/testing/hyperv/fedora-coreos-38.20230830.dev.0-hyperv.x86_64.vhdx.zip"
1. `./test/tools/build/ginkgo.exe -vv  --tags "remote exclude_graphdriver_btrfs btrfs_noversion exclude_graphdriver_devicemapper containers_image_openpgp remote" -timeout=90m --trace --no-color  pkg/machine/e2e/. `

Note: Add `--focus-file "basic_test.go" ` to only run basic test

### WSL
1. Open a powershell as a regular user
1. Build and copy win-sshproxy into bin/
1. `./test/tools/build/ginkgo.exe -vv  --tags "remote exclude_graphdriver_btrfs btrfs_noversion exclude_graphdriver_devicemapper containers_image_openpgp remote" -timeout=90m --trace --no-color  pkg/machine/e2e/. `

Note: Add `--focus-file "basic_test.go" ` to only run basic test

## MacOS

### Apple Hypervisor

1. `make podman-remote`
1. `make .install.ginkgo`
1. `export TMPDIR=/Users/<yourname>`
1. `export CONTAINERS_MACHINE_PROVIDER="applehv"`
1. `export MACHINE_IMAGE="https://fedorapeople.org/groups/podman/testing/applehv/arm64/fedora-coreos-38.20230925.dev.0-applehv.aarch64.raw.gz"`
1. `./test/tools/build/ginkgo -vv  --tags "remote exclude_graphdriver_btrfs btrfs_noversion exclude_graphdriver_devicemapper containers_image_openpgp remote" -timeout=90m --trace --no-color  pkg/machine/e2e/.`

Note: Add `--focus-file "basic_test.go" ` to only run basic test
