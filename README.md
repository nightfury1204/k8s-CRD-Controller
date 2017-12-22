# k8s-CRD-Controller

It's create same number of pods mantioned in replicas

```
kubectl create -f ./yaml/crd.yaml

kubectl create -f ./yaml/example.yaml

go run main.go create

```

### Custom Resource Defination
```
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: podwatchs.nahid.try.com
spec:
  # group name to use for REST API: /apis/<group>/<version>
  group: nahid.try.com
  # version name to use for REST API: /apis/<group>/<version>
  version: v1alpha1
  # either Namespaced or Cluster
  scope: Namespaced
  names:
    # plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: podwatchs
    # singular name to be used as an alias on the CLI and for display
    singular: podwatch
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: PodWatch
```
### example.yaml
```apple js
apiVersion: "nahid.try.com/v1alpha1"
kind: PodWatch
metadata:
  name: my-podwatch
spec:
  replicas: 2
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: ubuntu
        image: ubuntu:latest
        # Just spin & wait forever
        command: [ "/bin/bash", "-c", "--" ]
        args: [ "while true; do sleep 30; done;" ]
```