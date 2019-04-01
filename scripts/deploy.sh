#!/bin/bash
#

cat << EOF >./scripts/deployment.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    run: ${1}
  name: ${1}
spec:
  replicas: 1
  selector:
    matchLabels:
      run: ${1}
  template:
    metadata:
      labels:
        run: ${1}
    spec:
      containers:
      - image: ${2}
        imagePullPolicy: Always
        env:
          - name: CONSUL_ADDR
            value: 10.0.91.169
          - name: CONSUL_PORT
            value: "8500"
        name: ${1}
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF

kubectl apply -f ./scripts/deployment.yaml