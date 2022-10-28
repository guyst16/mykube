#!/bin/bash

virsh destroy myFedoraVM; virsh undefine --remove-all-storage myFedoraVM

virt-install -n myFedoraVM --description "my test Fedora vm" --os-variant=fedora36 --ram=2048 --vcpus=2 --disk path=/var/lib/libvirt/images/myFedoraVM.img,bus=virtio,size=10 --graphics none --location /root/Fedora-Server-dvd-x86_64-36-1.5.iso --initrd-inject /root/k8s-installer/ks.cfg --extra-args='inst.ks=file:/ks.cfg console=tty0 console=ttyS0,115200n8 serial' --noautoconsole

echo hi
sleep 10
echo hiiii
sleep 60

ansible-playbook install-k8.yaml -e "ansible_password=qwe123" -i $(virsh domifaddr --domain myFedoraVM | grep ':' | awk '{print $4}' | cut -d'/' -f1),
