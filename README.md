# k8s-installer

### Steps

> All steps are made on `Linux fedora 5.19.4-200.fc36.x86_64`

1. Install Libvirtd
```
$ sudo dnf install @virtualization;
$ sudo systemctl start libvirtd;
$ sudo systemctl enable libvirtd;
$ lsmod | grep kvm;
kvm_amd               114688  0
kvm                   831488  1 kvm_amd
```
2. Get Fedora-36 iso for the vm
```
wget https://download.fedoraproject.org/pub/fedora/linux/releases/36/Server/x86_64/iso/Fedora-Server-dvd-x86_64-36-1.5.iso
```
