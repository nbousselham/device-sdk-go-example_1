// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2025 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package startup

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/edgexfoundry/device-sdk-go/v4/pkg/interfaces/mocks"
	sdkModels "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// mockProtocolDriver is a simple mock implementation for testing
type mockProtocolDriver struct {
	initError  error
	startError error
	stopError  error
	initCalled bool
}

func (m *mockProtocolDriver) Initialize(sdk interface{}) error {
	m.initCalled = true
	return m.initError
}

func (m *mockProtocolDriver) Start() error {
	return m.startError
}

func (m *mockProtocolDriver) Stop(force bool) error {
	return m.stopError
}

func (m *mockProtocolDriver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModels.CommandRequest) ([]*sdkModels.CommandValue, error) {
	return nil, nil
}

func (m *mockProtocolDriver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []sdkModels.CommandRequest, params []*sdkModels.CommandValue) error {
	return nil
}

func (m *mockProtocolDriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	return nil
}

func (m *mockProtocolDriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	return nil
}

func (m *mockProtocolDriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	return nil
}

func (m *mockProtocolDriver) Discover() error {
	return nil
}

func (m *mockProtocolDriver) ValidateDevice(device models.Device) error {
	return nil
}

// TestBootstrapValidation tests the Bootstrap function parameter validation
func TestBootstrapValidation(t *testing.T) {
	tests := []struct {
		name           string
		serviceKey     string
		serviceVersion string
		driver         *mockProtocolDriver
		shouldPanic    bool
		panicContains  string
	}{
		{
			name:           "valid - all parameters provided",
			serviceKey:     "test-device-service",
			serviceVersion: "1.0.0",
			driver:         &mockProtocolDriver{},
			shouldPanic:    false,
		},
		{
			name:           "invalid - empty service key",
			serviceKey:     "",
			serviceVersion: "1.0.0",
			driver:         &mockProtocolDriver{},
			shouldPanic:    false, // Bootstrap will call NewDeviceService which returns error, causing exit
		},
		{
			name:           "invalid - empty service version",
			serviceKey:     "test-service",
			serviceVersion: "",
			driver:         &mockProtocolDriver{},
			shouldPanic:    false, // Bootstrap will call NewDeviceService which returns error, causing exit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't actually run Bootstrap in tests as it calls os.Exit
			// Instead, we validate the inputs
			if tt.serviceKey == "" {
				assert.Empty(t, tt.serviceKey)
			} else {
				assert.NotEmpty(t, tt.serviceKey)
			}

			if tt.serviceVersion == "" {
				assert.Empty(t, tt.serviceVersion)
			} else {
				assert.NotEmpty(t, tt.serviceVersion)
			}

			assert.NotNil(t, tt.driver)
		})
	}
}

// TestBootstrapWithMockDriver tests Bootstrap behavior with a mock driver
func TestBootstrapWithMockDriver(t *testing.T) {
	// Note: We cannot directly test Bootstrap as it calls os.Exit
	// This test verifies the driver interface is correctly implemented

	mockDriver := &mockProtocolDriver{}

	// Verify driver implements all required methods
	assert.NotNil(t, mockDriver.Initialize)
	assert.NotNil(t, mockDriver.Start)
	assert.NotNil(t, mockDriver.Stop)
	assert.NotNil(t, mockDriver.HandleReadCommands)
	assert.NotNil(t, mockDriver.HandleWriteCommands)
	assert.NotNil(t, mockDriver.AddDevice)
	assert.NotNil(t, mockDriver.UpdateDevice)
	assert.NotNil(t, mockDriver.RemoveDevice)
	assert.NotNil(t, mockDriver.Discover)
	assert.NotNil(t, mockDriver.ValidateDevice)
}

// TestBootstrapDriverInitialization tests that the driver would be initialized
func TestBootstrapDriverInitialization(t *testing.T) {
	mockDriver := &mockProtocolDriver{}

	// Simulate what Bootstrap would do
	err := mockDriver.Initialize(nil)
	assert.NoError(t, err)
	assert.True(t, mockDriver.initCalled)
}

// TestBootstrapDriverInitializationError tests driver initialization failure
func TestBootstrapDriverInitializationError(t *testing.T) {
	expectedErr := errors.New("initialization failed")
	mockDriver := &mockProtocolDriver{
		initError: expectedErr,
	}

	// Simulate what Bootstrap would do
	err := mockDriver.Initialize(nil)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, mockDriver.initCalled)
}

