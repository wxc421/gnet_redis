package poolx

const defaultBufferSize = 4 * 1024

var SliceBuffer = NewPoolNormal[[]byte](func() []byte {
	return make([]byte, defaultBufferSize)
})
