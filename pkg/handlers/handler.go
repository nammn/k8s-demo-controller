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

package handlers

import (
	"github.com/Sirupsen/logrus"
	"github.com/nammn/k8s-demo-controller/pkg/controller"
	"k8s.io/api/core/v1"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	// Init initializes handler configuration
	// loads from viper/key store and setup connection
	Init() error
	// this is responsible for relaying the information to a specific backend
	Relay(event controller.Event) error
}

// Map maps each event handler function to a name for easily lookup
var Map = map[string]interface{}{
	"cloudant": &Cloudant{},
	"aurora":   &Aurora{},
	"local":    &Local{},
}

type Cloudant struct {
}

func (d *Cloudant) Init() error {
	return nil
}

func (d *Cloudant) Relay(event controller.Event) error {
	return nil
}

type Aurora struct {
}

func (a *Aurora) Init() error {
	return nil
}

func (a *Aurora) Relay(event controller.Event) error {
	return nil

}

/**
The Local Handler is only responsible to take the current obj and formats it into proper JSON to dump this into the log.
Mainly purpose is to show that the relay action is actually working for the Event case.
*/
type Local struct {
}

func (a *Local) Init() error {
	return nil
}

func (a *Local) Relay(event controller.Event) error {
	logrus.WithField("pkg", v1.Event{}).Infof("Processing %+v", event)

	return nil

}
