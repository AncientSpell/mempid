package mempid

import (
	"os"
	"path"
)

var (
	_pm *ProgMutex
)

func NewProgMutex(keyName string) *ProgMutex {
	key := genKey(keyName)
	return &ProgMutex{AppName: path.Base(os.Args[0]), key: key}
}

/*
GetPid return value of int:
	0 	not running.
	>0 	progress is running, and the return value is the PID of progress.
*/
func GetPid() (int, error) {
	return _pm.GetPid()
}

// LockProg release share memory
func LockProg() error {
	return _pm.LockProg()
}

func UnLockProg() {
	_pm.UnLockProg()
}

func init() {
	_pm = NewProgMutex(path.Base(os.Args[0]))
}
