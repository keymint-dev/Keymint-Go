package keymint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ─── Garbage Detection ──────────────────────────────────────────────────

var garbageStrings = []string{
	"ffffffffffffffffffffffffffffffff",
	"03000200040005000006000700080009",
	"defaultstring",
	"tobefilledbyoem",
	"notapplicable",
	"notspecified",
	"systemserialnum",
	"none",
}

var garbageRegexes = []*regexp.Regexp{
	regexp.MustCompile(`^0+$`),
	regexp.MustCompile(`^f+$`),
}

var normalizeRegex = regexp.MustCompile(`[-:\s._]`)

func isGarbageID(id string) bool {
	normalized := strings.ToLower(normalizeRegex.ReplaceAllString(id, ""))
	for _, re := range garbageRegexes {
		if re.MatchString(normalized) {
			return true
		}
	}
	for _, garbage := range garbageStrings {
		if normalized == garbage || strings.Contains(normalized, garbage) {
			return true
		}
	}
	return false
}

func hashID(raw string) string {
	h := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(raw))))
	return hex.EncodeToString(h[:])
}

// ─── Fingerprint Layers ────────────────────────────────────────────────

// getBiosUUID attempts to read the BIOS/Hardware UUID.
func getBiosUUID() string {
	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("powershell.exe", "-Command",
			"(Get-CimInstance Win32_ComputerSystemProduct).UUID").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	case "darwin":
		out, err := exec.Command("bash", "-c",
			"ioreg -rd1 -c IOPlatformExpertDevice | grep IOPlatformUUID | awk -F'\"' '{print $4}'").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	case "linux":
		data, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	return ""
}

// getOSMachineID attempts to read the OS-level persistent machine ID.
func getOSMachineID() string {
	switch runtime.GOOS {
	case "windows":
		out, err := exec.Command("powershell.exe", "-Command",
			"(Get-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Cryptography' -Name MachineGuid).MachineGuid").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	case "darwin":
		out, err := exec.Command("bash", "-c",
			"ioreg -rd1 -c IOPlatformExpertDevice | grep IOPlatformSerialNumber | awk -F'\"' '{print $4}'").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	case "linux":
		for _, path := range []string{"/etc/machine-id", "/var/lib/dbus/machine-id"} {
			data, err := os.ReadFile(path)
			if err == nil {
				val := strings.TrimSpace(string(data))
				if val != "" {
					return val
				}
			}
		}
	}
	return ""
}

// getPrimaryMAC returns the MAC address of the first non-loopback interface.
func getPrimaryMAC() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range interfaces {
		// Skip loopback, virtual, and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		name := strings.ToLower(iface.Name)
		if strings.HasPrefix(name, "veth") || strings.HasPrefix(name, "docker") || strings.HasPrefix(name, "br-") {
			continue
		}
		mac := iface.HardwareAddr.String()
		if mac != "" && mac != "00:00:00:00:00:00" {
			return mac
		}
	}
	return ""
}

// ─── Public API ─────────────────────────────────────────────────────────

// GetMachineID returns a best-effort hardware fingerprint as a SHA-256 hex string.
//
// It attempts to read the BIOS UUID, then falls back through OS-level IDs
// and network MAC addresses. May return different values after hardware
// changes or OS reinstalls. May fail or collide on cheap/virtualized hardware.
//
// Use this for logging, display, or secondary validation.
// For activation HostID, prefer GetOrCreateInstallationID.
//
// Returns an empty string if every layer failed.
func GetMachineID() string {
	layers := []func() string{getBiosUUID, getOSMachineID, getPrimaryMAC}

	for _, layer := range layers {
		raw := layer()
		if raw != "" && len(raw) > 4 && !isGarbageID(raw) {
			return hashID(raw)
		}
	}

	return ""
}

// GetOrCreateInstallationID returns a guaranteed-unique, guaranteed-stable
// installation identifier. On first call, it generates a UUIDv4 seeded with
// whatever hardware info is available and persists it to disk. Every subsequent
// call returns the same value — even across reboots, app updates, and
// hardware upgrades.
//
// This is the recommended value to pass as HostID when activating a license key.
//
// storagePath: Optional custom path for the persistence file.
// Defaults to ~/.keymint/installation-id.
func GetOrCreateInstallationID(storagePath string) (string, error) {
	if storagePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		storagePath = filepath.Join(home, ".keymint", "installation-id")
	}

	// 1. If the file exists, trust it
	if data, err := os.ReadFile(storagePath); err == nil {
		stored := strings.TrimSpace(string(data))
		if stored != "" {
			h := sha256.Sum256([]byte(stored))
			return hex.EncodeToString(h[:]), nil
		}
	}

	// 2. Generate a new installation ID, anchored to hardware when possible
	hardwareAnchor := GetMachineID()
	newUUID := uuid.New().String()
	compositeID := fmt.Sprintf("%s:%s:%d", newUUID, hardwareAnchor, time.Now().UnixMilli())

	// 3. Persist it
	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create directory %s: %w", dir, err)
	}
	if err := os.WriteFile(storagePath, []byte(compositeID), 0600); err != nil {
		return "", fmt.Errorf("cannot write installation ID: %w", err)
	}

	h := sha256.Sum256([]byte(compositeID))
	return hex.EncodeToString(h[:]), nil
}
