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
	types "github.com/nammn/k8s-demo-controller/pkg/common"
	"k8s.io/apimachinery/pkg/util/json"
)

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	// Init initializes handler configuration
	// loads from viper/key store and setup connection
	Init() error
	GetType() types.BackendTypes
	// this is responsible for relaying the information to a specific backend
	Relay(event types.RelayEvent) error
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

func (d *Cloudant) GetType() types.BackendTypes {
	return types.Cloudant
}

func (d *Cloudant) Relay(event types.RelayEvent) error {
	return nil
}

type Aurora struct {
}

func (a *Aurora) Init() error {
	return nil
}

func (a *Aurora) GetType() types.BackendTypes {
	return types.Aurora
}

func (a *Aurora) Relay(event types.RelayEvent) error {
	return nil

}

/**
The Local Handler is only responsible to take the current obj and formats it into proper JSON to dump this into the log.
Mainly purpose is to show that the relay action is actually working for the RelayEvent case.
*/
type Local struct {
}

func (a *Local) Init() error {
	return nil
}

func (a *Local) GetType() types.BackendTypes {
	return types.Local
}

func (a *Local) Relay(event types.RelayEvent) error {
	b, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Error trying to convert event to a json")
		return err
	}
	logrus.WithField("pkg", types.RelayEvent{}).Infof("Relaying to: %s the following information: %s", a.GetType(), b)

	return nil

}
