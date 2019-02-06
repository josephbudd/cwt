package notjs

import "syscall/js"

// GetEventTarget gets an event's target attribute which is an html element.
func (notjs *NotJS) GetEventTarget(event js.Value) js.Value {
	return event.Get("target")
}

// SetOnClick sets an element's onclick.
func (notjs *NotJS) SetOnClick(element js.Value, cb js.Callback) {
	element.Set("onclick", cb)
}

// SetOnChange sets an element's onchange.
func (notjs *NotJS) SetOnChange(element js.Value, cb js.Callback) {
	element.Set("onchange", cb)
}

// SetOnScroll sets an element's onscroll.
func (notjs *NotJS) SetOnScroll(element js.Value, cb js.Callback) {
	element.Set("onscroll", cb)
}

// SetOnMouseEnter sets an element's onmouseenter.
func (notjs *NotJS) SetOnMouseEnter(element js.Value, cb js.Callback) {
	element.Set("onmouseenter", cb)
}

// SetOnMouseLeave sets an element's onmouseleave.
func (notjs *NotJS) SetOnMouseLeave(element js.Value, cb js.Callback) {
	element.Set("onmouseleave", cb)
}

// SetOnMouseDown sets an element's onmousedown.
func (notjs *NotJS) SetOnMouseDown(element js.Value, cb js.Callback) {
	element.Set("onmousedown", cb)
}

// SetOnMouseUp sets an element's onmouseup.
func (notjs *NotJS) SetOnMouseUp(element js.Value, cb js.Callback) {
	element.Set("onmouseup", cb)
}
