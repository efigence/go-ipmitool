package ipmitool

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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
	// meta-status that says if last connection to server was successful
	Connected          bool   `json:"connected"`
	Power              bool   `json:"power"`
	PowerRestorePolicy string `json:"power_restore_policy"`
	LastPowerEvent     string `json:"last_power_event"`
	DriveFault         bool   `json:"drive_fault"`
	CoolingFault       bool   `json:"cooling_fault"`
	Identify           bool   `json:"identify"`
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
			st.Connected = true
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

func (i Instance) Identify(state bool) (err error) {
	var idState string
	if state {
		idState = "force"
	} else {
		idState = "0"
	}
	out, err := i.Cmd([]string{"chassis", "identify", idState})
	if len(out) > 0 {
		re := regexp.MustCompile("Chassis identify interval")
		if !re.MatchString(out[0]) {
			return fmt.Errorf("Error: %s,%+v", err, out)
		}
		i.identify = true
	} else {
		return fmt.Errorf("Error: %s,%+v", err, out)
	}
	return err
}

//Chassis Power Control: Cycle

func (i Instance) PowerOff() (err error) {
	out, err := i.Cmd([]string{"chassis", "power", "off"})
	if !strings.Contains(out[0], "Chassis Power Control") || !strings.Contains(out[0], "Off") {
		return fmt.Errorf("unexpected ipmitool output: %+v\nerr: %s", out, err)
	}
	return err
}

func (i Instance) PowerOn() (err error) {
	out, err := i.Cmd([]string{"chassis", "power", "on"})
	if !strings.Contains(out[0], "Chassis Power Control") || !strings.Contains(out[0], "On") {
		return fmt.Errorf("unexpected ipmitool output: %+v\nerr: %s", out, err)
	}
	return err
}

func (i Instance) PowerCycle() (err error) {
	out, err := i.Cmd([]string{"chassis", "power", "cycle"})
	if !strings.Contains(out[0], "Chassis Power Control") || !strings.Contains(out[0], "Cycle") || err != nil {
		return fmt.Errorf("unexpected ipmitool output: %+v\nerr: %s", out, err)
	}
	return err
}
