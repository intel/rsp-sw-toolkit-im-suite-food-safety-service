/*
 * INTEL CONFIDENTIAL
 * Copyright (2019) Intel Corporation.
 *
 * The source code contained or described herein and all documents related to the source code ("Material")
 * are owned by Intel Corporation or its suppliers or licensors. Title to the Material remains with
 * Intel Corporation or its suppliers and licensors. The Material may contain trade secrets and proprietary
 * and confidential information of Intel Corporation and its suppliers and licensors, and is protected by
 * worldwide copyright and trade secret laws and treaty provisions. No part of the Material may be used,
 * copied, reproduced, modified, published, uploaded, posted, transmitted, distributed, or disclosed in
 * any way without Intel/'s prior express written permission.
 * No license under any patent, copyright, trade secret or other intellectual property right is granted
 * to or conferred upon you by disclosure or delivery of the Materials, either expressly, by implication,
 * inducement, estoppel or otherwise. Any license under such intellectual property rights must be express
 * and approved by Intel in writing.
 * Unless otherwise agreed by Intel in writing, you may not remove or alter this notice or any other
 * notice embedded in Materials by Intel or Intel's suppliers or licensors in any way.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	log "github.com/sirupsen/logrus"
	"github.impcloud.net/RSP-Inventory-Suite/food-safety-sample/app/config"
	"github.impcloud.net/RSP-Inventory-Suite/food-safety-sample/app/notification"
	"github.impcloud.net/RSP-Inventory-Suite/food-safety-sample/app/tag"
	"github.impcloud.net/RSP-Inventory-Suite/utilities/go-metrics"
	reporter "github.impcloud.net/RSP-Inventory-Suite/utilities/go-metrics-influxdb"
)

const (
	serviceKey = "food-safety-sample"
)

const (
	inventoryEvent   = "inventory_event"
	temperatureEvent = "Temperature"
)

var (
	// Filter data by value descriptors (aka device resource name)
	valueDescriptors = []string{
		inventoryEvent,
		temperatureEvent,
	}
	// Global temperature value that is updated from EdgeX
	currentTempValue float32
)

type params struct {
	tagsMap map[string]interface{}
}

func main() {
	mConfigurationError := metrics.GetOrRegisterGauge("food-safety-sample.Main.ConfigurationError", nil)

	// Load config variables
	err := config.InitConfig()
	fatalErrorHandler("unable to load configuration variables", err, &mConfigurationError)

	// Initialize metrics reporting
	initMetrics()

	setLoggingLevel(config.AppConfig.LoggingLevel)

	log.WithFields(log.Fields{
		"Method": "main",
		"Action": "Start",
	}).Info("Starting Food Safety Sample...")

	// Register a subscriber to EdgeX notification service
	emails := strings.Split(config.AppConfig.EmailSubscribers, ",")
	if err := notification.RegisterSubscriber(emails, config.AppConfig.NotificationServiceURL); err != nil {
		log.Fatalf("Unable to register subscriber in EdgeX: %s", err)
	}

	// Build map of tracking epcs
	tags := strings.Split(config.AppConfig.TrackingEPCs, ",")
	tagsMap := make(map[string]interface{}, len(tags))
	for _, val := range tags {
		tagsMap[val] = nil
	}

	// Connect to EdgeX core data
	receiveZMQEvents(tagsMap)

	log.WithField("Method", "main").Info("Completed.")

}

func receiveZMQEvents(trackingEPCs map[string]interface{}) {

	param := params{tagsMap: trackingEPCs}

	//Initialized EdgeX apps functionSDK
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v", err))
		os.Exit(-1)
	}

	// Filter data by inventory_event from Inventory Suite and Temperature from temperature sensor
	edgexSdk.SetFunctionsPipeline(
		transforms.NewFilter(valueDescriptors).FilterByValueDescriptor,
		param.processEvents,
	)

	err := edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

}

func (param params) processEvents(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	if len(params) < 1 {
		return false, nil
	}

	event := params[0].(models.Event)
	if len(event.Readings) < 1 {
		return false, nil
	}

	for _, reading := range event.Readings {

		// RSP events
		switch reading.Name {
		case inventoryEvent:
			log.Debugf("inventory-event data received: %s", string(reading.Value))

			var invData tag.InventoryEvent
			if err := json.Unmarshal([]byte(reading.Value), &invData); err != nil {
				log.Errorf("Unable to unmarshal inventory event. %s", err)
			}

			// Check if the tag(s) reached its destination zone
			var tagsInDestination []tag.Tag
			for _, tagData := range invData.Data {
				if reached := tag.ReachedFreezer(tagData, config.AppConfig.FreezerReaderName, param.tagsMap); reached {
					tagsInDestination = append(tagsInDestination, tagData)
				}
			}

			// Send notification to EdgeX
			if len(tagsInDestination) > 0 {
				bodyContent := notification.CreateBodyContent(tagsInDestination, currentTempValue, config.AppConfig.FreezerReaderName)
				if err := notification.PostNotification(bodyContent, config.AppConfig.NotificationServiceURL); err != nil {
					log.Errorf("Unable to send notification to EdgeX. %s", err)
				}
			}
			break

		// Temperature sensor events
		case temperatureEvent:

			if reading.Device == config.AppConfig.TemperatureSensor {
				// Maintaining current value of temperature sensor in global variable
				currentTempValue = base64TemperatureToFloat32(reading.Value)
			}
			break
		}

	}

	return false, nil
}

func initMetrics() {
	// setup metrics reporting
	if config.AppConfig.TelemetryEndpoint != "" {
		go reporter.InfluxDBWithTags(
			metrics.DefaultRegistry,
			time.Second*10, //cfg.ReportingInterval,
			config.AppConfig.TelemetryEndpoint,
			config.AppConfig.TelemetryDataStoreName,
			"",
			"",
			nil,
		)
	}
}
