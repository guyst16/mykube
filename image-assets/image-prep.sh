#!/bin/bash

# Set SELinux in permissive mode (effectively disabling it)
sudo setenforce 0
sudo sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

# This overwrites any existing configuration in /etc/yum.repos.d/kubernetes.repo
cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://pkgs.k8s.io/core:/stable:/v1.28/rpm/
enabled=1
gpgcheck=1
gpgkey=https://pkgs.k8s.io/core:/stable:/v1.28/rpm/repodata/repomd.xml.key
exclude=kubelet kubeadm kubectl cri-tools kubernetes-cni
EOF

#Disable swap 
sudo dnf remove -y zram-generator

sudo swapoff -a

#Enable network parameters
echo "br_netfilter" | sudo tee /etc/modules-load.d/br_netfilter.conf
echo "net.ipv4.ip_forward = 1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
sudo dnf install -y iproute-tc


#Install kubelet, kubeadm and kubectl, and enable kubelet to ensure it's automatically started on startup:
sudo yum install -y kubelet kubeadm kubectl --disableexcludes=kubernetes
sudo systemctl enable --now kubelet

#Install containerd
sudo yum install -y containerd
sudo systemctl enable --now containerd
sudo sh -c "containerd config default > /etc/containerd/config.toml" 
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/'  /etc/containerd/config.toml

#Install docker engine
sudo dnf -y install dnf-plugins-core
sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
sudo dnf install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl enable docker

#################################################################################################
# TODO: Document about virt-customize the base image                                            #
# virt-customize --run /home/guy/Projects/mykube/image-assets/image-prep.sh -a Base-image.qcow2 #
#################################################################################################