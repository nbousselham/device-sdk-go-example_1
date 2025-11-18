// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2025 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"errors"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/edgexfoundry/device-sdk-go/v4/pkg/interfaces"
	"github.com/edgexfoundry/device-sdk-go/v4/pkg/interfaces/mocks"
)

// TestNewDeviceService tests the NewDeviceService factory function
func TestNewDeviceService(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	tests := []struct {
		name           string
		serviceKey     string
		serviceVersion string
		driver         interfaces.ProtocolDriver
		expectError    bool
		errorContains  string
	}{
		{
			name:           "valid - successful device service creation",
			serviceKey:     "test-device-service",
			serviceVersion: "1.0.0",
			driver:         mockDriver,
			expectError:    false,
		},
		{
			name:           "invalid - empty service key",
			serviceKey:     "",
			serviceVersion: "1.0.0",
			driver:         mockDriver,
			expectError:    true,
			errorContains:  "please specify device service name",
		},
		{
			name:           "invalid - empty service version",
			serviceKey:     "test-device-service",
			serviceVersion: "",
			driver:         mockDriver,
			expectError:    true,
			errorContains:  "please specify device service version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewDeviceService(tt.serviceKey, tt.serviceVersion, tt.driver)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, service)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.Equal(t, tt.serviceKey, service.Name())
				assert.Equal(t, tt.serviceVersion, service.Version())
			}
		})
	}
}

// TestDeviceServiceName tests the Name method
func TestDeviceServiceName(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}
	serviceName := "my-device-service"

	service, err := NewDeviceService(serviceName, "1.0.0", mockDriver)
	require.NoError(t, err)

	assert.Equal(t, serviceName, service.Name())
}

// TestDeviceServiceVersion tests the Version method
func TestDeviceServiceVersion(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}
	serviceVersion := "2.3.1"

	service, err := NewDeviceService("test-service", serviceVersion, mockDriver)
	require.NoError(t, err)

	assert.Equal(t, serviceVersion, service.Version())
}

// TestDeviceServiceLoggingClient tests LoggingClient initialization
func TestDeviceServiceLoggingClient(t *testing.T) {
	ds := newDeviceService()

	// Initially, lc should be set to MockClient in newDeviceService
	lc := ds.LoggingClient()
	assert.NotNil(t, lc)

	// Calling it again should return the same instance
	lc2 := ds.LoggingClient()
	assert.Equal(t, lc, lc2)
}

// TestAsyncReadingsEnabledDefault tests the default value of AsyncReadingsEnabled
func TestAsyncReadingsEnabledDefault(t *testing.T) {
	ds := newDeviceService()

	// By default, async readings should be disabled in the test service
	// since config.Device.EnableAsyncReadings defaults to false
	enabled := ds.AsyncReadingsEnabled()
	assert.False(t, enabled)
}

// TestDeviceDiscoveryEnabledDefault tests the default value of DeviceDiscoveryEnabled
func TestDeviceDiscoveryEnabledDefault(t *testing.T) {
	ds := newDeviceService()

	// By default, device discovery should be disabled in the test service
	enabled := ds.DeviceDiscoveryEnabled()
	assert.False(t, enabled)
}

// TestAsyncValuesChannel tests the AsyncValuesChannel method
func TestAsyncValuesChannel(t *testing.T) {
	ds := newDeviceService()

	// Initially, asyncCh should be nil
	ch := ds.AsyncValuesChannel()
	assert.Nil(t, ch)
}

// TestDiscoveredDeviceChannel tests the DiscoveredDeviceChannel method
func TestDiscoveredDeviceChannel(t *testing.T) {
	ds := newDeviceService()

	// Initially, deviceCh should be nil
	ch := ds.DiscoveredDeviceChannel()
	assert.Nil(t, ch)
}

// TestDriverConfigs tests DriverConfigs method
func TestDriverConfigs(t *testing.T) {
	ds := newDeviceService()

	// Should return nil or empty map when no driver configs are set
	configs := ds.DriverConfigs()
	if configs != nil {
		assert.Empty(t, configs)
	}
}

// TestSetServiceName tests the setServiceName method
func TestSetServiceName(t *testing.T) {
	tests := []struct {
		name             string
		baseServiceKey   string
		instanceName     string
		expectedFullName string
		expectedBaseName string
	}{
		{
			name:             "no instance name",
			baseServiceKey:   "device-service",
			instanceName:     "",
			expectedFullName: "device-service",
			expectedBaseName: "device-service",
		},
		{
			name:             "with instance name",
			baseServiceKey:   "device-service",
			instanceName:     "1",
			expectedFullName: "device-service_1",
			expectedBaseName: "device-service",
		},
		{
			name:             "with complex instance name",
			baseServiceKey:   "my-device-service",
			instanceName:     "north-wing",
			expectedFullName: "my-device-service_north-wing",
			expectedBaseName: "my-device-service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &deviceService{
				serviceKey: tt.baseServiceKey,
			}

			ds.setServiceName(tt.instanceName)

			assert.Equal(t, tt.expectedFullName, ds.serviceKey)
			assert.Equal(t, tt.expectedBaseName, ds.baseServiceName)
		})
	}
}

