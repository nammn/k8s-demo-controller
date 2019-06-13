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

// Handler is implemented by any handler.
// The Handle method is used to process event
type Handler interface {
	Init() error
	Relay(obj interface{}) error
}

// Map maps each event handler function to a name for easily lookup
var Map = map[string]interface{}{
	"cloudant": &Cloudant{},
	"aurora":   &Aurora{},
}

// Cloudant handler implements Handler interface,
// print each event with JSON format
type Cloudant struct {
}

// Init initializes handler configuration
// loads from viper/key store and setup connection
func (d *Cloudant) Init() error {
	return nil
}

//TODO: this is responsible for relaying the information
func (d *Cloudant) Relay(obj interface{}) error {
	return nil
}

// Cloudant handler implements Handler interface,
// print each event with JSON format
type Aurora struct {
}

// Init initializes handler configuration
// loads from viper/key store and setup connection
func (a *Aurora) Init() error {
	return nil
}

//TODO: this is responsible for relaying the information
func (a *Aurora) Relay(obj interface{}) error {
	return nil

}
