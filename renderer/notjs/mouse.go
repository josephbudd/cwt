package notjs

import "syscall/js"

// SetOnMouseEnter sets an element's onmouseenter.
func (notjs *NotJS) SetOnMouseEnter(element js.Value, cb js.Func) {
	element.Set("onmouseenter", cb)
}

// SetOnMouseLeave sets an element's onmouseleave.
func (notjs *NotJS) SetOnMouseLeave(element js.Value, cb js.Func) {
	element.Set("onmouseleave", cb)
}

// SetOnMouseDown sets an element's onmousedown.
func (notjs *NotJS) SetOnMouseDown(element js.Value, cb js.Func) {
	element.Set("onmousedown", cb)
}

// SetOnMouseUp sets an element's onmouseup.
func (notjs *NotJS) SetOnMouseUp(element js.Value, cb js.Func) {
	element.Set("onmouseup", cb)
}
