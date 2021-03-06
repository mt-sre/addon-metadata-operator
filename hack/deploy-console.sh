#!/bin/bash

source ${PWD}/hack/_helpers.sh

# If release is missing, assume latest
release="${1:-latest}"

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: openshift-console
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: openshift-console
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: openshift-console
  namespace: kube-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: openshift-console
  name: openshift-console
  namespace: kube-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: openshift-console
  template:
    metadata:
      labels:
        app: openshift-console
    spec:
      serviceAccountName: openshift-console
      containers:
      - name: openshift-console
        image: quay.io/openshift/origin-console:${release}
        env:
        - name: BRIDGE_USER_AUTH
          value: "disabled"
---
apiVersion: v1
kind: Service
metadata:
  name: openshift-console
  namespace: kube-system
spec:
  selector:
    app: openshift-console
  ports:
  - port: 9000
    targetPort: 9000
EOF
