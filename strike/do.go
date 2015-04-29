package strike

import (
	"github.com/blinkat/blinks/strike/css"
	"github.com/blinkat/blinks/strike/html"
	"github.com/blinkat/blinks/strike/js"
	"io/ioutil"
	//"github.com/blinkat/blinks/strike/json"
)

// ---------- [ text ] -----------
func StrikeJsText(text string) string {
	return js.Strike(text)
}

func StrikeCssText(text string) string {
	return css.Strike(text)
}

func StrikeHtmlText(text string) string {
	return html.Strike(text)
}

// -----------[ file ] ------------
func StrikeJsPath(path string) string {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return StrikeJsText(string(buf))
}

func StrikeCssPath(path string) string {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return StrikeCssText(string(buf))
}

func StrikeHtmlPath(path string) string {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return StrikeHtmlText(string(buf))
}
