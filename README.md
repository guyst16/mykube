<p align="center">
  <img src="https://user-images.githubusercontent.com/100173467/201345993-2ecc594a-d066-42b2-bc76-76d743e45e2f.png">
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
