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
)
