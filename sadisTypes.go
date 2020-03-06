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
	ID         string             `json:"id"`
	NasPortID  string             `json:"nasPortId"`
	CircuitID  string             `json:"circuitId"`
	RemoteID   string             `json:"remoteId"`
	UniTagList []*sadisUnitaginfo `json:"uniTagList"`
}

/*
  XOS RCORD subscriber format
*/
type subscriber struct {
	ID              int    `json:"id"`
	OnuSerialNumber string `json:"onu_device"`
	NasPortID       string `json:"nas_port_id"`
	CircuitID       string `json:"circuit_id"`
	RemoteID        string `json:"remote_id"`
	UniTagListId    []int  `json:"unitaglist_ids"`
}

type subscribers struct {
	Subscribers []*subscriber `json:"items"`
}

/*
  ONOS SADIS UNI Tag information format
*/
type sadisUnitaginfo struct {
	//FIXME: which fields can be omitted??
	UniTagMatch                int16  `json:"uniTagMatch"`
	PonCTag                    int16  `json:"ponCTag"`
	PonSTag                    int16  `json:"ponSTag"`
	UsPonCTagPriority          int    `json:"usPonCTagPriority"`
	UsPonSTagPriority          int    `json:"usPonSTagPriority"`
	DsPonCTagPriority          int    `json:"dsPonCTagPriority"`
	DsPonSTagPriority          int    `json:"dsPonSTagPriority"`
	TechnologyProfileID        int    `json:"technologyProfileId"`
	UpstreamBandwidthProfile   string `json:"upstreamBandwidthProfile"`
	DownstreamBandwidthProfile string `json:"downstreamBandwidthProfile"`
	ServiceName                string `json:"serviceName"`
	EnableMacLearning          bool   `json:"enableMacLearning"`
	ConfiguredMacAddress       *string `json:"configuredMacAddress,omitempty"`
	IsDhcpRequired             bool   `json:"isDhcpRequired"`
	IsIgmpRequired             bool   `json:"isIgmpRequired"`
}

type sadisUnitaginfolist struct {
	SadisUniTagList []*sadisUnitaginfo `json:"items"`
}

/*
  XOS RCORD UNI Tag information format
*/

type unitaginfo struct {
	//FIXME: which fields can be empty and must be NOT be propagated? (for example see configured_mac_address)
	ID                         int    `json:"id"`
	UniTagMatch                int16  `json:"uni_tag_match"`
	PonCTag                    int16  `json:"pon_c_tag"`
	PonSTag                    int16  `json:"pon_s_tag"`
	UsPonCTagPriority          int    `json:"us_pon_ctag_priority"`
	UsPonSTagPriority          int    `json:"us_pon_stag_priority"`
	DsPonCTagPriority          int    `json:"ds_pon_ctag_priority"`
	DsPonSTagPriority          int    `json:"ds_pon_stag_priority"`
	TechnologyProfileID        int    `json:"tech_profile_id"`
	UpstreamBandwidthProfile   int    `json:"upstream_bps_id"`
	DownstreamBandwidthProfile int    `json:"downstream_bps_id"`
	ServiceName                string `json:"service_name"`
	EnableMacLearning          bool   `json:"enable_mac_learning"`
	ConfiguredMacAddress       *string `json:"configured_mac_address"`
	IsDhcpRequired             bool   `json:"is_dhcp_required"`
	IsIgmpRequired             bool   `json:"is_igmp_required"`
	Subscriber                 int    `json:"subscriber_id"`
}

/*
  XOS BandwidthProfile format
*/
type bandwidthprofile struct {
	Name string `json:"name"`
	Cir  int    `json:"cir"`
	Cbs  int    `json:"cbs"`
	Eir  int    `json:"eir"`
	Ebs  int    `json:"ebs"`
	Air  int    `json:"air"`
}

type bandwidthprofiles struct {
	Profiles []*bandwidthprofile `json:"items"`
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

/*
  ONOS SADIS bandwidth profile format
*/

type sadisBandwidthProfile struct {
	ID  string `json:"id"`
	Cir int    `json:"cir"`
	Cbs int    `json:"cbs"`
	Eir int    `json:"eir"`
	Ebs int    `json:"ebs"`
	Air int    `json:"air"`
}
