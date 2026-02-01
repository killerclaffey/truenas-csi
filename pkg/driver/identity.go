package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IdentityServer implements the CSI Identity service.
type IdentityServer struct {
	driver *Driver
	csi.UnimplementedIdentityServer
}

// NewIdentityServer creates a new CSI identity service.
func NewIdentityServer(driver *Driver) *IdentityServer {
	return &IdentityServer{
		driver: driver,
	}
}

// GetPluginInfo returns the driver name and version.
func (s *IdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	s.driver.Log().V(LogLevelDebug).Info("GetPluginInfo called")

	if s.driver.name == "" {
		return nil, status.Error(codes.Unavailable, "driver name not configured")
	}

	if s.driver.version == "" {
		return nil, status.Error(codes.Unavailable, "driver version not configured")
	}

	return &csi.GetPluginInfoResponse{
		Name:          s.driver.name,
		VendorVersion: s.driver.version,
	}, nil
}

// GetPluginCapabilities returns the driver's capabilities.
func (s *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	s.driver.Log().V(LogLevelDebug).Info("GetPluginCapabilities called")

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: s.driver.pluginCaps,
	}, nil
}

// Probe checks if the driver is healthy by verifying the TrueNAS connection status.
// We only check if the client is connected rather than actively pinging TrueNAS
// on every probe call. The client has its own background ping loop that handles
// connection health monitoring and automatic reconnection.
func (s *IdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	s.driver.Log().V(LogLevelDebug).Info("Probe called")

	// Check if the client is connected - don't actively ping on every probe
	// as this can cause timeouts during high load or network latency
	if !s.driver.client.Connected() {
		s.driver.Log().Error(nil, "Health check failed", "reason", "not connected to TrueNAS")
		return nil, status.Error(codes.FailedPrecondition, "TrueNAS connection failed")
	}

	return &csi.ProbeResponse{}, nil
}
