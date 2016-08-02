package ipmitool

import (
	"fmt"
	"os/exec"
	"strings"
)

type Instance struct {
	IP       string
	AuthType string
	User     string
	Password string
	identify bool
}

func (i Instance) Cmd(cmdArgs []string) ([]string, error) {
	var args []string
	var env []string
	// auth
	args = append(args, "-H", i.IP)
	args = append(args, "-I", i.AuthType)
	if i.AuthType != AuthNone {
		// -E means "take password from env IPMI_PASSWORD"
		args = append(args, "-E")
		env = append(env, fmt.Sprintf("IPMITOOL_PASSWORD=%s", i.Password))
	}
	args = append(args, "-U", i.User)
	args = append(args, "-I", "lan")
	args = append(args, cmdArgs...)
	cmd := exec.Command("ipmitool", args...)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	return strings.Split(string(out), "\n"), err
}
