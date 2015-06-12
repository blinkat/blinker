+ function($) {
	function Area(element, option) {
		this.shell = element.addClass('area');
		this.option = option;

		var divs = element.children();
		var area = this;

		if (divs.length == 1 && divs[0].tagName == 'DIV') {
			this.content = divs.addClass('area-content');
		} else {
			var html = element.html();
			element.html('');
			this.content = $.blink.div('area-content').html(html);
			element.append(this.content);
		}

		if (option.horizontal != 'hide') {
			element.append((this.horizontal = $.blink.div().height(option.size).width(element.width() - option.size).addClass('horizontal hide')));
			this.horizontal.scrollbar({
				style: 'horizontal'
			});
			this.horizontal.data('blink.area.used', false);

			this.horizontal.scrollbar('scroll', function(scr) {
				if (!area.horizontal.hasClass('hide')) {
					var pos = scr.position();
					area.content.css('margin-left', -1 * ((area.content.width() - element.width() + option.size) * pos));
				}
			});

			this.horizontal.scrollbar('dropend', function() {
				if (!area.shell.data('blink.area.hover')) {
					hide();
				}
			});
		}

		if (option.vertical != 'hide') {
			element.append((this.vertical = $.blink.div().width(option.size).height(element.height() - option.size).addClass('vertical hide')));
			this.vertical.scrollbar({
				style: 'vertical'
			});
			this.vertical.data('blink.area.used', false);

			this.vertical.scrollbar('scroll', function(scr) {
				if (!area.horizontal.hasClass('hide')) {
					var pos = scr.position();
					area.content.css('margin-top', -1 * ((area.content.height() - element.height() + option.size) * pos));
				}
			});

			this.vertical.scrollbar('dropend', function() {
				if (!area.shell.data('blink.area.hover')) {
					hide();
				}
			});
		}

		if (this.horizontal && this.vertical) {
			this.vertical.scrollbar('register', this.shell);
		} else if (this.horizontal) {
			this.horizontal.scrollbar('register', this.shell);
		} else if (this.vertical) {
			this.vertical.scrollbar('register', this.shell);
		}

		this.content.sizeChange(function() {
			area.checksize();
		})
		this.shell.sizeChange(function () {
			area.checksize();
		})

		function hide(argument) {
			if (area.horizontal.data('blink.area.used')) {
				area.horizontal.addClass('hide');
			}
			if (area.vertical.data('blink.area.used')) {
				area.vertical.addClass('hide');
			}
		}

		this.shell.hover(function() {
			if (area.horizontal.data('blink.area.used')) {
				area.horizontal.removeClass('hide');
			}
			if (area.vertical.data('blink.area.used')) {
				area.vertical.removeClass('hide');
			}
			area.shell.data('blink.area.hover', true);
		}, function() {
			if ((area.horizontal && !area.horizontal.scrollbar().droping) &&
				(area.vertical && !area.vertical.scrollbar().droping)) {
				hide();
			}
			area.shell.data('blink.area.hover', false);
		});

		this.checksize();
	}

	Area.prototype = $.blink.__blink__;

	Area.prototype.checksize = function() {
		if (this.vertical) {
			var ch = this.content.height() - 8;
			var eh = this.shell.height();
			if (ch > eh) {
				if (this.vertical.hasClass('hide')) {
					this.vertical.data('blink.area.used', true);
				}
				this.vertical.scrollbar('correction', eh / ch);
			} else if (ch <= eh) {
				this.vertical.addClass('hide');
				this.vertical.data('blink.area.used', false);
			}
		}
		if (this.horizontal) {
			var ch = this.content.width() - 8;
			var eh = this.shell.width();
			if (ch > eh) {
				if (this.horizontal.hasClass('hide')) {
					this.horizontal.data('blink.area.used', true);
				}
				this.horizontal.scrollbar('correction', eh / ch);
			} else if (ch <= eh) {
				this.horizontal.addClass('hide');
				this.horizontal.data('blink.area.used', false);
			}
		}
	}

	Area.prototype.html = function(t) {
		if (t) {
			return this.content.html(t);
		}
		return this.content.html();
	}

	Area.prototype.append = function() {
		return $.fn.append.apply(this.content, arguments);
	}

	Area.VERSION = '0.0.1';
	Area.DEFAULT = {
		// when show scroll, 'auto'/'aways'
		show: 'auto',
		// 'auto'/'hide'/'show'
		horizontal: 'auto',
		vertical: 'auto',
		size: 8
	}

	function Plugin(option) {
		var $this = $(this);
		var data = $this.data('blink.area');
		var opt = typeof option == 'object' && option;
		if (!data) $this.data('blink.area', (data = new Area($this, $.extend({}, Area.DEFAULT, option))));
		if (option != null && typeof option == 'string') {
			return data.doit.apply(data, arguments) || data;
		}
		return data;
	}

	$.fn.area = Plugin;
}(jQuery)