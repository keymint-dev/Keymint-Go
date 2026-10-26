package keymint

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

type offlineFixtureLicense struct {
	LicenseKey  string `json:"licenseKey"`
	ProductID   string `json:"productId"`
	MachineCode string `json:"machineCode"`
	IssuedAt    *int64 `json:"issuedAt"`
	ExpiresAt   int64  `json:"expiresAt"`
	KeyID       string `json:"keyId"`
}

type offlineFixtures struct {
	PublicKeys []string `json:"publicKeys"`
	KeyIDs     []string `json:"keyIds"`
	Cases      []struct {
		Name    string `json:"name"`
		Input   string `json:"input"`
		Options struct {
			ProductID   string   `json:"productId"`
			MachineCode string   `json:"machineCode"`
			Now         int64    `json:"now"`
			ClockSkew   *float64 `json:"clockSkew"`
		} `json:"options"`
		Expect struct {
			Valid   bool                   `json:"valid"`
			Error   string                 `json:"error"`
			License *offlineFixtureLicense `json:"license"`
		} `json:"expect"`
	} `json:"cases"`
}

func loadOfflineFixtures(t *testing.T) offlineFixtures {
	t.Helper()
	raw, err := os.ReadFile("testdata/offline-license-fixtures.json")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	var f offlineFixtures
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("parse fixtures: %v", err)
	}
	return f
}

func TestVerifyOfflineLicenseKeyIDs(t *testing.T) {
	f := loadOfflineFixtures(t)
	for i, pemKey := range f.PublicKeys {
		kid, err := OfflineKeyID(pemKey)
		if err != nil {
			t.Fatalf("OfflineKeyID(%d): %v", i, err)
		}
		if kid != f.KeyIDs[i] {
			t.Errorf("key %d: kid %q, want %q", i, kid, f.KeyIDs[i])
		}
	}
}

