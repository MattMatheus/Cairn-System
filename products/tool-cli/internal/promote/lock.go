package promote

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type MemoryRootLock struct {
	path string
	file *os.File
}

func AcquireMemoryRootLock(root string) (*MemoryRootLock, error) {
	if err := os.MkdirAll(root, 0o775); err != nil {
		return nil, err
	}
	lockPath := filepath.Join(root, ".promote.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o664)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("acquire promotion lock: %w", err)
	}
	return &MemoryRootLock{path: lockPath, file: f}, nil
}

func (l *MemoryRootLock) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	if err := syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN); err != nil {
		_ = l.file.Close()
		return err
	}
	return l.file.Close()
}
