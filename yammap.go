package yammap

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Mmap holds our in-memory data
type Mmap struct {
	sync.RWMutex
	fd     *os.File
	offset int64
	data   []byte
}

// pageSize is the size of a page in the system.
var pageSize int

const (
	SEEK_START   int = 0 // seek relative to the origin of the file
	SEEK_CURRENT int = 1 // seek relative to the current offset
	SEEK_END     int = 2 // seek relative to the end

	MREMAP_MAYMOVE = 0x1

	maxSize = 0xFFFFFFFFFFFF // Maximum allocation size, 2^48 bytes for x86_64
)

func init() {
	pageSize = unix.Getpagesize()
}

// Open opens or creates the named file as memmory-mapped.
func OpenFile(name string, flag int, perm os.FileMode) (*Mmap, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	m := new(Mmap)
	err = m.fileMap(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return m, nil
}

// Close closes the memory-mapped file, rendering it unusable for I/O.
func (m *Mmap) Close() (err error) {
	m.Lock()
	defer m.Unlock()
	addr := unsafe.Pointer(&m.data[0])
	_, _, errno := unix.Syscall(unix.SYS_MUNMAP, uintptr(addr), uintptr(len(m.data)), 0)
	if errno != 0 {
		err = fmt.Errorf("SYS_MUNMAP: %v", errno)
	}
	err = m.fd.Close()
	if err != nil {
		return err
	}
	m = nil
	return err
}

// Sync flushes changes made to a file that was mapped into memory using mmap back to the filesystem.
func (m *Mmap) Sync() (err error) {
	m.Lock()
	addr := unsafe.Pointer(&m.data[0])
	_, _, errno := unix.Syscall(unix.SYS_MSYNC, uintptr(addr), uintptr(len(m.data)), uintptr(unix.MS_SYNC))
	if errno != 0 {
		err = fmt.Errorf("SYS_MSYNC: %v", errno)
	}
	m.Unlock()
	return err
}

// Read reads up to len(b) bytes from the File. It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (m *Mmap) Read(b []byte) (n int, err error) {
	m.Lock()
	defer m.Unlock()
	if m.offset >= int64(len(m.data)) {
		return n, io.EOF
	}
	n = copy(b, m.data[m.offset:])
	m.offset += int64(n)
	return n, nil
}

// ReadAt reads len(b) bytes from the File starting at byte offset off. It returns the number of bytes read and the error, if any.
// ReadAt always returns a non-nil error when n < len(b). At end of file, that error is io.EOF.
func (m *Mmap) ReadAt(b []byte, off int64) (n int, err error) {
	m.RLock()
	defer m.RUnlock()
	if off >= int64(len(m.data)) {
		return 0, io.EOF
	}
	n = copy(b, m.data[off:])
	if n < len(b) {
		err = io.EOF
	}
	return n, err
}

// Size returns the size of the file.
func (m *Mmap) Size() int64 {
	m.RLock()
	size := int64(len(m.data))
	m.RUnlock()
	return size
}

// Name returns the name of the file as presented to Open.
func (m *Mmap) Name() string {
	return m.fd.Name()
}

// Offset returns the current offset.
func (m *Mmap) Offset() int64 {
	m.RLock()
	offset := m.offset
	m.RUnlock()
	return offset
}

// Seek sets the offset for the next Read or Write on file to offset, interpreted according to whence:
// 0 means relative to the origin of the file,
// 1 means relative to the current offset,
// and 2 means relative to the end. It returns the new offset and an error, if any.
func (m *Mmap) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	m.Lock()
	defer m.Unlock()
	switch whence {
	case SEEK_START:
		abs = offset
	case SEEK_CURRENT:
		abs = m.offset + offset
	case SEEK_END:
		abs = int64(len(m.data)) + offset
	default:
		return 0, errors.New("invalid whence value")
	}
	if abs < 0 {
		return 0, errors.New("negative position")
	}
	if abs > int64(len(m.data)) {
		return 0, errors.New("offset goes beyond the end of file")
	}
	m.offset = abs
	return abs, nil
}

