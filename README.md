# mempid
## Purpose
Record pid to shared memory to create a mutex progress
## Usage
```go
func main() {
	pid, err := mempid.GetPid()
	if err != nil {
		fmt.Println("fail to get pid.", err)
		return
	}
	if pid != 0 {
		fmt.Println("master process id is", pid)
		return
	}
	mempid.LockProg()
	defer mempid.UnLockProg()
	...
}
```
