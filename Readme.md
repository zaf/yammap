# yammap

Package yammap provides an interface to memory mapped files.

WIP - Don't use in production

## Constants

```golang
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

    PROT_READ  = 0x1 // page can be read
    PROT_WRITE = 0x2 // page can be written

    MAP_SHARED   = 0x01 // share changes
    MAP_PRIVATE  = 0x02 // changes are private
    MAP_LOCKED   = 0x2000
    MAP_POPULATE = 0x8000

    MREMAP_MAYMOVE   = 0x1 // may move the mapping
    MREMAP_FIXED     = 0x2 // map at a fixed address
    MREMAP_DONTUNMAP = 0x4 // don't unmap the mapping on close

    SEEK_START   = 0x0 // seek relative to the origin of the file
    SEEK_CURRENT = 0x1 // seek relative to the current offset
    SEEK_END     = 0x2 // seek relative to the end

)
```

## Types

### type [Mmap](/yammap.go#L27)

`type Mmap struct { ... }`

Mmap holds our in-memory file data

#### func [Create](/yammap.go#L93)

`func Create(name string, size int64, flag int, perm uint32) (*Mmap, error)`

Create creates the named file of specified size as memmory-mapped.

#### func [OpenFile](/yammap.go#L73)

`func OpenFile(name string, flag int, perm uint32) (*Mmap, error)`

Open opens or creates the named file as memmory-mapped.

#### func (*Mmap) [Close](/yammap.go#L109)

`func (m *Mmap) Close() (err error)`

Close closes the memory-mapped file, rendering it unusable for I/O.

#### func (*Mmap) [Name](/yammap.go#L174)

`func (m *Mmap) Name() string`

Name returns the name of the file as presented to Open.

#### func (*Mmap) [Offset](/yammap.go#L179)

`func (m *Mmap) Offset() int64`

Offset returns the current offset.

#### func (*Mmap) [Read](/yammap.go#L139)

`func (m *Mmap) Read(b []byte) (n int, err error)`

Read reads up to len(b) bytes from the File. It returns the number of bytes read and any error encountered.
At end of file, Read returns 0, io.EOF.

#### func (*Mmap) [ReadAt](/yammap.go#L152)

`func (m *Mmap) ReadAt(b []byte, off int64) (n int, err error)`

ReadAt reads len(b) bytes from the File starting at byte offset off. It returns the number of bytes read and the error, if any.
ReadAt always returns a non-nil error when n < len(b). At end of file, that error is io.EOF.

#### func (*Mmap) [Seek](/yammap.go#L190)

`func (m *Mmap) Seek(offset int64, whence int) (int64, error)`

Seek sets the offset for the next Read or Write on file to offset, interpreted according to whence:
0 means relative to the origin of the file,
1 means relative to the current offset,
and 2 means relative to the end. It returns the new offset and an error, if any.

#### func (*Mmap) [Size](/yammap.go#L166)

`func (m *Mmap) Size() int64`

Size returns the size of the file.

#### func (*Mmap) [Sync](/yammap.go#L126)

`func (m *Mmap) Sync() (err error)`

Sync flushes changes made to a file that was mapped into memory using mmap back to the filesystem.

#### func (*Mmap) [Truncate](/yammap.go#L254)

`func (m *Mmap) Truncate(size int64) error`

Truncate changes the size of the file. It does not change the I/O offset.

#### func (*Mmap) [Write](/yammap.go#L216)

`func (m *Mmap) Write(b []byte) (n int, err error)`

Write writes len(b) bytes to the File. It returns the number of bytes written and an error, if any.
Write returns a non-nil error when n != len(b).

#### func (*Mmap) [WriteAt](/yammap.go#L236)

`func (m *Mmap) WriteAt(b []byte, off int64) (n int, err error)`

WriteAt writes len(b) bytes to the File starting at byte offset off. It returns the number of bytes written and an error, if any.
WriteAt returns a non-nil error when n != len(b).
