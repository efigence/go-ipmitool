package ipmitool

import (
	"errors"
	"fmt"
	"regexp"
)

// Example output of chassis status:
// System Power         : on
// Power Overload       : false
// Power Interlock      : inactive
// Main Power Fault     : false
// Power Control Fault  : false
// Power Restore Policy : always-on
// Last Power Event     :
// Chassis Intrusion    : inactive
// Front-Panel Lockout  : inactive
// Drive Fault          : false
// Cooling/Fan Fault    : false

type ChassisStatus struct {
	Power              bool   `json:"power"`
	PowerRestorePolicy string `json:"power_restore_policy"`
	LastPowerEvent     string `json:"last_power_event"`
	DriveFault         bool   `json:"drive_fault"`
	CoolingFault       bool   `json:"cooling_fault"`
}

func (i Instance) GetChassisStatus() (ChassisStatus, error) {
	var st ChassisStatus
	out, err := i.Cmd([]string{"chassis", "status"})
	splitRe := regexp.MustCompile(`(.*\S)\s+?:\s+(\S.*)`)
	if err != nil {
		return st, err
	}
	for _, line := range out {
		m := splitRe.FindStringSubmatch(line)
		if len(m) < 3 {
			continue
		}
		switch {
		case m[1] == "System Power":
			if m[2] == "on" {
				st.Power = true
			} else if m[2] == "off" {
				st.Power = false
			} else {
				return st, errors.New(fmt.Sprintf("Unknown system power state: [%s]"))
			}
		case m[1] == "Power Restore Policy":
			st.PowerRestorePolicy = m[2]
		case m[1] == "Last Power Event":
			st.LastPowerEvent = m[2]
		case m[1] == "Drive Fault":
			if m[2] == "true" {
				st.DriveFault = true
			}
		case m[1] == "Cooling/Fan Fault":
			if m[2] == "true" {
				st.DriveFault = true
			}
		}
	}
	return st, err
}
