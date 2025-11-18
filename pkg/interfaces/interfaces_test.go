// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2025 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package interfaces_test

import (
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/edgexfoundry/device-sdk-go/v4/pkg/interfaces"
	"github.com/edgexfoundry/device-sdk-go/v4/pkg/interfaces/mocks"
	sdkModels "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
)

// TestProtocolDriverInterface tests the ProtocolDriver interface contract
func TestProtocolDriverInterface(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}
	mockSDK := &mocks.DeviceServiceSDK{}

	t.Run("Initialize", func(t *testing.T) {
		mockDriver.On("Initialize", mockSDK).Return(nil).Once()
		err := mockDriver.Initialize(mockSDK)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Start", func(t *testing.T) {
		mockDriver.On("Start").Return(nil).Once()
		err := mockDriver.Start()
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Stop", func(t *testing.T) {
		mockDriver.On("Stop", false).Return(nil).Once()
		err := mockDriver.Stop(false)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("HandleReadCommands", func(t *testing.T) {
		deviceName := "test-device"
		protocols := map[string]models.ProtocolProperties{
			"tcp": {"host": "localhost", "port": "502"},
		}
		reqs := []sdkModels.CommandRequest{
			{DeviceResourceName: "temperature"},
		}
		expectedValues := []*sdkModels.CommandValue{
			{DeviceResourceName: "temperature", Type: "Float64"},
		}

		mockDriver.On("HandleReadCommands", deviceName, protocols, reqs).Return(expectedValues, nil).Once()
		values, err := mockDriver.HandleReadCommands(deviceName, protocols, reqs)
		assert.NoError(t, err)
		assert.Equal(t, expectedValues, values)
		mockDriver.AssertExpectations(t)
	})

	t.Run("HandleWriteCommands", func(t *testing.T) {
		deviceName := "test-device"
		protocols := map[string]models.ProtocolProperties{
			"tcp": {"host": "localhost", "port": "502"},
		}
		reqs := []sdkModels.CommandRequest{
			{DeviceResourceName: "setpoint"},
		}
		params := []*sdkModels.CommandValue{
			{DeviceResourceName: "setpoint", Type: "Float64", Value: float64(25.5)},
		}

		mockDriver.On("HandleWriteCommands", deviceName, protocols, reqs, params).Return(nil).Once()
		err := mockDriver.HandleWriteCommands(deviceName, protocols, reqs, params)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("AddDevice", func(t *testing.T) {
		deviceName := "new-device"
		protocols := map[string]models.ProtocolProperties{
			"tcp": {"host": "192.168.1.100", "port": "502"},
		}
		adminState := models.Unlocked

		mockDriver.On("AddDevice", deviceName, protocols, adminState).Return(nil).Once()
		err := mockDriver.AddDevice(deviceName, protocols, adminState)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("UpdateDevice", func(t *testing.T) {
		deviceName := "test-device"
		protocols := map[string]models.ProtocolProperties{
			"tcp": {"host": "192.168.1.101", "port": "502"},
		}
		adminState := models.Locked

		mockDriver.On("UpdateDevice", deviceName, protocols, adminState).Return(nil).Once()
		err := mockDriver.UpdateDevice(deviceName, protocols, adminState)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("RemoveDevice", func(t *testing.T) {
		deviceName := "test-device"
		protocols := map[string]models.ProtocolProperties{
			"tcp": {"host": "localhost", "port": "502"},
		}

		mockDriver.On("RemoveDevice", deviceName, protocols).Return(nil).Once()
		err := mockDriver.RemoveDevice(deviceName, protocols)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Discover", func(t *testing.T) {
		mockDriver.On("Discover").Return(nil).Once()
		err := mockDriver.Discover()
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("ValidateDevice", func(t *testing.T) {
		device := models.Device{
			Name: "test-device",
			Protocols: map[string]models.ProtocolProperties{
				"tcp": {"host": "localhost", "port": "502"},
			},
		}

		mockDriver.On("ValidateDevice", device).Return(nil).Once()
		err := mockDriver.ValidateDevice(device)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
	})
}

// TestDeviceServiceSDKInterface tests the DeviceServiceSDK interface contract
func TestDeviceServiceSDKInterface(t *testing.T) {
	mockSDK := &mocks.DeviceServiceSDK{}

	t.Run("Name", func(t *testing.T) {
		expectedName := "test-device-service"
		mockSDK.On("Name").Return(expectedName).Once()
		name := mockSDK.Name()
		assert.Equal(t, expectedName, name)
		mockSDK.AssertExpectations(t)
	})

	t.Run("Version", func(t *testing.T) {
		expectedVersion := "1.0.0"
		mockSDK.On("Version").Return(expectedVersion).Once()
		version := mockSDK.Version()
		assert.Equal(t, expectedVersion, version)
		mockSDK.AssertExpectations(t)
	})

	t.Run("AsyncReadingsEnabled", func(t *testing.T) {
		mockSDK.On("AsyncReadingsEnabled").Return(true).Once()
		enabled := mockSDK.AsyncReadingsEnabled()
		assert.True(t, enabled)
		mockSDK.AssertExpectations(t)
	})

	t.Run("DeviceDiscoveryEnabled", func(t *testing.T) {
		mockSDK.On("DeviceDiscoveryEnabled").Return(true).Once()
		enabled := mockSDK.DeviceDiscoveryEnabled()
		assert.True(t, enabled)
		mockSDK.AssertExpectations(t)
	})

	t.Run("AddDevice", func(t *testing.T) {
		device := models.Device{Name: "test-device"}
		expectedID := "device-123"
		mockSDK.On("AddDevice", device).Return(expectedID, nil).Once()
		id, err := mockSDK.AddDevice(device)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockSDK.AssertExpectations(t)
	})

	t.Run("GetDeviceByName", func(t *testing.T) {
		deviceName := "test-device"
		expectedDevice := models.Device{Name: deviceName}
		mockSDK.On("GetDeviceByName", deviceName).Return(expectedDevice, nil).Once()
		device, err := mockSDK.GetDeviceByName(deviceName)
		assert.NoError(t, err)
		assert.Equal(t, expectedDevice, device)
		mockSDK.AssertExpectations(t)
	})

	t.Run("UpdateDevice", func(t *testing.T) {
		device := models.Device{Name: "test-device"}
		mockSDK.On("UpdateDevice", device).Return(nil).Once()
		err := mockSDK.UpdateDevice(device)
		assert.NoError(t, err)
		mockSDK.AssertExpectations(t)
	})

	t.Run("RemoveDeviceByName", func(t *testing.T) {
		deviceName := "test-device"
		mockSDK.On("RemoveDeviceByName", deviceName).Return(nil).Once()
		err := mockSDK.RemoveDeviceByName(deviceName)
		assert.NoError(t, err)
		mockSDK.AssertExpectations(t)
	})

	t.Run("DeviceExistsForName", func(t *testing.T) {
		deviceName := "test-device"
		mockSDK.On("DeviceExistsForName", deviceName).Return(true).Once()
		exists := mockSDK.DeviceExistsForName(deviceName)
		assert.True(t, exists)
		mockSDK.AssertExpectations(t)
	})

	t.Run("DriverConfigs", func(t *testing.T) {
		expectedConfigs := map[string]string{"key1": "value1", "key2": "value2"}
		mockSDK.On("DriverConfigs").Return(expectedConfigs).Once()
		configs := mockSDK.DriverConfigs()
		assert.Equal(t, expectedConfigs, configs)
		mockSDK.AssertExpectations(t)
	})
}

// TestAutoEventManagerInterface tests the AutoEventManager interface contract
func TestAutoEventManagerInterface(t *testing.T) {
	mockManager := &mocks.AutoEventManager{}

	t.Run("StartAutoEvents", func(t *testing.T) {
		mockManager.On("StartAutoEvents").Once()
		mockManager.StartAutoEvents()
		mockManager.AssertExpectations(t)
	})

	t.Run("RestartForDevice", func(t *testing.T) {
		deviceName := "test-device"
		mockManager.On("RestartForDevice", deviceName).Once()
		mockManager.RestartForDevice(deviceName)
		mockManager.AssertExpectations(t)
	})

	t.Run("StopForDevice", func(t *testing.T) {
		deviceName := "test-device"
		mockManager.On("StopForDevice", deviceName).Once()
		mockManager.StopForDevice(deviceName)
		mockManager.AssertExpectations(t)
	})
}

// TestAuthenticationConstants tests the Authentication type constants
func TestAuthenticationConstants(t *testing.T) {
	t.Run("Unauthenticated constant", func(t *testing.T) {
		assert.Equal(t, interfaces.Authentication(false), interfaces.Unauthenticated)
	})

	t.Run("Authenticated constant", func(t *testing.T) {
		assert.Equal(t, interfaces.Authentication(true), interfaces.Authenticated)
	})
}

// TestUpdatableConfigInterface tests the UpdatableConfig interface
func TestUpdatableConfigInterface(t *testing.T) {
	mockConfig := &mocks.UpdatableConfig{}

	t.Run("UpdateFromRaw", func(t *testing.T) {
		rawConfig := make(map[string]any)
		mockConfig.On("UpdateFromRaw", rawConfig).Return(true).Once()
		updated := mockConfig.UpdateFromRaw(rawConfig)
		assert.True(t, updated)
		mockConfig.AssertExpectations(t)
	})

	t.Run("UpdateWritableFromRaw", func(t *testing.T) {
		rawWritable := make(map[string]any)
		mockConfig.On("UpdateWritableFromRaw", rawWritable).Return(true).Once()
		updated := mockConfig.UpdateWritableFromRaw(rawWritable)
		assert.True(t, updated)
		mockConfig.AssertExpectations(t)
	})
}

// TestMockInteractions tests interactions between mocked interfaces
func TestMockInteractions(t *testing.T) {
	t.Run("Driver initialization with SDK", func(t *testing.T) {
		mockDriver := &mocks.ProtocolDriver{}
		mockSDK := &mocks.DeviceServiceSDK{}

		// Setup expectations
		mockSDK.On("Name").Return("test-service")
		mockSDK.On("Version").Return("1.0.0")
		mockDriver.On("Initialize", mockSDK).Return(nil)

		// Execute
		serviceName := mockSDK.Name()
		serviceVersion := mockSDK.Version()
		err := mockDriver.Initialize(mockSDK)

		// Verify
		assert.Equal(t, "test-service", serviceName)
		assert.Equal(t, "1.0.0", serviceVersion)
		assert.NoError(t, err)
		mockDriver.AssertExpectations(t)
		mockSDK.AssertExpectations(t)
	})

	t.Run("Device lifecycle management", func(t *testing.T) {
		mockDriver := &mocks.ProtocolDriver{}
		mockSDK := &mocks.DeviceServiceSDK{}

		deviceName := "lifecycle-device"
		protocols := map[string]models.ProtocolProperties{
			"modbus": {"address": "127.0.0.1", "port": "502"},
		}
		device := models.Device{
			Name:      deviceName,
			Protocols: protocols,
		}

		// Add device
		mockSDK.On("AddDevice", device).Return("device-id-123", nil).Once()
		mockDriver.On("AddDevice", deviceName, protocols, models.Unlocked).Return(nil).Once()

		id, err := mockSDK.AddDevice(device)
		assert.NoError(t, err)
		assert.Equal(t, "device-id-123", id)

		err = mockDriver.AddDevice(deviceName, protocols, models.Unlocked)
		assert.NoError(t, err)

		// Update device
		mockSDK.On("UpdateDevice", device).Return(nil).Once()
		mockDriver.On("UpdateDevice", deviceName, protocols, models.Locked).Return(nil).Once()

		err = mockSDK.UpdateDevice(device)
		assert.NoError(t, err)

		err = mockDriver.UpdateDevice(deviceName, protocols, models.Locked)
		assert.NoError(t, err)

		// Remove device
		mockSDK.On("RemoveDeviceByName", deviceName).Return(nil).Once()
		mockDriver.On("RemoveDevice", deviceName, protocols).Return(nil).Once()

		err = mockSDK.RemoveDeviceByName(deviceName)
		assert.NoError(t, err)

		err = mockDriver.RemoveDevice(deviceName, protocols)
		assert.NoError(t, err)

		mockDriver.AssertExpectations(t)
		mockSDK.AssertExpectations(t)
	})
}

// TestErrorHandling tests error scenarios with mocked interfaces
func TestErrorHandling(t *testing.T) {
	t.Run("Driver initialization failure", func(t *testing.T) {
		mockDriver := &mocks.ProtocolDriver{}
		mockSDK := &mocks.DeviceServiceSDK{}

		expectedErr := assert.AnError
		mockDriver.On("Initialize", mockSDK).Return(expectedErr).Once()

		err := mockDriver.Initialize(mockSDK)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Device not found", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		deviceName := "non-existent-device"
		expectedErr := assert.AnError

		mockSDK.On("GetDeviceByName", deviceName).Return(models.Device{}, expectedErr).Once()

		device, err := mockSDK.GetDeviceByName(deviceName)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, device.Name)
		mockSDK.AssertExpectations(t)
	})

	t.Run("HandleReadCommands failure", func(t *testing.T) {
		mockDriver := &mocks.ProtocolDriver{}

		deviceName := "failing-device"
		protocols := map[string]models.ProtocolProperties{}
		reqs := []sdkModels.CommandRequest{{DeviceResourceName: "sensor"}}
		expectedErr := assert.AnError

		mockDriver.On("HandleReadCommands", deviceName, protocols, reqs).Return([]*sdkModels.CommandValue(nil), expectedErr).Once()

		values, err := mockDriver.HandleReadCommands(deviceName, protocols, reqs)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, values)
		mockDriver.AssertExpectations(t)
	})
}

// TestChannelOperations tests channel-based operations
func TestChannelOperations(t *testing.T) {
	t.Run("AsyncValuesChannel", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}
		asyncCh := make(chan *sdkModels.AsyncValues, 1)

		mockSDK.On("AsyncValuesChannel").Return(asyncCh).Once()

		ch := mockSDK.AsyncValuesChannel()
		assert.NotNil(t, ch)
		assert.Equal(t, asyncCh, ch)
		mockSDK.AssertExpectations(t)
	})

	t.Run("DiscoveredDeviceChannel", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}
		deviceCh := make(chan []sdkModels.DiscoveredDevice, 1)

		mockSDK.On("DiscoveredDeviceChannel").Return(deviceCh).Once()

		ch := mockSDK.DiscoveredDeviceChannel()
		assert.NotNil(t, ch)
		assert.Equal(t, deviceCh, ch)
		mockSDK.AssertExpectations(t)
	})
}

// TestDeviceResourceAndCommand tests device resource and command retrieval
func TestDeviceResourceAndCommand(t *testing.T) {
	t.Run("DeviceResource found", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		deviceName := "test-device"
		resourceName := "temperature"
		expectedResource := models.DeviceResource{Name: resourceName}

		mockSDK.On("DeviceResource", deviceName, resourceName).Return(expectedResource, true).Once()

		resource, found := mockSDK.DeviceResource(deviceName, resourceName)
		assert.True(t, found)
		assert.Equal(t, expectedResource, resource)
		mockSDK.AssertExpectations(t)
	})

	t.Run("DeviceResource not found", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		deviceName := "test-device"
		resourceName := "non-existent"

		mockSDK.On("DeviceResource", deviceName, resourceName).Return(models.DeviceResource{}, false).Once()

		resource, found := mockSDK.DeviceResource(deviceName, resourceName)
		assert.False(t, found)
		assert.Empty(t, resource.Name)
		mockSDK.AssertExpectations(t)
	})

	t.Run("DeviceCommand found", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		deviceName := "test-device"
		commandName := "read-all"
		expectedCommand := models.DeviceCommand{Name: commandName}

		mockSDK.On("DeviceCommand", deviceName, commandName).Return(expectedCommand, true).Once()

		command, found := mockSDK.DeviceCommand(deviceName, commandName)
		assert.True(t, found)
		assert.Equal(t, expectedCommand, command)
		mockSDK.AssertExpectations(t)
	})
}