func TestVerifyOfflineLicenseFixtures(t *testing.T) {
	f := loadOfflineFixtures(t)
	if len(f.Cases) != 34 {
		t.Fatalf("expected 34 fixture cases, got %d", len(f.Cases))
	}
	for _, tc := range f.Cases {
		t.Run(tc.Name, func(t *testing.T) {
			opts := OfflineVerifyOptions{
				PublicKeys:  f.PublicKeys,
				ProductID:   tc.Options.ProductID,
				MachineCode: tc.Options.MachineCode,
				Now:         time.Unix(tc.Options.Now, 0),
			}
			if tc.Options.ClockSkew != nil {
				skew := time.Duration(*tc.Options.ClockSkew * float64(time.Second))
				opts.ClockSkew = &skew
			}
			got, err := VerifyOfflineLicense(tc.Input, opts)

			if !tc.Expect.Valid {
				var licErr *OfflineLicenseError
				if !errors.As(err, &licErr) {
					t.Fatalf("expected OfflineLicenseError %q, got license=%+v err=%v", tc.Expect.Error, got, err)
				}
				if string(licErr.Code) != tc.Expect.Error {
					t.Fatalf("code %q, want %q", licErr.Code, tc.Expect.Error)
				}
				if got != nil {
					t.Fatalf("expected nil license on failure")
				}
				if strings.Contains(licErr.Error(), "KM-ABCD-1234-EFGH-5678") {
					t.Fatalf("error message leaks the license key: %q", licErr.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected valid, got %v", err)
			}
			want := tc.Expect.License
			if got.LicenseKey != want.LicenseKey || got.ProductID != want.ProductID ||
				got.MachineCode != want.MachineCode || got.KeyID != want.KeyID {
				t.Fatalf("license %+v, want %+v", got, want)
			}
			if got.ExpiresAt.Unix() != want.ExpiresAt {
				t.Fatalf("ExpiresAt %d, want %d", got.ExpiresAt.Unix(), want.ExpiresAt)
			}
			switch {
			case want.IssuedAt == nil && got.IssuedAt != nil:
				t.Fatalf("IssuedAt %v, want nil", got.IssuedAt)
			case want.IssuedAt != nil && (got.IssuedAt == nil || got.IssuedAt.Unix() != *want.IssuedAt):
				t.Fatalf("IssuedAt %v, want %d", got.IssuedAt, *want.IssuedAt)
			}
		})
	}
}

func TestVerifyOfflineLicenseDefaultsNowToSystemClock(t *testing.T) {
	f := loadOfflineFixtures(t)
	// The fixture tokens expire in 2027; with the system clock well before
	// that, a valid fixture must verify without Now. Skip once that is false.
	valid := f.Cases[0]
	if time.Now().Unix() >= valid.Options.Now {
		t.Skip("system clock is past the fixture's fixed time")
	}
	_, err := VerifyOfflineLicense(valid.Input, OfflineVerifyOptions{
		PublicKeys:  f.PublicKeys,
		ProductID:   valid.Options.ProductID,
		MachineCode: valid.Options.MachineCode,
	})
	if err != nil {
		t.Fatalf("expected valid with default Now, got %v", err)
	}
}

func TestVerifyOfflineLicenseOptionErrors(t *testing.T) {
	f := loadOfflineFixtures(t)
	valid := f.Cases[0]

	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ecDER, err := x509.MarshalPKIXPublicKey(&ecKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	ecPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: ecDER}))

	dur := func(d time.Duration) *time.Duration { return &d }
	base := func() OfflineVerifyOptions {
		return OfflineVerifyOptions{
			PublicKeys:  f.PublicKeys,
			ProductID:   valid.Options.ProductID,
			MachineCode: valid.Options.MachineCode,
			Now:         time.Unix(valid.Options.Now, 0),
		}
	}

	cases := []struct {
		name   string
		mutate func(*OfflineVerifyOptions)
	}{
		{"no public keys", func(o *OfflineVerifyOptions) { o.PublicKeys = nil }},
		{"empty public keys", func(o *OfflineVerifyOptions) { o.PublicKeys = []string{} }},
		{"garbage PEM", func(o *OfflineVerifyOptions) { o.PublicKeys = []string{"not a pem"} }},
		{"PEM with garbage body", func(o *OfflineVerifyOptions) {
			o.PublicKeys = []string{"-----BEGIN PUBLIC KEY-----\nAAAA\n-----END PUBLIC KEY-----\n"}
		}},
		{"non-Ed25519 key", func(o *OfflineVerifyOptions) { o.PublicKeys = []string{ecPEM} }},
		{"one bad key among good", func(o *OfflineVerifyOptions) { o.PublicKeys = append([]string{ecPEM}, f.PublicKeys...) }},
		{"empty product", func(o *OfflineVerifyOptions) { o.ProductID = "" }},
		{"empty machine", func(o *OfflineVerifyOptions) { o.MachineCode = "" }},
		{"negative skew", func(o *OfflineVerifyOptions) { o.ClockSkew = dur(-time.Second) }},
		{"skew over 1h", func(o *OfflineVerifyOptions) { o.ClockSkew = dur(time.Hour + time.Second) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := base()
			tc.mutate(&opts)
			got, err := VerifyOfflineLicense(valid.Input, opts)
			if !errors.Is(err, ErrInvalidOfflineOptions) {
				t.Fatalf("expected ErrInvalidOfflineOptions, got license=%v err=%v", got, err)
			}
			var licErr *OfflineLicenseError
			if errors.As(err, &licErr) {
				t.Fatalf("option error must not be an OfflineLicenseError")
			}
		})
	}

	for _, skew := range []time.Duration{0, time.Hour} {
		opts := base()
		opts.ClockSkew = dur(skew)
		if _, err := VerifyOfflineLicense(valid.Input, opts); err != nil {
			t.Errorf("skew %v should be accepted, got %v", skew, err)
		}
	}

	if _, err := OfflineKeyID(ecPEM); !errors.Is(err, ErrInvalidOfflineOptions) {
		t.Errorf("OfflineKeyID(non-Ed25519) = %v, want ErrInvalidOfflineOptions", err)
	}
}
