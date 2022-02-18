package ui

import (
	"log"
)

// Capacity of each context.
// Might be more flexible to take a capacity input
const ctxWindowCap = 50

// Only one context can be active at the moment.
// Not sure why one would want multiple
var ctx *Context

type Context struct {
	renderBuf    renderBuffer
	winBuf       [ctxWindowCap]winNode
	head         *winNode
	actives      [ctxWindowCap]*Window
	currentFocus int
	count        int
	input        inputData

	cursorShapeCallback func(s CursorShape)
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
	c.renderBuf = newRenderBuffer(1000)
	c.freeAllWindows()
	return c
}

func MakeContextCurrent(c *Context) {
	ctx = c
}

func (c *Context) SetCursorShapeCallback(cb func(CursorShape)) {
	c.cursorShapeCallback = cb
}

// Add a window (a copy of the one given as argument)
// to the current context.
//
// WARNING: A context must be set to current before trying to add windows
func AddWindow(w Window) WinHandle {
	// Pop the head of the list
	node := ctx.head
	if node == nil {
		log.SetPrefix("[UI Fatal Error]: ")
		log.Fatalln("Current Context is out of memory")
	}
	// Set the next node at the top of the list
	ctx.head = node.next
	handle := WinHandle{
		id:  node.win.handle.id,
		gen: node.win.handle.gen + 1,
	}
	node.win = w
	node.win.handle = handle
	// Set to nil for safety.
	node.next = nil

	// Push the new window onto the active
	// window array and initialize it
	ctx.actives[ctx.count] = &node.win
	FocusWindow(handle)
	ctx.actives[ctx.count].initWindow()
	ctx.count += 1
	return handle
}

// Delete a Window with the given handle from the current context.
//
// Note: Also removes all the child nodes
func DeleteWindow(h WinHandle) {
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

func getWindow(h WinHandle) *Window {
	node := &ctx.winBuf[h.id]
	if node.win.handle.gen != h.gen {
		return nil
	}
	return &node.win
}

func FocusWindow(h WinHandle) {
	for i := 0; i < ctx.count; i += 1 {
		win := ctx.actives[i]
		if win.handle.id == h.id && win.handle.gen == h.gen {
			win.zIndex = 0
		} else {
			win.zIndex += 1
		}
	}
}

// Function used internally!
// Reset the memory and all the handles
func (c *Context) freeAllWindows() {
	c.head = nil
	for i := 0; i < ctxWindowCap; i += 1 {
		node := &c.winBuf[i]
		node.next = c.head
		node.win = Window{
			handle: WinHandle{id: i, gen: 0},
		}
		c.head = node
	}
	c.count = 0
}

func (c *Context) UpdateUI(data Input) {
	// Update all the input supplied
	{
		c.input.previousmPos = c.input.mPos
		c.input.previousmLeft = c.input.mLeft
		c.input.mPos = data.MPos
		c.input.mLeft = data.MLeft
		c.input.previousKeys = c.input.keys
		c.input.keys[keyEsc] = data.Esc
		c.input.keys[keyEnter] = data.Enter
		c.input.keys[keyDelete] = data.Del
		c.input.keys[keyCtlr] = data.Ctrl
		c.input.keys[keyShift] = data.Shift
		c.input.keys[keySpace] = data.Space
		c.input.keys[keyTab] = data.Tab
		c.input.keys[keyUp] = data.Up
		c.input.keys[keyDown] = data.Down
		c.input.keys[keyLeft] = data.Left
		c.input.keys[keyRight] = data.Right

		for i := range c.input.keyCounts {
			if c.input.keys[i] {
				c.input.keyCounts[i] += 1
			}
			if !c.input.keys[i] && c.input.previousKeys[i] {
				c.input.keyCounts[i] = 0
			}
		}
	}
	for i := 0; i < ctx.count; i += 1 {
		c.actives[i].update()
	}
	c.input.pressedCharsCount = 0
}

func (c *Context) DrawUI() []RenderEntry {
	for i := 0; i < ctx.count; i += 1 {
		c.actives[i].draw(&c.renderBuf)
	}
	return c.renderBuf.flushBuffer()
}
