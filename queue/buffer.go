package queue

// import (
// 	"runtime"

// 	"gitlab.com/around25/products/matching-engine/engine"
// )

// // Buffer structure to process messages in a buffered ring
// type Buffer struct {
// 	// The padding members 1 to 5 below are here to ensure each item is on a separate cache line.
// 	// This prevents false sharing and hence improves performance.
// 	padding1           [8]uint64
// 	lastCommittedIndex uint64
// 	padding2           [8]uint64
// 	nextFreeIndex      uint64
// 	padding3           [8]uint64
// 	readerIndex        uint64
// 	padding4           [8]uint64
// 	contents           []engine.Event
// 	padding5           [8]uint64
// 	queueSize          uint64
// 	indexMask          uint64
// }

// // NewBuffer ring
// func NewBuffer(queueSize uint64) *Buffer {
// 	return &Buffer{
// 		lastCommittedIndex: 0, nextFreeIndex: 1, readerIndex: 1, queueSize: queueSize, indexMask: queueSize - 1,
// 		contents: make([]engine.Event, queueSize),
// 	}
// }

// // Write a new value to the ring buffer
// func (buffer *Buffer) Write(value engine.Event) {
// 	buffer.nextFreeIndex++
// 	var myIndex = buffer.nextFreeIndex - 1
// 	//Wait for reader to catch up, so we don't clobber a slot which it is (or will be) reading
// 	for myIndex > (buffer.readerIndex + buffer.queueSize - 2) {
// 		runtime.Gosched()
// 	}
// 	//Write the item into it's slot
// 	buffer.contents[myIndex&buffer.indexMask] = value
// 	//Increment the lastCommittedIndex so the item is available for reading
// 	buffer.lastCommittedIndex = myIndex
// 	// for !atomic.CompareAndSwapUint64(&self.lastCommittedIndex, myIndex-1, myIndex) {
// 	// 	runtime.Gosched()
// 	// }
// }

// // Read next element from the buffer
// func (buffer *Buffer) Read() engine.Event {
// 	buffer.readerIndex++
// 	var myIndex = buffer.readerIndex - 1
// 	//If reader has out-run writer, wait for a value to be committed
// 	for myIndex > buffer.lastCommittedIndex {
// 		runtime.Gosched()
// 	}
// 	return buffer.contents[myIndex&buffer.indexMask]
// }
