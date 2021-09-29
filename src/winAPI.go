package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32DLL               = windows.NewLazyDLL("user32.dll")
	procGetMessage          = user32DLL.NewProc("GetMessageW")
	procGetKeyState         = user32DLL.NewProc("GetKeyState")
	procCallNextHookEx      = user32DLL.NewProc("CallNextHookEx")
	procSetWindowsHookEx    = user32DLL.NewProc("SetWindowsHookExA")
	procUnhookWindowsHookEx = user32DLL.NewProc("UnhookWindowsHookEx")
	keyboardHook            HHOOK
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 256
	WM_KEYUP       = 257
)

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	UINT      uint32
	ULONG_PTR uintptr
	LPMSG     MSG
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type KBDLLHOOKSTRUCT struct {
	vkCode      DWORD
	scanCode    DWORD
	flags       DWORD
	time        DWORD
	dwExtraInfo ULONG_PTR
}

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(lpMsg *MSG, hWnd HWND, wMsgFilterMin UINT, wMsgFilterMax UINT) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr(hWnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax),
	)

	return int(ret)
}
