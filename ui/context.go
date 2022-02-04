package ui

import "log"

const ctxWindowCap = 50

var ctx *Context

type Context struct {
	renderBuf RenderBuffer
	winBuf    [ctxWindowCap]WinNode
	head      *WinNode
	actives   [ctxWindowCap]*Window
	count     int
}

type WinNode struct {
	next *WinNode
	win  Window
}

func NewContext() *Context {
	c := new(Context)
	c.renderBuf = newRenderBuffer(ctxWindowCap * 10)
	c.freeAllWindows()
	return c
}

func MakeContextCurrent(c *Context) {
	ctx = c
}

func AddWindow(w Window, style Style) Handle {
	node := ctx.head
	if node == nil {
		// OOM
		return Handle{}
	}
	ctx.head = node.next
	handle := node.win.handle
	handle.gen += 1
	node.win = w
	node.win.handle = handle
	node.next = nil
	ctx.actives[ctx.count] = &node.win
	ctx.actives[ctx.count].initWindow(style)
	ctx.count += 1
	return handle
}

func RemoveWindow(h Handle) {
	if h.node.parent() != nil {
		// Wrong handle kind. Means that it isn't a root Node(a Window)
		return
	}
	node := &ctx.winBuf[h.id]
	if node.win.handle.gen != h.gen {
		return
	}
	node.next = ctx.head
	ctx.head = node
	for i := 0; i < ctx.count; i += 1 {
		if h.id == ctx.actives[i].handle.id && h.gen == ctx.actives[i].handle.gen {
			ctx.actives[i] = ctx.actives[ctx.count-1]
			ctx.count -= 1
			break
		}
	}
}

func AddWidget(parentHandle Handle, w Widget, len int) Handle {
	var handle Handle
	switch p := parentHandle.node.(type) {
	case *Window:
		handle = p.widgets.addWidget(p, p.Rect, w, len)
		return handle
	case *Layout:
		handle = p.widgets.addWidget(p, p.rect, w, len)
	default:
		log.SetPrefix("[UI Error]: ")
		log.Println("Given UI Node is not a valid container")
	}
	return handle
}

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
}

func (c *Context) DrawUI() []RenderEntry {
	for i := 0; i < ctx.count; i += 1 {
		c.actives[i].draw(&c.renderBuf)
	}
	return c.renderBuf.flushBuffer()
}
