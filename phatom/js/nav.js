// =============================
// navigation control
// =============================

+ function($) {
	function Navigation(element, option) {
		this.element = element;
		this.option = option;

		// ul li 
		if (element[0].tagName == 'ul')
			this.ul = element;
		else
			this.ul = $.blink.doc('ul', 'navigation');
		element.append(this.ul);
		if (option.style == 'horizontal') this.ul.addClass('horizontal');
		else this.ul.addClass('vertical');
		this.append(option.data)
		this.listener();
	}

	Navigation.prototype.append = function(data) {
		var nav = this;

		function add(data, obj) {
			if (data != null) {
				var draw_icon = function(div, fn) {
					if (typeof fn == 'string') {
						var img = $.blink.doc('img').attr('src', fn);
						div.append(img);
					} else {
						fn.call(div);
					}
					return nav.option.icon.position == 'left' ? div : div.addClass('right');
				}
				for (var k in data) {
					var node = data[k];
					var li = $.blink.doc('li');
					obj.append(li);
					if (typeof node == 'string') {
						li.append($.blink.doc('a').attr('href', node).html(k));
					} else {
						if (node.href) li.append($.blink.doc('a').attr('href', node.href).html(k));
						else li.append($.blink.doc('a').attr('href', 'javascript:void(0);').html(k));
						if (node.nodes) {
							// --- add icon ---
							if (obj.hasClass('navigation')) {
								var icon = $.blink.div('navigation-icon');
								li.append(icon);
								draw_icon(icon, nav.option.icon.horizontal);
							} else {
								var icon = $.blink.div('navigation-icon');
								li.append(icon);
								draw_icon(icon, nav.option.icon.normal);
							}

							var ul = $.blink.doc('ul');
							li.append(ul);
							add(node.nodes, ul);
						}
					}
				};
			}
		}
		add(data, this.ul)
	}

	Navigation.prototype.listener = function() {
		var uls = this.element.find('ul');
		for (var i = uls.length - 1; i >= 0; i--) {
			var li = $(uls[i]).parent('li');
			li.hover(function() {
				var $this = $(this);
				if (!$this.hasClass('action')) {
					var ul = $this.children('ul');
					if ($this.parent('ul').hasClass('horizontal')) {
						$this.addClass('action');
					} else {
						$this.addClass('action');
						var ls = ul.children('li');
					}
					var fn = $this.data('blink.navigation.hover');
					if (fn) {
						fn.call($this);
					}
				}
			}, function() {
				var $this = $(this);
				if ($this.hasClass('action')) {
					var ul = $this.children('ul');
					if ($this.parent('ul').hasClass('horizontal')) {
						$this.removeClass('action');
					} else {
						$this.removeClass('action');
					}
					var fn = $this.data('blink.navigation.leave');
					if (fn) {
						fn.call($this);
					}
				}
			});

			li.click(function(event) {
				var fn = $this.data('blink.navigation.click');
				if (fn) {
					fn.call($this);
				}
			});
		};

	}

	Navigation.prototype.hover = function(fn) {
		if (typeof fn == 'function') {
			this.element.data('blink.navigation.hover', fn);
		}
		return this.element.data('blink.navigation.hover');
	}

	Navigation.prototype.leave = function(fn) {
		if (typeof fn == 'function') {
			this.element.data('blink.navigation.leave', fn);
		}
		return this.element.data('blink.navigation.leave');
	}

	Navigation.prototype.click = function() {
		if (typeof fn == 'function') {
			this.element.data('blink.navigation.click', fn);
		}
		return this.element.data('blink.navigation.click');
	}

	Navigation.VERSION = '0.0.1';
	Navigation.DEFAULT = {
		style: "horizontal",
		icon: {
			position: 'right',
			horizontal: function() {
				if ($.blink.html5) {
					var $this = $(this);
					var can_tag = $.blink.doc('canvas').attr('width', $this.width()).attr('height', $this.height());
					$this.append(can_tag);
					var context = can_tag.get(0).getContext('2d');
					var wide = $this.width(),
						high = $this.height();
					var ps = [{
						x: 0,
						y: 0
					}, {
						x: wide,
						y: 0,
					}, {
						x: wide / 2,
						y: high - (high / 3)
					}];

					context.fillStyle = $this.prev('a').css('color');

					for (var i = 0; i < ps.length; i++) {
						var p = ps[i];
						context.lineTo(p.x, p.y);
					};

					context.fill();
				}
			},
			normal: function() {
				if ($.blink.html5) {
					var $this = $(this);
					var can_tag = $.blink.doc('canvas').attr('width', $this.width()).attr('height', $this.height());
					$this.append(can_tag);
					var context = can_tag.get(0).getContext('2d');
					var wide = $this.width(),
						high = $this.height();
					var ps = [{
						x: 0,
						y: 0
					}, {
						x: wide - (wide / 3),
						y: high / 2,
					}, {
						x: 0,
						y: high
					}];

					context.fillStyle = $this.prev('a').css('color');

					for (var i = 0; i < ps.length; i++) {
						var p = ps[i];
						context.lineTo(p.x, p.y);
					};

					context.fill();
				}
			}
		}
	}

	function Plugin(option) {
		var $this = $(this);
		var data = $this.data('blink.navigation');
		var opt = typeof option == 'object' && option;
		if (!data) {
			var opr = $.extend({}, Navigation.DEFAULT, option);
			opr.icon = $.extend({}, Navigation.DEFAULT.icon, option.icon);
			$this.data('blink.navigation', (data = new Navigation($this, opr)))
		}
		if (option && typeof option == 'string') {
			return data.doit.apply(data, arguments) || data;
		}
		return data;
	}

	$.fn.nav = Plugin;
}(jQuery)