package models_ext_test

import (
	"testing"

	"strings"

	. "github.com/harrisbaird/dailyteedeals/models_ext"
	"github.com/nbio/st"
)

func TestTagQuery(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", []string{}, "tags && '{}'::text[] AND tags != '{}'"},
		{"valid", []string{"tag 1", "tag 2"}, "tags && '{tag1,tag2}'::text[] AND tags != '{}'"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			q := TagQuery("tags", tt.input, NormalizeTags)
			st.Expect(t, q, tt.want)
		})

	}
}

func TestNormalizeTags(t *testing.T) {
	tags := []string{"hello", "HELLO", "Hello", "World", "hello world!", "World-"}
	want := []string{"hello", "helloworld", "world"}
	st.Expect(t, NormalizeTags(tags), want)
}

func TestNormalizeUrls(t *testing.T) {
	urls := []string{
		"http://redbubble.com/people/arinesart",
		"https://redbubble.com/people/arinesart",
		"http://www.redbubble.com/people/arinesart",
		"http://www.redbubble.com/people/arinesart/",
		"http://www.redbubble.com/people/arinesart/shop",
		"http://www.redbubble.com/people/arinesart/shop?a=b",
		"http://www.redbubble.com/people/arinesart/shop?#test",
		"http://www.REDBUBBLE.com/people/arinesart",
	}
	want := []string{
		"http://redbubble.com/people/arinesart",
		"http://redbubble.com/people/arinesart/shop",
		"http://redbubble.com/people/arinesart/shop?a=b"}

	st.Expect(t, NormalizeUrls(urls), want)
}

func TestMakeSlug(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{"1", "Test Slug", "test-slug"},
		{"2", "25", "25"},
		{"3", "日本", "ri-ben-"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			slug := MakeSlug(tt.input)
			st.Expect(t, strings.Contains(slug, tt.output), true)
			st.Expect(t, ValidSlug.MatchString(slug), true)
		})
	}
}