// Write writes len(b) bytes to the File. It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (m *Mmap) Write(b []byte) (n int, err error) {
	m.Lock()
	if m.offset+int64(len(b)) > int64(len(m.data)) {
		err = m.mremap(m.offset + int64(len(b)))
		if err != nil {
			m.Unlock()
			return 0, err
		}
	}
	m.fd.Seek(m.offset+int64(len(b)-1), SEEK_START)
	m.fd.Write([]byte("\x00"))
	n = copy(m.data[m.offset:], b)
	m.offset += int64(n)
	m.Unlock()
	if n != len(b) {
		err = io.ErrShortWrite
	}
	return n, err
}

// WriteAt writes len(b) bytes to the File starting at byte offset off. It returns the number of bytes written and an error, if any.
// WriteAt returns a non-nil error when n != len(b).
func (m *Mmap) WriteAt(b []byte, off int64) (n int, err error) {
	m.Lock()
	if off+int64(len(b)) > int64(len(m.data)) {
		err = m.mremap(off + int64(len(b)))
		if err != nil {
			m.Unlock()
			return 0, err
		}
	}
	m.fd.Seek(off+int64(len(b)-1), SEEK_START)
	m.fd.Write([]byte("\x00"))
	n = copy(m.data[off:], b)
	m.Unlock()
	if n != len(b) {
		err = io.ErrShortWrite
	}
	return n, err
}

// Truncate changes the size of the file. It does not change the I/O offset.
func (m *Mmap) Truncate(size int64) error {
	size = int64(align(size))
	m.Lock()
	defer m.Unlock()
	err := m.fd.Truncate(size)
	if err != nil {
		return err
	}
	err = m.mremap(size)
	return err
}

// slice is the runtime representation of a Go slice.
type slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// Map file to memory
func (m *Mmap) fileMap(f *os.File) error {
	m.fd = f
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	size := align(stat.Size())
	mmapAddr, _, errno := unix.Syscall6(
		unix.SYS_MMAP,
		0,
		uintptr(size),
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_SHARED,
		uintptr(f.Fd()),
		0,
	)
	if errno != 0 {
		return fmt.Errorf("mmap failed with errno: %v", err)
	}
	header := (*slice)(unsafe.Pointer(&m.data))
	header.Data = unsafe.Pointer(mmapAddr)
	header.Cap = int(size)
	header.Len = int(size)
	runtime.KeepAlive(mmapAddr)
	return err
}

// Use mremap to increase the size of allocated memory
func (m *Mmap) mremap(size int64) error {
	size = align(size)
	header := (*slice)(unsafe.Pointer(&m.data))
	mmapAddr, mmapSize, err := unix.Syscall6(
		unix.SYS_MREMAP,
		uintptr(header.Data),
		uintptr(header.Len),
		uintptr(size),
		uintptr(MREMAP_MAYMOVE),
		0,
		0,
	)
	if err != 0 {
		return fmt.Errorf("mremap failed with errno: %v", err)
	}
	if mmapSize != uintptr(size) {
		return fmt.Errorf("mremap size mismatch: requested: %d got: %d", size, mmapSize)
	}
	header.Data = unsafe.Pointer(mmapAddr)
	header.Cap = int(size)
	header.Len = int(size)
	runtime.KeepAlive(mmapAddr)
	return nil
}

// Align to page boundries
func align(size int64) int64 {
	var aligned int64
	if size == 0 {
		aligned = int64(pageSize)
	} else if (size % int64(pageSize)) != 0 {
		aligned = (size/int64(pageSize) + 1) * int64(pageSize)
	} else {
		aligned = size
	}
	if aligned > maxSize {
		aligned = maxSize
	}
	return aligned
}
