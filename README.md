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

## What is missing?
*  Kobra configuration to set which backend to actually plug with to 
*  Concrete Backend integration
*  Only on can run at a time 
*  The controller only handles `Events`
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

## Auditing Backend
The current Events are saved in the following format
```
// Event indicate the informerEvent
type Event struct {
	key            string
	reason         string
	message        string
	firstTimestamp meta_v1.Time
	lastTimestamp  meta_v1.Time
	eventType      string
	namespace      string
	resourceType   string
}
``` 
and each `Event` from K8 will be handled by a pluggable backend.
The usual case is that the `Event` struct will be marshalled into a JSON and send to the corresponding API.

### Running as standalone binary
Just run `./audit-controller`. 
