package main

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name     string
		contract Contract
		valid    bool
	}{
		{"service identifier", Contract{Service: "mastodon.demo-1"}, true},
		{"missing service", Contract{}, false},
		{"multiline service", Contract{Service: "mastodon\nsecond-line"}, false},
		{"invalid port", Contract{Service: "matrix", Network: Network{Inbound: []int{65536}}}, false},
		{"unknown federation flow", Contract{Service: "matrix", Network: Network{Federation: []string{"unknown"}}}, false},
		{"known flows", Contract{Service: "matrix", Network: Network{Inbound: []int{443, 8448}, Outbound: []string{"dns", "https"}, Federation: []string{"matrix"}}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := validate(tc.contract); (err == nil) != tc.valid {
				t.Fatalf("validate error = %v, want valid %t", err, tc.valid)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	rules := generate(Contract{Service: "matrix", Network: Network{
		Inbound: []int{8448, 443}, Outbound: []string{"dns", "https"}, Federation: []string{"matrix"},
	}})
	for _, line := range []string{"tcp dport 443 accept", "tcp dport 8448 accept", "udp dport 53 accept", "tcp dport 53 accept"} {
		if !strings.Contains(rules, line) {
			t.Errorf("missing rule %q", line)
		}
	}
	if strings.Count(rules, "policy drop;") != 2 {
		t.Error("input and output chains must both use a default drop policy")
	}
}
