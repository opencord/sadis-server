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
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	lkh "github.com/gfremex/logrus-kafka-hook"
)

const appName = "SADISSERVER"

type Config struct {
	Port         int    `default:"8000" desc:"port on which to listen for requests"`
	Xos		     string `default:"127.0.0.1:8181" desc:"connection string with which to connect to XOS"`
	Username     string `default:"admin@opencord.org" desc:"username with which to connect to XOS"`
	Password     string `default:"letmein" desc:"password with which to connect to XOS"`
	LogLevel     string `default:"info" envconfig:"LOG_LEVEL" desc:"detail level for logging"`
	LogFormat    string `default:"text" envconfig:"LOG_FORMAT" desc:"log output format, text or json"`
	KafkaBroker  string `default:"" desc:"url of the kafka broker"`

	connect string
}

var logger = logrus.New()
var log *logrus.Entry
var appFlags = flag.NewFlagSet("", flag.ContinueOnError)
var config Config

func init()  {
	config = Config{}
	appFlags.Usage = func() {
		envconfig.Usage(appName, &config)
	}
	if err := appFlags.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			os.Exit(1)
		} else {
			return
		}
	}

	err := envconfig.Process(appName, &config)
	if err != nil {
		logger.Fatalf("[ERROR] Unable to parse configuration options : %s", err)
	}

	if len(config.KafkaBroker) > 0 {
		logger.Debug("Setting up kafka integration")
		hook, err := lkh.NewKafkaHook(
			"kh",
			[]logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel},
			&logrus.JSONFormatter{
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "@timestamp",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
				},
			},
			[]string{config.KafkaBroker},
		)

		if err != nil {
			logger.Error(err)
		}

		logger.Hooks.Add(hook)
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.WarnLevel
	}
	logger.Level = level

	switch config.LogFormat {
	case "json":
		logger.Formatter = &logrus.JSONFormatter{}
	default:
		logger.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		}
	}

	log = logger.WithField("topics", []string{"sadis-server.log"})
}

func main() {



	log.WithFields(logrus.Fields{
		"PORT": config.Port,
		"XOS": config.Xos,
		"USERNAME": config.Username,
		"PASSWORD": config.Password,
		"LOG_LEVEL": config.LogLevel,
		"LOG_FORMAT": config.LogFormat,
		"KAFKA_BROKER": config.KafkaBroker,
	}).Infof(`Sadis-server started`)

	router := mux.NewRouter()
	router.HandleFunc("/subscriber/{id}", config.getSubscriberHandler)
	http.Handle("/", router)

	connectStringFormat := "http://%s:%s@%s"
	config.connect = fmt.Sprintf(connectStringFormat, config.Username, config.Password, config.Xos)

	panic(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
