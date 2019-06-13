# Audit-Controller

The audit controller builds upon the simplest k8 controller I could find, which itself is based on kubewatch.
I used especially them because they do not work with CRD, which I do not need right now.

```
`demo-controller` is the simplest, yet fully valid, kubernetes controller I could up come with. When I wanted to learn how to build k8s controllers, I search the net and found only some general ideas or already quite complicated examples, that were actually really doing "something" or were using Custom Resource Definitions (CRDs).
https://github.com/piontec/k8s-demo-controller
```
## What does this controller do?

This controller  observes the audit logs and event streams of Kubernetes, 
and relays them out to another service outside the cluster.
Main goal is to expose the audit history of a cluster for analysis by another service.

## Building
You need `dep`. Get and install it here: [https://github.com/golang/dep](https://github.com/golang/dep). Then run,
```
# to fetch dependencies
dep ensure
# to build the whole thing
make
```

## Running
Make sure your `kubectl` is working. 

### Running as standalone binary
Just run `./audit-controller`. 

### Running as pod in a cluster
*  set `DOCKER_REPO` variable in [`Makefile`](Makefile) to point to a docker registry where you can store the image
*  run `make build-image` to build locally a docker image with your controller
*  run `make push-image` to push it to the registry
*  edit [`demo-controller.yaml`](demo-controller.yaml) and change `image: YOUR_URL:TAG` to pint to your image registry and the version tag you want to deploy
*  run `kubectl create -f demo-controller.yaml`

