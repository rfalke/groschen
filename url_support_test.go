package groschen

import (
	. "groschen/hamcrest"
	"testing"
)

func TestMakeLinkAbsolute(t *testing.T) {
	base := "http://www.example.com/abc/dec"
	Assert(t).That(MakeLinkAbsolute(base, "http://www.otherhost/foo"), Equals("http://www.otherhost/foo"))
	Assert(t).That(MakeLinkAbsolute(base, "/foo/bar"), Equals("http://www.example.com/foo/bar"))
	Assert(t).That(MakeLinkAbsolute(base, "foo/bar"), Equals("http://www.example.com/abc/foo/bar"))

	Assert(t).That(MakeLinkAbsolute(base, "http://www.otherhost/foo#xyz"), Equals("http://www.otherhost/foo"))
	Assert(t).That(MakeLinkAbsolute(base, "/foo/bar#xyz"), Equals("http://www.example.com/foo/bar"))
	Assert(t).That(MakeLinkAbsolute(base, "foo/bar#xyz"), Equals("http://www.example.com/abc/foo/bar"))

	Assert(t).That(MakeLinkAbsolute(base, "/foo/././bar"), Equals("http://www.example.com/foo/bar"))
	Assert(t).That(MakeLinkAbsolute(base, "/foo/."), Equals("http://www.example.com/foo/DOT"))

	Assert(t).That(MakeLinkAbsolute(base, "http://otherhost"), Equals("http://otherhost/"))

	Assert(t).That(MakeLinkAbsolute(base, "//newhost/some/path"), Equals("http://newhost/some/path"))
	Assert(t).That(MakeLinkAbsolute(base, "//newhost"), Equals("http://newhost/"))

	Assert(t).That(MakeLinkAbsolute(base, "/a/b/c/./../../g?foo=.."), Equals("http://www.example.com/a/g?foo=.."))
	Assert(t).That(MakeLinkAbsolute(base, "mid/content=5/../6"), Equals("http://www.example.com/abc/mid/6"))

	Assert(t).That(MakeLinkAbsolute(base, "/foo/bar/.."), Equals("http://www.example.com/foo"))
	Assert(t).That(MakeLinkAbsolute(base, "/foo/bar/../.."), Equals("http://www.example.com"))

	Assert(t).That(MakeLinkAbsolute(base, "foo/bar/.."), Equals("http://www.example.com/abc/foo"))
	Assert(t).That(MakeLinkAbsolute(base, "foo/bar/../.."), Equals("http://www.example.com/abc"))
	Assert(t).That(MakeLinkAbsolute(base, "foo/bar/../../.."), Equals("http://www.example.com"))
	Assert(t).That(MakeLinkAbsolute(base, "foo/bar/../../../.."), Equals("http://www.example.com"))
	Assert(t).That(MakeLinkAbsolute(base, "foo/.."), Equals("http://www.example.com/abc"))
	Assert(t).That(MakeLinkAbsolute(base, ".."), Equals("http://www.example.com"))
	Assert(t).That(MakeLinkAbsolute(base, "../.."), Equals("http://www.example.com"))

	Assert(t).That(MakeLinkAbsolute("http://www.example.com", "../.."), Equals("http://www.example.com"))
	Assert(t).That(MakeLinkAbsolute("http://www.example.com ", "abc"), Equals("http://www.example.com/abc"))
}

func TestSupportedUrl(t *testing.T) {
	Assert(t).That(SupportedUrl("#"), Equals(false))
	Assert(t).That(SupportedUrl("mailto:abc"), Equals(false))
	Assert(t).That(SupportedUrl("javascript:abc"), Equals(false))
	Assert(t).That(SupportedUrl("http:///www.regulations.gov/"), Equals(false))
	Assert(t).That(SupportedUrl("ftp://ftp.lrz.de:21/"), Equals(false))

	Assert(t).That(SupportedUrl("/abc"), Equals(true))
	Assert(t).That(SupportedUrl("http://www.example.com/abc"), Equals(true))
}
