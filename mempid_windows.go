package mempid

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

type ProgMutex struct {
	AppName string
	key     *uint16
	handle  syscall.Handle
}

const (
	PAGE_READWRITE = 0x0004
	FILE_MAP_READ  = 0x0004
	FILE_MAP_WRITE = 0x0008
	int_size       = strconv.IntSize
)

func genKey(name string) *uint16 {
	sum := sha1.Sum([]byte(name))
	fmt.Println(name)
	share, _ := syscall.UTF16PtrFromString(hex.EncodeToString(sum[:])) // name of the share memory
	return share
}

func (pm *ProgMutex) GetPid() (int, error) {
	handle, err := syscall.CreateFileMapping(syscall.InvalidHandle, nil, syscall.PAGE_READONLY, 0, int_size, pm.key)
	if err != nil {
		return 0, err
	}
	defer syscall.CloseHandle(handle)
	ptr, err := syscall.MapViewOfFile(handle, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		return 0, err
	}
	pid := *(*int)(unsafe.Pointer(ptr))
	if _, err := os.FindProcess(pid); err != nil {
		pid = 0
	}
	return pid, syscall.UnmapViewOfFile(ptr)
}

func (pm *ProgMutex) LockProg() error {
	var err error
	if pm.handle != 0 && pm.handle != syscall.InvalidHandle {
		return errors.New("you should unlock first.")
	}
	pm.handle, err = syscall.CreateFileMapping(syscall.InvalidHandle, nil, syscall.PAGE_READWRITE, 0, int_size, pm.key)
	if err != nil {
		return err
	}
	ptr, err := syscall.MapViewOfFile(pm.handle, syscall.FILE_MAP_WRITE, 0, 0, 0)
	if err != nil {
		return err
	}
	*(*int)(unsafe.Pointer(ptr)) = os.Getpid()
	return syscall.UnmapViewOfFile(ptr)
}

func (pm *ProgMutex) UnLockProg() {
	if pm.handle != 0 && pm.handle != syscall.InvalidHandle {
		syscall.CloseHandle(pm.handle)
		pm.handle = 0
	}
}
