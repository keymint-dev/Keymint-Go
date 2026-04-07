package keymint

import (
	"testing"
)

func TestIdentity(t *testing.T) {
	machineID := GetMachineID()
	if machineID == "" {
		t.Log("GetMachineID returned empty string (expected on some systems but not ideal)")
	} else {
		t.Logf("Machine ID: %s", machineID)
		if len(machineID) != 64 {
			t.Errorf("Expected 64 character hex string, got %d chars", len(machineID))
		}
	}

	installID, err := GetOrCreateInstallationID("")
	if err != nil {
		t.Fatalf("Failed to generate installation ID: %v", err)
	}
	t.Logf("Installation ID: %s", installID)
	
	if len(installID) != 64 {
		t.Errorf("Expected 64 character hex string, got %d chars", len(installID))
	}
	
	// Test persistence
	installID2, err := GetOrCreateInstallationID("")
	if err != nil {
		t.Fatalf("Failed to second installation ID: %v", err)
	}
	if installID != installID2 {
		t.Errorf("Persistence failed: %s != %s", installID, installID2)
	} else {
		t.Log("Persistence test passed.")
	}
}