// TestProfileAndWatcherManagement tests profile and watcher management operations
func TestProfileAndWatcherManagement(t *testing.T) {
	t.Run("AddDeviceProfile", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		profile := models.DeviceProfile{Name: "test-profile"}
		expectedID := "profile-123"

		mockSDK.On("AddDeviceProfile", profile).Return(expectedID, nil).Once()

		id, err := mockSDK.AddDeviceProfile(profile)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockSDK.AssertExpectations(t)
	})

	t.Run("AddProvisionWatcher", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		watcher := models.ProvisionWatcher{Name: "test-watcher"}
		expectedID := "watcher-123"

		mockSDK.On("AddProvisionWatcher", watcher).Return(expectedID, nil).Once()

		id, err := mockSDK.AddProvisionWatcher(watcher)
		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockSDK.AssertExpectations(t)
	})

	t.Run("GetProfileByName", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		profileName := "test-profile"
		expectedProfile := models.DeviceProfile{Name: profileName}

		mockSDK.On("GetProfileByName", profileName).Return(expectedProfile, nil).Once()

		profile, err := mockSDK.GetProfileByName(profileName)
		assert.NoError(t, err)
		assert.Equal(t, expectedProfile, profile)
		mockSDK.AssertExpectations(t)
	})
}

// TestAutoEventOperations tests auto-event management
func TestAutoEventOperations(t *testing.T) {
	t.Run("AddDeviceAutoEvent", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		deviceName := "test-device"
		event := models.AutoEvent{SourceName: "temperature", Interval: "30s"}

		mockSDK.On("AddDeviceAutoEvent", deviceName, event).Return(nil).Once()

		err := mockSDK.AddDeviceAutoEvent(deviceName, event)
		assert.NoError(t, err)
		mockSDK.AssertExpectations(t)
	})

	t.Run("RemoveDeviceAutoEvent", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		deviceName := "test-device"
		event := models.AutoEvent{SourceName: "temperature", Interval: "30s"}

		mockSDK.On("RemoveDeviceAutoEvent", deviceName, event).Return(nil).Once()

		err := mockSDK.RemoveDeviceAutoEvent(deviceName, event)
		assert.NoError(t, err)
		mockSDK.AssertExpectations(t)
	})
}

