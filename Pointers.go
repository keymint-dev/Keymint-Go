package keymint

// String returns a pointer to s, for optional string fields such as
// ActivateKeyParams.HostID.
func String(s string) *string { return &s }

// Int returns a pointer to i, for optional integer fields.
func Int(i int) *int { return &i }

// Bool returns a pointer to b, for optional boolean fields.
func Bool(b bool) *bool { return &b }
