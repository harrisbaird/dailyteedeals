package utils_test

import (
	"sort"
	"testing"

	. "github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
)

func TestStringSlicesDiff(t *testing.T) {
	have := StringSlicesDiff([]string{"hello", "world", "test"},
		[]string{"hello", "world", "hello world"})
	want := []string{"hello world", "test"}
	sort.Strings(have)
	sort.Strings(want)
	st.Expect(t, have, want)
}

func TestStringSliceUnique(t *testing.T) {
	have := StringSliceUnique([]string{"hello", "world", "hello", "test", "TEST"})
	want := []string{"hello", "world", "test", "TEST"}
	sort.Strings(have)
	sort.Strings(want)
	st.Expect(t, have, want)
}

func TestStringSliceCleanup(t *testing.T) {
	have := StringSliceCleanup([]string{"", "  ", "hello", "world", "test", "TEST"})
	want := []string{"hello", "world", "test", "TEST"}
	sort.Strings(have)
	sort.Strings(want)
	st.Expect(t, have, want)
}
