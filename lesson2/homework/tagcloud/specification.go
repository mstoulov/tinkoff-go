package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	dict map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() TagCloud {
	return TagCloud{make(map[string]int)}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (cloud *TagCloud) AddTag(tag string) {
	cloud.dict[tag]++
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (cloud *TagCloud) TopN(n int) []TagStat {
	var list []TagStat
	for key, value := range cloud.dict {
		list = append(list, TagStat{Tag: key, OccurrenceCount: value})
	}
	sort.Slice(list, func(l, r int) bool {
		return list[l].OccurrenceCount > list[r].OccurrenceCount
	})
	N := min(n, len(cloud.dict))
	return list[:N]
}
