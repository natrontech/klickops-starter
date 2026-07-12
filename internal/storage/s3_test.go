package storage

import "testing"

func TestNormalizeEndpoint(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"rook-ceph-rgw:80", "http://rook-ceph-rgw:80"},
		{"http://localhost:8333", "http://localhost:8333"},
		{"https://s3.amazonaws.com", "https://s3.amazonaws.com"},
	}
	for _, tt := range tests {
		if got := NormalizeEndpoint(tt.in); got != tt.want {
			t.Errorf("NormalizeEndpoint(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
