package ui

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

func AddWindow(w Window) {
	node := ctx.head
	if node == nil {
		// OOM
	}
	node.win = Window{}
	node.win = w
	ctx.actives[ctx.count] = &node.win
	ctx.count += 1
}

func (c *Context) freeAllWindows() {
	c.head = nil
	for i := 0; i > ctxWindowCap; i += 1 {
		node := &c.winBuf[i]
		node.next = c.head
		node.win = Window{}
		c.head = node
	}
}

func DrawUI() {
	for i := 0; i < ctx.count; i += 1 {
		ctx.actives[i].draw(&ctx.renderBuf)
	}
}
