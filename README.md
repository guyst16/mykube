<p align="center">
  <img src="https://user-images.githubusercontent.com/100173467/201345993-2ecc594a-d066-42b2-bc76-76d743e45e2f.png">
</p>


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
(Optional)
3. Create volume for the vm
Create storage pool on the current disk
```
$ mkdir ~/guest_images
$ cat <<EOF > guest_pool.xml 
<pool type='fs'>
  <name>guest_images_fs</name>
  <source>
    <device path='$(df . | grep dev | cut -d" " -f1)'/>
  </source>
  <target>
    <path>$HOME/guest_images</path>
  </target>
</pool> 
EOF

$ sudo virsh pool-create guest_pool.xml
$ sudo virsh pool-list --all
 Name              State    Autostart
---------------------------------------
 guest_images_fs   active   no
```

Create volume in the same pool
```
$ sudo virsh vol-create-as guest_images_fs volume1 40G
```
4. Create the <i>Kickstart</i> file
```
# System timezone
timezone Africa/Bissau
```

5. Create the vm
```
$ virt-install -n myFedoraVM --description "my test Fedora vm" --os-variant=fedora36 --ram=2048 --vcpus=2 --disk path=/var/lib/libvirt/images/myFedoraVM.img,bus=virtio,size=10 --graphics none --location /root/Fedora-Server-dvd-x86_64-36-1.5.iso --initrd-inject /root/test-ks.cfg --extra-args='inst.ks=file:/test-ks.cfg console=tty0 console=ttyS0,115200n8 serial'
```

### Delete and stop vm
Run the next 2 commands
```
$ virsh destroy myFedoraVM; virsh undefine --remove-all-storage myFedoraVM
```

### Quick ssh to vm
Run the command:
```
$ sshpass -p qwe123 ssh liveuser@$(virsh domifaddr --domain myFedoraVM | grep ':' | awk '{print $4}' | cut -d'/' -f1)
```
