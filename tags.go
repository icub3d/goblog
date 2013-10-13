// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"sort"
)

// Tag is a representation of a tag and it's associated entries. It is
// used as a storage mechanism for the tags page.
type Tag struct {
	// The name of the Tag.
	Name string

	// The list of Entries associated with this tag.
	Entries []*BlogEntry
}

// Add links the given BlogEntry to this Tag.
func (t *Tag) Add(e *BlogEntry) {
	// If it needs to be initialized, do that now.
	if t.Entries == nil {
		t.Entries = make([]*BlogEntry, 0, 0)
	}

	t.Entries = append(t.Entries, e)
}

// Tags is a map of Tags structures with some methods for easily
// adding blog entries. It also has the ability to export the entries
// as a sorted list for output.
type Tags map[string]*Tag

// GetTags builds a Tags from the given list of blogs.
func GetTags(entries []*BlogEntry) Tags {
	t := make(Tags)

	for _, blog := range entries {
		t.Add(blog)
	}

	return t
}

// Add links the given BlogEntry to all of it's tags.
func (t Tags) Add(e *BlogEntry) {
	for _, tag := range e.Tags {
		f, ok := t[tag]
		if !ok {
			// It wasn't found, so create one.
			t[tag] = &Tag{
				Name:    tag,
				Entries: []*BlogEntry{e},
			}
		} else {
			// Add it to the one we found.
			f.Add(e)
		}
	}
}

// Slice returns the Tags as a slice. The
// list is in sorted order.
func (t Tags) Slice() TagSlice {
	s := make(TagSlice, 0, len(t))

	for _, t := range t {
		s = append(s, t)
	}

	// Sort the slice.
	sort.Sort(s)

	return s
}

// TagsSlice is a slice of Tags this is returns by the Slice()
// function for Tags. It implements the sorting interface for go's
// sort package.
type TagSlice []*Tag

// Len returns the length of the TagSlice.
func (t TagSlice) Len() int {
	return len(t)
}

// Less returns true if the value at i is less than the value at j.
func (t TagSlice) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Swap switches the elemens at i and j.
func (t TagSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
