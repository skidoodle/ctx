//go:build windows

package clipboard

import (
	"fmt"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	openClipboard    = user32.NewProc("OpenClipboard")
	closeClipboard   = user32.NewProc("CloseClipboard")
	emptyClipboard   = user32.NewProc("EmptyClipboard")
	setClipboardData = user32.NewProc("SetClipboardData")
	globalAlloc      = kernel32.NewProc("GlobalAlloc")
	globalLock       = kernel32.NewProc("GlobalLock")
	globalUnlock     = kernel32.NewProc("GlobalUnlock")
	globalFree       = kernel32.NewProc("GlobalFree")
)

const (
	cfHDrop = 15
	gHnd    = 0x0042
)

type dropFiles struct {
	pFiles uint32
	pt     struct{ x, y int32 }
	fNC    int32
	fWide  int32
}

func CopyFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	pathUTF16, err := syscall.UTF16FromString(absPath)
	if err != nil {
		return err
	}
	pathUTF16 = append(pathUTF16, 0)

	dropSize := uint32(unsafe.Sizeof(dropFiles{}))
	dataSize := uint32(len(pathUTF16) * 2)
	totalSize := uintptr(dropSize + dataSize)

	hMem, _, err := globalAlloc.Call(gHnd, totalSize)
	if hMem == 0 {
		return fmt.Errorf("GlobalAlloc failed: %v", err)
	}

	ptrVal, _, _ := globalLock.Call(hMem)
	if ptrVal == 0 {
		_, _, _ = globalFree.Call(hMem)
		return fmt.Errorf("GlobalLock failed")
	}

	df := (*dropFiles)(unsafe.Pointer(ptrVal))
	df.pFiles = dropSize
	df.fWide = 1

	targetPtr := unsafe.Pointer(ptrVal + uintptr(dropSize))

	srcSlice := unsafe.Slice((*uint16)(unsafe.Pointer(&pathUTF16[0])), len(pathUTF16))
	dstSlice := unsafe.Slice((*uint16)(targetPtr), len(pathUTF16))
	copy(dstSlice, srcSlice)

	_, _, _ = globalUnlock.Call(hMem)

	var openSuccess bool
	for range 10 {
		ret, _, _ := openClipboard.Call(0)
		if ret != 0 {
			openSuccess = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if !openSuccess {
		_, _, _ = globalFree.Call(hMem)
		return fmt.Errorf("clipboard locked")
	}
	defer func() {
		_, _, _ = closeClipboard.Call()
	}()

	_, _, _ = emptyClipboard.Call()

	if ret, _, err := setClipboardData.Call(uintptr(cfHDrop), hMem); ret == 0 {
		_, _, _ = globalFree.Call(hMem)
		return fmt.Errorf("SetClipboardData failed: %v", err)
	}

	return nil
}
