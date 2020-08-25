package base

import (
	"bytes"
	"container/list"
)

type BufferChunk struct {
	data []byte
	next *BufferChunk
}

type MultiChunkBuffer struct {
	curr *BufferChunk
	pos  int
	last *BufferChunk
	eof  bool
}

func NewMultiChunkBuffer() *MultiChunkBuffer {
	return &MultiChunkBuffer{
		curr: nil,
		pos:  0,
		last: nil,
		eof:  false,
	}
}

func (buffer *MultiChunkBuffer) Write(data []byte) {
	if buffer.eof {
		panic("MultiChunkBuffer with eof!")
	}
	if len(data) == 0 {
		return
	}
	chunk := &BufferChunk{data, nil}
	if buffer.curr == nil {
		buffer.curr = chunk
		buffer.pos = 0
	} else {
		buffer.last.next = chunk
	}
	buffer.last = chunk
}

func (buffer *MultiChunkBuffer) WriteEOF() {
	buffer.eof = false
}

func (buffer *MultiChunkBuffer) ReadLine() []byte {
	if buffer.curr == nil {
		return nil
	}
	i := bytes.IndexByte(buffer.curr.data[buffer.pos:], '\n')
	if i >= 0 {
		lineBytes := buffer.curr.data[buffer.pos : buffer.pos+i]
		lineBytes = trimLine(lineBytes)
		buffer.move(buffer.curr, buffer.pos+i+1)
		return lineBytes
	}
	chunkList := list.New()
	chunk := buffer.curr.next
	for chunk != nil {
		i = bytes.IndexByte(chunk.data, '\n')
		if i >= 0 {
			lineBuf := bytes.NewBuffer(nil)
			lineBuf.Write(buffer.curr.data[buffer.pos:])
			for e := chunkList.Front(); e != nil; e = e.Next() {
				c, _ := e.Value.(*BufferChunk)
				lineBuf.Write(c.data)
			}
			lineBuf.Write(chunk.data[:i])
			lineBytes := trimLine(lineBuf.Bytes())
			buffer.move(chunk, i+1)
			return lineBytes
		}
		chunkList.PushBack(chunk)
		chunk = chunk.next
	}
	if buffer.eof {
		lineBuf := bytes.NewBuffer(nil)
		lineBuf.Write(buffer.curr.data[buffer.pos:])
		for e := chunkList.Front(); e != nil; e = e.Next() {
			c, _ := e.Value.(*BufferChunk)
			lineBuf.Write(c.data)
		}
		lineBytes := trimLine(lineBuf.Bytes())
		buffer.move(nil, 0)
		return lineBytes
	}
	return nil
}

func (buffer *MultiChunkBuffer) move(expectChunk *BufferChunk, expectPos int) {
	if expectPos >= len(expectChunk.data) {
		buffer.curr = expectChunk.next
		buffer.pos = 0
	} else {
		buffer.curr = expectChunk
		buffer.pos = expectPos
	}
	if buffer.curr == nil {
		buffer.last = nil
	}
}

func trimLine(lineBytes []byte) []byte {
	if len(lineBytes) > 0 && lineBytes[len(lineBytes)-1] == '\r' {
		lineBytes = lineBytes[:len(lineBytes)-1]
	}
	return lineBytes
}
