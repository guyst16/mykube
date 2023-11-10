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
sudo modprobe br_netfilter
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
sudo systemctl restart containerd

# #Install K8S cluster
# sudo kubeadm init --pod-network-cidr=10.244.0.0/16

# until sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf taint node hal9000 node-role.kubernetes.io/control-plane-;
# do
#     echo "Try untaint node..."
#     sleep 1
# done;

# until sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml;
# do
#     echo "Try install flannel..."
#     sleep 1
# done;

# until sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml;
# do
#     echo "Try install K8S console..."
#     sleep 1
# done;

# until [ -f /opt/cni/bin/flannel ];
# do 
#     echo "Waiting for flannel bin creation...";
#     sleep 1;
# done
# sudo cp /opt/cni/bin/* /usr/libexec/cni/
# sudo cp /usr/libexec/cni/* /opt/cni/bin/

# sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf patch svc kubernetes-dashboard -n kubernetes-dashboard --type='merge' -p '{"spec":{"type":"NodePort"}}'
# sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf patch svc kubernetes-dashboard -n kubernetes-dashboard --type='json' -p='[{"op": "replace", "path": "/spec/ports/0/nodePort", "value": 31000}]'