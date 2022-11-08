dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo

dnf install -y containerd.io git go

containerd config default > /etc/containerd/config.toml;
sed -i  "s/SystemdCgroup = false/SystemdCgroup = true/" /etc/containerd/config.toml;
systemctl restart containerd;
systemctl enable containerd;
systemctl daemon-reload;
    
systemctl restart containerd

mkdir /tmp/networking-plugins
git clone https://github.com/containernetworking/plugins.git /tmp/networking-plugins
/tmp/networking-plugins/test_linux.sh



cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
exclude=kubelet kubeadm kubectl
EOF

setenforce 0

sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

modprobe -a overlay br_netfilter

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sysctl --system

swapoff -a;
sed -i '/swap/d' /etc/fstab

systemctl restart dbus;
systemctl restart firewalld;

firewall-cmd --permanent --add-port=6443/tcp;
firewall-cmd --permanent --add-port=2379-2380/tcp;
firewall-cmd --permanent --add-port=10250/tcp;
firewall-cmd --permanent --add-port=10251/tcp;
firewall-cmd --permanent --add-port=10252/tcp;
firewall-cmd --permanent --add-port=10255/tcp;
firewall-cmd --permanent --add-port=8472/udp;
firewall-cmd --add-masquerade --permanent;
firewall-cmd --permanent --add-port=30000-32767/tcp;
systemctl status firewalld;
systemctl restart firewalld;

dnf install -y kubelet kubeadm kubectl --disableexcludes=kubernetes

kubeadm init --pod-network-cidr=10.244.0.0/16

systemctl enable --now kubelet

sleep 30

export KUBECONFIG=/etc/kubernetes/admin.conf

kubectl taint node fedora node-role.kubernetes.io/control-plane-
kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml

cp -r /tmp/networking-plugins/bin/* /opt/cni/bin

