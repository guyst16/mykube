![MyKube Logo](https://user-images.githubusercontent.com/100173467/202854244-a0b4d1c7-27a5-45f0-a2cb-b93615993c11.png)





# k8s-installer - One-click k8s single-node cluster installation on your own device.



### About


![GitHub](https://img.shields.io/github/license/guyst16/mykube) ![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/guyst16/mykube) 
MyKube is a new easy-to-use open source tool for creating your own virtual machine with k8s installed only by one click.

All the dependencies are included, which means there are **no** previous steps that need to be taken.



### How to Install



> All steps are made on `Linux fedora 5.19.4-200.fc36.x86_64`



Run the next commands for deploying k8s single-node cluster installation:

```

$ git clone https://github.com/guyst16/mykube.git

$ cd mykube/ansible

$ ./start.sh

```



And you done.



* The script `start.sh` is reuseable, which means it will destroy the vm and create a new one if re-executed



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

- [x] ![Fedora](https://img.shields.io/badge/Fedora-294172?style=for-the-badge&logo=fedora&logoColor=white) - `6.0.5-200.fc36.x86_64`

- [ ] ![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white) 

- [ ] ![Debian](https://img.shields.io/badge/Debian-D70A53?style=for-the-badge&logo=debian&logoColor=white)





## Feel free to open issues and suggestions

