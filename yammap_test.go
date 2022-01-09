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
	"os"
	"strconv"
	"testing"
	"time"
)

// Generate random temp names
func tmpname() string {
	prefix := "yammap_test_"
	dir := os.TempDir()
	rand := uint32(time.Now().UnixNano() + int64(os.Getpid()))
	rand = rand*1664525 + 1013904223
	return dir + "/" + prefix + strconv.Itoa(int(1e9 + rand%1e9))[1:]
}

func TestReadWrite(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(name)
	msg := "Hello World!"
	_, err = m.Write([]byte(msg))
	if err != nil {
		t.Error(err)
	}
	err = m.Sync()
	if err != nil {
		t.Error(err)
	}
	err = m.Close()
	if err != nil {
		t.Error(err)
	}
	m2, err := OpenFile(name, O_RDONLY, 0644)
	if err != nil {
		t.Error(err)
	}
	b := make([]byte, len(msg))
	_, err = m2.Read(b)
	if err != nil {
		t.Error(err)
	}
	err = m2.Close()
	if err != nil {
		t.Error(err)
	}
	if string(b) != msg {
		t.Error("Wrong data")
	}
}

func TestTurncate(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer m.Close()
	defer os.Remove(name)
	newsize := int64(16384)
	err = m.Truncate(newsize)
	if err != nil {
		t.Error(err)
	}
	if newsize != m.Size() {
		t.Error("Wrong size")
	}
	err = m.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestOpenFile(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	m.Close()
	defer os.Remove(name)
	f, err := os.Stat(name)
	if err != nil {
		t.Error(err)
	}
	if f.Size() != int64(pageSize) {
		t.Error("Wrong size")
	}
}

func TestCreate(t *testing.T) {
	name := tmpname()
	m, err := Create(name, int64(pageSize), O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	m.Close()
	defer os.Remove(name)
	f, err := os.Stat(name)
	if err != nil {
		t.Error(err)
	}
	if f.Size() != int64(pageSize) {
		t.Error("Wrong size")
	}
}

func TestName(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer m.Close()
	defer os.Remove(name)
	if name != m.Name() {
		t.Error("Wrong name")
	}
}

func TestSeek(t *testing.T) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer m.Close()
	defer os.Remove(name)
	var position int64 = 1024
	_, err = m.Seek(position, SEEK_START)
	if err != nil {
		t.Error(err)
	}
	if position != m.Offset() {
		t.Error("Wrong offset")
	}
	current := m.Offset()
	_, err = m.Seek(position, SEEK_CURRENT)
	if err != nil {
		t.Error(err)
	}
	if current+position != m.Offset() {
		t.Error("Wrong offset")
	}
	position = -position
	_, err = m.Seek(position, SEEK_END)
	if err != nil {
		t.Error(err)
	}
	if m.Size()+position != m.Offset() {
		t.Error("Wrong offset")
	}
	_, err = m.Seek(1024, SEEK_END)
	if err.Error() != "offset goes beyond the end of file" {
		t.Error("Allowed to seek beyond the end of file")
	}
	_, err = m.Seek(-1024, SEEK_START)
	if err.Error() != "negative position" {
		t.Error("Allowed to seek with negative position")
	}
}

func BenchmarkWrite(b *testing.B) {
	name := tmpname()
	m, err := OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer m.Close()
	defer os.Remove(name)
	data := make([]byte, 512)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Write(data)
		//m.Sync()
	}
}

func BenchmarkOsWrite(b *testing.B) {
	name := tmpname()
	f, err := os.OpenFile(name, O_RDWR|O_CREATE, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(name)
	data := make([]byte, 512)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Write(data)
		//f.Sync()
	}
}
