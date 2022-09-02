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
	"bytes"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate random data
func rndmessage(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return b
}

// Generate random temp names
func tmpname() string {
	prefix := "yammap_test_"
	dir := os.TempDir()
	rand := uint32(time.Now().UnixNano() + int64(os.Getpid()))
	rand = rand*1664525 + 1013904223
	return dir + "/" + prefix + strconv.Itoa(int(1e9 + rand%1e9))[1:]
}

// Generate files with random data
func rndfile(size int) (string, error) {
	name := tmpname()
	m, err := Create(name, int64(size), O_RDWR|O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer m.Close()
	msg := rndmessage(size)
	_, err = m.Write(msg)
	if err != nil {
		os.Remove(name)
		return "", err
	}
	err = m.Sync()
	if err != nil {
		os.Remove(name)
		return "", err
	}
	return name, nil
}

func TestOpenFile(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	size := m.Size()
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size() != size {
		t.Fatal("wrong size of created file")
	}
}

func TestCreate(t *testing.T) {
	name := tmpname()
	m, err := Create(name, int64(os.Getpagesize()), O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	size := m.Size()
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size() != size {
		t.Fatal("wrong size of created file")
	}
}

func TestMadvise(t *testing.T) {
	name := tmpname()
	m, err := Create(name, int64(os.Getpagesize()), O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	err = m.Madvise(MADV_SEQUENTIAL)
	if err != nil {
		t.Fatal(err)
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTruncate(t *testing.T) {
	name := tmpname()
	msg := rndmessage(os.Getpagesize() * 2)
	m, err := Create(name, int64(len(msg)), O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	n, err := m.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(msg) {
		t.Error("wrong number of bytes written")
	}
	newsize := int64(os.Getpagesize())
	err = m.Truncate(newsize)
	if err != nil {
		t.Fatal(err)
	}
	if newsize != m.Size() {
		t.Error("wrong size when shrinking")
	}
	newsize = int64(4 * os.Getpagesize())
	err = m.Truncate(newsize)
	if err != nil {
		t.Fatal(err)
	}
	if newsize != m.Size() {
		t.Error("wrong size when growing")
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size() != newsize {
		t.Error("wrong file size after closing")
	}
}

func TestTruncateToZero(t *testing.T) {
	name, err := rndfile(os.Getpagesize())
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	m, err := OpenFile(name, O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = m.Truncate(0)
	if err != nil {
		t.Fatal(err)
	}
	if m.Size() != 0 {
		t.Error("wrong size when truncating to zero")
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size() != 0 {
		t.Error("wrong file size after truncating to zero")
	}
}

func TestName(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	if name != m.Name() {
		t.Error("wrong name of created file")
	}
}

func TestSeek(t *testing.T) {
	name := tmpname()
	m, err := Create(name, int64(os.Getpagesize()), O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	var position int64 = 1024
	_, err = m.Seek(position, SEEK_SET)
	if err != nil {
		t.Fatal(err)
	}
	if position != m.Offset() {
		t.Error("wrong offset")
	}
	current := m.Offset()
	_, err = m.Seek(position, SEEK_CUR)
	if err != nil {
		t.Fatal(err)
	}
	if current+position != m.Offset() {
		t.Error("wrong offset")
	}
	position = -position
	_, err = m.Seek(position, SEEK_END)
	if err != nil {
		t.Fatal(err)
	}
	if m.Size()+position != m.Offset() {
		t.Error("wrong offset")
	}
	_, err = m.Seek(1024, SEEK_END)
	if err.Error() != "offset goes beyond the end of file" {
		t.Error("allowed to seek beyond the end of file")
	}
	_, err = m.Seek(-1024, SEEK_SET)
	if err.Error() != "negative position" {
		t.Error("allowed to seek with negative position")
	}
}

func TestReadWrite(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	msg := rndmessage(os.Getpagesize() * 2)
	n, err := m.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(msg) {
		t.Error("wrong number of bytes written")
	}
	if m.offset != int64(len(msg)) {
		t.Error("wrong offset after write")
	}
	err = m.Sync()
	if err != nil {
		t.Fatal(err)
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
	m2, err := OpenFile(name, O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, len(msg))
	n, err = m2.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(b) {
		t.Error("wrong number of bytes read")
	}
	if m2.offset != int64(len(msg)) {
		t.Error("wrong offset after read")
	}
	err = m2.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, msg) {
		t.Error("wrong data read")
	}
}

func TestReadAtWriteAt(t *testing.T) {
	name := tmpname()
	offset := int64(512)
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	msg := rndmessage(os.Getpagesize() * 2)
	n, err := m.WriteAt([]byte(msg), offset)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(msg) {
		t.Error("wrong number of bytes written")
	}
	err = m.Sync()
	if err != nil {
		t.Fatal(err)
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
	m2, err := OpenFile(name, O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, len(msg))
	n, err = m2.ReadAt(b, offset)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(b) {
		t.Error("wrong number of bytes read")
	}
	err = m2.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, msg) {
		t.Error("wrong data read")
	}
	m3, err := OpenFile(name, O_RDWR|O_APPEND, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = m3.WriteAt([]byte(msg), offset)
	if err.Error() != "invalid use of WriteAt on file opened with O_APPEND" {
		t.Error("allowed to write at offset in append mode")
	}
}

func TestAppend(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(name)
	msg := rndmessage(os.Getpagesize() * 2)
	n, err := m.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(msg) {
		t.Error("wrong number of bytes written")
	}
	if m.offset != int64(len(msg)) {
		t.Error("wrong offset after write")
	}
	err = m.Sync()
	if err != nil {
		t.Fatal(err)
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}

	m2, err := OpenFile(name, O_RDWR|O_APPEND, 0644)
	if err != nil {
		t.Fatal(err)
	}
	n, err = m2.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(msg) {
		t.Error("wrong number of bytes written")
	}
	err = m2.Sync()
	if err != nil {
		t.Fatal(err)
	}
	err = m2.Close()
	if err != nil {
		t.Fatal(err)
	}

	m3, err := OpenFile(name, O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 2*len(msg))
	n, err = m3.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(b) {
		t.Error("wrong number of bytes read")
	}
	if m3.offset != int64(2*len(msg)) {
		t.Error("wrong offset after read")
	}
	err = m3.Close()
	if err != nil {
		t.Fatal(err)
	}
	msg = append(msg, msg...)
	if !bytes.Equal(b, msg) {
		t.Error("wrong data read")
	}
}

func TestBigFiles(t *testing.T) {
	var size int64 = 1 << 31 // 2GB
	msg := rndmessage(os.Getpagesize())
	name := tmpname()
	m, err := Create(name, size, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal("Failed to create large file", err)
	}
	defer m.Close()
	defer os.Remove(name)
	_, err = m.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	err = m.Sync()
	if err != nil {
		t.Fatal(err)
	}
	_, err = m.Seek(size-int64(len(msg)), 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = m.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	err = m.Sync()
	if err != nil {
		t.Fatal(err)
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}
	m2, err := OpenFile(name, O_RDONLY, 0644)
	if err != nil {
		t.Fatal("Failed to open large file", err)
	}
	err = m2.Close()
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(name)
}

func BenchmarkWrite(b *testing.B) {
	testSize := os.Getpagesize() * 1024
	name := tmpname()
	m, err := Create(name, int64(testSize), O_RDWR|O_CREATE, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	m.Madvise(MADV_SEQUENTIAL)
	data := rndmessage(testSize)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Write(data)
		m.Seek(0, 0)
	}
}

func BenchmarkOSWrite(b *testing.B) {
	testSize := os.Getpagesize() * 1024
	name := tmpname()
	f, err := os.OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(name)
	data := rndmessage(testSize)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Write(data)
		f.Seek(0, 0)
	}
}

func BenchmarkRead(b *testing.B) {
	testSize := os.Getpagesize() * 1024
	name, err := rndfile(testSize)
	if err != nil {
		b.Fatal(err)
	}
	m, err := OpenFile(name, O_RDONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	m.Madvise(MADV_SEQUENTIAL)
	data := make([]byte, testSize)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Read(data)
		m.Seek(0, 0)
	}
}

func BenchmarkOSRead(b *testing.B) {
	testSize := os.Getpagesize() * 1024
	name, err := rndfile(testSize)
	if err != nil {
		b.Fatal(err)
	}
	f, err := os.OpenFile(name, O_RDONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(name)
	data := make([]byte, testSize)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Read(data)
		f.Seek(0, 0)
	}
}
