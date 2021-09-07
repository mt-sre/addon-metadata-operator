#!/bin/bash

source ${PWD}/hack/_helpers.sh

echo "Installing hive-operator..."
kubectl create ns hive
cat <<EOF | kubectl apply -f -
---
apiVersion: operators.coreos.com/v1alpha2
kind: OperatorGroup
metadata:
  name: hive-operator
  namespace: hive
spec:
  targetNamespaces:
  - hive
---
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: hive-operator
  namespace: hive
spec:
  channel: alpha
  name: hive-operator
  source: operatorhubio-catalog
  sourceNamespace: olm
  installPlanApproval: Automatic
EOF

echo "Waiting for deployment/hive-operator and HiveConfig CRD..."
kubectl wait --for=condition=available deployment/hive-operator -n hive --timeout=120s
kubectl wait --for=condition=established crd/hiveconfigs.hive.openshift.io --timeout=120s

echo "Deploying remaining hive components (hive-controllers, hive-clustersync, hiveadmission)..."
cat <<EOF | kubectl apply -f -
apiVersion: hive.openshift.io/v1
kind: HiveConfig
metadata:
  name: hive
spec:
  logLevel: debug
  targetNamespace: hive
EOF

echo "Waiting for deployment/hive-controllers..."
kubectl wait --for=condition=available deployment/hive-controllers -n hive --timeout=240s