// TestBulkOperations tests bulk retrieval operations
func TestBulkOperations(t *testing.T) {
	t.Run("Devices", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		expectedDevices := []models.Device{
			{Name: "device1"},
			{Name: "device2"},
		}

		mockSDK.On("Devices").Return(expectedDevices).Once()

		devices := mockSDK.Devices()
		assert.Len(t, devices, 2)
		assert.Equal(t, expectedDevices, devices)
		mockSDK.AssertExpectations(t)
	})

	t.Run("DeviceProfiles", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		expectedProfiles := []models.DeviceProfile{
			{Name: "profile1"},
			{Name: "profile2"},
		}

		mockSDK.On("DeviceProfiles").Return(expectedProfiles).Once()

		profiles := mockSDK.DeviceProfiles()
		assert.Len(t, profiles, 2)
		assert.Equal(t, expectedProfiles, profiles)
		mockSDK.AssertExpectations(t)
	})

	t.Run("ProvisionWatchers", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		expectedWatchers := []models.ProvisionWatcher{
			{Name: "watcher1"},
			{Name: "watcher2"},
		}

		mockSDK.On("ProvisionWatchers").Return(expectedWatchers).Once()

		watchers := mockSDK.ProvisionWatchers()
		assert.Len(t, watchers, 2)
		assert.Equal(t, expectedWatchers, watchers)
		mockSDK.AssertExpectations(t)
	})
}

