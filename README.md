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

> All steps are made on `Linux fedora 5.19.4-200.fc36.x86_64`

Run the next commands for deploying k8s single-node cluster installation:

1. Go to _MyKube_ [releases](https://github.com/guyst16/mykube/releases/tag/v0.0.1-alpha)
2. Download & install the desired release package
3. Run the next command for creating the cluster:
```
$ mykube create --domain <CLUSTER-NAME>
```

And you done.

### Need some help?

Ask for help:
![Screenshot from 2023-01-22 20-54-50](https://user-images.githubusercontent.com/100173467/213934360-3f867824-674e-4ab1-9415-832f3bf203e2.png)

### Delete cluster

Run the next command:

```
$ mykube delete --domain <CLUSTER-NAME>
```

### List clusters

Run the next command:

```
$ mykube list
```

### Get login credentials

Run the command:

```
$ mykube connect --domain <CLUSTER-NAME>
```


#### (Compatibility - currently) Linux kernel distributions:

- [x] ![Fedora](https://img.shields.io/badge/Fedora-294172?style=for-the-badge&logo=fedora&logoColor=white)

- [x] ![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white) 

- [x] ![Debian](https://img.shields.io/badge/Debian-D70A53?style=for-the-badge&logo=debian&logoColor=white)





## Feel free to open issues and suggestions