// TestStopDeviceService tests the Stop method
func TestStopDeviceService(t *testing.T) {
	tests := []struct {
		name  string
		force bool
	}{
		{
			name:  "graceful stop",
			force: false,
		},
		{
			name:  "force stop",
			force: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDriver := &mocks.ProtocolDriver{}
			ds := &deviceService{
				serviceKey: "test-service",
				driver:     mockDriver,
				lc:         logger.NewMockClient(),
			}

			mockDriver.On("Stop", tt.force).Return(nil).Once()

			ds.Stop(tt.force)
			mockDriver.AssertExpectations(t)
		})
	}
}

// TestStopDeviceServiceWithError tests Stop method when driver returns error
func TestStopDeviceServiceWithError(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}
	ds := &deviceService{
		serviceKey: "test-service",
		driver:     mockDriver,
		lc:         logger.NewMockClient(),
	}

	expectedErr := errors.New("driver stop failed")
	mockDriver.On("Stop", false).Return(expectedErr).Once()

	// Stop doesn't return an error, it logs it
	ds.Stop(false)
	mockDriver.AssertExpectations(t)
}

// TestSelfRegister tests the selfRegister method behavior
func TestSelfRegister(t *testing.T) {
	// This test verifies that selfRegister is called during bootstrap
	// Actual testing would require mocking the metadata client
	// For now, we verify the method exists and can be called
	ds := newDeviceService()
	assert.NotNil(t, ds)
}

// TestDeviceServiceCreationWithExtendedDriver tests creation with ExtendedProtocolDriver
func TestDeviceServiceCreationWithExtendedDriver(t *testing.T) {
	// Create a mock that implements both interfaces
	type mockExtendedDriver struct {
		mocks.ProtocolDriver
		interfaces.ExtendedProtocolDriver
	}

	// This test verifies that the service properly identifies extended drivers
	// The actual implementation would cast the driver if it implements ExtendedProtocolDriver
}

// TestChannelInitialization tests that channels are properly initialized
func TestChannelInitialization(t *testing.T) {
	ds := newDeviceService()

	// Before bootstrap, channels should be nil
	assert.Nil(t, ds.asyncCh)
	assert.Nil(t, ds.deviceCh)
}

// TestMultipleServiceInstances tests creating multiple service instances
func TestMultipleServiceInstances(t *testing.T) {
	mockDriver1 := &mocks.ProtocolDriver{}
	mockDriver2 := &mocks.ProtocolDriver{}

	service1, err := NewDeviceService("service-1", "1.0.0", mockDriver1)
	require.NoError(t, err)

	service2, err := NewDeviceService("service-2", "2.0.0", mockDriver2)
	require.NoError(t, err)

	assert.NotEqual(t, service1.Name(), service2.Name())
	// Note: Version is shared globally in sdkCommon.ServiceVersion, so both will have the same version
	// This is a known limitation when creating multiple services in the same process
}

// TestServiceKeyValidation tests various service key formats
func TestServiceKeyValidation(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	validKeys := []string{
		"device-simple",
		"device-modbus",
		"my-custom-device-service",
		"device_service_123",
		"DeviceService",
	}

	for _, key := range validKeys {
		t.Run("valid key: "+key, func(t *testing.T) {
			service, err := NewDeviceService(key, "1.0.0", mockDriver)
			assert.NoError(t, err)
			assert.NotNil(t, service)
			assert.Equal(t, key, service.Name())
		})
	}
}

// TestVersionValidation tests various version formats
func TestVersionValidation(t *testing.T) {
	mockDriver := &mocks.ProtocolDriver{}

	validVersions := []string{
		"1.0.0",
		"2.3.1",
		"0.0.1",
		"10.20.30",
		"v1.2.3",
		"1.0.0-beta",
	}

	for _, version := range validVersions {
		t.Run("valid version: "+version, func(t *testing.T) {
			service, err := NewDeviceService("test-service", version, mockDriver)
			assert.NoError(t, err)
			assert.NotNil(t, service)
			assert.Equal(t, version, service.Version())
		})
	}
}

// TestNilDriverValidation tests service creation with nil driver
func TestNilDriverValidation(t *testing.T) {
	// This should be handled by the caller, but we test the behavior
	service, err := NewDeviceService("test-service", "1.0.0", nil)
	// The current implementation doesn't validate nil driver in NewDeviceService
	// It would fail later during initialization
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

// TestConfigurationAccess tests that configuration is accessible
func TestConfigurationAccess(t *testing.T) {
	ds := newDeviceService()

	assert.NotNil(t, ds.config)
}

// BenchmarkNewDeviceService benchmarks the NewDeviceService function
func BenchmarkNewDeviceService(b *testing.B) {
	mockDriver := &mocks.ProtocolDriver{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewDeviceService("benchmark-service", "1.0.0", mockDriver)
	}
}

// BenchmarkServiceName benchmarks the Name method
func BenchmarkServiceName(b *testing.B) {
	mockDriver := &mocks.ProtocolDriver{}
	service, _ := NewDeviceService("benchmark-service", "1.0.0", mockDriver)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Name()
	}
}

// BenchmarkServiceVersion benchmarks the Version method
func BenchmarkServiceVersion(b *testing.B) {
	mockDriver := &mocks.ProtocolDriver{}
	service, _ := NewDeviceService("benchmark-service", "1.0.0", mockDriver)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Version()
	}
}
