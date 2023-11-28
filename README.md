![MyKube Logo](https://user-images.githubusercontent.com/100173467/202854244-a0b4d1c7-27a5-45f0-a2cb-b93615993c11.png)





# k8s-installer - One-click k8s single-node cluster installation on your own device.

<div align="center">
  <img src="https://img.shields.io/github/license/guyst16/mykube">
  <img src="https://img.shields.io/github/languages/code-size/guyst16/mykube"> 
  <img src="https://github.com/guyst16/mykube/workflows/Lint%20Code%20Base/badge.svg">
</div>
MyKube is a new easy-to-use open source tool for creating your own virtual machine with k8s installed only by one click.

All the dependencies are included, which means there are **no** previous steps that need to be taken.

### How to Use

##### Prerequisites
Install `Libvirt`
For RPM OS:
```
$ yum install qemu-kvm libvirt bridge-utils
```

Run the next commands for deploying k8s single-node cluster installation:

1. Go to _MyKube_ [releases](https://github.com/guyst16/mykube/releases/tag/v0.0.1-alpha)
2. Download & install the desired release package
3. Run the next command for creating the cluster:
```
$ mykube create --domain <NAME>
```
![image](https://github.com/guyst16/mykube/assets/100173467/4ac2ebb4-ce5b-4305-bab1-c659abebfc5d)


And you done.

### Need some help?

Ask for help:
```
$ mykube --help
```
### Destroy vm

Run the next command:

```
$ mykube delete --domain <NAME>
```

### Get connection details

Run the command:

```
$ mykube connect --domain <NAME>
```


#### (Compatibility) OS:

- ✔️ ![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)

- ⏲️ ![Mac OS](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white) 

- ⏲️ ![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)





## Feel free to open issues and suggestions

