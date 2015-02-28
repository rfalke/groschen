package groschen

import (
	"regexp"
	"strings"
)

func is_mailto(url string) bool {
	return strings.HasPrefix(url, "mailto:")
}

func is_javascript(url string) bool {
	return strings.HasPrefix(url, "javascript:")
}

func is_ftp(url string) bool {
	return strings.HasPrefix(url, "ftp:")
}

func is_bad_http(url string) bool {
	return strings.HasPrefix(url, "http:///")
}

func SupportedUrl(url string) bool {
	return url != "" && url != "#" && !(is_mailto(url) || is_javascript(url) || is_ftp(url) || is_bad_http(url))
}
func remove_anchor(url string) string {
	if url == "#" {
		return url
	}
	return strings.Split(url, "#")[0]
}

type FullUrl struct {
	proto string
	host  string
	port  string
	path  string
	query string
}

func split_full(url string) *FullUrl {
	//	fmt.Printf("  split_full('%s')\n", url)
	result := regexp.MustCompile("^(https?://)([a-zA-Z0-9-.]+)(:[0-9]+)?([^?]*)(.*)$").FindStringSubmatch(url)
	if result == nil {
		//		fmt.Printf("  split_full  = not full url\n")
		return nil
	}

	//	fmt.Printf("    split_full(%q) = %q\n", url, result[1:])
	Proto, Host, Port, Path, Query := result[1], result[2], result[3], result[4], result[5]
	if Path == "" {
		Path = "/"
	}
	r := new(FullUrl)
	r.proto = Proto
	r.host = Host
	r.port = Port
	r.path = Path
	r.query = Query

	return r
}

func split_url(url string) []string {
	//	fmt.Printf("  split_url(%q)\n", url)
	decomposed := split_full(url)
	var result []string = nil
	if decomposed != nil {
		Path, Name := rsplit(decomposed.path, "/")
		result = []string{"full", decomposed.proto + decomposed.host + decomposed.port, Path, Name + decomposed.query}
	} else {
		if strings.HasPrefix(url, "//") {
			if regexp.MustCompile("^//[a-zA-Z0-9-.]+(:[0-9]+)?/.*$").MatchString(url) {
				result = []string{"withoutProto", url}
			} else {
				result = []string{"withoutProto", url + "/"}
			}
		} else {
			if strings.HasPrefix(url, "/") {
				result = []string{"absolute"}
			} else {
				result = []string{"relative"}
			}
		}
	}
	//	fmt.Printf("  split_url  = %q\n", result)
	return result
}

func rsplit(s, sep string) (base, tail string) {
	value := strings.SplitN(Reverse(s), sep, 2)
	if len(value) != 2 {
		panic("foobar")
	}
	tail, base_ := value[0], value[1]
	return Reverse(base_) + "/", Reverse(tail)
}

func normalize_url(url string) string {
	url = replace_dot_at_end(url)
	decomposed := split_full(url)
	if decomposed == nil {
		return url
	} else {
		return decomposed.proto + decomposed.host + decomposed.port + removeDotSegmentsInPath(decomposed.path) + decomposed.query
	}
}

// http://tools.ietf.org/html/rfc3986#section-5.2.4
func removeDotSegmentsInPath(path string) string {
	result := make([]string, 0)
	for _, part := range strings.Split(path, "/") {
		if part != "." {
			if part == ".." {
				result = result[:len(result)-1]
			} else {
				result = append(result, part)
			}
		}
	}
	return strings.Join(result, "/")
}

func replace_dot_at_end(url string) string {
	if strings.HasSuffix(url, "/.") {
		return strings.TrimSuffix(url, "/.") + "/DOT"
	}
	return url
}

func MakeLinkAbsolute(baseUrl string, url string) string {
	//	fmt.Printf("MakeLinkAbsolute(base='%s', url='%s')\n", baseUrl, url)
	tmp := split_url(baseUrl)
	urlType, host, path, _ := tmp[0], tmp[1], tmp[2], tmp[3]
	if urlType != "full" {
		panic("foobar")
	}
	url = remove_anchor(url)
	parts := split_url(url)
	result := ""
	if parts[0] == "full" {
		result = parts[1] + parts[2] + parts[3]
	} else if parts[0] == "withoutProto" {
		result = get_proto(baseUrl) + parts[1]
	} else if parts[0] == "absolute" {
		result = host + url
	} else if parts[0] == "relative" {
		result = host + path + url
	} else {
		panic("foobar")
	}
	result = normalize_url(result)
	//	fmt.Printf("MakeLinkAbsolute  = '%s'\n", result)
	return result
}

func get_proto(url string) string {
	return strings.TrimSuffix(split_full(url).proto, "//")
}

func IsFullUrl(url string) bool {
	return split_url(url) != nil
}
