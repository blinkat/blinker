//================================
// size change event
//================================
(function() {
    var event_name = "blink.sizechange";

    // transition-duartion sub to ms
    // 123ms return 123
    // 3s return 3000
    function sub_unit(str) {
        var ret = $.blink.getnum(str);
        if (ret.unit == "s") {
            ret.num *= 1000;
        }
        return ret.num;
    }

    // get transition time (ms)
    function get_times($div) {
        if ($.blink.html5) {
            var res = {};
            var trans = $div.transition();
            if (trans != null) {
                for (var k in trans) {
                    if (k == 'width' || k == 'height') {
                        var item = trans[k];
                        res[k] = item.delay + item.time;
                    }
                }
            }
            return res;
        }
        return {
            width: 0,
            height: 0
        };
    }

    // size change event shell func
    function size_event_shell(ele) {
        var now = $.now(),
            $this = $(ele),
            timer = $this.data(event_name + "_timer"),
            fns = $this.data(event_name);
        if (fns) {
            fns.call($this, $this.width(), $this.height())
        }

        if (now - timer.start >= timer.time) {
            clearInterval(timer.id)
            console.log('event end timer:' + timer.id);
        }
    }

    // add size change event
    $.fn.sizeChange = function(fn) {
        return this.each(function() {
            var $this = $(this);
            if (fn == null) {
                return $this.data(event_name);
            }
            $this.data(event_name, fn);

            $this.on('sizechange', function(e, need) {
                var ele = $(this);
                var old_timer = $this.data(event_name + "_timer");
                if (ele.data(event_name) == null) {
                    return;
                } else if (old_timer) {
                    clearInterval(old_timer.id);
                }

                var width = $this.width(),
                    height = $this.height(),
                    times = need ? get_times(ele) : {
                        width: 0,
                        height: 0
                    };

                if ((times.width == 0 && times.height == 0) || (times.width == null && times.height == null)) {
                    $this.data(event_name).call(this, width, height);
                } else {
                    var time = times.width > times.height ? times.width : times.height,
                        id = setInterval(size_event_shell, $.blink.option.interval, this);
                    console.log("event timer start:" + id);
                    $this.data(event_name + "_timer", {
                        id: id,
                        time: time,
                        start: $.now(),
                    });
                }

            });
            return fn;
        });
    };

    // jquery funcs
    var old = {
        height: $.fn.height,
        width: $.fn.width,
        css: $.fn.css,
        addClass: $.fn.addClass,
        removeClass: $.fn.removeClass,
        attr: $.fn.attr,
        removeAttr: $.fn.removeAttr,
        html: $.fn.html,
        append: $.fn.append,
        appendTo: $.fn.appendTo,
    }


    function call_old(ele, fn, args, tr, need) {
        var r = fn.apply(ele, args);
        if ($(ele).data(event_name) && ((typeof tr == "number" && args.length >= tr) || (typeof tr == "function" && tr.apply(this, args))))
            $(ele).trigger("sizechange", !need);
        return r;
    }

    // overload function
    $.fn.height = function() {
        return call_old(this, old.height, arguments, 1);
    }

    $.fn.width = function() {
        return call_old(this, old.width, arguments, 1);
    }

    $.fn.css = function() {
        return call_old(this, old.css, arguments, function() {
            return arguments[1] && (arguments[0] == "height" || arguments[0] == "width")
        });
    }

    $.fn.addClass = function() {
        return call_old(this, old.addClass, arguments, 1);
    }

    $.fn.removeClass = function() {
        return call_old(this, old.removeClass, arguments, 1);
    }

    $.fn.attr = function() {
        return call_old(this, old.attr, arguments, function() {
            return arguments[0] == "style"
        })
    }

    $.fn.removeAttr = function() {
        return call_old(this, old.removeAttr, arguments, function() {
            return arguments[0] == "style"
        })
    }

    $.fn.html = function() {
        return call_old(this, old.html, arguments, 1, true);
    }

    $.fn.append = function() {
        return call_old(this, old.append, arguments, 1, true);
    }

    $.fn.appendTo = function() {
        return call_old(this, old.appendTo, arguments, 1, true);
    }
})();

