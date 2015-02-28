package groschen

import (
	"fmt"
	"regexp"
)

func findAll(body string, pattern string) []string {
	matches := regexp.MustCompile("(?i)"+pattern).FindAllStringSubmatch(body, -1)
	if matches == nil {
		return make([]string, 0)
	}
	result := make([]string, 0)
	for _, match := range matches {
		if len(match) != 2 {
			panic("foobar")
		}
		result = append(result, match[1])
	}
	return result
}

func ExtractLinks(baseUrl string, body string) []string {
	links := append(findAll(body, "(?:src)=\"([^\"]+)\""),
		findAll(body, "(?:href)=\"([^\"]+)\"")...)
	links = append(links, findAll(body, "(?:background)=\"([^\"]+)\"")...)
	//	fmt.Printf("got the following raw links: %q\n", links)
	absoluteLinks := make([]string, 0)
	for _, link := range links {
		if SupportedUrl(link) {
			absoluteLinks = append(absoluteLinks, MakeLinkAbsolute(baseUrl, link))
		}
	}
	goodLinks := make([]string, 0)
	badLinks := make([]string, 0)
	for _, link := range absoluteLinks {
		if IsFullUrl(link) {
			goodLinks = append(goodLinks, link)
		} else {
			badLinks = append(badLinks, link)
		}
	}
	if len(badLinks) > 0 {
		fmt.Printf("Ignore links which are not valid %q\n", badLinks)
	}
	return unique(goodLinks)
}

func unique(input []string) []string {
	asMap := make(map[string]bool, 0)
	for _, i := range input {
		asMap[i] = true
	}
	result := make([]string, 0)
	for i, _ := range asMap {
		result = append(result, i)
	}
	return result
}
