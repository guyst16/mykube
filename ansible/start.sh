#!/bin/bash

# Source variables
# shellcheck source=/dev/null
source ENV.sh

display_logo()
{
    echo '                                                        '
    echo ' ||\      /|| \\  //      || //  ||   ||  ||==\\  ||===='
    echo ' ||\\    //||  \\//       ||//   ||   ||  ||   || ||____'
    echo ' || \\  // ||   ||        ||\\   ||   ||  ||== // ||    '
    echo ' ||  \\//  ||   ||        || \\  \\===//  ||___)) ||===='
    echo '                                                        '
    echo '                             ^                          '
    echo '                           xxxxx                        '
    echo '                         xxxxxxxxx                      '
    echo '                       xxxxxxxxxxxxx                    '
    echo '                     xxxxxxxxxxxxxxxxx                  '
    echo '                    x  xxxxxxxxxxxxx  x                 '
    echo '                    xxx  xxxxxxxxx  xxx                 '
    echo '                    xxxxx  xxxxx   xxxx                 '
    echo '                    xxxxxxx  x   xxxxxx                 '
    echo '                    xxxxxxxxx xxxxxxxxx                 '
    echo '                     xxxxxxxx xxxxxxxx                  '
    echo '                       xxxxxx xxxxxx                    '
    echo '                         xxxx xxxx                      '
    echo '                           xx xx                        '
}

display_help()
{
    # Display Help
    display_logo
    echo
    echo "MyKube is a new easy-to-use tool for creating your own virtual machine with k8s installed only by one click."
    echo
    echo "Syntax: ./start [-h|--help|--no-console-deployment|--destroy|--connect]"
    echo
    echo "options:"
    echo "--no-console-deployment  Disable console deployment."
    echo "--destroy                Destroy existing vms"
    echo "--connect                Connect to vm"
    echo "--help|-h                Print this Help."
    echo
}

# Destroy existing vms
destroy_vms()
{
    virsh destroy "$VM_NAME";
    virsh undefine --remove-all-storage "$VM_NAME";
}

# Connect to vm
connect_to_vm()
{
    sshpass -p qwe123 \
        ssh liveuser@"$(virsh domifaddr --domain $VM_NAME | grep ':' | awk '{print $4}' | cut -d'/' -f1)"
}

# Options
if [[ $1 = "--help" ]] || [[ $1 = "-h" ]];
then
    display_help
    exit 0;
elif [[ $1 = "--no-console-deployment" ]];
then
    K8S_CONSOLE_DEPLOYMENT="false"
elif [[ $1 = "--destroy" ]];
then
    destroy_vms
    exit 0;
elif [[ $1 = "--connect" ]];
then
    connect_to_vm
    exit 0;
elif [[ $1 != "" ]];
then
    display_help
    echo
    echo "argument $1 not found!"
    exit 0;
fi

# Declare vars for the ansible playbook
ANSIBLE_EXTRA_VARS="{'k8s_console_deployment':'$K8S_CONSOLE_DEPLOYMENT'}"

echo -e "\nDownloading packages..."
yum install -y "$HOST_PACKAGES"

echo -e "Start libvirtd service..."
systemctl start libvirtd

echo -e "\nValidate that $OS_ISO_SHORT_NAME image is ready..."
if [ -f "$OS_ISO_FULL_PATH" ]; then
    echo -e "$OS_ISO_SHORT_NAME iso file exists"
else
    echo -e "File not exists\nStart download $OS_ISO_SHORT_NAME iso file..."
    wget "$OS_ISO_URL" -P "$OS_ISO_PATH"
    echo -e "Iso file is ready to be used"
fi

echo -e "\nValidate that $OS_ISO_SHORT_NAME os exists in osdb-info..."
if (osinfo-query os | grep -iq "$OS_ISO_SHORT_NAME"); then
    echo -e "$OS_ISO_SHORT_NAME os exists"
else
    echo -e "$OS_ISO_SHORT_NAME os does not exists\nStart  updating OS..."
    wget "$OS_INFO_DB_URL"
    osinfo-db-import "$OS_INFO_DB_FILE"
    echo -e "OS db updated"
fi


echo -e "\nTry deleting '$VM_NAME' if exists..."
destroy_vms

echo -e "\nCheck if default network is activated"
if virsh net-info --network default | grep Active | grep -q yes; then
    echo -e "default network is activated"
else
    echo -e "default network is not activated\nActivating default network"
    virsh net-start default;
fi

echo -e "\nStart deploying the new vm..."
virt-install -n "$VM_NAME" \
    --description "my test $OS_ISO_SHORT_NAME vm" \
    --os-variant="$VM_OS_VARIANT" --ram="$VM_MEMORY" \
    --vcpus="$VM_VCPUS" \
    --disk path="$VM_DISK_PATH",bus=virtio,size="$VM_DISK_SIZE" \
    --graphics=none \
    --location="$OS_ISO_FULL_PATH" \
    --initrd-inject=../ks.cfg \
    --extra-args='inst.ks=file:/ks.cfg console=tty0 console=ttyS0,115200n8' \
    --noautoconsole \
    --wait=-1

# Waiting for IP address
echo -e "\nWait 20 seconds for IP address to get assigned..."
sleep 20

# VM IP address
echo -e "Find IP address"
VM_IP_ADDRESS=$(virsh domifaddr --domain "$VM_NAME" |
                grep ':' |
                awk '{print $4}' |
                cut -d'/' -f1)
echo -e "IP address is: $VM_IP_ADDRESS"

# Delete ssh fingerprint if exists
echo -e "Delete fingerprint from ~/.ssh/known_hosts if exists..."
ssh-keygen -f ~/.ssh/known_hosts -R "$VM_IP_ADDRESS"

echo -e "\nInstall k8s module for ansible"
ansible-galaxy collection install kubernetes.core

echo "$ANSIBLE_EXTRA_VARS"

echo -e "\nRun ansible-playbook for deploying k8s..."
ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook install-k8.yaml -b \
    -e "ansible_password=qwe123" \
    -i "$VM_IP_ADDRESS", \
    -e "$ANSIBLE_EXTRA_VARS"