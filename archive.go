// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"sort"
	"time"
)

// GetEntriesByDate sorts the given entries by date.
func GetEntriesByDate(entries []*Entry) EntriesByDate {

	es := make(EntriesByDate, 0, len(entries))
	for _, e := range entries {
		es = append(es, e)
	}

	sort.Sort(es)

	return es
}

// EntriesByDate is a slice of entries which is sortable using go's
// sort package by date, newest being first.
type EntriesByDate []*Entry

// Len returns the length of the EntriesByDate.
func (e EntriesByDate) Len() int {
	return len(e)
}

// Less returns true if the value at i is less than the value at j.
func (e EntriesByDate) Less(i, j int) bool {
	return e[i].Created.After(e[j].Created)
}

// Swap switches the elemens at i and j.
func (e EntriesByDate) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// GetArchives formats the given entries sorted by year in a slice of
// YearEntries suitable for making the archvie page.
func GetArchives(es EntriesByDate) []*YearEntries {
	yes := []*YearEntries{}
	curYear := -1
	var curMonth time.Month = -1

	for _, e := range es {
		// Create the Year and Month as necessary.
		if e.Created.Year() != curYear {
			yes = append(yes, &YearEntries{
				Year:   fmt.Sprintf("%v", e.Created.Year()),
				Months: []*MonthEntries{},
			})
			curYear = e.Created.Year()
			curMonth = -1
		}
		if e.Created.Month() != curMonth {
			yes[len(yes)-1].Months = append(yes[len(yes)-1].Months,
				&MonthEntries{
					Month:   e.Created.Month().String(),
					Entries: []*Entry{},
				})
			curMonth = e.Created.Month()
		}

		// Append this entry to the last month in the last year.
		yes[len(yes)-1].Months[len(yes[len(yes)-1].Months)-1].Entries =
			append(yes[len(yes)-1].Months[len(yes[len(yes)-1].Months)-1].Entries, e)
	}

	return yes
}

// MonthEntries is a list of entries and their associated month. It
// implements the sort interface for sorting the entries by date
// descending.
type MonthEntries struct {
	// The name of the month.
	Month string

	// A list of blog entries for this month.
	Entries []*Entry
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
