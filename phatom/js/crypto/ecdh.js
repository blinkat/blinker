+ function() {

	var one = new sjcl.bn(1)

	function randFieldElement(c) {
		var b = sjcl.random.randomWords(Math.floor(c.r.bitLength() / 8) + 8),
			k = new sjcl.bn(b),
			n = c.r.sub(one)
		k = k.mod(n).add(one)
		return k
	}

	blink.crypto.generateKey = function(c) {
		var k = randFieldElement(c)

		var priv = {
			publicKey: {
				curve: c,
				point: c.mult(k.toBits(), c.G)
			},
			d: k,
		}
		return priv
	}
}();