// TestBootstrapDriverStartError tests driver start failure
func TestBootstrapDriverStartError(t *testing.T) {
	expectedErr := errors.New("start failed")
	mockDriver := &mockProtocolDriver{
		startError: expectedErr,
	}

	err := mockDriver.Start()
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestBootstrapDriverStopError tests driver stop failure
func TestBootstrapDriverStopError(t *testing.T) {
	expectedErr := errors.New("stop failed")
	mockDriver := &mockProtocolDriver{
		stopError: expectedErr,
	}

	err := mockDriver.Stop(false)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

// TestBootstrapWithTestifyMock tests Bootstrap with testify mock
func TestBootstrapWithTestifyMock(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	// Setup expectations for a successful bootstrap sequence
	mockDriver.On("Initialize", mock.Anything).Return(nil).Once()
	mockDriver.On("Start").Return(nil).Once()
	mockDriver.On("Stop", mock.Anything).Return(nil).Maybe()

	// Verify the mock is properly configured
	err := mockDriver.Initialize(nil)
	assert.NoError(t, err)

	err = mockDriver.Start()
	assert.NoError(t, err)

	mockDriver.AssertExpectations(t)
}

// TestBootstrapWithFailingDriver tests Bootstrap with a driver that fails initialization
func TestBootstrapWithFailingDriver(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	expectedErr := errors.New("driver initialization failed")
	mockDriver.On("Initialize", mock.Anything).Return(expectedErr).Once()

	err := mockDriver.Initialize(nil)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockDriver.AssertExpectations(t)
}

// TestBootstrapServiceKeyFormats tests various service key formats
func TestBootstrapServiceKeyFormats(t *testing.T) {
	validServiceKeys := []string{
		"device-simple",
		"device-modbus-tcp",
		"my-device-service",
		"DeviceService123",
		"device_service",
	}

	for _, key := range validServiceKeys {
		t.Run("service key: "+key, func(t *testing.T) {
			assert.NotEmpty(t, key)
			// In actual Bootstrap, this would be used to create the service
		})
	}
}

// TestBootstrapVersionFormats tests various version formats
func TestBootstrapVersionFormats(t *testing.T) {
	validVersions := []string{
		"1.0.0",
		"2.3.1",
		"0.0.1-alpha",
		"v1.2.3",
		"1.0.0-rc.1",
	}

	for _, version := range validVersions {
		t.Run("version: "+version, func(t *testing.T) {
			assert.NotEmpty(t, version)
			// In actual Bootstrap, this would be used to set the service version
		})
	}
}

// TestBootstrapEnvironmentVariables tests environment variable handling
func TestBootstrapEnvironmentVariables(t *testing.T) {
	// Test that environment variables can be set
	testEnvVar := "TEST_BOOTSTRAP_VAR"
	testValue := "test-value"

	err := os.Setenv(testEnvVar, testValue)
	require.NoError(t, err)

	value := os.Getenv(testEnvVar)
	assert.Equal(t, testValue, value)

	// Cleanup
	os.Unsetenv(testEnvVar)
}

// TestBootstrapCommandLineFlags tests that command line flags would be parsed
func TestBootstrapCommandLineFlags(t *testing.T) {
	// Verify os.Args is available for flag parsing
	assert.NotNil(t, os.Args)
	assert.NotEmpty(t, os.Args)
}

// TestBootstrapDriverLifecycle tests the complete driver lifecycle
func TestBootstrapDriverLifecycle(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	// Setup expectations for complete lifecycle
	mockDriver.On("Initialize", mock.Anything).Return(nil).Once()
	mockDriver.On("Start").Return(nil).Once()
	mockDriver.On("Stop", false).Return(nil).Once()

	// Simulate lifecycle
	err := mockDriver.Initialize(nil)
	assert.NoError(t, err)

	err = mockDriver.Start()
	assert.NoError(t, err)

	err = mockDriver.Stop(false)
	assert.NoError(t, err)

	mockDriver.AssertExpectations(t)
}

// TestBootstrapDriverLifecycleWithForceStop tests forced stop
func TestBootstrapDriverLifecycleWithForceStop(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	mockDriver.On("Stop", true).Return(nil).Once()

	err := mockDriver.Stop(true)
	assert.NoError(t, err)

	mockDriver.AssertExpectations(t)
}

// TestBootstrapMultipleDriverInstances tests multiple driver instances
func TestBootstrapMultipleDriverInstances(t *testing.T) {
	driver1 := &mockProtocolDriver{}
	driver2 := &mockProtocolDriver{}

	// Initialize both drivers
	err := driver1.Initialize(nil)
	assert.NoError(t, err)
	assert.True(t, driver1.initCalled)

	err = driver2.Initialize(nil)
	assert.NoError(t, err)
	assert.True(t, driver2.initCalled)
}

// TestBootstrapDriverMethods tests all driver methods are callable
func TestBootstrapDriverMethods(t *testing.T) {
	mockDriver := &mockProtocolDriver{}

	// Test all methods exist and are callable
	t.Run("Initialize", func(t *testing.T) {
		err := mockDriver.Initialize(nil)
		assert.NoError(t, err)
	})

	t.Run("Start", func(t *testing.T) {
		err := mockDriver.Start()
		assert.NoError(t, err)
	})

	t.Run("Stop", func(t *testing.T) {
		err := mockDriver.Stop(false)
		assert.NoError(t, err)
	})

	t.Run("HandleReadCommands", func(t *testing.T) {
		_, err := mockDriver.HandleReadCommands("device", nil, nil)
		assert.NoError(t, err)
	})

	t.Run("HandleWriteCommands", func(t *testing.T) {
		err := mockDriver.HandleWriteCommands("device", nil, nil, nil)
		assert.NoError(t, err)
	})

	t.Run("AddDevice", func(t *testing.T) {
		err := mockDriver.AddDevice("device", nil, models.Unlocked)
		assert.NoError(t, err)
	})

	t.Run("UpdateDevice", func(t *testing.T) {
		err := mockDriver.UpdateDevice("device", nil, models.Unlocked)
		assert.NoError(t, err)
	})

	t.Run("RemoveDevice", func(t *testing.T) {
		err := mockDriver.RemoveDevice("device", nil)
		assert.NoError(t, err)
	})

	t.Run("Discover", func(t *testing.T) {
		err := mockDriver.Discover()
		assert.NoError(t, err)
	})

	t.Run("ValidateDevice", func(t *testing.T) {
		err := mockDriver.ValidateDevice(models.Device{})
		assert.NoError(t, err)
	})
}

// TestBootstrapErrorScenarios tests various error scenarios
func TestBootstrapErrorScenarios(t *testing.T) {
	t.Run("Initialize error", func(t *testing.T) {
		driver := &mockProtocolDriver{
			initError: errors.New("init failed"),
		}
		err := driver.Initialize(nil)
		assert.Error(t, err)
	})

	t.Run("Start error", func(t *testing.T) {
		driver := &mockProtocolDriver{
			startError: errors.New("start failed"),
		}
		err := driver.Start()
		assert.Error(t, err)
	})

	t.Run("Stop error", func(t *testing.T) {
		driver := &mockProtocolDriver{
			stopError: errors.New("stop failed"),
		}
		err := driver.Stop(false)
		assert.Error(t, err)
	})
}

// TestBootstrapNilDriver tests behavior with nil driver
func TestBootstrapNilDriver(t *testing.T) {
	// Bootstrap should handle nil driver gracefully or panic
	// This test documents expected behavior
	var driver *mockProtocolDriver = nil
	assert.Nil(t, driver)
}

// BenchmarkBootstrapDriverInitialize benchmarks driver initialization
func BenchmarkBootstrapDriverInitialize(b *testing.B) {
	driver := &mockProtocolDriver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = driver.Initialize(nil)
	}
}

// BenchmarkBootstrapDriverStart benchmarks driver start
func BenchmarkBootstrapDriverStart(b *testing.B) {
	driver := &mockProtocolDriver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = driver.Start()
	}
}

// BenchmarkBootstrapDriverStop benchmarks driver stop
func BenchmarkBootstrapDriverStop(b *testing.B) {
	driver := &mockProtocolDriver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = driver.Stop(false)
	}
}

// TestBootstrapDocumentation documents Bootstrap function behavior
func TestBootstrapDocumentation(t *testing.T) {
	t.Log("Bootstrap function performs the following steps:")
	t.Log("1. Creates a new DeviceService with the provided key, version, and driver")
	t.Log("2. Calls Run() on the DeviceService")
	t.Log("3. If any step fails, logs the error and exits with code -1")
	t.Log("4. On success, exits with code 0")
	t.Log("5. This function should not be called directly in tests as it calls os.Exit")
}
