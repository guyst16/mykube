#!/bin/bash

#Install K8S cluster
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
 
until sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf taint node hal9000 node-role.kubernetes.io/control-plane-;
do
    echo "Try untaint node..."
    sleep 1
done;
 
until sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml;
do
    echo "Try install flannel..."
    sleep 1
done;
 
until sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml;
do
    echo "Try install K8S console..."
    sleep 1
done;
 
until [ -f /opt/cni/bin/flannel ];
do 
    echo "Waiting for flannel bin creation...";
    sleep 1;
done
sudo cp /opt/cni/bin/* /usr/libexec/cni/
sudo cp /usr/libexec/cni/* /opt/cni/bin/
 
sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf patch svc kubernetes-dashboard -n kubernetes-dashboard --type='merge' -p '{"spec":{"type":"NodePort"}}'
sudo kubectl --kubeconfig=/etc/kubernetes/admin.conf patch svc kubernetes-dashboard -n kubernetes-dashboard --type='json' -p='[{"op": "replace", "path": "/spec/ports/0/nodePort", "value": 31000}]'