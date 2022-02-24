//go:build windows
// +build windows

package clipboard

import (
	"runtime"
	"syscall"
	"unsafe"
)

const (
	cfUnicodeText = 13
	// gmemMoveable  = 0x0002
)

var (
	user32                     = syscall.MustLoadDLL("user32")
	isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable")
	openClipboard              = user32.MustFindProc("OpenClipboard")
	closeClipboard             = user32.MustFindProc("CloseClipboard")
	emptyClipboard             = user32.MustFindProc("EmptyClipboard")
	getClipboardData           = user32.MustFindProc("GetClipboardData")
	setClipboardData           = user32.MustFindProc("SetClipboardData")

	kernel32     = syscall.NewLazyDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
)

func tryOpenClipboard() error {
	ok, _, err := openClipboard.Call(0)
	if ok != 0 {
		return nil
	}
	return err
}

func readClipboard() (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if ok, _, err := isClipboardFormatAvailable.Call(cfUnicodeText); ok == 0 {
		return "", err
	}

	err := tryOpenClipboard()
	if err != nil {
		return "", err
	}

	clipPtr, _, err := getClipboardData.Call(cfUnicodeText)
	if clipPtr == 0 {
		return "", err
	}

	data, _, err := globalLock.Call(clipPtr)
	if data == 0 {
		closeClipboard.Call()
		return "", err
	}

	t := (*[1 << 20]uint16)(unsafe.Pointer(data))
	clipText := syscall.UTF16ToString(t[:])

	ok, _, err := globalUnlock.Call()
	if ok == 0 {
		closeClipboard.Call()
		return "", err
	}

	ok, _, err = closeClipboard.Call()
	if ok == 0 {
		return "", err
	}

	return clipText, nil
}
