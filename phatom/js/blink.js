+ function() {
	var blink = {
			crypto: {},
		} // namesapce

	window.blink = blink

	Number.prototype.byte = Number.prototype.uint8 = function() {
		return this & 255
	}
}()