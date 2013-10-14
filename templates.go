// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"
)

// Templates is a set of goblog templates.
type Templates map[string]*template.Template

// SiteData is a struct that contains all of the information necessary
// for generating a site page.
type SiteData struct {
	Helper
	Title       string
	Description string
	Author      string
	Content     string
	Languages   []string
	AtHome      bool
	AtTags      bool
	AtArchives  bool
	AtAbout     bool
}

// Helper is included with each template data. It allows the methods
// associated with this value to be run within the template.
type Helper bool

// Exec runs the given command and returns the combined output.
func (h Helper) Exec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// MakeAbout creates a complleted about HTML page and puts it into the
// given directory. It uses the template from about.html and will fill
// in the following values:
//
//        .CDate - The date the page was created.
//
// The results of that templating are then used as the content for
// calling MakeWebPage.
func (t Templates) MakeAbout(dir string) error {

	// Make the data that will be passed to the templater.
	data := struct {
		Helper
		CDate string
	}{
		CDate: time.Now().Format("2006-01-02"),
	}

	// Perform the templating
	content, err := ExecTemplate(t["about"], data)
	if err != nil {
		return err
	}

	// Make the pages with the siteData Helper Function
	return t.MakeWebPage(path.Join(dir, "about.html"), &SiteData{
		Title:      "About",
		Content:    content,
		AtHome:     false,
		AtTags:     false,
		AtArchives: false,
		AtAbout:    true,
	})

}

// MakeArchive creates a completed archives HTML page and puts it into
// the given directory. It uses the template from tags.html and will
// fill in the following values:
//
//      .CDate   - The date the page was created.
//      .Years   - A slice of Years that contain blog entries. Each one
//	               contains:
//        .Year   - The name of the Year (e.g. 2013).
//        .Months - A slice of months for this year that contains blog
//                  entries. Each one contains:
//          .Month   - The name of the month (e.g. January).
//          .Entries - A slice of blog entries for the given month of
//                     the given year. Each one contains:
//            .CDate   - The date of the blog entry.
//            .Url     - The url of the blog entry.
//            .Title   - The title of the blog entry.
//
// The results of that templating are then used as the content for
// calling MakeWebPage.
func (t Templates) MakeArchive(dir string, a []*YearEntries) error {

	// Make the data that will be passed to the templater.
	data := struct {
		Helper
		Years []*YearEntries
		CDate string
	}{
		Years: a,
		CDate: time.Now().Format("2006-01-02"),
	}

	// Perform the templating
	content, err := ExecTemplate(t["archive"], data)
	if err != nil {
		return err
	}

	// Make the pages with the siteData Helper Function
	return t.MakeWebPage(path.Join(dir, "archives.html"), &SiteData{
		Title:      "Archives",
		Content:    content,
		AtHome:     false,
		AtTags:     false,
		AtArchives: true,
		AtAbout:    false,
	})

}

// MakeIndex creates a completed index HTML page and puts it into the
// given directory. It uses the template from tags.html and will fill
// in the following values:
//
//      .Entries - A list of entries to display. Each one contains:
//        .CDate   - The date the entry was created.
//        .Title   - The title of the entry.
//        .UDate   - If the entry has changed since it's original
//                   creation, this will be the most recent update
//                   date.
//        .Content - The HTML formated Content of blog entry.
//        .Tags    - A list of tags (strings) for the blog entry.
//
// The results of that templating are then used as the content for
// calling MakeWebPage.
func (t Templates) MakeIndex(dir string, b []*Entry) error {

	// Make the HTML for each entry.
	entries := struct {
		Helper
		Entries []struct {
			*Entry
			Content string
		}
	}{
		Entries: []struct {
			*Entry
			Content string
		}{},
	}

	// Generate the entries list.
	languages := []string{}
	for _, blog := range b {
		// Make the content.
		c, err := blog.Parse()
		if err != nil {
			fmt.Println("parsing blog", blog, ":", err)
			os.Exit(1)
		}

		// Store the languages.
		for _, l := range blog.Languages {
			languages = append(languages, l)
		}

		// Save the blog and content.
		entries.Entries = append(entries.Entries, struct {
			*Entry
			Content string
		}{
			blog,
			c,
		})

	}

	// Get a unique list of languages.
	languages = removeDuplicates(languages)

	// Perform the templating
	content, err := ExecTemplate(t["entries"], entries)
	if err != nil {
		return err
	}

	// Make the pages with the siteData Helper Function
	return t.MakeWebPage(path.Join(dir, "index.html"), &SiteData{
		Title:      "Index",
		Content:    content,
		Languages:  languages,
		AtHome:     true,
		AtTags:     false,
		AtArchives: false,
		AtAbout:    false,
	})

}

