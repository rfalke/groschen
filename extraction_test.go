package groschen

import (
	. "groschen/hamcrest"
	"testing"
)

func TestExtraction(t *testing.T) {
	Assert(t).ThatStringSlice(ExtractLinks("http://www.example.com/abc/dec",
		"HREF=\"link1\" href=\"mailto:foo@bar.com\" <img src=\"http://www.foobar.com/pixel.gif\" /> <body background=\"/background.png\""),
		ContainsExactly("http://www.example.com/abc/link1", "http://www.example.com/background.png", "http://www.foobar.com/pixel.gif"))
	Assert(t).ThatStringSlice(ExtractLinks("http://www.example.com/abc/dec", "HREF=\"link1\" href=\"link1\""),
		ContainsExactly("http://www.example.com/abc/link1"))
	Assert(t).ThatStringSlice(ExtractLinks("http://www.example.com/abc/dec", "href=\"mailto:foo@bar.com\""), IsEmpty())
	Assert(t).ThatStringSlice(ExtractLinks("http://www.example.com/abc/dec", "href=\"javascript:void(0)\""), IsEmpty())
	Assert(t).ThatStringSlice(ExtractLinks("http://www.example.com/abc/dec", "href=\"#\""), IsEmpty())
	Assert(t).ThatStringSlice(ExtractLinks("http://www.example.com/abc/dec", "some text hr_ef=\"link1\""), IsEmpty())
}