// TestMultipleCallsToSameMethod tests calling the same method multiple times
func TestMultipleCallsToSameMethod(t *testing.T) {
	mockManager := &mocks.AutoEventManager{}

	deviceName1 := "device1"
	deviceName2 := "device2"

	mockManager.On("RestartForDevice", deviceName1).Once()
	mockManager.On("RestartForDevice", deviceName2).Once()

	mockManager.RestartForDevice(deviceName1)
	mockManager.RestartForDevice(deviceName2)

	mockManager.AssertExpectations(t)
}

// TestEdgeCases tests edge case scenarios
func TestEdgeCases(t *testing.T) {
	t.Run("Empty protocols map", func(t *testing.T) {
		mockDriver := &mocks.ProtocolDriver{}

		deviceName := "test-device"
		emptyProtocols := map[string]models.ProtocolProperties{}
		reqs := []sdkModels.CommandRequest{{DeviceResourceName: "sensor"}}

		mockDriver.On("HandleReadCommands", deviceName, emptyProtocols, reqs).Return([]*sdkModels.CommandValue{}, nil).Once()

		values, err := mockDriver.HandleReadCommands(deviceName, emptyProtocols, reqs)
		assert.NoError(t, err)
		assert.Empty(t, values)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Empty command requests", func(t *testing.T) {
		mockDriver := &mocks.ProtocolDriver{}

		deviceName := "test-device"
		protocols := map[string]models.ProtocolProperties{"tcp": {"host": "localhost"}}
		emptyReqs := []sdkModels.CommandRequest{}

		mockDriver.On("HandleReadCommands", deviceName, protocols, emptyReqs).Return([]*sdkModels.CommandValue{}, nil).Once()

		values, err := mockDriver.HandleReadCommands(deviceName, protocols, emptyReqs)
		assert.NoError(t, err)
		assert.Empty(t, values)
		mockDriver.AssertExpectations(t)
	})

	t.Run("Device with no name", func(t *testing.T) {
		mockSDK := &mocks.DeviceServiceSDK{}

		device := models.Device{Name: ""}
		mockSDK.On("AddDevice", device).Return("", assert.AnError).Once()

		id, err := mockSDK.AddDevice(device)
		assert.Error(t, err)
		assert.Empty(t, id)
		mockSDK.AssertExpectations(t)
	})

	t.Run("ValidateDevice with invalid protocols", func(t *testing.T) {
		mockDriver := &mocks.ProtocolDriver{}

		invalidDevice := models.Device{
			Name:      "invalid-device",
			Protocols: map[string]models.ProtocolProperties{},
		}

		mockDriver.On("ValidateDevice", invalidDevice).Return(assert.AnError).Once()

		err := mockDriver.ValidateDevice(invalidDevice)
		assert.Error(t, err)
		mockDriver.AssertExpectations(t)
	})
}

// TestMethodWithAnyMatcher tests using mock.Anything matcher
func TestMethodWithAnyMatcher(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	// Use mock.Anything to match any argument
	mockDriver.On("AddDevice", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := mockDriver.AddDevice("any-device", map[string]models.ProtocolProperties{}, models.Unlocked)
	assert.NoError(t, err)
	mockDriver.AssertExpectations(t)
}
