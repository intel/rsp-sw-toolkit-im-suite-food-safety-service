/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"github.com/pkg/errors"
	"github.impcloud.net/RSP-Inventory-Suite/utilities/configuration"
)

type (
	variables struct {
		ServiceName, LoggingLevel, TrackingEPCs, TemperatureSensor  string
		TelemetryEndpoint, TelemetryDataStoreName                   string
		FreezerReaderName, NotificationServiceURL, EmailSubscribers string
	}
)

// AppConfig exports all config variables
var AppConfig variables

// InitConfig loads application variables
// nolint :gocyclo
func InitConfig() error {
	AppConfig = variables{}

	config, err := configuration.NewConfiguration()
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.ServiceName, err = config.GetString("serviceName")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	// Set "debug" for development purposes. Nil or "" for Production.
	AppConfig.LoggingLevel, err = config.GetString("loggingLevel")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.TelemetryEndpoint, err = config.GetString("telemetryEndpoint")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.TelemetryDataStoreName, err = config.GetString("telemetryDataStoreName")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.FreezerReaderName, err = config.GetString("freezerReaderName")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.NotificationServiceURL, err = config.GetString("notificationServiceURL")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.EmailSubscribers, err = config.GetString("emailSubscribers")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.TrackingEPCs, err = config.GetString("trackingEPCs")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	AppConfig.TemperatureSensor, err = config.GetString("temperatureSensor")
	if err != nil {
		return errors.Wrapf(err, "Unable to load config variables: %s", err.Error())
	}

	return nil
}
