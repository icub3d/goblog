// Copyright 2013 Joshua Marsh. All rights reserved.  Use of this
// source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"text/template"
	"time"
)

var sfeed = `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <lastBuildDate>{{.CreateDate}}</lastBuildDate> 
{{.ChannelContent}} 
{{range .Blogs}}    <item>
      <title>{{.Title}}</title>
      <link>%s{{.Url}}</link>
      <description>{{.Description}}</description>
      <pubDate>{{.PubDate}}</pubDate>
{{range .Tags}}      <category>{{.}}</category>
{{end}}    </item>
{{end}}
  </channel>
</rss>`

// MakeRss creates a completed feed.rss xml document and puts it into
// the given directory. It uses the template from channel.rss to
// populated the channel values except for the <item>s.
func MakeRss(entries []*Entry, url, tdir, dir string) error {

	// Get the channel data.
	channelContent, err := ioutil.ReadFile(path.Join(tdir, "channel.rss"))
	if err != nil {
		return err
	}

	// Make the feed template.
	// We need to get the URL to for the <links>
	if url == "" {
		// Try to get it from the channel.rss <link>
		re := regexp.MustCompile("<link>([^<]*)</link>")
		found := re.FindSubmatch(channelContent)
		if len(found) > 1 {
			url = string(found[1])
		}
	}

	feed := fmt.Sprintf(sfeed, url)

	var tmplt = template.Must(template.New("rss").Parse(feed))

	// Make the data that will be passed to the templater.
	data := struct {
		Blogs          []*Entry
		CreateDate     string
		ChannelContent string
	}{
		Blogs:          entries,
		CreateDate:     time.Now().Format(time.RFC822),
		ChannelContent: string(channelContent),
	}

	// Perform the templating
	sw := new(bytes.Buffer)
	err = tmplt.Execute(sw, data)
	if err != nil {
		return err
	}

	// Write out the file.
	err = ioutil.WriteFile(path.Join(dir, "feed.rss"), sw.Bytes(), 0644)

	return err
}
