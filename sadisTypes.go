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

/*
  {
    "id": "PONSIM",
    "cTag": 333,
    "sTag": 444,
    "nasPortId": "PON 1/1/03/1:1.1.1",
    "circuitId": "foo",
    "remoteId": "bar"
  }
*/
/*
  ONOS SADIS subscriber format
*/
type sadisSubscriber struct {
	ID        string `json:"id"`
	CTag      int16  `json:"cTag"`
	STag      int16  `json:"sTag"`
	NasPortID string `json:"nasPortId"`
	CircuitID string `json:"circuitId"`
	RemoteID  string `json:"remoteId"`
}

/*
  XOS RCORD subscriber format
*/
type subscriber struct {
	CTag            int16  `json:"c_tag"`
	STag            int16  `json:"s_tag"`
	OnuSerialNumber string `json:"onu_device"`
	NasPortID       string `json:"nas_port_id"`
	CircuitID       string `json:"circuit_id"`
	RemoteID        string `json:"remote_id"`
}

type subscribers struct {
	Subscribers []*subscriber `json:"items"`
}

/*
  {
    "id" : "10.1.1.1:9191",
    "hardwareIdentifier" : "de:ad:be:ef:ba:11",
    "uplinkPort" : 128
  }
*/
/*
  ONOS SADIS device format
*/
type sadisDevice struct {
	ID         string `json:"id"`
	HardwareID string `json:"hardwareIdentifier"`
	Uplink     int    `json:"uplinkPort"`
	IPAddress  string `json:"ipAddress"`
	NasID      string `json:"nasId"`
}

/*
  XOS vOLT device format
*/
type oltDevice struct {
	Uplink       string `json:"uplink"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	NasID        string `json:"nas_id"`
	SerialNumber string `json:"serial_number"`
}

type oltDevices struct {
	OltDevices []*oltDevice `json:"items"`
}
