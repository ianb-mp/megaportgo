//go:build integration

// Copyright 2020 Megaport Pty Ltd
//
// Licensed under the Mozilla Public License, Version 2.0 (the
// "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//       https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package port

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/megaport/megaportgo/config"
	"github.com/megaport/megaportgo/mega_err"
	"github.com/megaport/megaportgo/service/authentication"
	"github.com/megaport/megaportgo/service/location"
	"github.com/megaport/megaportgo/shared"
	"github.com/megaport/megaportgo/types"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_LOCATION_ID_A = 19 // 	Interactive 437 Williamstown
	MEGAPORTURL        = "https://api-staging.megaport.com/"
)

var logger *config.DefaultLogger
var cfg config.Config

func TestMain(m *testing.M) {
	logger = config.NewDefaultLogger()
	logger.SetLevel(config.DebugLevel)

	clientID := os.Getenv("MEGAPORT_ACCESS_KEY")
	clientSecret := os.Getenv("MEGAPORT_SECRET_KEY")
	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel != "" {
		logger.SetLevel(config.StringToLogLevel(logLevel))
	}

	if clientID == "" {
		logger.Error("MEGAPORT_ACCESS_KEY environment variable not set.")
		os.Exit(1)
	}

	if clientSecret == "" {
		logger.Error("MEGAPORT_SECRET_KEY environment variable not set.")
		os.Exit(1)
	}

	cfg = config.Config{
		Log:      logger,
		Endpoint: MEGAPORTURL,
	}

	auth := authentication.New(&cfg, clientID, clientSecret)
	token, loginErr := auth.Login()

	if loginErr != nil {
		logger.Errorf("LoginError: %s", loginErr.Error())
	}

	cfg.SessionToken = token
	os.Exit(m.Run())
}

// TestSinglePort tests the creation of a LAG Port, then passes the id to PortScript to finalise lifecycle testing.
func TestSinglePort(t *testing.T) {
	port := New(&cfg)
	loc := location.New(&cfg)

	testLocation, err := loc.GetLocationByID(TEST_LOCATION_ID_A)

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	portsListInitial, err := port.GetPorts()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	portId, portErr := testCreatePort(port, types.SINGLE_PORT, testLocation)

	if !assert.NoError(t, portErr) && !assert.True(t, shared.IsGuid(portId)) {
		cfg.PurchaseError(portId, portErr)
		t.FailNow()
	}

	portCreated, err := port.WaitForPortProvisioning(portId)

	if !assert.NoError(t, err) || !portCreated {
		t.FailNow()
	}

	portsListPostCreate, err := port.GetPorts()
	if err != nil {
		logger.Debug("Failed to get ports list: ", err.Error())
		t.FailNow()
	}

	portIsActuallyNew := true
	for _, p := range portsListInitial {
		if p.UID == portId {
			portIsActuallyNew = false
		}
	}

	if !portIsActuallyNew {
		logger.Debug("Failed to find port we just created in ports list: ", portId)
		t.FailNow()
	}

	foundNewPort := false
	for _, p := range portsListPostCreate {
		if p.UID == portId {
			foundNewPort = true
		}
	}

	if !foundNewPort {
		logger.Debug("Failed to find port we just created in ports list: ", portId)
		t.FailNow()
	}

	testModifyPort(port, portId, types.SINGLE_PORT, t)
	testLockPort(port, portId, t)
	testCancelPort(port, portId, types.SINGLE_PORT, t)
	testDeletePort(port, portId, types.SINGLE_PORT, t)

}

// TestLAGPort tests the creation of a LAG Port, then passes the id to PortScript to finalise lifecycle testing.
func TestLAGPort(t *testing.T) {
	port := New(&cfg)
	loc := location.New(&cfg)

	testLocation, err := loc.GetLocationByID(TEST_LOCATION_ID_A)

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	portsListInitial, err := port.GetPorts()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	portId, portErr := testCreatePort(port, types.LAG_PORT, testLocation)

	if !assert.NoError(t, portErr) && !assert.True(t, shared.IsGuid(portId)) {
		cfg.PurchaseError(portId, portErr)
		t.FailNow()
	}

	portCreated, err := port.WaitForPortProvisioning(portId)

	if !assert.NoError(t, err) || !portCreated {
		t.FailNow()
	}

	portsListPostCreate, err := port.GetPorts()
	if err != nil {
		logger.Debug("Failed to get ports list: ", err.Error())
		t.FailNow()
	}

	portIsActuallyNew := true
	for _, p := range portsListInitial {
		if p.UID == portId {
			portIsActuallyNew = false
		}
	}

	if !portIsActuallyNew {
		logger.Debug("Failed to find port we just created in ports list: ", portId)
		t.FailNow()
	}

	foundNewPort := false
	for _, p := range portsListPostCreate {
		if p.UID == portId {
			foundNewPort = true
		}
	}

	if !foundNewPort {
		logger.Debug("Failed to find port we just created in ports list: ", portId)
		t.FailNow()
	}

	testModifyPort(port, portId, types.LAG_PORT, t)
	testCancelPort(port, portId, types.LAG_PORT, t)
}

