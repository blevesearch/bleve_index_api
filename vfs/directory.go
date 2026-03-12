//  Copyright (c) 2025 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vfs

import (
	"io"
	"io/fs"
	"os"
	"time"
)

// Directory abstracts filesystem operations for storage backends.
// This interface enables segments to use different storage backends
// (local filesystem, S3, GCS, Azure Blob, etc.) without modification.
//
// Implementations must be safe for concurrent use by multiple goroutines.
//
// This is a low-level interface used by index implementations (like Scorch).
// Most users should use higher-level index APIs instead of working with
// Directory implementations directly.
type Directory interface {
	// Open opens the named file for reading. The caller must close the
	// returned ReadCloser when done.
	Open(name string) (io.ReadCloser, error)

	// Create creates or truncates the named file for writing. If the file
	// already exists, it is truncated. The caller must close the returned
	// WriteCloser when done.
	Create(name string) (WriteCloser, error)

	// Remove removes the named file.
	Remove(name string) error

	// Rename renames (moves) oldpath to newpath. If newpath already exists
	// and is not a directory, Rename replaces it.
	Rename(oldpath, newpath string) error

	// Stat returns FileInfo describing the named file.
	Stat(name string) (FileInfo, error)

	// ReadDir reads the named directory and returns a list of directory entries.
	ReadDir(name string) ([]FileInfo, error)

	// MkdirAll creates a directory named path, along with any necessary
	// parents, and returns nil, or else returns an error.
	MkdirAll(path string, perm fs.FileMode) error

	// Sync commits the current contents of the directory to stable storage.
	// This is a hint that implementations can use to optimize durability.
	Sync() error

	// Lock acquires an exclusive lock on the directory. This is used to
	// prevent multiple processes from opening the same index simultaneously.
	// Must be called before any other operations.
	Lock() error

	// Unlock releases the lock acquired by Lock.
	Unlock() error

	// OpenAt opens the named file for random access reading. This is used
	// for memory-mapped segments and other use cases requiring io.ReaderAt.
	// The caller must close the returned ReaderAtCloser when done.
	OpenAt(name string) (ReaderAtCloser, error)
}

// FileInfo describes a file and is returned by Stat and ReadDir.
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes
	Mode() fs.FileMode  // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
}

// WriteCloser extends io.WriteCloser with a Sync method for durability.
type WriteCloser interface {
	io.WriteCloser
	Sync() error
}

// ReaderAtCloser extends io.ReaderAt and io.Closer with a method to get
// the underlying file descriptor for memory-mapped operations.
// This enables zero-copy segment reading through mmap.
type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer

	// AsFd returns the underlying file descriptor if available, for mmap support.
	// Returns 0 if not backed by a real file (e.g., in-memory buffer).
	// Callers should check the return value before attempting mmap.
	AsFd() uintptr
}

// FileReaderAtCloser wraps an *os.File to implement ReaderAtCloser.
// This is a convenience helper for implementations that use os.File.
type FileReaderAtCloser struct {
	*os.File
}

// AsFd returns the underlying file descriptor for mmap operations.
func (f *FileReaderAtCloser) AsFd() uintptr {
	return f.File.Fd()
}

// NewFileReaderAtCloser wraps an *os.File in a ReaderAtCloser.
func NewFileReaderAtCloser(f *os.File) ReaderAtCloser {
	return &FileReaderAtCloser{File: f}
}
