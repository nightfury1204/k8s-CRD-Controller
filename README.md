# k8s-CRD-Controller

It's create same number of pods mantioned in replicas.
Initializer used for adding finalizer

```
#for initializer
kubectl create -f ./yaml/initializerConfigure.yaml

#For custom resource defination
kubectl create -f ./yaml/crd.yaml

kubectl create -f ./yaml/example_crd_obj.yaml

#run controller
go run main.go create

```

### initializerConfigure.yaml
```apple js
apiVersion: admissionregistration.k8s.io/v1alpha1
kind: InitializerConfiguration
metadata:
  name: podwatch-initializer-config
initializers:
  #the name needs to be fully qualified, i.e., containing at least two "."
  - name: podwatchinit.nahid.try.com
    rules:
      #apiGroups, apiVersion, resources all support wildcard "*".
      #"*" cannot be mixed with non-wildcard.
      - apiGroups:
          - nahid.try.com
        apiVersions:
          - v1alpha1
        resources:
          - podwatchs
```
### Custom Resource Definition
```apple js
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  #name must match the spec fields below, and be in the form: <plural>.<group>
  name: podwatchs.nahid.try.com
spec:
  #group name to use for REST API: /apis/<group>/<version>
  group: nahid.try.com
  #version name to use for REST API: /apis/<group>/<version>
  version: v1alpha1
  #either Namespaced or Cluster
  scope: Namespaced
  names:
    #plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: podwatchs
    #singular name to be used as an alias on the CLI and for display
    singular: podwatch
    #kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: PodWatch
```
### example_crd_obj.yaml
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
        #Just spin & wait forever
        command: [ "/bin/bash", "-c", "--" ]
        args: [ "while true; do sleep 30; done;" ]
```