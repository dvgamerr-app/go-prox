package main

import (
	"syscall"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procPostMessage = user32.NewProc("PostMessageW")
)

const (
	HWND_BROADCAST  = 0xFFFF // Broadcast to all windows
	WM_SYSCOMMAND   = 0x0112 // System command message
	SC_MONITORPOWER = 0xF170 // Monitor power command
	MONITOR_OFF     = 2      // Command to turn off the monitor
	MONITOR_ON      = -1     // Command to turn on the monitor
)

func PostMessage(hWnd uintptr, msg uint32, wParam, lParam uintptr) bool {
	ret, _, _ := procPostMessage.Call(
		hWnd,
		uintptr(msg),
		wParam,
		lParam,
	)
	return ret != 0
}