// removeDuplicates is a helper function for the MakeIndex page. It
// removes duplicate languages.
func removeDuplicates(a []string) []string {
	result := []string{}
	seen := map[string]string{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

// MakeTags creates a completed tags HTML page and puts it into the
// given directory. It uses the template from tags.html and will fill
// in the following values:
//
//      .CDate - The date the page was created.
//      .Tags - A list of tags for the blog entry. Each one contains:
//         .Name - The name of the tag.
//         .Entries - A slice of blog entries for with the given tag.
//                    Each one contains:
//            .Url   - The url of the blog entry.
//            .Title - The title of the blog entry.
//
// The results of that templating are then used as the content for
// calling MakeWebPage.
func (t Templates) MakeTags(dir string, ta []*Tag) error {

	// Make the data that will be passed to the templater.
	data := struct {
		Helper
		Tags  []*Tag
		CDate string
	}{
		Tags:  ta,
		CDate: time.Now().Format("2006-01-02"),
	}

	// Perform the templating
	content, err := ExecTemplate(t["tags"], data)
	if err != nil {
		return err
	}

	// Make the pages with the siteData Helper Function
	return t.MakeWebPage(path.Join(dir, "tags.html"), &SiteData{
		Title:      "Tags",
		Content:    content,
		AtHome:     false,
		AtTags:     true,
		AtArchives: false,
		AtAbout:    false,
	})

}

// MakeEntry creates a completed HTML page of the given blog entry
// and puts it in the given directory. It uses the template from
// entry.html and will fill in the following values:
//
//      .CDate - The date the entry was created.
//      .Title   - The title of the entry.
//      .UDate   - If the entry has changed since it's original
//                 creation, this will be the most recent update
//                 date.
//      .Content - The HTML formated Content of blog entry.
//      .Tags    - A list of tags (strings) for the blog entry.
//
// The results of that templating are then used as the content for
// calling MakeWebPage.
func (t Templates) MakeEntry(dir string, blog *Entry,
	contents string) error {

	// Get the inner HTML.
	inner, err := t.makeBlogHelper(blog, contents)
	if err != nil {
		return nil
	}

	// Make the pages with the siteData Helper Function
	return t.MakeWebPage(path.Join(dir, blog.Url), &SiteData{
		Title:       blog.Title,
		Description: blog.Description,
		Author:      blog.Author,
		Content:     inner,
		Languages:   blog.Languages,
		AtHome:      false,
		AtTags:      false,
		AtArchives:  false,
		AtAbout:     false,
	})
}

// MakeWebPage write the resutls of using the site template on the
// given SiteData to the given file. This is the main function that
// makes complete web pages. It uses site.html to render the page and
// passes the following values:
//
//      .Title       - The title to use for this page.
//      .Description - The description of this page.
//      .Author      - The author of this page.
//      .Content     - The pages content.
//      .Languages   - A list of languages (string) used by the page.
//      .AtHome      - If true, the page is the index.html page.
//      .AtTags      - If true, the page is the index.html page.
//      .AtArchives  - If true, the page is the index.html page.
//      .AtAbout     - If true, the page is the index.html page.
func (t Templates) MakeWebPage(file string, sd *SiteData) error {
	// Get a file handle to write the contents to.
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t["site"].Execute(f, sd)
	if err != nil {
		return err
	}

	return nil
}

// makeBLogHelper is a helper function that generates the main content
// of a blog entry from the entry.html template.
func (t Templates) makeBlogHelper(blog *Entry,
	contents string) (string, error) {

	// Make the data that will be passed to the templater.
	templateData := struct {
		Helper
		*Entry
		Content string
	}{
		false,
		blog,
		contents,
	}

	// Perform the templating
	return ExecTemplate(t["entry"], templateData)
}

// LoadTemplates reads templates from the given directory and returns
// them as a map where the key is the template name and the value is
// the template itself. It currently looks for and loads the following
// templates with the following known dot values:
//
//  about.html - The about page of the site.
//    Variables:
//      .CDate   - The date the page was created.
//  archive.html - The archive page of the site.
//    Variables:
//  entries.html - Displays a list of entries.
//    Variables:
//      .Entries - A list of entries to display. Each one contains:
//        .CDate   - The date the entry was created.
//        .Title   - The title of the entry.
//        .UDate   - If the entry has changed since it's original
//                   creation, this will be the most recent update
//                   date.
//        .Content - The HTML formated Content of blog entry.
//        .Tags    - A list of tags (strings) for the blog entry.
//  entry.html - Display a single entry.
//    Variables:
//  site.html - The sites main template. All pages derive from this
//              template.
//    Variables:
//  tags.html - The sites list of tags.
//    Variables:
//
// All of the templates must exist for this to succeed.
func LoadTemplates(dir string) (Templates, error) {
	// This will be our return value.
	ret := make(Templates)

	// This is the list of templates to look for
	templates := []string{
		"about",
		"archive",
		"entries",
		"entry",
		"site",
		"tags",
	}

	// Process each template.
	for _, t := range templates {
		filename := path.Join(dir, t+".html")

		// Get the contents.
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		// Generate the template.
		tmplt, err := template.New(t).Parse(string(contents))
		if err != nil {
			return nil, err
		}

		// Save the template to the map.
		ret[t] = tmplt
	}

	return ret, nil
}

// ExecTemplate calls the Execute function on the given template and
// saves the results to the string. The given set of args should be a
// map of arguments within the template and their values.
func ExecTemplate(t *template.Template, args interface{}) (string,
	error) {

	sw := new(bytes.Buffer)

	err := t.Execute(sw, args)
	if err != nil {
		return "", err
	}

	return sw.String(), nil
}
