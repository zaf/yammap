# yammap

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/zaf/yammap)

Package yammap provides an interface to memory mapped files.

WIP - Don't use in production

## Constants

```golang
const (
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

    // Mapping advice, refer to madvise(2) manual page.
    MADV_NORMAL      = 0x0  // no special treatment.  This is the default.
    MADV_RANDOM      = 0x1  // expect random page references.
    MADV_SEQUENTIAL  = 0x2  // expect sequential page references.
    MADV_WILLNEED    = 0x3  // will need these pages.
    MADV_DONTNEED    = 0x4  // don't need these pages.
    MADV_FREE        = 0x8  // pages can be freed.
    MADV_REMOVE      = 0x9  // remove these pages from the mappings.
    MADV_DONTFORK    = 0xa  // do not inherit across fork.
    MADV_DOFORK      = 0xb  // inherit across fork.
    MADV_MERGEABLE   = 0xc  // enable Kernel Samepage Merging (KSM) for the pages
    MADV_UNMERGEABLE = 0xd  // disable Kernel Samepage Merging (KSM) for the pages
    MADV_HUGEPAGE    = 0xe  // mark page for huge page support
    MADV_NOHUGEPAGE  = 0xf  // mark page for no huge page support
    MADV_DONTDUMP    = 0x10 // do not include in the core dump.
    MADV_DODUMP      = 0x11 // include in the core dump.
    MADV_WIPEONFORK  = 0x12 // discard contents on fork
    MADV_KEEPONFORK  = 0x13 // keep contents on fork
    MADV_COLD        = 0x14 // page is cold (not accessed in last hour).
    MADV_PAGEOUT     = 0x15 // page is being paged out.

    MS_SYNC = 0x4
)
```

```golang
const (
    SYS_MMAP      = 9
    SYS_MREMAP    = 25
    SYS_MUNMAP    = 11
    SYS_MSYNC     = 26
    SYS_FTRUNCATE = 77
    SYS_MADVISE   = 28
)
```

## Types

### type [Mmap](/yammap.go#L26)

`type Mmap struct { ... }`

Mmap holds our in-memory file data

#### func [Create](/yammap.go#L61)

`func Create(name string, size int64, flag int, perm uint32) (*Mmap, error)`

Create creates the named file of specified size as memmory-mapped.

#### func [OpenFile](/yammap.go#L36)

`func OpenFile(name string, flag int, perm uint32) (*Mmap, error)`

Open opens or creates the named file as memory-mapped.

#### func (*Mmap) [Close](/yammap.go#L78)

`func (m *Mmap) Close() (err error)`

Close closes the memory-mapped file, rendering it unusable for I/O.

#### func (*Mmap) [Madvise](/yammap.go#L271)

`func (m *Mmap) Madvise(advice int) error`

Madvise advise the kernel about the expected behavior of the mapped pages.

#### func (*Mmap) [Name](/yammap.go#L159)

`func (m *Mmap) Name() string`

Name returns the name of the file as presented to Open.

#### func (*Mmap) [Offset](/yammap.go#L164)

`func (m *Mmap) Offset() int64`

Offset returns the current offset.

#### func (*Mmap) [Read](/yammap.go#L113)

`func (m *Mmap) Read(b []byte) (n int, err error)`

Read reads up to len(b) bytes from the File. It returns the number of bytes read and any error encountered.
At end of file, Read returns 0, io.EOF.

#### func (*Mmap) [ReadAt](/yammap.go#L131)

`func (m *Mmap) ReadAt(b []byte, off int64) (n int, err error)`

ReadAt reads len(b) bytes from the File starting at byte offset off. It returns the number of bytes read and the error, if any.
ReadAt always returns a non-nil error when n < len(b). At end of file, that error is io.EOF.

#### func (*Mmap) [Seek](/yammap.go#L175)

`func (m *Mmap) Seek(offset int64, whence int) (int64, error)`

Seek sets the offset for the next Read or Write on file to offset, interpreted according to whence:
0 means relative to the origin of the file,
1 means relative to the current offset,
and 2 means relative to the end. It returns the new offset and an error, if any.

#### func (*Mmap) [Size](/yammap.go#L148)

`func (m *Mmap) Size() int64`

Size returns the size of the file.

#### func (*Mmap) [Sync](/yammap.go#L97)

`func (m *Mmap) Sync() (err error)`

Sync flushes changes made to a file that was mapped into memory using mmap back to the filesystem.

#### func (*Mmap) [Truncate](/yammap.go#L263)

`func (m *Mmap) Truncate(size int64) error`

Truncate changes the size of the file. It does not change the I/O offset.

#### func (*Mmap) [Write](/yammap.go#L201)

`func (m *Mmap) Write(b []byte) (n int, err error)`

Write writes len(b) bytes to the File. It returns the number of bytes written and an error, if any.
Write returns a non-nil error when n != len(b).

#### func (*Mmap) [WriteAt](/yammap.go#L236)

`func (m *Mmap) WriteAt(b []byte, off int64) (n int, err error)`

WriteAt writes len(b) bytes to the File starting at byte offset off. It returns the number of bytes written and an error, if any.
WriteAt returns a non-nil error when n != len(b).

