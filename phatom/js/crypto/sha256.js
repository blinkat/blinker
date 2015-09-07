+ function() {
	var chunk = 64,
		init0 = 0x6A09E667,
		init1 = 0xBB67AE85,
		init2 = 0x3C6EF372,
		init3 = 0xA54FF53A,
		init4 = 0x510E527F,
		init5 = 0x9B05688C,
		init6 = 0x1F83D9AB,
		init7 = 0x5BE0CD19,

		size = 32,
		blocksize = 64

	var _K = [
		0x428a2f98,
		0x71374491,
		0xb5c0fbcf,
		0xe9b5dba5,
		0x3956c25b,
		0x59f111f1,
		0x923f82a4,
		0xab1c5ed5,
		0xd807aa98,
		0x12835b01,
		0x243185be,
		0x550c7dc3,
		0x72be5d74,
		0x80deb1fe,
		0x9bdc06a7,
		0xc19bf174,
		0xe49b69c1,
		0xefbe4786,
		0x0fc19dc6,
		0x240ca1cc,
		0x2de92c6f,
		0x4a7484aa,
		0x5cb0a9dc,
		0x76f988da,
		0x983e5152,
		0xa831c66d,
		0xb00327c8,
		0xbf597fc7,
		0xc6e00bf3,
		0xd5a79147,
		0x06ca6351,
		0x14292967,
		0x27b70a85,
		0x2e1b2138,
		0x4d2c6dfc,
		0x53380d13,
		0x650a7354,
		0x766a0abb,
		0x81c2c92e,
		0x92722c85,
		0xa2bfe8a1,
		0xa81a664b,
		0xc24b8b70,
		0xc76c51a3,
		0xd192e819,
		0xd6990624,
		0xf40e3585,
		0x106aa070,
		0x19a4c116,
		0x1e376c08,
		0x2748774c,
		0x34b0bcb5,
		0x391c0cb3,
		0x4ed8aa4a,
		0x5b9cca4f,
		0x682e6ff3,
		0x748f82ee,
		0x78a5636f,
		0x84c87814,
		0x8cc70208,
		0x90befffa,
		0xa4506ceb,
		0xbef9a3f7,
		0xc67178f2,
	]

	function block(dig, p) {
		var w = [],
			h0 = dig.h[0],
			h1 = dig.h[1],
			h2 = dig.h[2],
			h3 = dig.h[3],
			h4 = dig.h[4],
			h5 = dig.h[5],
			h6 = dig.h[6],
			h7 = dig.h[7]

		while (p.length >= chunk) {
			for (var i = 0; i < 16; i++) {
				var j = i * 4
				w[i] = p[j] << 24 | p[j + 1] << 16 | p[j + 2] << 8 | p[j + 3]
			}

			for (var i = 16; i < 64; i++) {
				var v1 = w[i - 2],
					t1 = (v1 >>> 17 | v1 << (32 - 17)) ^ (v1 >>> 19 | v1 << (32 - 19)) ^ (v1 >>> 10),
					v2 = w[i - 15],
					t2 = (v2 >>> 7 | v2 << (32 - 7)) ^ (v2 >>> 18 | v2 << (32 - 18)) ^ (v2 >>> 3)
				w[i] = t1 + w[i - 7] + t2 + w[i - 16]
			}

			var a = h0,
				b = h1,
				c = h2,
				d = h3,
				e = h4,
				f = h5,
				g = h6,
				h = h7

			for (var i = 0; i < 64; i++) {
				var t1 = h + ((e >>> 6 | e << (32 - 6)) ^ (e >>> 11 | e << (32 - 11)) ^ (e >>> 25 | e << (32 - 25))) + ((e & f) ^ (~e & g)) + _K[i] + w[i]
				var t2 = ((a >>> 2 | a << (32 - 2)) ^ (a >>> 13 | a << (32 - 13)) ^ (a >>> 22 | a << (32 - 22))) + ((a & b) ^ (a & c) ^ (b & c))
				h = g
				g = f
				f = e
				e = d + t1
				d = c
				c = b
				b = a
				a = t1 + t2
			}

			h0 = (a + h0).uint32()
			h1 = (b + h1).uint32()
			h2 = (c + h2).uint32()
			h3 = (d + h3).uint32()
			h4 = (e + h4).uint32()
			h5 = (f + h5).uint32()
			h6 = (g + h6).uint32()
			h7 = (h + h7).uint32()
			p = p.slice(chunk)
		}

		dig.h[0] = h0
		dig.h[1] = h1
		dig.h[2] = h2
		dig.h[3] = h3
		dig.h[4] = h4
		dig.h[5] = h5
		dig.h[6] = h6
		dig.h[7] = h7
	}

	function checksum(d) {
		var len = d.len,
			tmp = Array.new(64, 0)
		tmp[0] = 0x80
		if (len % 64 < 56)
			d.write(tmp.slice(0, 56 - len % 64))
		else
			d.write(tmp.slice(0, 64 + 56 - len % 64))
		len <<= 3
		for (var i = 0; i < 8; i++)
			tmp[i] = (len >>> (56 - 8 * i)).byte()
		d.write(tmp.slice(0, 8))
		if (d.nx != 0) throw "d.nx != 0"

		var h = d.h.slice()
		var digest = []
		for (var i = 0; i < h.length; i++) {
			var s = h[i]
			digest[i * 4] = (s >>> 24).byte()
			digest[i * 4 + 1] = (s >>> 16).byte()
			digest[i * 4 + 2] = (s >>> 8).byte()
			digest[i * 4 + 3] = s.byte()
		}
		return digest
	}

	function Digest() {
		this.h = Array(8)
		this.x = Array(chunk)
		this.nx = 0
		this.len = 0
	}

	Digest.prototype.reset = function() {
		this.h[0] = init0
		this.h[1] = init1
		this.h[2] = init2
		this.h[3] = init3
		this.h[4] = init4
		this.h[5] = init5
		this.h[6] = init6
		this.h[7] = init7

		this.nx = 0
		this.len = 0
	}

	Digest.prototype.size = function() {
		return size
	}

	Digest.prototype.blocksize = function() {
		return blocksize
	}

	Digest.prototype.write = function(p) {
		var nn = p.length,
			d = this
		this.len += nn
		if (this.nx > 0) {
			var n = d.x.copy(p, d.nx)
			d.nx += n
			if (d.nx == chunk) {
				block(d, d.x.slice())
				d.nx = 0
			}
			p = p.slice(n)
		}
		if (p.length >= chunk) {
			var n = p.length & ~(chunk - 1)
			block(d, p.slice(0, n))
			p = p.slice(n)
		}
		if (p.length > 0) {
			d.nx = d.x.copy(p)
		}
		return nn
	}

	Digest.prototype.sum = function(inp) {
		var hash = checksum(this)
		return inp.concat(hash)
	};

	blink.crypto.sha256 = function() {
		var d = new Digest()
		d.reset()
		return d
	}
}()