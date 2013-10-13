// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"sort"
	"strconv"
)

// MonthEntries is a list of entries and their associated month. It
// implements the sort interface for sorting the entries by date
// descending.
type MonthEntries struct {
	// The name of the month.
	Month string

	// A list of blog entries for this month.
	Entries []*BlogEntry
}

// Add appends the given BlogEntry to the Entries list. It doesn't
// check to see if the given entry was actually in this month.
func (me *MonthEntries) Add(e *BlogEntry) {
	if me.Entries == nil {
		me.Entries = make([]*BlogEntry, 0, 0)
	}

	me.Entries = append(me.Entries, e)
}

// Len returns the length of the MonthEntries.
func (me *MonthEntries) Len() int {
	return len(me.Entries)
}

// Less returns true if the entry at j is newer than the entry at i.
func (me *MonthEntries) Less(i, j int) bool {
	return me.Entries[j].Created.Before(me.Entries[i].Created)
}

// Swap switches the elements at i and j.
func (me *MonthEntries) Swap(i, j int) {
	me.Entries[i], me.Entries[j] = me.Entries[j], me.Entries[i]
}

type YearEntriesSlice []*YearEntries

// Len returns the length of the YearEntriesSlice.
func (ye YearEntriesSlice) Len() int {
	return len(ye)
}

// Less returns true if the value at i is newer than the value at j.
func (ye YearEntriesSlice) Less(i, j int) bool {
	first, err := strconv.Atoi(ye[j].Year)
	if err != nil {
		return true
	}

	second, err := strconv.Atoi(ye[i].Year)
	if err != nil {
		return true
	}

	return first < second
}

// Swap switches the elements at i and j.
func (ye YearEntriesSlice) Swap(i, j int) {
	ye[i], ye[j] = ye[j], ye[i]
}

// This is used to simplify the Less function.
var mstrs map[string]int = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

// YearEntries is a list of entries and their associated year. It
// implements the sort interface for sorting the MonthEntries by date
// descending.
type YearEntries struct {
	// The name of the year.
	Year string

	// A list of MonthEntries for this year.
	Months []*MonthEntries
}

// Add appends the given BlogEntry to the Months list for the given
// month. It doesn't check to see if the given entry was actually in
// this year/month.
func (ye *YearEntries) Add(month string, e *BlogEntry) {
	// Make it if we don't have one.
	if ye.Months == nil {
		ye.Months = make([]*MonthEntries, 0, 0)
	}

	// Find the right month.
	var index int = -1
	for k, v := range ye.Months {
		if v.Month == month {
			index = k
			break
		}
	}

	if index == -1 {
		// We didn't find a month, so let's make it
		ye.Months = append(ye.Months, &MonthEntries{
			Month:   month,
			Entries: []*BlogEntry{e},
		})
	} else {
		ye.Months[index].Add(e)
	}
}

// Len returns the length of the YearEntries.
func (ye YearEntries) Len() int {
	return len(ye.Months)
}

// Less returns true if the value at i is newer than the value at j.
func (ye YearEntries) Less(i, j int) bool {
	return mstrs[ye.Months[j].Month] < mstrs[ye.Months[i].Month]
}

// Swap switches the elements at i and j.
func (ye YearEntries) Swap(i, j int) {
	ye.Months[i], ye.Months[j] = ye.Months[j], ye.Months[i]
}

// DateEntries is a map that stores blog entries by year and month.
type DateEntries map[string]*YearEntries

// GetMost Recent returns up to the max most recent entries from the
// given (and hopefully sorted) list of YearEntries.
func GetMostRecent(y []*YearEntries, max int) []*BlogEntry {
	b := make([]*BlogEntry, 0, max)

	for _, year := range y {
		if len(b) == max {
			break
		}

		for _, month := range year.Months {
			if len(b) == max {
				break
			}

			for _, entry := range month.Entries {
				if len(b) == max {
					break
				}

				b = append(b, entry)
			}
		}
	}

	return b
}

// ParseBlogs creates a DateEntries from the given list of blogs.
func ParseBlogs(entries []*BlogEntry) DateEntries {
	t := make(DateEntries)

	for _, blog := range entries {
		year := blog.Created.Format("2006")
		month := blog.Created.Format("January")
		t.Add(year, month, blog)
	}

	return t
}

// Add stores the given BlogEntry under the given year and month.
func (de DateEntries) Add(year, month string, e *BlogEntry) {
	y, ok := de[year]
	if !ok {
		// We need to create it.
		de[year] = &YearEntries{
			Year: year,
			Months: []*MonthEntries{
				&MonthEntries{
					Month:   month,
					Entries: []*BlogEntry{e},
				},
			},
		}

		return
	}

	y.Add(month, e)
}

// Slice returns the Year, Month, and BlogEntries as a
// slice which is suitable for transformation in the
// templates. The year and months are sorted.
func (de DateEntries) Slice() YearEntriesSlice {
	s := make(YearEntriesSlice, 0, len(de))

	for _, y := range de {
		// Sort the months in the year.
		sort.Sort(y)

		// Sort each entry in each month.
		for _, m := range y.Months {
			sort.Sort(m)
		}

		// Add the year to our list.
		s = append(s, y)
	}

	// Now sort our year.
	sort.Sort(s)

	return s
}
