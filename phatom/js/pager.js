+ function($) {
	function Pages(element, option) {
		this.element = element;
		this.option = option;
		this.content = $.blink.div('pages-content');
		element.data('blink.pages.current', option.current);
		var arr = [];
		element.data('blink.pages.array', arr);
		element.append(this.content);
		element.data('blink.pages.active', true);

		var page = this;
		this.content.append($.blink.doc('a').attr('href', 'javascript:void(0);').html(option.text.home).click(function() {
			page.page(0);
		}));

		this.content.append($.blink.doc('a').attr('href', 'javascript:void(0);').html(option.text.prev).click(function() {
			var cur = element.data('blink.pages.current') - 1;
			page.page(cur);
		}));

		// add numbers
		this.number = $.blink.doc('span')
		this.content.append(this.number);
		this.number.tophide = $.blink.doc('span', 'pages-hides hide').html('...');
		this.number.append(this.number.tophide);
		for (var i = 0; i < option.total; i++) {
			var p = $.blink.doc('a').attr('href', 'javascript:void(0);').html(i + 1).click(function() {
				page.page($(this).data('blink.pages.number'));
			}).data('blink.pages.number', i).addClass('pages-number hide');
			arr[i] = p;
			if (i == element.data('blink.pages.current')) {
				p.addClass('action');
			}
			this.number.append(p);
		};
		this.number.bothide = $.blink.doc('span', 'pages-hides hide').html('...');
		this.number.append(this.number.bothide);
		this.resizePage();

		this.content.append($.blink.doc('a').attr('href', 'javascript:void(0);').html(option.text.next).click(function() {
			var cur = element.data('blink.pages.current') + 1;
			page.page(cur);
		}));

		this.content.append($.blink.doc('a').attr('href', 'javascript:void(0);').html(option.text.end).click(function() {
			page.page(option.total - 1);
		}));

		this.content.append($.blink.doc('input').attr('type', 'text').attr('placeholder', 'page').bind('keypress', function(event) {
			if (event.keyCode == "13") {
				var $this = $(this);
				page.page(parseInt($this.val()) - 1);
				$this.val('');
				$this.blur();
			}
		}));

		this.content.append($.blink.doc('a').attr('href', 'javascript:void(0);').html(option.text.jump).click(function() {
			var inp = element.find('input');
			page.page(parseInt(inp.val()) - 1);
			inp.val('');
			inp.blur();
		}));
	}
	Pages.prototype = $.blink.__blink__;

	Pages.prototype.turning = function(fn) {
		this.element.data('blink.tab.event', fn);
	}

	Pages.prototype.page = function(num) {
		if (typeof num == 'number' && num != NaN && this.element.data('blink.pages.active')) {
			num = num < 0 ? 0 : num;
			num = num > this.option.total - 1 ? this.option.total - 1 : num;

			var cur = this.element.data('blink.pages.current');
			if (cur != num) {
				var arr = this.element.data('blink.pages.array');
				arr[cur].removeClass('action');
				this.element.data('blink.pages.current', num);
				arr[num].addClass('action');
				var fn = this.element.data('blink.pages.event');
				if (fn) fn.call(this, num);
				this.resizePage();
			}
		}
	}

	Pages.prototype.resizePage = function() {
		var arr = this.element.data('blink.pages.array');
		var cur = this.element.data('blink.pages.current');
		var num = this.number;

		var flag = 0;
		var cha = this.option.max / 2;
		var b = false;
		if (cur <= cha) {
			if (this.option.total > this.option.max) {
				this.number.tophide.addClass('hide');
				this.number.bothide.removeClass('hide');
			}

			for (var i = 0; i < this.option.total; i++) {
				if (flag < this.option.max) {
					arr[i].removeClass('hide');
					flag++;
				} else {
					arr[i].addClass('hide');
				}
			};
		} else if (this.option.total - cur <= cha) {
			if (this.option.total > this.option.max) {
				this.number.tophide.removeClass('hide');
				this.number.bothide.addClass('hide');
			}
			for (var i = this.option.total - 1; i >= 0; i--) {
				if (flag < this.option.max) {
					arr[i].removeClass('hide');
					flag++;
				} else {
					arr[i].addClass('hide');
				}
			};
		} else {
			if (this.option.total > this.option.max) {
				this.number.tophide.removeClass('hide');
				this.number.bothide.removeClass('hide');
			}
			var max = cur + cha;
			var min = cur - cha;
			for (var i = 0; i < this.option.total; i++) {
				if (i <= max && i > min) arr[i].removeClass('hide');
				else arr[i].addClass('hide');
			}
		}
	}

	Pages.prototype.active = function(active) {
		if (arguments.length == 0) {
			active = !this.element.data('blink.pages.active');
		}
		this.element.data('blink.pages.active', active);
	}

	Pages.VERSION = '0.0.1'
	Pages.DEFAULT = {
		text: {
			home: "首页",
			prev: "上一页",
			next: '下一页',
			end: '尾页',
			jump: '跳页'
		},
		current: 0,
		max: 6,
		total: 0
	}

	function Plugin(option) {
		var $this = $(this);
		var data = $this.data('blink.pages');
		var opt = typeof option == 'object' && option;
		if (!data) {
			var opr = $.extend({}, Pages.DEFAULT, opt);
			opr.text = $.extend({}, Pages.DEFAULT.text, opt.text);
			$this.data('blink.pages', (data = new Pages($this, opr)));
		}
		if (option != null && typeof option == 'string') {
			return data.doit.apply(data, arguments) || data;
		}
		return data;
	}
	$.fn.pager = Plugin;
}(jQuery)