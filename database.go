package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Database interface {
	Put(key string, value string) bool
	Get(key string) (string, bool)
	Delete(key string) bool
}

// MapDatabase is a database that is backed by a simple map as a data structure
// only for testing purposes
type MapDatabase struct {
	data map[string]string
}

func NewMapDatabase() *MapDatabase {
	return &MapDatabase{data: make(map[string]string)}
}

func (kv *MapDatabase) Put(key string, value string) bool {
	_, contains := kv.data[key]
	kv.data[key] = value
	return contains
}

func (kv *MapDatabase) Get(key string) (string, bool) {
	result, contains := kv.data[key]
	return result, contains
}

func (kv *MapDatabase) Delete(key string) bool {
	_, contains := kv.data[key]
	delete(kv.data, key)
	return contains
}

type KvDatabase struct {
	backingFile *os.File
}

func NewKvDatabase(fileName string) *KvDatabase {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	return &KvDatabase{backingFile: f}
}

func (kv *KvDatabase) getValue(key string) (string, int, error) {
	kv.backingFile.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(kv.backingFile)

	var currentLine int
	// Scan through the lines of the file
	for scanner.Scan() {
		lineText := scanner.Text()
		separator := strings.Index(lineText, " ")
		keyLength, err := strconv.Atoi(lineText[:separator])
		if err != nil {
			return "", currentLine, err
		}

		lineKey := lineText[separator+1 : separator+keyLength+1]
		lineValue := lineText[separator+1+keyLength:]

		if key == lineKey {
			return lineValue, currentLine, nil
		}

		currentLine += 1
	}

	return "", currentLine, ErrNoSuchKey
}

func (kv *KvDatabase) writeLine(line int, key string, value string) error {
	kv.backingFile.Seek(0, io.SeekStart)
	dataText, err := io.ReadAll(kv.backingFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(dataText), "\n")
	if line >= len(lines) {
		lines = append(lines, fmt.Sprintf("%d %s%s", len(key), key, value))
	} else {
		lines[line] = fmt.Sprintf("%d %s%s", len(key), key, value)
	}

	kv.backingFile.Seek(0, io.SeekStart)
	kv.backingFile.Truncate(0)
	if _, err = kv.backingFile.Write([]byte(strings.Join(lines, "\n"))); err != nil {
		return err
	}

	return nil
}

func (kv *KvDatabase) deleteLine(line int) error {
	kv.backingFile.Seek(0, io.SeekStart)
	dataText, err := io.ReadAll(kv.backingFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(dataText), "\n")
	if line < 0 || line >= len(lines) {
		return nil
	}

	lines = append(lines[:line], lines[line+1:]...)

	kv.backingFile.Seek(0, io.SeekStart)
	kv.backingFile.Truncate(0)
	if _, err = kv.backingFile.Write([]byte(strings.Join(lines, "\n"))); err != nil {
		return err
	}

	return nil
}

func (kv *KvDatabase) Put(key string, value string) bool {
	_, line, err := kv.getValue(key)
	exists := true
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			exists = false
		} else {
			panic(err)
		}
	}

	kv.writeLine(line, key, value)
	return !exists
}

func (kv *KvDatabase) Get(key string) (string, bool) {
	value, _, err := kv.getValue(key)
	exists := true
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			exists = false
		} else {
			panic(err)
		}
	}

	return value, exists
}

func (kv *KvDatabase) Delete(key string) bool {
	_, line, err := kv.getValue(key)
	exists := true
	if err != nil {
		if errors.Is(err, ErrNoSuchKey) {
			exists = false
		} else {
			panic(err)
		}
	}

	if err := kv.deleteLine(line); err != nil {
		panic(err)
	}

	return exists
}
