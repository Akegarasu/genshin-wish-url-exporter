package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

const (
	InternetSettings = `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`
)

func SetProxy(enable bool, server string, ignore string) error {
	var err error
	if enable {
		err = setProxyServer(server)
		if err != nil {
			return err
		}
		err = setIgnoreProxy(ignore)
		if err != nil {
			return err
		}
		err = enableProxy()
		if err != nil {
			return err
		}
	} else {
		err = disableProxy()
		if err != nil {
			return err
		}
	}
	return nil
}

func setProxyServer(server string) error {
	return set("ProxyServer", "REG_SZ", server)
}

func setIgnoreProxy(ignore string) error {
	return set("ProxyOverride", "REG_SZ", ignore)
}

func enableProxy() error {
	return set("ProxyEnable", "REG_DWORD", "1")
}

func disableProxy() error {
	return set("ProxyEnable", "REG_DWORD", "0")
}

func set(key string, typ string, value string) error {
	_, err := execCmd(`reg`, `add`, InternetSettings, `/v`, key, `/t`, typ, `/d`, value, `/f`)
	return err
}

func execCmd(name string, arg ...string) (string, error) {
	c := exec.Command(name, arg...)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%q: %w: %q", strings.Join(append([]string{name}, arg...), " "), err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
