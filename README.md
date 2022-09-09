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

3. Create volume for the vm
Create storage pool on the current disk
```
$ mkdir ~/guest_images
$ cat <<EOF >> guest_pool.xml 
<pool type='fs'>
  <name>guest_images_fs</name>
  <source>
    <device path='$(df . | grep dev | cut -d" " -f1)'/>
  </source>
  <target>
    <path>~/guest_images</path>
  </target>
</pool> 
EOF

$ virsh pool-create guest_images.xml
```

Use the xml file for the volume
```
<volume>
  <name>volume1</name>
  <allocation>0</allocation>
  <capacity>20G</capacity>
  <target>
    <path>/var/lib/virt/images/sparse.img</path>
  </target>
</volume> 
```
Create the volume
```
$ virsh vol-create guest_images_dir guest_volume.xml
```

4. Create the vm
```
```
