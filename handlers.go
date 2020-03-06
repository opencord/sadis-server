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
	"net"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

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
				NasPortID: sub.NasPortID,
				CircuitID: sub.CircuitID,
				RemoteID:  sub.RemoteID,
			}

			log.Debugf("Fetching UNI Tag list for subscriber %s", sub.OnuSerialNumber)

			unitaglist := sadisUnitaginfolist{}
			for _, unitagid := range sub.UniTagListId {
				utinfo := unitaginfo{}
				err = c.getOneUniTagInfo(unitagid, &utinfo)
				if err != nil {
					log.Errorf("Cannot fetch UNI tag information%s for subscriber %s", strconv.Itoa(unitagid), sub.OnuSerialNumber)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if (unitaginfo{}) == utinfo {
					// it's empty
					log.WithFields(logrus.Fields{
						"UniTagInfoId": unitagid,
						"Subscriber":   sub.OnuSerialNumber,
						"sadisId":      sadisRequestID,
					}).Error("UNI Tag info not found in XOS")
					http.Error(w, "UNI Tag info not found in XOS", http.StatusInternalServerError)
					return
				}
				sadisUnitaginfo := sadisUnitaginfo{
					UniTagMatch:          utinfo.UniTagMatch,
					PonCTag:              utinfo.PonCTag,
					PonSTag:              utinfo.PonSTag,
					UsPonCTagPriority:    utinfo.UsPonCTagPriority,
					UsPonSTagPriority:    utinfo.UsPonSTagPriority,
					DsPonCTagPriority:    utinfo.DsPonCTagPriority,
					DsPonSTagPriority:    utinfo.DsPonSTagPriority,
					TechnologyProfileID:  utinfo.TechnologyProfileID,
					ServiceName:          utinfo.ServiceName,
					EnableMacLearning:    utinfo.EnableMacLearning,
					ConfiguredMacAddress: utinfo.ConfiguredMacAddress,
					IsDhcpRequired:       utinfo.IsDhcpRequired,
					IsIgmpRequired:       utinfo.IsIgmpRequired,
				}

				log.Debugf("Fetching bandwidth profiles for subscriber %s and unitagid: %s", sub.OnuSerialNumber, strconv.Itoa(utinfo.ID))

				dsBandwidthprofile := bandwidthprofile{}
				err = c.getOneBandwidthProfileHandler(utinfo.DownstreamBandwidthProfile, &dsBandwidthprofile)
				if err != nil {
					log.Errorf("Cannot fetch downstream bandwidth profile %s for subscriber %s", strconv.Itoa(utinfo.DownstreamBandwidthProfile), sub.OnuSerialNumber)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if (bandwidthprofile{}) == dsBandwidthprofile {
					// it's empty
					log.WithFields(logrus.Fields{
						"DownstreamBandwidthProfile": utinfo.DownstreamBandwidthProfile,
						"Subscriber":                 sub.OnuSerialNumber,
						"sadisId":                    sadisRequestID,
					}).Error("Downstream bandwidth profile not found in XOS")
					http.Error(w, "Downstream bandwidth profile not found in XOS", http.StatusInternalServerError)
					return
				}
				sadisUnitaginfo.DownstreamBandwidthProfile = dsBandwidthprofile.Name

				usBandwidthprofile := bandwidthprofile{}
				err = c.getOneBandwidthProfileHandler(utinfo.UpstreamBandwidthProfile, &usBandwidthprofile)
				if err != nil {
					log.Errorf("Cannot fetch upstream bandwidth profile %s for subscriber %s", strconv.Itoa(utinfo.UpstreamBandwidthProfile), sub.OnuSerialNumber)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if (bandwidthprofile{}) == usBandwidthprofile {
					// it's empty
					log.WithFields(logrus.Fields{
						"UpstreamBandwidthProfile": utinfo.UpstreamBandwidthProfile,
						"Subscriber":               sub.OnuSerialNumber,
						"sadisId":                  sadisRequestID,
					}).Error("Upstream bandwidth profile not found in XOS")
					http.Error(w, "Upstream bandwidth profile not found in XOS", http.StatusInternalServerError)
					return
				}
				sadisUnitaginfo.UpstreamBandwidthProfile = usBandwidthprofile.Name

				log.WithFields(logrus.Fields{
					"Subscriber":                 sub.OnuSerialNumber,
					"UniTagInfo":                 utinfo.ID,
					"UpstreamBandwidthProfile":   usBandwidthprofile.Name,
					"DownstreamBandwidthProfile": dsBandwidthprofile.Name,
					"sadisId": sadisRequestID,
				}).Debug("Bandwidth profiles for subscriber/unitaginfo")
				unitaglist.SadisUniTagList = append(unitaglist.SadisUniTagList, &sadisUnitaginfo)
			}
			sadisSubscriber.UniTagList = unitaglist.SadisUniTagList

			sadisjson, e := json.Marshal(&sadisSubscriber)
			if e != nil {
				log.Errorf("Unable to marshal JSON: %s", e)
				http.Error(w, e.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(sadisjson)
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
			log.Debugf("ID: %s, Uplink: %d, IPAddress: %s, NasID: %s", devID, toInt(device.Uplink), ipaddr, device.NasID)
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

func (c *Config) getBandwidthProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sadisRequestID := vars["id"]

	log.WithFields(logrus.Fields{
		"sadisId": sadisRequestID,
	}).Infof("Fetching BandwidthProfiles from XOS database")
	defer r.Body.Close()

	bandwidthprofiles := bandwidthprofiles{}

	err := c.fetch("/xosapi/v1/rcord/bandwidthprofiles", &bandwidthprofiles)

	if err != nil {
		log.Errorf("Unable to retrieve bandwidth profiles information from XOS: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, profile := range bandwidthprofiles.Profiles {
		profileID := profile.Name
		if profileID == sadisRequestID {
			sadisProfile := sadisBandwidthProfile{
				ID:  profile.Name,
				Cir: profile.Cir,
				Cbs: profile.Cbs,
				Eir: profile.Eir,
				Ebs: profile.Ebs,
				Air: profile.Air,
			}
			json, e := json.Marshal(&sadisProfile)
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

func (c *Config) getOneBandwidthProfileHandler(id int, data interface{}) error {
	err := c.fetch("/xosapi/v1/rcord/bandwidthprofiles/"+strconv.Itoa(id), &data)
	if err != nil {
		log.Errorf("Unable to retrieve bandwidth profile information from XOS: %s", err)
		return err
	}
	return nil
}

func (c *Config) getOneUniTagInfo(id int, data interface{}) error {
	err := c.fetch("/xosapi/v1/rcord/rcordunitaginformations/"+strconv.Itoa(id), &data)
	if err != nil {
		log.Errorf("Unable to retrieve UNI Tag information from XOS: %s", err)
		return err
	}
	return nil
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
	log.WithFields(logrus.Fields{
		"path":   path,
		"resp":   data,
		"status": resp.Status,
	}).Debug("Received data from XOS")
	return err
}
