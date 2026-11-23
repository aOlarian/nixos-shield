package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func showEnterprise() {
	stat, err := os.Stderr.Stat()
	if err != nil || stat.Mode()&os.ModeCharDevice == 0 {
		return
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	dir = filepath.Join(dir, "nixos-shield")
	if os.MkdirAll(dir, 0700) != nil {
		return
	}
	marker, err := os.OpenFile(filepath.Join(dir, "enterprise-intro-v1"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	_ = marker.Close()
	fmt.Fprintln(os.Stderr, `NixOS Shield Enterprise edition · Built to order
We can develop a central console for firewall policy across NixOS servers,
Git and deployment-pipeline integrations, and role-based approvals with an audit log.
Enterprise and custom development: https://verobots.eu/ or https://olarian.com/`)
}
