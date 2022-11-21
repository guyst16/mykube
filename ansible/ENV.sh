#!/bin/bash

HOST_PACKAGES="ansible qemu-kvm qemu-img libvirt python3-libvirt libvirt-client virt-install virt-viewer bridge-utils"

OS_ISO_PATH="/var/lib/libvirt/images"
OS_ISO_FULL_PATH="/var/lib/libvirt/images/Fedora-Server-dvd-x86_64-36-1.5.iso"
OS_ISO_URL="https://download.fedoraproject.org/pub/fedora/linux/releases/36/Server/x86_64/iso/Fedora-Server-dvd-x86_64-36-1.5.iso"
OS_ISO_SHORT_NAME="Fedora36"

OS_INFO_DB_URL="https://releases.pagure.org/libosinfo/osinfo-db-20221018.tar.xz"
OS_INFO_DB_FILE="osinfo-db-20221018.tar.xz"

VM_NAME="myFedoraVM"
VM_OS_VARIANT="fedora36"
VM_MEMORY="2048"
VM_VCPUS="2"
VM_DISK_PATH="/var/lib/libvirt/images/$VM_NAME.img"
VM_DISK_SIZE="20"

K8S_CONSOLE_DEPLOYMENT="true"

