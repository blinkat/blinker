+ function($) {
	function Tab(element, option) {
		this.element = element;
		this.option = option;
		var list = new $.blink.linkList();
		this.element.data('blink.tab.list', list);
		this.element.data('blink.tab.map', {});

		var divs = element.children();
		element.empty();

		this.header = $.blink.div('tab-header');
		element.append(this.header);
		this._control = $.blink.div('tab-control');
		this.header.append(this._control);
		this.content = $.blink.div('tab-content').height(element.height() - this.header.height() - 1);
		element.append(this.content);

		for (var i = 0; i < divs.length; i++) {
			var d = $(divs[i]);
			var c = d.data('close');
			this.append(d, c == null ? true : c);
		};

		this.__init = true;

		if (list.length() > 0) {
			this.action(list.get(0));
			this._control.addClass('show')
		}

		// init control
		var tab = this;
		this._control.icon = $.blink.div('tab-control-icon');
		this._control.append(this._control.icon);
		option.control.draw.call(this._control.icon);
		this._control.hover(function() {
			option.control.draw.call(tab._control.icon);
		}, function() {
			option.control.draw.call(tab._control.icon);
		});
		this._control.click(function() {
			tab.control();
		});
		this._control.btns = $.blink.doc('ul');
		this._control.list = $.blink.doc('ul');
		this._control.page = $.blink.div('tab-control-list');
		this.header.append(this._control.page);
		this._control.page.append(this._control.btns).append(this._control.list);
		this.addcontrol('关闭所有', function() {
			var map = element.data('blink.tab.map');
			tab.__init = false;
			for (var k in map) {
				var item = map[k];
				if (item.data.close) tab.remove(map[k]);
			}
			tab.__init = true;
			tab.resize();
		})
		this.resize();

		// add event
		this.element.sizeChange(function () {
			tab.resize();
			tab.data('blink.tab.current').area('checksize');
		})
	}
	Tab.prototype = $.blink.__blink__;

	Tab.prototype.append = function(div, close) {
		var tab = this;
		var uuid = $.blink.uuid();
		var head = $.blink.div('tab-block').html(div.data('name') || this.option.default_name).data('blink.tab.uuid', uuid);
		var list = this.element.data('blink.tab.list');
		var map = this.element.data('blink.tab.map');
		// add content
		this.header.append(head);
		this.content.append(div);

		// head drop
		head.mousedown(function(event) {
			var $this = $(this);
			// if down in close icon
			function in_close() {
				var icon = $this.find('.tab-block-icon');
				if (icon.length) {
					var io = icon.offset();
					var m = {
						x: event.pageX,
						y: event.pageY
					}
					return m.x >= io.left && m.x <= io.left + icon.outerWidth() &&
						m.y >= io.top && m.y <= io.top + icon.outerHeight();
				}
			}
			if (!in_close()) {
				tab.action($(this).data('blink.tab.uuid'));
				$('html').attr('onselectstart', 'return false;');
				$this.addClass('drop');
				var first = false;
				var offset = $this.offset().left - event.pageX;
				var node = map[$this.data('blink.tab.uuid')];
				node.lock = true;
				var value = node.data.header.cssval('left').num;
				var frame = tab.element.offset().left;
				var wide = node.data.header.outerWidth();
				var tolerance = tab.option.offset ? 2 : tab.option.offset;

				function enddrop() {
					$this.removeClass('drop');
					delete node.lock;
					tab.element.unbind('mousemove mouseleave mouseup');
					tab.resize();
				}

				tab.element.mousemove(function(e) {
					if (first) {
						var x = e.pageX;
						var head = node.data.header;

						var left = x - frame;
						//console.log(offset + left);
						if (left < value - tolerance && node.prev && node.prev.data.header.hasClass('show') && node.data.close == node.prev.data.close) {
							list.swap(node, node.prev);
							tab.resize();
							value -= tab.option.offset + wide;
						} else if (left > value + wide + tolerance && node.next && node.next.data.header.hasClass('show') && node.data.close == node.next.data.close) {
							list.swap(node, node.next);
							tab.resize();
							value += tab.option.offset + wide;
						}
						head.css('left', left + offset);
					} else first = true;
				});

				tab.element.mouseleave(enddrop);
				tab.element.mouseup(enddrop);
			}
		});

		var ret = {
			header: head,
			content: div,
			close: close,
			id: uuid
		}

		if (!close) {
			map[uuid] = list.front(ret);
		} else {
			// if is has close button
			// have close arrange in after
			// add icon
			var icon = $.blink.div('tab-block-icon');
			head.append(icon);
			this.option.icon.draw.call(icon);
			// icon event
			icon.click(function(event) {
				var id = $(this).parent('.tab-block').data('blink.tab.uuid');
				tab.remove(map[id]);
			});
			icon.hover(function() {
				tab.option.icon.hover.call($(this));
			}, function() {
				tab.option.icon.draw.call($(this));
			});

			head.hover(function() {
				tab.option.icon.draw.call($(this).find('.tab-block-icon'));
			}, function() {
				tab.option.icon.draw.call($(this).find('.tab-block-icon'));
			});

			// add to list
			map[uuid] = list.append(ret);
		}
		div.area();

		if (this.__init) {
			this.resize();
		}

		// if divs length == 0
		if (list.length() == 1 && this.__init) {
			this.action(uuid);
			this._control.addClass('show');
		}
	}

	Tab.prototype.remove = function(index) {
		var list = this.element.data('blink.tab.list');
		var map = this.element.data('blink.tab.map');
		var node = typeof index == 'number' ? list.get(index) : typeof index == 'string' ? map[index] : index;
		if (node) {
			list.remove(node);
			node.data.header.remove();
			node.data.content.remove();
			delete map[node.data.id];
			//this.element.data('blink.tab.map', map);
			if (node.data.id == this.element.data('blink.tab.current').data.id) {
				if (node.next) {
					this.action(node.next);
				} else if (node.prev) {
					this.action(node.prev);
				} else {
					this.element.data('blink.tab.current', null);
				}
			}
			if (this.__init)
				this.resize();
			if (list.length() == 0) {
				this._control.removeClass('show');
			}
		}
	}

	Tab.prototype.action = function(index) {
		var tf = typeof index;
		var list = this.element.data('blink.tab.list');
		var node = tf == 'number' ? list.get(index) : tf == 'string' ? this.element.data('blink.tab.map')[index] : index;
		var current = this.element.data('blink.tab.current');
		if (node && (!current || current.data.id != node.data.id)) {
			if (current) {
				current.data.header.removeClass('action');
				current.data.content.removeClass('show');
				this.option.icon.draw.call(current.data.header.find('.tab-block-icon'));
			}
			this.element.data('blink.tab.current', node);
			node.data.header.addClass('action');
			node.data.content.addClass('show');
			node.data.content.area('checksize');
			this.option.icon.draw.call(node.data.header.find('.tab-block-icon'));
			// if head not show
			if (!node.data.header.hasClass('show')) {
				var head = node.data.header;
				//head.addClass('show');
				var last;
				list.reverse(function(index) {
					if (this.data.close == node.data.close && this.data.header.hasClass('show') && this.data.id != node.data.id) {
						last = this;
						return true;
					}
				});
				// swap nodes
				if (last) {
					list.swap(node, last);
					this.resize();
				}
			}
		}
	}

	Tab.prototype.resize = function() {
		var tab = this;
		var list = this.element.data('blink.tab.list');
		var wide = this.header.width(); //- this.control.width();
		var b = true;
		var total = 0;
		var first = true;
		list.foreach(function(index) {
			var head = this.data.header;
			if (head.hasClass('first')) head.removeClass('first');
			if (head.hasClass('last')) head.removeClass('last');
			var left = index == 0 ? 0 : (total += tab.option.offset);
			if (b) {
				head.addClass('show');
				if (index == 0) {
					head.addClass('first');
				}
				if (!this.data.lock) {
					head.css('left', left);
				}
			} else {
				head.removeClass('show');
				head.css('left', 0);
			}
			total += head.outerWidth();

			if (total >= wide) {
				b = false;
				head.removeClass('show');
				head.css('left', 0);
			}
			if (first && !b) {
				first = false;
				if (this.prev) {
					this.prev.data.header.addClass('last');
				}
				tab._control.css('left', left);
				tab._control.page.css('left', left)
			} else if (first && this == list.last()) {
				first = false;
				this.data.header.addClass('last');
				left += head.outerWidth() + tab.option.offset;
				tab._control.css('left', left);
				tab._control.page.css('left', left)
			}
		});
	}

	// control
	Tab.prototype.addcontrol = function(text, fn) {
		var btns = this._control.btns;
		var li = $.blink.doc('li').data('blink.tab.btn.event', fn).html(text);
		var tab = this;
		btns.append(li);
		li.click(function() {
			var fn = $(this).data('blink.tab.btn.event');
			if (fn) {
				fn.call(tab, this);
			}
		});
	}

	Tab.prototype.control = function() {
		var page = this._control.page;
		if (page.hasClass('show')) {
			page.removeClass('show');
			this._control.removeClass('action');
			this.element.unbind('click');
		} else {
			// add list
			var list = this.element.data('blink.tab.list');
			var ul = this._control.list;
			ul.empty();
			var tab = this;
			list.foreach(function() {
				var head = this.data.header;
				var li = $.blink.doc('li').data('blink.tab.control.head', this.data.id).html(head.text());
				if (head.hasClass('action')) {
					li.addClass('action');
				}
				li.click(function(event) {
					tab.action($(this).data('blink.tab.control.head'));
				});
				ul.append(li);
				if (this.data.close) {
					var icon = $.blink.div('tab-icon');
					li.append(icon);
					tab.option.icon.draw.call(icon);
					li.hover(function() {
						tab.option.icon.draw.call(icon);
					}, function() {
						tab.option.icon.draw.call(icon);
					});
					icon.hover(function() {
						var $this = $(this);
						tab.option.icon.mousehover.call($this);
					}, function() {
						var $this = $(this);
						tab.option.icon.draw.call($this);
					});
					icon.click(function(event) {
						var $this = $(this);
						tab.remove($this.parent().data('blink.tab.control.head'));
					});
				}
			})

			page.addClass('show');
			this._control.addClass('action');
			var tab = this;
			var b = false;
			this.element.click(function() {
				if (b) {
					tab.control();
					tab.element.unbind('click');
				}
				b = true;
			});
		}
		this.option.control.draw.call(this._control.icon);
	}


	Tab.VERSION = '0.0.2';
	Tab.DEFAULT = {
		default_name: 'tab',
		offset: -1,
		// close icon
		icon: {
			draw: function() {
				$.blink.clearDraw.call(this);
				$.blink.drawFork.call(this)
			},
			hover: function() {
				var $this = $(this);
				//$this.empty();
				$.blink.clearDraw.call($this);
				var col = $this.css('color');
				var rgb = $.blink.strToRGB(col);
				if (rgb.r <= 100) {
					col = $.blink.colorFilter(col, 60);
				} else {
					col = $.blink.colorFilter(col, -60);
				}
				$.blink.drawFork.call(this, col);
			}
		},
		control: {
			draw: function() {
				$.blink.clearDraw.call(this);
				$.blink.drawArrow.call(this, 'bottom');
			}
		}

	}

	function Plugin(option) {
		var $this = $(this);
		var data = $this.data('blink.tab');
		var opt = typeof option == 'object' && option;
		if (!data) {
			var opr = $.extend({}, Tab.DEFAULT, opt);
			opr.icon = $.extend({}, Tab.DEFAULT.icon, opt.icon);
			$this.data('blink.tab', (data = new Tab($this, opr)));
		}
		if (option != null && typeof option == 'string') {
			return data.doit.apply(data, arguments) || data;
		}
		return data;
	}

	$.fn.tab = Plugin;
}(jQuery)