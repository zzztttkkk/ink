package utils

import (
	"bytes"
	"sync"
)

type _BytesBufferPoolNamespace struct{}

var (
	bufPool = sync.Pool{New: func() any { return bytes.NewBuffer(nil) }}
)

func (_BytesBufferPoolNamespace) Get() *bytes.Buffer { return bufPool.Get().(*bytes.Buffer) }

func (_BytesBufferPoolNamespace) Put(v *bytes.Buffer) { bufPool.Put(v) }

var BytesBufferPool _BytesBufferPoolNamespace
