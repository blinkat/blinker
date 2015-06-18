//====================================
// core script
//====================================
(function() {
	// if jquery is not load, throw error
	if ($ == null) {
		throw "use this control must have jquery"
	}

	// some function namespace
	$.blink = {
		html5: document.createElement('canvas').getContext != null,
		// rsa
		rsa: {
			pem: function(str) {
				var ret = Base64.re.exec(str);
				return ret && ret.length > 1 ? ret[1] : str;
			},
			decrypt: function(text, prk) {
				if (text && prk) {
					var rsa = $.blink.rsa;
					text = rsa.pem(text);
					var dec = new JSEncrypt();
					dec.setPrivateKey(prk);
					return dec.decrypt(text);
				}
				throw "decryption error";
			},
			encrypt: function(text, puk) {
				if (text && puk) {
					var rsa = $.blink.rsa;
					text = rsa.pem(text);
					var dec = new JSEncrypt();
					dec.setPublicKey(puk);
					return "-----BEGIN TEXT-----" + dec.encrypt(text) + "-----END TEXT-----";
				}
				throw "encryption error";
			}
		},
		md5: $.md5,

		doc: function(tag, cls, id) {
			var ret = $(document.createElement(tag));
			if (id != null)
				ret.attr('id', id);
			if (cls != null)
				ret.addClass(cls);
			return ret;
		},
		div: function(cls, id) {
			return $.blink.doc("div", cls, id);
		},

		// transition-duartion sub to ms
		// 123ms return 123
		// 3s return 3000
		getnum: function(str) {
			str = str.trim();
			var num = 0,
				unit = '';
			var i = 0;
			for (; i < str.length; i++) {
				var ch = str.charAt(i);
				if (/\d/.test(ch) || ch == ".") {
					num += ch;
				} else {
					break;
				}
			}
			num = parseFloat(num);
			if (num == NaN) {
				num = 0;
			}
			unit = str.substr(i);
			return {
				num: num,
				unit: unit
			};
		},

		mouseup: function(fn) {
			if (fn) {
				$(document).mouseup(function() {
					if (fn) {
						fn.call(this);
						fn = null;
					}
				});
			}
		},

		mousemove: function(fn) {
			if (fn) {
				$(document).mousemove(function(e) {
					if (fn) {
						b = fn.call(this, e);
						if (b) {
							fn = null;
						}
					}
				});
			}
		},

		mouseclick: function(fn) {
			if (fn) {
				$(document).click(function(e) {
					if (fn) {
						fn.call(this, e);
						fn = null;
					}
				});
			}
		},

		strToRGB: function(str) {
			if (/#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})$/.test(str)) {
				var sR = "";
				var sG = "";
				var sB = "";
				var nR = 0;
				var nG = 0;
				var nB = 0;
				str = str.substring(1);
				if (str.length == 3) {
					str = str.toLowerCase()
					sR = str.substring(0, 1) + str.substring(0, 1);
					sG = str.substring(1, 2) + str.substring(1, 2);
					sB = str.substring(2) + str.substring(2);
				} else if (str.length <= 6) {
					sR = str.substring(0, 2);
					sG = str.substring(2, 4);
					sB = str.substring(4);
				}
				nR = parseInt(sR, 16);
				nG = parseInt(sG, 16);
				nB = parseInt(sB, 16);
				return {
					r: nR,
					g: nG,
					b: nB
				}
			} else if (/^(rgb|RGB)/.test(str)) {
				col = str.replace(/(?:\(|\)|rgb|RGB)*/g, "").split(',');
				return {
					r: parseInt(col[0]),
					g: parseInt(col[1]),
					b: parseInt(col[2])
				}
			}

		},

		rgbToHex: function(r, g, b) {
			if (arguments == 1) {
				var rgb = r;
				r = rgb.r;
				g = rgb.g;
				b = rgb.b;
			}
			var ir = r.toString(16);
			var ig = g.toString(16);
			var ib = b.toString(16);
			return "#" + ir + ig + ib;
		},

		colorFilter: function(color, depth) {
			if (depth != 0) {
				var rgb = $.blink.strToRGB(color);
				var rgb2 = depth > 0 ? {
					r: 255,
					g: 255,
					b: 255
				} : {
					r: 0,
					g: 0,
					b: 0
				};
				depth = depth > 0 ? depth : -depth;

				var dis = 100;
				var rd = rgb.r - rgb2.r;
				var gd = rgb.g - rgb2.g;
				var bd = rgb.b - rgb2.b;
				var r, g, b;
				r = Math.floor(rgb2.r + (rd / dis) * depth);
				g = Math.floor(rgb2.g + (gd / dis) * depth);
				b = Math.floor(rgb2.b + (bd / dis) * depth);

				return $.blink.rgbToHex(r, g, b);
			}
			return color;
		},
		// draw canvas
		drawFork: function(color) {
			if ($.blink.html5) {
				if (arguments.length == 0) {
					color = this.css('color');
				}
				var canvas = this.find('canvas');
				if (canvas.length == 0) {
					canvas = $.blink.doc('canvas').attr('width', this.width()).attr('height', this.height());
					this.append(canvas);
				}

				var context = canvas[0].getContext('2d');
				var wide = this.width(),
					high = this.height();
				context.moveTo(0, 0);
				context.lineTo(wide, high);
				context.moveTo(wide, 0);
				context.lineTo(0, wide);

				context.strokeStyle = color;
				context.lineWidth = 1;
				context.stroke();
			}
		},
		drawRound: function(color) {
			if ($.blink.html5) {
				if (arguments.length == 0) {
					color = this.css('background-color');
					color = $.blink.colorFilter(50);
				}
				var canvas = this.find('canvas');
				if (canvas.length == 0) {
					canvas = $.blink.doc('canvas').attr('width', this.width()).attr('height', this.height());
					this.append(canvas);
				}

				var context = canvas[0].getContext('2d');
				var wide = this.width(),
					high = this.height();
				context.beginPath();
				context.arc(wide / 2, high / 2, wide < high ? wide / 2 : high / 2, 0, 2 * Math.PI);
				context.fillStyle = color;
				context.fill();
			}
		},
		drawArrow: function(angle, color) {
			if ($.blink.html5) {
				if (color == null) color = this.css('color');
				if (angle == null) angle = 'right';
				var wide = this.width(),
					high = this.height();
				var ps;
				switch (angle) {
					case 'right':
						ps = [{
							x: wide / 4,
							y: 0
						}, {
							x: wide - wide / 4,
							y: high / 2,
						}, {
							x: wide / 4,
							y: high
						}];
						break;
					case 'bottom':
						ps = [{
							x: 0,
							y: high / 4
						}, {
							x: wide / 2,
							y: high - high / 4,
						}, {
							x: wide,
							y: high / 4
						}];
						break;
				}
				var canvas = this.find('canvas');
				if (canvas.length == 0) {
					canvas = $.blink.doc('canvas').attr({
						width: wide,
						height: high
					});
					this.append(canvas);
				}
				var content = canvas[0].getContext('2d');
				for (var i = 0; i < ps.length; i++) {
					var p = ps[i];
					if (i == 0) content.moveTo(p.x, p.y);
					else content.lineTo(p.x, p.y);
				};

				content.strokeStyle = color;
				content.lineWidth = 2;
				content.stroke();
			}
		},
		clearDraw: function() {
			if ($.blink.html5) {
				var canvas = this.find('canvas');
				if (canvas.length > 0)
					canvas[0].getContext('2d').clearRect(0, 0, canvas.width(), canvas.height());
			}
		},

		// uuid
		uuid: function(len, radix) {
			var uuid = [],
				i,
				chars = $.blink.option.uuid.chars;
			radix = radix || chars.length;
			len = len || $.blink.option.uuid.length;

			if (len) {
				for (i = 0; i < len; i++) uuid[i] = chars[0 | Math.random() * radix];
			} else {
				var r;
				uuid[8] = uuid[13] = uuid[18] = uuid[23] = '-';
				uuid[14] = '4';
				for (var i = 0; i < 36; i++) {
					if (!uuid[i]) {
						r = 0 | Math.random() * 16;
						uuid[i] = chars[(i == 19) ? (r & 0x3) | 0x8 : r];
					}
				};
			}
			return uuid.join('');
		},

		// setting
		option: {
			// timer interval
			interval: 100,
			uuid: {
				chars: '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'.split(''),
				length: 16,
			}
		}
	};

	// extend jquery
	$.fn.cssval = function() {
		var val;
		if (arguments.length == 1) {
			val = this.css(arguments[0]);
		} else if (arguments.length >= 2) {
			val = this.css(arguments[0], arguments[1]);
		}
		return $.blink.getnum(val);
	}

	// get transition val
	// res[transition-property] = { 
	//			time =transition-duration, 
	//			timing = transition-timing-function 
	//			delay = transition-delay
	//}
	$.fn.transition = function() {
		if ($.blink.html5) {
			var transition = this.css('transition-property'); // get all names
			transition = transition != null ? transition.split(',') : null;

			var durations = this.css('transition-duration');
			durations = durations != null ? durations.split(',') : null;

			var timings = this.css('transition-timing-function');
			timings = timings != null ? timings.split(',') : null;

			var delays = this.css('transition-delay');
			delays = delays != null ? delays.split(',') : null;

			var ret = [];
			if (transition != null) {
				for (var k = 0; k < transition.length; k++) {
					var dur = $.blink.getnum(durations[k]);
					var del = $.blink.getnum(delays[k]);
					ret[transition[k]] = {
						time: dur.unit == 's' ? dur.num * 1000 : dur.num,
						timing: timings[k],
						delay: del.num == 0 ? 0 : del.unit == 's' ? del.num * 1000 : del.num,
					}
				}
			}
			return ret;
		}
	}

	// extend string
	String.prototype.trim = function() {　　
		return this.replace(/(^\s*)|(\s*$)/g, "");　　
	}

	// ------- base class -----------
	$.blink.__blink__ = {
		"doit": function(opt) {
			var arg = Array.prototype.slice.call(arguments, 1);
			var fn = this[opt];
			if (fn != null) {
				return fn.apply(this, arg);
			}
			return this;
		}
	}

	// blink link list
	$.blink.linkList = function() {
		var head, last, length = 0;

		this.append = function(data) {
			var ret = {
				data: data,
			}
			if (!head) {
				head = last = ret;
				//ret.position = 0;
			} else {
				last.next = ret;
				ret.prev = last;
				//ret.position = last.position + 1;
				last = ret;
			}
			length++;
			return ret;
		}

		this.front = function(data) {
			var ret = {
				data: data
			}
			var old = head;
			head = ret;
			ret.next = old;
			if (old) old.prev = ret;
			length++;
			return ret;
		}

		this.insertAfter = function(index, data) {
			if (this.isEmpty() || index == length - 1 || index == null) {
				return this.append(data);
				/*} else if (index == 0) {
					return this.front(data);*/
			} else {
				var node = this.get(index);
				if (node) {
					ret = {
						data: data
					};
					node.next.prev = ret;
					ret.prev = node;
					ret.next = node.next;
					node.next = ret;
					length++;
					/*for (var i = ret.next.position, next = ret; i < length; i++, next = next.next) {
						next.position = i;
					};*/
					return ret;
				} else {
					return this.append(data);
				}
			}
		}

		this.foreach = function(fn) {
			for (var i = 0, node = head; i < length; i++, node = node.next) {
				if (fn.call(node, i)) {
					break;
				}
			};
		}

		this.reverse = function(fn) {
			for (var i = length - 1, node = last; i >= 0; i--, node = node.prev) {
				if (fn.call(node, i)) {
					break;
				}
			};
		}

		this.insertBefore = function(index, data) {
			if (this.isEmpty() || index == length - 1) {
				return this.append(data);
			} else if (index == 0) {
				return this.front(data);
			} else {
				var node = this.get(index);
				if (node) {
					ret = {
						data: data
					}
					node.prev.next = ret;
					ret.prev = node.prev;
					ret.next = node;
					node.prev = ret;
					length++;
					return ret;
				} else {
					return this.append(data);
				}
			}
		}

		this.get = function(index) {
			if (!index) {
				return head;
			}
			for (var i = 0, node = head; i < length; i++, node = node.next) {
				if (i == index) {
					return node;
				}
			};
		}

		this.swap = function(n1, n2) {

			if (n1 != n2) {
				var node1 = typeof n1 == 'number' ? this.get(n1) : n1;
				var node2 = typeof n2 == 'number' ? this.get(n2) : n2;
				var prev1 = node1.prev,
					prev2 = node2.prev;
				var next1 = node1.next,
					next2 = node2.next;
				if (prev1 && prev1 == node2) {
					if (node2 == head) head = node1;
					if (node1 == last) last = node2;
					node1.next = node2;
					node2.prev = node1;

					node1.prev = prev2;
					if (prev2) prev2.next = node1;

					node2.next = next1;
					if (next1) next1.prev = node2;


				} else if (prev2 && prev2 == node1) {
					if (node1 == head) head = node2;
					if (node2 == last) last = node1;
					node1.prev = node2;
					node2.next = node1;

					node1.next = next2;
					if (next2) next2.prev = node1;

					node2.prev = prev1;
					if (prev1) prev1.next = node2;
				} else {
					node1.next = next2;
					node2.next = next1;
					if (next2) next2.prev = node1;
					if (next1) next1.prev = node2;
					node1.prev = prev2;
					if (prev2) prev2.next = node1;
					node2.prev = prev1;
					if (prev1) prev1.next = node2
				}
			}
		}

		this.remove = function(index) {
			var node = typeof index == 'number' ? node = this.get(index) : index;
			if (node) {
				if (--length == 0) head = last = null;
				else if (node == head) head = node.next;
				else if (node == last) last = node.prev;
				else {
					if (node.prev) node.prev.next = node.next;
					if (node.next) node.next.prev = node.prev;
				}
				return node;
			}
		}

		this.pop = function() {
			return this.remove(length - 1);
		}

		this.shift = function() {
			return this.remove(0);
		}

		this.isEmpty = function() {
			return length == 0;
		}

		this.head = function() {
			return head;
		}
		this.last = function() {
			return last;
		}
		this.length = function() {
			return length;
		}
	}
})();