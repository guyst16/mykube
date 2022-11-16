#!/bin/bash

echo -e "\nDownloading packages..."
yum install ansible qemu-kvm qemu-img libvirt python3-libvirt libvirt-client virt-install virt-viewer bridge-utils

echo -e "\nStart sshd service..."
systemctl start sshd

echo -e "\nStart libvirtd service..."
systemctl start libvirtd

echo -e "\nValidate that fedora image is ready..."
if [ -f /var/lib/libvirt/images/Fedora-*.iso ]; then 
	echo -e "\nFedora iso file exists"
else 
	echo -e "\nFile not exists\nStart download Fedora iso file..."
	wget https://download.fedoraproject.org/pub/fedora/linux/releases/36/Server/x86_64/iso/Fedora-Server-dvd-x86_64-36-1.5.iso -P /var/lib/libvirt/images
	echo -e "\nIso file is ready for use"
fi

echo -e "\nValidate that fedora os exists in osdb-info..."
if (osinfo-query os | grep -q fedora36); then
        echo -e "\nFedora os exists"
else
        echo -e "\nFedor os does not exists\nStart  updating OS..."
        wget https://releases.pagure.org/libosinfo/osinfo-db-20221018.tar.xz
	osinfo-db-import osinfo-db-20221018.tar.xz
        echo -e "\nOS db updated"
fi


echo -e "\nTry deleting 'myFedoraVM' if exists..."
virsh destroy myFedoraVM; virsh undefine --remove-all-storage myFedoraVM

echo -e "\nCheck if default network is activated"
if virsh net-info --network default | grep Active | grep -q yes; then
	echo -e "\ndefault network is activated"
else
	echo -e "\ndefault network is not activated\nActivating default network"
	virsh net-start default;
fi

echo -e "\nStart deploying the new vm..."
virt-install -n myFedoraVM --description "my test Fedora vm" --os-variant=fedora36 --ram=2048 --vcpus=2 --disk path=/var/lib/libvirt/images/myFedoraVM.img,bus=virtio,size=20 --graphics none --location /var/lib/libvirt/images/Fedora-Server-dvd-x86_64-36-1.5.iso --initrd-inject ../ks.cfg --extra-args='inst.ks=file:/ks.cfg console=tty0 console=ttyS0,115200n8' --noautoconsole --wait=-1

# Waiting for IP address
echo -e "\nWait 20 seconds for IP address to get assigned..."
sleep 20

# Delete ssh fingerprint if exists
echo -e "\nDelete fingerprint from ~/.ssh/known_hosts if exists..."
ssh-keygen -f ~/.ssh/known_hosts -R $(virsh domifaddr --domain myFedoraVM | grep ':' | awk '{print $4}' | cut -d'/' -f1)

echo -e "\nInstall k8s module for ansible"
ansible-galaxy collection install kubernetes.core

echo -e "\nRun ansible-playbook for deploying k8s..."
ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook install-k8.yaml -e "ansible_password=qwe123" -i $(virsh domifaddr --domain myFedoraVM | grep ':' | awk '{print $4}' | cut -d'/' -f1), -b
