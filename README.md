# mempid
## Purpose
Record pid to shared memory to create a mutex progress
## Usage
```go
func main() {
	pid, _ := mempid.GetPid()
	if pid != 0 {
		fmt.Println("master process id is", pid)
		return
	}
	mempid.LockProg()
	defer mempid.UnLockProg()
	...
}
```
