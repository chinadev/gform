package gform

import (
    "w32"
    "w32/user32"
)

type W32Control struct {
    ControlBase

    originalWndProc uintptr
    isMouseLeft bool
}

func (this *W32Control) init(className string, parent Controller, exstyle, style uint) {
    this.hwnd = CreateWindow(className, parent, exstyle, style)
    if this.hwnd == 0 {
        panic("cannot create window for " + className)
    }
    this.isMouseLeft = true
    this.originalWndProc = user32.SetWindowLongPtr(this.hwnd, w32.GWLP_WNDPROC, GeneralWndprocCallBack)
    this.ControlBase.init(parent)
}

func (this *W32Control) attach(parent Controller, dlgItemID int) {
    if parent == nil {
        panic("parent cannot be nil")
    }

    if this.hwnd = user32.GetDlgItem(parent.Handle(), dlgItemID); this.hwnd == 0 {
        panic("hwnd cannot be nil")
    }

    this.isMouseLeft = true
    this.originalWndProc = user32.SetWindowLongPtr(this.hwnd, w32.GWLP_WNDPROC, GeneralWndprocCallBack)
    this.ControlBase.init(parent)
}

func (this *W32Control) WndProc(msg uint, wparam, lparam uintptr) uintptr {
    switch msg {
    case w32.WM_CREATE:
        internalTrackMouseEvent(this.hwnd)
        this.onCreate.Fire(this)
    case w32.WM_MOUSEMOVE:
        if this.isMouseLeft {
            this.onMouseHover.Fire(this)
            internalTrackMouseEvent(this.hwnd)
            this.isMouseLeft = false
        }
    case w32.WM_MOUSELEAVE:
        this.onMouseLeave.Fire(this)
        this.isMouseLeft = true
    }
    return user32.CallWindowProc(this.originalWndProc, this.hwnd, msg, wparam, lparam)
}
