//go:build linux && (amd64 || arm64 || mips64 || mips64le || ppc64 || ppc64le || s390x)
// +build linux
// +build amd64 arm64 mips64 mips64le ppc64 ppc64le s390x

/*
	Copyright (C) 2022, Lefteris Zafiris <zaf@fastmail.com>
	This program is free software, distributed under the terms of
	the GNU GPL v3 License. See the LICENSE file
	at the top of the source tree.
*/

package yammap

const (
	SYS_MMAP   = 9
	SYS_MREMAP = 25
	SYS_MUNMAP = 11
	SYS_MSYNC  = 26

	maxSize = 0xFFFFFFFFFFFF // maximum allocation size, 2^48 bytes for x86_64
)