func testCreatePort(port *Port, portType string, location types.Location) (string, error) {
	var portId string
	var portErr error

	logger.Debug("Buying Port:", portType)
	if portType == types.LAG_PORT {
		portId, portErr = port.BuyLAGPort("Buy Port (LAG) Test", 1, 10000, location.ID, location.Market, 2, true)
	} else {
		portId, portErr = port.BuySinglePort("Buy Port (Single) Test", 1, 10000, location.ID, location.Market, true)
	}

	logger.Debugf("Port Purchased: %s", portId)
	return portId, portErr
}

func testModifyPort(port *Port, portId string, portType string, t *testing.T) {
	portInfo, _ := port.GetPortDetails(portId)

	newPortName := fmt.Sprintf("Buy Port (%s) [Modified]", portType)

	logger.Debugf("Modifying Port: %s", portType)
	_, modifyErr := port.ModifyPort(portId, newPortName, "", portInfo.MarketplaceVisibility)
	assert.NoError(t, modifyErr)

	portInfo, _ = port.GetPortDetails(portId)
	assert.EqualValues(t, newPortName, portInfo.Name)
}

// PortScript tests the remaining lifecycle for a Port (not dependant on port-type), Go-Live, Modification,
// and Soft/Hard Deletes.
func testCancelPort(port *Port, portId string, portType string, t *testing.T) {
	// Soft Delete
	logger.Debugf("Scheduling %s Port for deletion (30 days).", portType)
	softDeleteStatus, deleteErr := port.DeletePort(portId, false)
	assert.True(t, softDeleteStatus)
	assert.NoError(t, deleteErr)

	portInfo, _ := port.GetPortDetails(portId)
	assert.EqualValues(t, types.STATUS_CANCELLED, portInfo.ProvisioningStatus)

	logger.Debugf("Status is: '%s'", portInfo.ProvisioningStatus)
	restoreStatus, restoreErr := port.RestorePort(portId)
	assert.True(t, restoreStatus)
	assert.NoError(t, restoreErr)
}

func testDeletePort(port *Port, portId string, portType string, t *testing.T) {
	// Hard Delete
	logger.Debugf("Deleting %s Port now.", portType)
	hardDeleteStatus, deleteErr := port.DeletePort(portId, true)
	assert.True(t, hardDeleteStatus)
	assert.NoError(t, deleteErr)

	portInfo, _ := port.GetPortDetails(portId)
	assert.EqualValues(t, types.STATUS_DECOMMISSIONED, portInfo.ProvisioningStatus)
	logger.Debugf("Status is: %s", portInfo.ProvisioningStatus)
}

func testLockPort(port *Port, portId string, t *testing.T) {
	logger.Debug("Locking Port now.")
	lockStatus, lockErr := port.LockPort(portId)
	assert.True(t, lockStatus)
	assert.NoError(t, lockErr)

	portInfo, _ := port.GetPortDetails(portId)
	assert.EqualValues(t, true, portInfo.Locked)

	logger.Debug("Test lock of an already locked port.")
	lockStatus, lockErr = port.LockPort(portId)
	assert.True(t, lockStatus)
	assert.Error(t, errors.New(mega_err.ERR_PORT_ALREADY_LOCKED), lockErr)

	logger.Debug("Unlocking Port now.")
	unlockStatus, unlockErr := port.UnlockPort(portId)
	assert.True(t, unlockStatus)
	assert.NoError(t, unlockErr)

	logger.Debug("Test unlocking of a port that doesn't have a lock.")
	unlockStatus, unlockErr = port.UnlockPort(portId)
	assert.True(t, unlockStatus)
	assert.Error(t, errors.New(mega_err.ERR_PORT_NOT_LOCKED), unlockErr)
}
