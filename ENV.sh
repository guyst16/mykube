#!/bin/bash

export OS_ISO_PATH="/var/lib/libvirt/images"
export OS_ISO_FULL_PATH="/var/lib/libvirt/images/Fedora-Server-dvd-x86_64-36-1.5.iso"
export OS_ISO_URL="https://download.fedoraproject.org/pub/fedora/linux/releases/36/Server/x86_64/iso/Fedora-Server-dvd-x86_64-36-1.5.iso"
export OS_ISO_SHORT_NAME="Fedora36"
 
export OS_INFO_DB_URL="https://releases.pagure.org/libosinfo/osinfo-db-20221018.tar.xz"
export OS_INFO_DB_FILE="osinfo-db-20221018.tar.xz"
 
export VM_NAME="myFedoraVM"
export VM_OS_VARIANT="fedora36"
export VM_MEMORY="2048"
export VM_VCPUS="2"
export VM_DISK_PATH="/var/lib/libvirt/images/$VM_NAME.img"
export VM_DISK_SIZE="20"

export K8S_CONSOLE_DEPLOYMENT="true"
export K8S_CONSOLE_NODE_PORT="32000"

