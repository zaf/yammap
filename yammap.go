/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

/*
Package yammap provides an interface to memory mapped files.
*/

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

// Mmap holds our in-memory file data
type Mmap struct {
	sync.RWMutex
	fd     *os.File
	offset int64
	data   []byte
	append bool
}

// pageSize is the size of a page in the system.
var pageSize int

const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY = 0x0 // open the file read-only
	O_WRONLY = 0x1 // open the file write-only
	O_RDWR   = 0x2 // open the file read-write
	// The remaining values may be or'ed in to control behavior.
	O_APPEND = 0x400    // append data to the file when writing
	O_CREATE = 0x40     // create a new file if none exists
	O_EXCL   = 0x80     // used with O_CREATE, file must not exist
	O_SYNC   = 0x101000 // open for synchronous I/O
	O_TRUNC  = 0x200    // truncate to zero length
	// Page protections modes
	PROT_NONE  = 0x0 // page protection: no access
	PROT_READ  = 0x1 // page protection: read-only
	PROT_WRITE = 0x2 // page protection: read-write
	PROT_EXEC  = 0x4 // page protection: read-execute

	MAP_SHARED          = 0x1    // share changes
	MAP_PRIVATE         = 0x2    // changes are private
	MAP_SHARED_VALIDATE = 0x3    // share changes, but validate
	MAP_LOCKED          = 0x2000 // pages are locked to RAM
	MAP_POPULATE        = 0x8000 // populate (prefault) pagetables

	MREMAP_MAYMOVE   = 0x1 // may move the mapping
	MREMAP_FIXED     = 0x2 // map at a fixed address
	MREMAP_DONTUNMAP = 0x4 // don't unmap the mapping on close

	SEEK_SET = 0x0 // seek relative to the origin of the file
	SEEK_CUR = 0x1 // seek relative to the current offset
	SEEK_END = 0x2 // seek relative to the end
)

func init() {
	pageSize = os.Getpagesize()
}

// Open opens or creates the named file as memmory-mapped.
func OpenFile(name string, flag int, perm uint32) (*Mmap, error) {
	f, err := os.OpenFile(name, flag, os.FileMode(perm))
	if err != nil {
		return nil, err
	}
	m := new(Mmap)
	m.fd = f
	m.append = flag&O_APPEND != 0
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		os.Remove(name)
		return nil, err
	}
	err = m.mmap(stat.Size(), flag)
	if err != nil {
		f.Close()
		os.Remove(name)
		return nil, err
	}
	return m, nil
}

// Create creates the named file of specified size as memmory-mapped.
func Create(name string, size int64, flag int, perm uint32) (*Mmap, error) {
	f, err := os.OpenFile(name, flag, os.FileMode(perm))
	if err != nil {
		return nil, err
	}
	m := new(Mmap)
	m.fd = f
	err = m.mmap(size, flag)
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
	_, _, errno := unix.Syscall(SYS_MUNMAP, uintptr(addr), uintptr(len(m.data)), 0)
	if errno != 0 {
		err = fmt.Errorf("SYS_MUNMAP: %s", errno.Error())
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
	_, _, errno := unix.Syscall(SYS_MSYNC, uintptr(addr), uintptr(len(m.data)), uintptr(unix.MS_SYNC))
	if errno != 0 {
		err = fmt.Errorf("SYS_MSYNC: %s", errno.Error())
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
	case SEEK_SET:
		abs = offset
	case SEEK_CUR:
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
	if m.append {
		m.offset = int64(len(m.data))
	}
	if m.offset+int64(len(b)) > int64(len(m.data)) {
		err = m.mremap(int64(len(m.data) * 2))
		if err != nil {
			m.Unlock()
			return 0, err
		}
	}
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
	if m.append {
		return 0, errors.New("invalid use of WriteAt on file opened with O_APPEND")
	}
	m.Lock()
	if off+int64(len(b)) > int64(len(m.data)) {
		err = m.mremap(int64(len(b) * 2))
		if err != nil {
			m.Unlock()
			return 0, err
		}
	}
	n = copy(m.data[off:], b)
	m.Unlock()
	if n != len(b) {
		err = io.ErrShortWrite
	}
	return n, err
}

// Truncate changes the size of the file. It does not change the I/O offset.
func (m *Mmap) Truncate(size int64) error {
	m.Lock()
	err := m.mremap(size)
	m.Unlock()
	return err
}

// slice is the runtime representation of a Go slice.
type slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// Map file to memory
func (m *Mmap) mmap(size int64, flag int) error {
	size = align(size)
	var protection int
	mapping := MAP_SHARED | MAP_POPULATE
	if flag&O_WRONLY != 0 {
		protection = PROT_WRITE
	} else if flag&O_RDWR != 0 {
		protection = PROT_READ | PROT_WRITE
	} else {
		protection = PROT_READ
	}
	if protection != PROT_READ {
		err := m.fd.Truncate(int64(size))
		if err != nil {
			return err
		}
	}
	mmapAddr, _, errno := unix.Syscall6(
		SYS_MMAP,
		0,
		uintptr(size),
		uintptr(protection),
		uintptr(mapping),
		m.fd.Fd(),
		0,
	)
	if errno != 0 {
		return fmt.Errorf("mmap failed: %s", errno.Error())
	}
	header := (*slice)(unsafe.Pointer(&m.data))
	header.Data = unsafe.Pointer(mmapAddr)
	header.Cap = int(size)
	header.Len = int(size)
	runtime.KeepAlive(mmapAddr)
	return nil
}

// Use mremap to increase the size of allocated memory
func (m *Mmap) mremap(size int64) error {
	size = align(size)
	err := m.fd.Truncate(size)
	if err != nil {
		return err
	}
	header := (*slice)(unsafe.Pointer(&m.data))
	mmapAddr, _, errno := unix.Syscall6(
		SYS_MREMAP,
		uintptr(header.Data),
		uintptr(header.Len),
		uintptr(size),
		uintptr(MREMAP_MAYMOVE),
		0,
		0,
	)
	if errno != 0 {
		return fmt.Errorf("mremap failed: %v", errno.Error())
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
