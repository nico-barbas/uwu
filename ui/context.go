package ui

import "log"

// Capacity of each context.
// Might be more flexible to take a capacity input
const ctxWindowCap = 50

// Only one context can be active at the moment.
// Not sure why one would want multiple
var ctx *Context

type Context struct {
	renderBuf renderBuffer
	winBuf    [ctxWindowCap]winNode
	head      *winNode
	actives   [ctxWindowCap]*Window
	count     int
	input     inputData
}

// Internal data used for the window free list.
// It has to be stored as a header like that
// since it isn't possible cast a pointer to an
// arbitrary pointer type in Go.
type winNode struct {
	next *winNode
	win  Window
}

// Allocate a new Context and return it
func NewContext() *Context {
	c := new(Context)
	c.renderBuf = newRenderBuffer(ctxWindowCap * 10)
	c.freeAllWindows()
	return c
}

func MakeContextCurrent(c *Context) {
	ctx = c
}

// Add a window (a copy of the one given as argument)
// to the current context.
//
// WARNING: A context must be set to current before trying to add windows
func AddWindow(w Window) Handle {
	// Pop the head of the list
	node := ctx.head
	if node == nil {
		log.SetPrefix("[UI Fatal Error]: ")
		log.Fatalln("Current Context is out of memory")
	}
	// Set the next node at the top of the list
	ctx.head = node.next
	handle := Handle{
		node: node.win.handle.node,
		id:   node.win.handle.id,
		gen:  node.win.handle.gen + 1,
	}
	node.win = w
	node.win.handle = handle
	// Set to nil for safety.
	node.next = nil

	// Push the new window onto the active
	// window array and initialize it
	ctx.actives[ctx.count] = &node.win
	ctx.actives[ctx.count].initWindow()
	ctx.count += 1
	return handle
}

// Delete a Window with the given handle from the current context.
//
// Note: Also removes all the child nodes
func DeleteWindow(h Handle) {
	// Wrong handle kind. Mean that it isn't a root Node(a Window)
	if h.node.parent() != nil {
		log.SetPrefix("[UI Error]: ")
		log.Println("Given Handle does not refer to a Window")
		return
	}

	// Push the node on top of the free list
	node := &ctx.winBuf[h.id]
	if node.win.handle.gen != h.gen {
		return
	}
	node.next = ctx.head
	ctx.head = node

	// Linear search in the active windows array and remove the window
	for i := 0; i < ctx.count; i += 1 {
		if h.id == ctx.actives[i].handle.id && h.gen == ctx.actives[i].handle.gen {
			ctx.actives[i] = ctx.actives[ctx.count-1]
			ctx.count -= 1
			break
		}
	}
}

// Try to add the given widget as a child of the Node (referenced by the Handle)
// Can fail if the Node is not a valid receiver.
//
// Valid receivers are Windows and Layouts.
func AddWidget(parentHandle Handle, w Widget, len int) Handle {
	var handle Handle
	switch p := parentHandle.node.(type) {
	case *Window:
		handle = p.widgets.addWidget(p, p.activeRect, w, len)
	case *Layout:
		handle = p.widgets.addWidget(p, p.rect, w, len)
	default:
		log.SetPrefix("[UI Error]: ")
		log.Println("Given UI Node is not a valid container")
	}
	return handle
}

// Function used internally!
// Reset the memory and all the handles
func (c *Context) freeAllWindows() {
	c.head = nil
	for i := 0; i < ctxWindowCap; i += 1 {
		node := &c.winBuf[i]
		node.next = c.head
		node.win = Window{
			handle: Handle{node: &node.win, id: i, gen: 0},
		}
		c.head = node
	}
	c.count = 0
}

func (c *Context) DrawUI() []RenderEntry {
	for i := 0; i < ctx.count; i += 1 {
		c.actives[i].draw(&c.renderBuf)
	}
	return c.renderBuf.flushBuffer()
}
