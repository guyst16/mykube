#!/bin/bash

virsh destroy myFedoraVM; virsh undefine --remove-all-storage myFedoraVM

virt-install -n myFedoraVM --description "my test Fedora vm" --os-variant=fedora36 --ram=2048 --vcpus=2 --disk path=/var/lib/libvirt/images/myFedoraVM.img,bus=virtio,size=10 --graphics none --location /root/Fedora-Server-dvd-x86_64-36-1.5.iso --initrd-inject /root/k8s-installer/ks.cfg --extra-args='inst.ks=file:/ks.cfg console=tty0 console=ttyS0,115200n8' --noautoconsole --wait=-1

# Waiting for IP address
sleep 20

# Delete ssh fingerprint if exists
ssh-keygen -f ~/.ssh/known_hosts -R $(virsh domifaddr --domain myFedoraVM | grep ':' | awk '{print $4}' | cut -d'/' -f1)

ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook install-k8.yaml -e "ansible_password=qwe123" -i $(virsh domifaddr --domain myFedoraVM | grep ':' | awk '{print $4}' | cut -d'/' -f1),
