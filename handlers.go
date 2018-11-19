// Copyright 2018 Open Networking Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (c *Config) getSubscriberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sadisRequestID := vars["id"]

	log.WithFields(logrus.Fields{
		"sadisId": sadisRequestID,
	}).Infof("Looking for object %s in XOS database", sadisRequestID)

	defer r.Body.Close()

	subscribers := subscribers{}

	log.Debug("Checking subscribers")

	err := c.fetch("/xosapi/v1/rcord/rcordsubscribers", &subscribers)
	if err != nil {
		log.Errorf("Unable to retrieve subscriber information from XOS: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	for _, sub := range subscribers.Subscribers {
		if sub.OnuSerialNumber == sadisRequestID {
			log.Infof("Found subscriber with ID %s", sub.OnuSerialNumber)
			sadisSubscriber := sadisSubscriber{
				ID:        sub.OnuSerialNumber,
				CTag:      sub.CTag,
				STag:      sub.STag,
				NasPortID: sub.NasPortID,
				CircuitID: sub.CircuitID,
				RemoteID:  sub.RemoteID,
			}

			json, e := json.Marshal(&sadisSubscriber)
			if e != nil {
				log.Errorf("Unable to marshal JSON: %s", e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(json)
			return
		}
	}

	log.Debug("Checking devices")

	devices := oltDevices{}

	err = c.fetch("/xosapi/v1/volt/oltdevices", &devices)
	if err != nil {
		log.Errorf("Unable to retrieve device information from XOS: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	for _, device := range devices.OltDevices {
		// NOTE if it's an OLT then sadisRequestID is the device serial number
		devID := device.SerialNumber
		if devID == sadisRequestID {
			log.Infof("Found OLT device with ID %s", devID)

			ipaddr := device.Host
			addr, err := net.ResolveIPAddr("ip", device.Host)
			if err != nil {
				log.Errorf("Resolution error: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				ipaddr = addr.String()
			}
			log.Debugf("ID: %s, Uplink: %s, IPAddress: %s, NasID: %s", devID, toInt(device.Uplink), ipaddr, device.NasID)
			sadisDevice := sadisDevice{
				ID:         devID,
				Uplink:     toInt(device.Uplink),
				HardwareID: "de:ad:be:ef:ba:11", // TODO do we really need to configure this?
				IPAddress:  ipaddr,
				NasID:      device.NasID,
			}

			json, e := json.Marshal(&sadisDevice)
			if e != nil {
				log.Errorf("Unable to marshal JSON: %s", e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(json)
			return
		}
	}

	log.WithFields(logrus.Fields{
		"sadisId": sadisRequestID,
	}).Infof("Couldn't find object %s in XOS database", sadisRequestID)

	http.NotFound(w, r)
}

func toInt(value string) int {
	r, _ := strconv.Atoi(value)
	return r
}

func (c *Config) fetch(path string, data interface{}) error {
	resp, err := http.Get(c.connect + path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)
	return err
}
