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
	"testing"
)

func TestReadWrite(t *testing.T) {
	m, err := OpenFile("/tmp/test.txt", O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove("/tmp/test.txt")
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
	m2, err := OpenFile("/tmp/test.txt", O_RDONLY, 0644)
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
	m, err := OpenFile("/tmp/test.txt", O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer m.Close()
	defer os.Remove("/tmp/test.txt")
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

func TestSeek(t *testing.T) {
	m, err := OpenFile("/tmp/test.txt", O_RDWR|O_CREATE, 0644)
	if err != nil {
		t.Error(err)
	}
	defer m.Close()
	defer os.Remove("/tmp/test.txt")
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
}
