package mempid

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"strconv"
	"syscall"
	"unsafe"
)

type ProgMutex struct {
	AppName string
	key     uintptr
}

const (
	IPC_RMID  = 0
	IPC_CREAT = 00001000
	IPC_EXCL  = 00002000
	int_size  = strconv.IntSize / 8
)

func (pm *ProgMutex) GetPid() (int, error) {
	// int shmget(key_t key, size_t size of Byte, int shmflg);
	shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, pm.key, int_size, IPC_EXCL)
	if errno != 0 {
		return 0, fmt.Errorf("syscall error: %v", errno)
	}

	// read pid from shared memory
	shmaddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	if errno != 0 {
		return 0, fmt.Errorf("syscall error: %v", errno)
	}
	pid := *(*int)(unsafe.Pointer(uintptr(shmaddr)))

	// test pid
	if f, err := os.Open(fmt.Sprintf("/proc/%d/cmdline", pid)); err == nil {
		buf := bufio.NewReader(f)
		pname, _ := buf.ReadString(0)
		pname = pname[0 : len(pname)-1]
		pname = path.Base(pname)
		f.Close()
		if pname == pm.AppName {
			return pid, nil
		}
	}
	return 0, nil
}

func (pm *ProgMutex) LockProg() error {
	// create the shared memory and put pid to it
	shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, pm.key, int_size, IPC_EXCL)
	if errno != 0 {
		shmid, _, errno = syscall.Syscall(syscall.SYS_SHMGET, pm.key, int_size, IPC_CREAT|IPC_EXCL)
		if errno != 0 {
			return fmt.Errorf("syscall error: %v", errno)
		}
	}
	shmaddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	if errno != 0 {
		return fmt.Errorf("syscall error: %v", errno)
	}
	*(*int)(unsafe.Pointer(uintptr(shmaddr))) = os.Getpid()
	syscall.Syscall(syscall.SYS_SHMDT, shmaddr, 0, 0)
	return nil
}

func (pm *ProgMutex) UnLockProg() {
	if shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, pm.key, int_size, IPC_EXCL); errno != 0 {
		syscall.Syscall6(syscall.SYS_SHMCTL, shmid, IPC_RMID, 0, 0, 0, 0)
	}
}

func genKey(name string) uintptr {
	sum := sha1.Sum([]byte(name))
	var key uint
	key = uint(sum[0])<<24 + uint(sum[1])<<16 + uint(sum[2])<<8 + uint(sum[3])
	return uintptr(key)
}
