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
	m.Close()
	defer os.Remove(name)
	f, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size() != int64(pageSize) {
		t.Fatal("wrong size of created file")
	}
}

func TestCreate(t *testing.T) {
	name := tmpname()
	m, err := Create(name, int64(pageSize), O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	m.Close()
	defer os.Remove(name)
	f, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if f.Size() != int64(pageSize) {
		t.Fatal("wrong size of created file")
	}
}

func TestTurncate(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	newsize := int64(16384)
	err = m.Truncate(newsize)
	if err != nil {
		t.Fatal(err)
	}
	if newsize != m.Size() {
		t.Error("wrong size")
	}
	err = m.Close()
	if err != nil {
		t.Fatal(err)
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
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
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
	msg := rndmessage(pageSize * 2)
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
	msg := rndmessage(pageSize * 2)
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
}

func BenchmarkWrite(b *testing.B) {
	testSize := pageSize / 8
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	data := rndmessage(testSize)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Write(data)
	}
}

func BenchmarkOSWrite(b *testing.B) {
	testSize := pageSize / 8
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
	}
}

func BenchmarkRead(b *testing.B) {
	testSize := pageSize / 8
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
	data := make([]byte, testSize)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Read(data)
	}
}

func BenchmarkOSRead(b *testing.B) {
	testSize := pageSize / 8
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
	}
}
