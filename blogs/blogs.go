// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// Package blog contains structures, methods and functions for
// manipulating blog entries.
package blogs

import (
	"bytes"
	"github.com/icub3d/goblog/fs"
	md "github.com/russross/blackfriday"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"time"
)

// BlogEntry is a representation of a blog entry. It contains the
// information necessary to generate the blogs contents and store
// information about the blog.
type BlogEntry struct {
	// Name is the name of the entry gleaned from the filename.
	Name string

	// Path is the path to the markdown file from the cwd.
	Path string

	// Aurhor is the name of the person who wrote the page.
	Author string

	// Title is the title of the Entry.
	Title string

	// Description is the description of the Entry.
	Description string

	// Url is the HTML file name of this entry (Name + ".html").
	Url string

	// Tags is a list of tags this blog entry contains. It is generated
	// when when the Parse method is called.
	Tags []string

	// Languages is a list of languages this blog entry contains. It is
	// generated when when the Parse method is called.
	Languages []string

	// Created is the date the blog entry was created. It is generated
	// when the Parse metod is called.
	Created time.Time

	// Updated is the date the blog entry was last updated. It is
	// generated when the Parse metod is called.
	Updated time.Time
}

// Parse reads the contents of the path for this BlogEntry. It gleans
// information from the file and saves it to this BlogEntry. It then
// formats the markdown to HTML and returns that.
func (be *BlogEntry) Parse() (string, error) {
	// Get the files contents.
	orgContents, err := ioutil.ReadFile(be.Path)
	if err != nil {
		return "", err
	}

	// Save some of the meta data.
	err = be.gleanInfo(string(orgContents))
	if err != nil {
		return "", err
	}

	// Return the markdown content.
	return string(md.MarkdownCommon(orgContents)), nil
}

// CDate is a helper function for the templating system that returns
// the Created date as a string or "" if there is no value.
func (be *BlogEntry) CDate() string {
	if be.Created.IsZero() {
		return ""
	}

	return be.Created.Format("2006-01-02")
}

// PubDate is a helper function for the templating system that returns
// the Created date as an RFC822 string or "" if there is no value.
func (be *BlogEntry) PubDate() string {
	if be.Created.IsZero() {
		return ""
	}

	return be.Created.Format("02 Jan 2006 15:04 MST")
}

// UDate is a helper function for the templating system that returns
// the Updated date as a string or "" if there is no value or if it's
// identical to the Created date.
func (be *BlogEntry) UDate() string {
	if be.Updated.IsZero() {
		return ""
	}

	if be.Created.Equal(be.Updated) {
		return ""
	}

	return be.Updated.Format("2006-01-02")
}

// GetBlogFiles looks in the given directory for blog entries and
// returns a list of them. Blog entries must have the '.md'
// extension. Entries are searched in the directory recursively. If a
// files is in a directory, the directory name is used as a prefix to
// the blog entries name concatenated with a '-'. The blog is not
// parsed or read. You should do that yourself elsewhere.
func GetBlogFiles(dir string) ([]*BlogEntry, error) {
	entries := []*BlogEntry{}

	// Read the list of entries for dir.
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p := path.Join(dir, file.Name())

		// If it's a directory, then recursively call this function and
		// merge the two slices.
		if file.IsDir() {
			blogs, err := GetBlogFiles(p)
			if err != nil {
				return nil, err
			}

			// Add the blog entries to this list.
			for _, blog := range blogs {
				newName, err := MakeBlogName(file.Name(), blog.Name)
				if err != nil {
					return nil, err
				}

				entries = append(entries, &BlogEntry{
					Name: newName,
					Url:  newName + ".html",
					Path: blog.Path,
				})

			}
		} else {
			// We only deal with .md (markdown) files.
			if path.Ext(p) != ".md" {
				continue
			}

			// Get the name of the entry.
			pieces := strings.Split(file.Name(), ".md")
			newName := pieces[0]

			// Just create the new entry.
			entries = append(entries, &BlogEntry{
				Name: newName,
				Url:  newName + ".html",
				Path: p,
			})
		}
	}

	return entries, nil
}

// MakeBlogName concatenates all of the given names with a dash and
// replaces illegal blog name characters with a dash.
func MakeBlogName(names ...string) (string, error) {
	re := regexp.MustCompile("[^-a-zA-Z0-9_]")
	buf := new(bytes.Buffer)

	end := len(names) - 1
	for i, name := range names {
		_, err := buf.Write(re.ReplaceAll([]byte(name), []byte("-")))
		if err != nil {
			return "", err
		}

		if i != end {
			_, err := buf.Write([]byte("-"))
			if err != nil {
				return "", err
			}

		}
	}

	return buf.String(), nil
}

// gleanInfo is a helper function that searches for various comments
// that contain useful information about the blog. It also uses fetchs
// the update and create dates.
func (be *BlogEntry) gleanInfo(contents string) error {
	/* These are the patterns we are searching for */
	var err error

	be.Title, err = regexSingle("Title", contents)
	if err != nil {
		return err
	}

	be.Author, err = regexSingle("Author", contents)
	if err != nil {
		return err
	}

	be.Description, err = regexSingle("Description", contents)
	if err != nil {
		return err
	}

	be.Languages, err = regexList("Languages", contents)
	if err != nil {
		return err
	}

	be.Tags, err = regexList("Tags", contents)
	if err != nil {
		return err
	}

	created, updated, err := fs.GetTimes(be.Path)
	if err != nil {
		return err
	}

	_created, err := regexSingle("Created", contents)
	be.Created, err = time.Parse("2006-01-02", _created)
	if err != nil {
		be.Created = created
	}

	_updated, err := regexSingle("Updated", contents)
	be.Updated, err = time.Parse("2006-01-02", _updated)
	if err != nil {
		be.Updated = updated
	}

	return nil
}

// regexList is a helper function that performs a regex search for an
// HTML comment with the given title. It returns the list (comma
// separated) of values.
func regexList(key, contents string) ([]string, error) {
	val, err := regexSingle(key, contents)
	if err != nil {
		return nil, err
	}

	if val == "" {
		return []string{}, nil
	}

	return strings.Split(val, ","), nil
}

// regexSingle is a helper function that performs a regex search for
// an HTML comment with the given title. It returns the value if it
// was found, or "" if it wasn't.
func regexSingle(key, contents string) (string, error) {
	pattern := "<!--[ ]*" + key + ":(.*)-->"

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(contents)
	if len(matches) < 2 {
		return "", nil
	}

	return strings.TrimSpace(matches[1]), nil
}
