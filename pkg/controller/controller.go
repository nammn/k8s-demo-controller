/*
Copyright 2016 Skippbox, Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	types "github.com/nammn/k8s-demo-controller/pkg/common"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nammn/k8s-demo-controller/pkg/handlers"
	"github.com/nammn/k8s-demo-controller/pkg/utils"

	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 5

var serverStartTime time.Time

// Controller object
type Controller struct {
	logger    *logrus.Entry
	clientset kubernetes.Interface
	queue     workqueue.RateLimitingInterface
	informer  cache.SharedIndexInformer
}

func Start() {
	var kubeClient kubernetes.Interface
	_, err := rest.InClusterConfig()
	if err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kubeClient.CoreV1().Events(meta_v1.NamespaceAll).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kubeClient.CoreV1().Events(meta_v1.NamespaceAll).Watch(options)
			},
		},
		&api_v1.Event{},
		0, //Skip resync
		cache.Indexers{},
	)

	c := newResourceController(kubeClient, informer, "event")
	stopCh := make(chan struct{})
	defer close(stopCh)

	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func updateEvent(event *types.RelayEvent, obj interface{}, resourceType string) {
	if resourceType == "event" {
		event.FirstTimestamp = obj.(*api_v1.Event).FirstTimestamp
		event.LastTimestamp = obj.(*api_v1.Event).LastTimestamp
		event.Reason = obj.(*api_v1.Event).Reason
		event.Message = obj.(*api_v1.Event).Message
	}

}

func newResourceController(client kubernetes.Interface, informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent types.RelayEvent
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newEvent.Key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.EventType = "create"
			newEvent.ResourceType = resourceType
			logrus.WithField("pkg", resourceType).Infof("Processing add to %v: %s", resourceType, newEvent.Key)
			if err == nil {
				updateEvent(&newEvent, obj, resourceType)
				queue.Add(newEvent)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newEvent.Key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.EventType = "update"
			newEvent.ResourceType = resourceType
			logrus.WithField("pkg", resourceType).Infof("Processing update to %v: %s", resourceType, newEvent.Key)
			if err == nil {
				updateEvent(&newEvent, new, resourceType)
				queue.Add(newEvent)
			}
		},
		DeleteFunc: func(obj interface{}) {
			newEvent.Key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			newEvent.EventType = "delete"
			newEvent.ResourceType = resourceType
			logrus.WithField("pkg", resourceType).Infof("Processing delete to %v: %s", resourceType, newEvent.Key)
			if err == nil {
				updateEvent(&newEvent, obj, resourceType)
				queue.Add(newEvent)
			}
		},
	})

	return &Controller{
		logger:    logrus.WithField("pkg", resourceType),
		clientset: client,
		informer:  informer,
		queue:     queue,
	}
}

// Run starts the controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("Starting controller")
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	c.logger.Info("Controller synced and ready")

	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(types.RelayEvent))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.logger.Errorf("Error processing %s (will retry): %v", newEvent.(types.RelayEvent).Key, err)
		c.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		c.logger.Errorf("Error processing %s (giving up): %v", newEvent.(types.RelayEvent).Key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}

	return true
}

/**
Processes the current item and relays it to a handler given the environment variable.
The handler is responsible to actually handle the payload by usually marshalling and sending it to a common place
*/

func (c *Controller) processItem(newEvent types.RelayEvent) error {
	_, _, err := c.informer.GetIndexer().GetByKey(newEvent.Key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", newEvent.Key, err)
	}

	handlerType := os.Getenv("BACKENDHANDLERTYPE")
	var eventHandler handlers.Handler
	switch types.BackendTypes(handlerType) {
	case types.Aurora:
		eventHandler = new(handlers.Aurora)
	case types.Local:
		eventHandler = new(handlers.Local)
	case types.Cloudant:
		eventHandler = new(handlers.Cloudant)
	default:
		eventHandler = new(handlers.Local)
	}
	//Now just takes the cloudant handler
	return eventHandler.Relay(newEvent)
}
