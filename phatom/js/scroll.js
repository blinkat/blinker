+ function($) {
	function Scroll(element, option) {
		this.element = element;
		this.option = option;

		this.element.addClass('scroll');
		this.element.append((this.slider = $.blink.div('slider')));
		this.correction();
		this.register(element);
		this.event = {};
		this.droping = false;

		var scr = this;
		this.slider.mousedown(function(e) {
			var b = false;
			$('html').attr('onselectstart', 'return false;');
			var mouse = {
				x: e.screenX,
				y: e.screenY
			}
			scr.droping = true;

			$.blink.mouseup(function(event) {
				$('html').removeAttr('onselectstart');
				b = true;
				scr.droping = false;
				var fn = scr.element.data('blink.scroll.dropend');
				if (fn) {
					fn.call(scr);
				}
			});

			$.blink.mousemove(function(e) {
				var m = {
					x: e.screenX,
					y: e.screenY,
				}
				var offset;
				if (option.style == 'horizontal') offset = m.x - mouse.x;
				else offset = m.y - mouse.y;
				mouse = m;
				scr.move(offset);
				return b;
			});
		});
	}

	Scroll.prototype = $.blink.__blink__;

	Scroll.prototype.correction = function(per) {
		var ele = this.element;
		var slider = this.slider;
		var opr = this.option;

		if (per != null) per = per > 1 ? 1: per;
		 else per = 1;

		if (opr.style == 'horizontal') {
			per = ele.width() * per;
			per = per < this.option.min ? this.option.min : per;
			slider.height(ele.height() - (opr.offset * 2)).css('margin-top', opr.offset).width(per);
		} else {
			per = ele.height() * per;
			per = per < this.option.min ? this.option.min : per;
			slider.width(ele.width() - (opr.offset * 2)).css('margin-left', opr.offset).height(per);
		}
	}

	Scroll.prototype.register = function(obj) {
		var scr = this;
		obj.mousewheel(function(e) {
			var opr = scr.option;
			scr.move(e.deltaY * scr.option.wheel * -1);
		})
	}

	Scroll.prototype.dropend = function (fn) {
		if (typeof fn == 'function') {
			this.element.data('blink.scroll.dropend', fn);
		}
		return this.element.data('blink.scroll.dropend');
	}

	Scroll.prototype.move = function(val) {
		if (val != null) {
			if (this.option.style == 'horizontal') {
				var offset = $.blink.getnum(this.slider.css('margin-left')).num + val;
				var max = this.element.width() - this.slider.width();
				if (offset > max) offset = max;
				if (offset < 0) offset = 0;
				this.slider.css('margin-left', offset);
			} else {
				var offset = $.blink.getnum(this.slider.css('margin-top')).num + val;
				var max = this.element.height() - this.slider.height();
				if (offset > max) offset = max;
				if (offset < 0) offset = 0;
				this.slider.css('margin-top', offset);
			}
			var fn = this.element.data('blink.scroll.scroll');
			if (fn) {
				fn.call(this.element, this);
			}
		}
	}

	// get position percent
	Scroll.prototype.position = function() {
		var pos;
		if (this.option.style == 'horizontal') {
			pos = $.blink.getnum(this.slider.css('margin-left')).num;
			pos = pos / (this.element.width() - this.slider.width());
		} else {
			pos = $.blink.getnum(this.slider.css('margin-top')).num;
			pos = pos / (this.element.height() - this.slider.height());
		}
		return pos;
	}

	Scroll.prototype.offset = function () {
		return this.option.style == 'horizontal' ? 
			$.blink.getnum(this.slider.css('margin-left')).num : $.blink.getnum(this.slider.css('margin-top')).num;
	}

	Scroll.prototype.scroll = function (fn) {
		if (typeof fn == 'function') {
			this.element.data('blink.scroll.scroll', fn);
		}
		return this.element.data('blink.scroll.scroll');
	}

	Scroll.VERSION = '0.0.1';
	Scroll.NAMESPACE = 'blink.scroll';
	Scroll.DEFAULT = {
		style: 'verictal',
		min: .1,
		percent: .5,
		offset: 1,
		wheel: 8,
	}

	function Plugin(option) {
		var $this = $(this);
		var data = $this.data(Scroll.NAMESPACE);
		var opt = typeof option == 'object' && option;
		if (!data)
			$this.data(Scroll.NAMESPACE, (data = new Scroll($this, $.extend({}, Scroll.DEFAULT, option))))
		if (option && typeof option == 'string') {
			return data.doit.apply(data, arguments) || data;
		}
		return data;
	}

	$.fn.scrollbar = Plugin;
}(jQuery)