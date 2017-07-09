package modext

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/purell"
	"github.com/harrisbaird/dailyteedeals/utils"
	"github.com/rainycape/unidecode"
)

// Normalizer is the type of Normalizer
type Normalizer func([]string) []string

// TagQuery queries and normalizes tags stored in a Postgres array column.
// Returns a sqlboiler QueryMod.
func TagQuery(field string, tags []string, fn Normalizer) string {
	tq := []string{fmt.Sprintf("%s && '{%s}'::text[]", field, strings.Join(fn(tags), ",")), fmt.Sprintf("%s != '{}'", field)}
	return strings.Join(tq, " AND ")
}

// NormalizeUrls attempts to normalize urls
// Returns a string slice of unique, normalized urls
func NormalizeUrls(urls []string) []string {
	output := []string{}
	for _, url := range urls {
		normalized, err := purell.NormalizeURLString(url, purell.FlagsAllGreedy)
		if err == nil && strings.Contains(normalized, "http") {
			output = append(output, normalized)
		}
	}
	output = utils.StringSliceUnique(output)
	sort.Strings(output)
	return output
}

// NormalizeTags attempts to normalize tags
// Returns a string slice of unique, normalized tags
func NormalizeTags(tags []string) []string {
	for i, tag := range tags {
		re := regexp.MustCompile(`\W`)
		tag = re.ReplaceAllString(tag, "")
		tag = strings.ToLower(tag)
		tags[i] = strings.TrimSpace(tag)
	}
	tags = utils.StringSliceUnique(tags)
	tags = utils.StringSliceCleanup(tags)
	sort.Strings(tags)
	return tags
}

// MakeSlug creates a slug in the format: 55555-my-slug-name.
// Fairly crude but ensures a unique slug without
// having to check the database.
func MakeSlug(slug string) string {
	slug = unidecode.Unidecode(slug)
	slug = strings.Replace(slug, " ", "-", -1)
	slug = strings.ToLower(slug)

	id := rangeIn(10000, 99999)

	return fmt.Sprintf("%d-%s", id, slug)
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
