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
$ sudo yum install qemu-kvm libvirt python3-libvirt libvirt-client bridge-utils
```
For DEB OS:
```
$ sudo sudo apt install qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils
```

Run the next commands for deploying k8s single-node cluster installation:

1. Go to _MyKube_ [releases](https://github.com/guyst16/mykube/releases/tag/v0.0.1-alpha)
2. Download & install the desired release package
3. Run the next command for creating the cluster:
```
$ sudo mykube create --domain <NAME>
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
$ sudo mykube delete --domain <NAME>
```

### Get connection details

Run the command:

```
$ sudo mykube connect --domain <NAME>
```

## How to build?
1. Run git clone:
```
$ git clone https://github.com/guyst16/mykube.git
$ cd mykube
```
2. Run go generate:
```
$ go generate pkg/embedfiles/util.go
```
3. Build:
```
$ go build
```
Done! Now you have a `mykube` binary file.

## So how does it work?
The Mykube procedure for creating new working K8S cluster is very simple, here are the steps:
1. Necessary directories for mykube are getting created.
1. A customized OS image downloaded for the virtual machine which the k8s will run above it (if not already downloaded).
2. A new virtual machine get deployed using cloud-init for automatic k8s installation.
3. Done!

#### (Compatibility) OS:

- ✔️ ![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)

- ⏲️ ![Mac OS](https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=apple&logoColor=white) 

- ⏲️ ![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)





## Feel free to open issues and suggestions

