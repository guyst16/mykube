<p align="center">
  <img src="https://user-images.githubusercontent.com/100173467/202853618-622951f8-86e6-4234-8352-c77bd498e840.png" width="615" height="600">
</p>


# k8s-installer - One-click k8s single-node cluster installation on your own device.

### Steps

> All steps are made on `Linux fedora 5.19.4-200.fc36.x86_64`

Run the next commands for deploying k8s single-node cluster installation:
```
$ git clone https://github.com/guyst16/mykube.git
$ cd mykube/ansible
$ ./start.sh
```

And you done.

* The script `start.sh` is reuseable which mean it will destrot the vm and create a new one if will executed again

### Delete and stop vm
Run the next 2 commands
```
$ virsh destroy myFedoraVM; virsh undefine --remove-all-storage myFedoraVM
```

### Quick ssh to vm
Run the command:
```
$ sshpass -p qwe123 \
ssh liveuser@$(virsh domifaddr --domain myFedoraVM | grep ':' | awk '{print $4}' | cut -d'/' -f1)
```

#### (Compatibility) Linux kernel distributions:
- [x] Fedora - 6.0.5-200.fc36.x86_64
- [ ] Ubuntu 
- [ ] Debian


## Feel free to open issues and suggestion
