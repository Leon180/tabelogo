package inject

// ControllerHandle Controller Handle
type ControllerHandle struct {
	// HealthCheckControllerHandle controller.HealthCheckControllerHandle
	// TestControllerHandle        controller.TestControllerHandle
}

func newControllerHandle(
// healthCheckControllerHandle controller.HealthCheckControllerHandle,
// testControllerHandle controller.TestControllerHandle,
) *ControllerHandle {
	return &ControllerHandle{
		// HealthCheckControllerHandle: healthCheckControllerHandle,
		// TestControllerHandle:        testControllerHandle,
	}
}
