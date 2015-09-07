+ function() {
	var blink = {
			crypto: {},
		} // namesapce

	window.blink = blink

	Number.prototype.byte = Number.prototype.uint8 = function() {
		return this & 255
	}

	Number.prototype.uint32 = function() {
		var a = new Uint32Array([this])
		return a[0]
	}

	Array.new = function(size, def) {
		def = def || 0
		size = size || 0
		var a = Array(size)
		for (var i = 0; i < a.length; i++) {
			a[i] = def
		}
		return a
	};

	Array.prototype.copy = function(a, begin) {
		var index = 0,
			i = begin || 0
		for (; i < this.length; i++, index++) {
			if (index == a.length) return index
			this[i] = a[index]
		}
		return index
	};

	Array.prototype.equal = function(a) {
		if (a.length != this.length) return false
		for (var i = 0; i < a.length; i++) {
			if (a[i] != this[i]) return false
		}
		return true
	};
}()