/*!
 * jQuery Mousewheel 3.1.12
 *
 * Copyright 2014 jQuery Foundation and other contributors
 * Released under the MIT license.
 * http://jquery.org/license
 */

(function(factory) {
    if (typeof define === 'function' && define.amd) {
        // AMD. Register as an anonymous module.
        define(['jquery'], factory);
    } else if (typeof exports === 'object') {
        // Node/CommonJS style for Browserify
        module.exports = factory;
    } else {
        // Browser globals
        factory(jQuery);
    }
}(function($) {

    var toFix = ['wheel', 'mousewheel', 'DOMMouseScroll', 'MozMousePixelScroll'],
        toBind = ('onwheel' in document || document.documentMode >= 9) ?
        ['wheel'] : ['mousewheel', 'DomMouseScroll', 'MozMousePixelScroll'],
        slice = Array.prototype.slice,
        nullLowestDeltaTimeout, lowestDelta;

    if ($.event.fixHooks) {
        for (var i = toFix.length; i;) {
            $.event.fixHooks[toFix[--i]] = $.event.mouseHooks;
        }
    }

    var special = $.event.special.mousewheel = {
        version: '3.1.12',

        setup: function() {
            if (this.addEventListener) {
                for (var i = toBind.length; i;) {
                    this.addEventListener(toBind[--i], handler, false);
                }
            } else {
                this.onmousewheel = handler;
            }
            // Store the line height and page height for this particular element
            $.data(this, 'mousewheel-line-height', special.getLineHeight(this));
            $.data(this, 'mousewheel-page-height', special.getPageHeight(this));
        },

        teardown: function() {
            if (this.removeEventListener) {
                for (var i = toBind.length; i;) {
                    this.removeEventListener(toBind[--i], handler, false);
                }
            } else {
                this.onmousewheel = null;
            }
            // Clean up the data we added to the element
            $.removeData(this, 'mousewheel-line-height');
            $.removeData(this, 'mousewheel-page-height');
        },

        getLineHeight: function(elem) {
            var $elem = $(elem),
                $parent = $elem['offsetParent' in $.fn ? 'offsetParent' : 'parent']();
            if (!$parent.length) {
                $parent = $('body');
            }
            return parseInt($parent.css('fontSize'), 10) || parseInt($elem.css('fontSize'), 10) || 16;
        },

        getPageHeight: function(elem) {
            return $(elem).height();
        },

        settings: {
            adjustOldDeltas: true, // see shouldAdjustOldDeltas() below
            normalizeOffset: true // calls getBoundingClientRect for each event
        }
    };

    $.fn.extend({
        mousewheel: function(fn) {
            return fn ? this.bind('mousewheel', fn) : this.trigger('mousewheel');
        },

        unmousewheel: function(fn) {
            return this.unbind('mousewheel', fn);
        }
    });


    function handler(event) {
        var orgEvent = event || window.event,
            args = slice.call(arguments, 1),
            delta = 0,
            deltaX = 0,
            deltaY = 0,
            absDelta = 0,
            offsetX = 0,
            offsetY = 0;
        event = $.event.fix(orgEvent);
        event.type = 'mousewheel';

        // Old school scrollwheel delta
        if ('detail' in orgEvent) {
            deltaY = orgEvent.detail * -1;
        }
        if ('wheelDelta' in orgEvent) {
            deltaY = orgEvent.wheelDelta;
        }
        if ('wheelDeltaY' in orgEvent) {
            deltaY = orgEvent.wheelDeltaY;
        }
        if ('wheelDeltaX' in orgEvent) {
            deltaX = orgEvent.wheelDeltaX * -1;
        }

        // Firefox < 17 horizontal scrolling related to DOMMouseScroll event
        if ('axis' in orgEvent && orgEvent.axis === orgEvent.HORIZONTAL_AXIS) {
            deltaX = deltaY * -1;
            deltaY = 0;
        }

        // Set delta to be deltaY or deltaX if deltaY is 0 for backwards compatabilitiy
        delta = deltaY === 0 ? deltaX : deltaY;

        // New school wheel delta (wheel event)
        if ('deltaY' in orgEvent) {
            deltaY = orgEvent.deltaY * -1;
            delta = deltaY;
        }
        if ('deltaX' in orgEvent) {
            deltaX = orgEvent.deltaX;
            if (deltaY === 0) {
                delta = deltaX * -1;
            }
        }

        // No change actually happened, no reason to go any further
        if (deltaY === 0 && deltaX === 0) {
            return;
        }

        // Need to convert lines and pages to pixels if we aren't already in pixels
        // There are three delta modes:
        //   * deltaMode 0 is by pixels, nothing to do
        //   * deltaMode 1 is by lines
        //   * deltaMode 2 is by pages
        if (orgEvent.deltaMode === 1) {
            var lineHeight = $.data(this, 'mousewheel-line-height');
            delta *= lineHeight;
            deltaY *= lineHeight;
            deltaX *= lineHeight;
        } else if (orgEvent.deltaMode === 2) {
            var pageHeight = $.data(this, 'mousewheel-page-height');
            delta *= pageHeight;
            deltaY *= pageHeight;
            deltaX *= pageHeight;
        }

        // Store lowest absolute delta to normalize the delta values
        absDelta = Math.max(Math.abs(deltaY), Math.abs(deltaX));

        if (!lowestDelta || absDelta < lowestDelta) {
            lowestDelta = absDelta;

            // Adjust older deltas if necessary
            if (shouldAdjustOldDeltas(orgEvent, absDelta)) {
                lowestDelta /= 40;
            }
        }

        // Adjust older deltas if necessary
        if (shouldAdjustOldDeltas(orgEvent, absDelta)) {
            // Divide all the things by 40!
            delta /= 40;
            deltaX /= 40;
            deltaY /= 40;
        }

        // Get a whole, normalized value for the deltas
        delta = Math[delta >= 1 ? 'floor' : 'ceil'](delta / lowestDelta);
        deltaX = Math[deltaX >= 1 ? 'floor' : 'ceil'](deltaX / lowestDelta);
        deltaY = Math[deltaY >= 1 ? 'floor' : 'ceil'](deltaY / lowestDelta);

        // Normalise offsetX and offsetY properties
        if (special.settings.normalizeOffset && this.getBoundingClientRect) {
            var boundingRect = this.getBoundingClientRect();
            offsetX = event.clientX - boundingRect.left;
            offsetY = event.clientY - boundingRect.top;
        }

        // Add information to the event object
        event.deltaX = deltaX;
        event.deltaY = deltaY;
        event.deltaFactor = lowestDelta;
        event.offsetX = offsetX;
        event.offsetY = offsetY;
        // Go ahead and set deltaMode to 0 since we converted to pixels
        // Although this is a little odd since we overwrite the deltaX/Y
        // properties with normalized deltas.
        event.deltaMode = 0;

        // Add event and delta to the front of the arguments
        args.unshift(event, delta, deltaX, deltaY);

        // Clearout lowestDelta after sometime to better
        // handle multiple device types that give different
        // a different lowestDelta
        // Ex: trackpad = 3 and mouse wheel = 120
        if (nullLowestDeltaTimeout) {
            clearTimeout(nullLowestDeltaTimeout);
        }
        nullLowestDeltaTimeout = setTimeout(nullLowestDelta, 200);

        return ($.event.dispatch || $.event.handle).apply(this, args);
    }

    function nullLowestDelta() {
        lowestDelta = null;
    }

    function shouldAdjustOldDeltas(orgEvent, absDelta) {
        // If this is an older event and the delta is divisable by 120,
        // then we are assuming that the browser is treating this as an
        // older mouse wheel event and that we should divide the deltas
        // by 40 to try and get a more usable deltaFactor.
        // Side note, this actually impacts the reported scroll distance
        // in older browsers and can cause scrolling to be slower than native.
        // Turn this off by setting $.event.special.mousewheel.settings.adjustOldDeltas to false.
        return special.settings.adjustOldDeltas && orgEvent.type === 'mousewheel' && absDelta % 120 === 0;
    }

}));