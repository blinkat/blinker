/* ==========================================================
 * flex box control
 * "action"  	change(open/close) control state
 * "content" 	return content jquery object
 * "html"	 	writer html to content
 * "append"  	append to content
 * "sizeChange"	listen size change to content
 * "resize"		tell the flexbox need listen controls
 * "title"		change title html
 * ========================================================== */

+ function($) {
	var FlexBox = function(element, option) {
		this.option = option;
		this.element = element;
		this._isopen = option.is_open;

		var content, header, arrow, title, open = this._isopen;

		element.addClass('flexbox');
		content = element.children();
		var con = element.data('content');
		if (content.length == 0) {
			var html = element.html();
			element.empty();
			content = $.blink.div('flexbox-content');
			if (html != "") {
				content.html(html);
			} else {
				content.html(con);
			}

			element.append(content);
		} else if (content.length == 1 && content.get(0).tagName == "DIV") {
			content.addClass('flexbox-content');
		} else {
			var div = $.blink.div('flexbox-content');
			content = content.wrapAll(div).parent();
		}
		// ---------- init title ------------
		var val = element.data('title');
		val = val == null ? option.title : val;
		header = $.blink.div('flexbox-header');
		element.prepend(header);
		// ---------- draw icon -------------
		arrow = $.blink.div('flexbox-icon');
		header.append(arrow);
		option.icon.call(arrow);
		if (open) {
			arrow.addClass('action');
		}
		// ---------- write title ------------
		title = $.blink.div('flexbox-title').html(val);
		header.append(title);

		// add event
		var $this = this;
		header.click(function() {
			$this._isopen = !$this._isopen;
			$this.do_action();
		});

		this.content = content;
		this.header = header;
		this.arrow = arrow;
		this.title = title;
	}

	FlexBox.prototype = $.blink.__blink__;

	// ----- prototypes ------
	FlexBox.prototype.action = function(b) {
		this._isopen = b == null ? !this._isopen : b;
		this.do_action();
		return this._isopen;
	};

	FlexBox.prototype.content = function() {
		return this.content;
	};

	FlexBox.prototype.html = function(con) {
		if (con != null) {
			return this.content.html(con);
		}
		return this.content.html();
	};

	FlexBox.prototype.append = function(con) {
		if (con != null) {
			return this.content.append(con)
		}
		return this.content;
	};
	FlexBox.prototype.sizeChange = function(fn) {
		return this.content.sizeChange(fn)
	};

	FlexBox.prototype.resize = function() {
		var arg = Array.prototype.slice.call(arguments);
		var flex = this;
		for (var i = arg.length - 1; i >= 0; i--) {
			var item = arg[i];
			if (item.Constructor == $.fn.Constructor) {
				item.sizeChange(function() {
					flex.do_action();
				})
			}
		};
	};

	FlexBox.prototype.title = function(title) {
		if (title != null) {
			this.title.html(title);
		}
	}

	FlexBox.prototype.do_action = function() {
		var ele = this.element;
		if (this._isopen) {
			ele.addClass('action')
			ele.height(this.header.height() + this.content.height() +
				(this.content.cssval('padding-top').num + this.content.cssval('padding-bottom').num));
		} else {
			ele.removeClass('action')
			ele.height(this.header.height());
		}
	}

	FlexBox.prototype._isopen = false;


	FlexBox.VERSION = '0.0.1';
	FlexBox.DEFAULT = {
		title: 'flexbox',
		icon: function() {
			$.blink.clearDraw.call(this);
			$.blink.drawArrow.call(this);
		}
	}

	function Plugin(option) {
		var $this = $(this);
		var data = $this.data('blink.flexbox');
		var opt = typeof option == 'object' && option;

		if (!data) $this.data('blink.flexbox', (data = new FlexBox($this, $.extend({}, FlexBox.DEFAULT, option))));

		if (option != null && typeof option == "string") {
			return data.doit.apply(data, arguments) || data;
		}
		return data;
	}

	$.fn.flexbox = Plugin;
}(jQuery)