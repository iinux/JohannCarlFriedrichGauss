package main

import "testing"

func TestLegacyNATType(t *testing.T) {
	tests := []struct {
		mapping   behavior
		filtering behavior
		noNAT     bool
		want      string
	}{
		{endpointIndependent, endpointIndependent, false, "full-cone-like"},
		{endpointIndependent, addressDependent, false, "address-restricted-cone-like"},
		{endpointIndependent, addressPortDependent, false, "port-restricted-cone-like"},
		{addressPortDependent, endpointIndependent, false, "symmetric NAT / endpoint-dependent mapping"},
		{inconclusive, inconclusive, false, "unknown"},
		{endpointIndependent, endpointIndependent, true, "no NAT (public address is assigned locally)"},
	}
	for _, test := range tests {
		if got := legacyNATType(test.mapping, test.filtering, test.noNAT); got != test.want {
			t.Errorf("legacyNATType(%q, %q, %v) = %q; want %q", test.mapping, test.filtering, test.noNAT, got, test.want)
		}
	}
}
