Goblog
======

A static blog generator implemented in Go. This all started when looking for a static blog tool. I asked around in my local Linux user group and one of them joked you could write your own over a weekend. Mission accomplished.

Installation
------------

You'll want go installed on your system, then you just need to 

    go get github.com/icub3d/goblog
	
This will install the binary `goblog` into `$GOPATH/bin`. In this case, you
have to [set](http://golang.org/doc/code.html#GOPATH) properly `GOPATH` and
`GOROOT`. Alternatively, you can download the source and do it yourself.

Usage
-----

All directory locations are configurable, but it is generally considered
wise to have a single place for your blog.

The default blog directory structure is this:

    $ ls
    blogs  public  static  templates

  * The `blogs` directory contains all of your blog entries. Each file
within that directory or any sub-directory that ends in `.md` will be
processed as a blog entry. The rest of extensions will be ignored. Goblog uses markdown
(like github), so feel free to mark down your blog.

  * The `public` directory is where your generated code will go. You'll want
to point your web server to that location.

  * The `static` directory contains static assets like CSS, JavaScript,
images, etc that your blog needs to function.

  * The `templates` directory contains a list of html templates to use
when generating the site. Each of the templates use Go's [templating system](http://golang.org/pkg/text/template/)
to display specific values. You can see [my own blog](https://github.com/icub3d/joshua.themarshians.com) for
an example.

Templates
---------

Each of the following templates are required. Without them, the system cannot generate the static pages.

  * The `site.html` template is the template for every page.
  * The `archive.html` template is used for printing a list of all your blog entries.
  * The `about.html` template is used for displaying information about yourself.
  * The `entries.html` template is used to display multiple blog entries on the *index.html* page.
  * The `entry.html` template renders a single blog entry.
  * The `tags.html` template renders all of the blog tags into a page.

Each template is rendered using Go's standard text/template library. When designing your templates, you can reference the documentation for the [templates package](http://godoc.org/github.com/icub3d/goblog/templates). For example, the _entry.html_ maps to the [MakeBlogEntry](http://godoc.org/github.com/icub3d/goblog/templates#Templates.MakeBlogEntry) function. In your _entry.html_ template, you'd put _{{.Title}}_ where you expect the title of the blog entry to go. You can see an example at my own [entry.html](https://github.com/icub3d/joshua.themarshians.com/blob/master/templates/entry.html).

As a special case, the templating engine has some helper functions:

  * _.Exec_ - you can include the output of a command line using this
    helper. For example, `{{.Exec "/bin/bash" "-c" "git log
    --pretty=oneline | wc -l "}}` would return the results of the
    executed command which might be something like _30_.

Blog Entry Meta Data
--------------------

You could specify meta data about a blog (title, created, etc.). Contrary to
other static blog generator which use YAML headers within markdown file for
setting fields, Goblog use XHTML comments for that. That is, if we want to set a tag with name `Foo` with `bar` value, we have to write:

    <!-- Foo: bar -->

An example can be found at [goblog.md](https://raw.github.com/icub3d/joshua.themarshians.com/master/blogs/goblog.md).
    
The meta data fields which Goblog recognizes are: 

  * `Title`: Title of the post, without quotes. Example: `Title: This if my first post`
  * `Author`: The author of the post. Example: `Author: Joshua Marsh`
  * `Description`: A brief description of the blog. This will be used by things like RSS feeds. Example: `Description: This is my first blog entry!`
  * `Languages`: This is the language the entry is in. This can be used to set html headers in your templates. Example: `Languages: en`
  * `Tags`: A list of tags. Example: `Tags: linux, oss, informatics`
  * `Created`: Data of creation of the post. The format of the date is YYYY-MM-DD. If this is not set, it will default to the timestamp of the file on the file system. Example: `Created: 2013-07-18`
  * `Updated`: Data of last update of the post. The format of the date is YYYY-MM-DD. If this is not set, it will default to the timestamp of the file on the file system. Example: `Updated: 2013-07-18`

All of these values are optional. They are mapped to your template. If you don't specify them in your blog entry but have them in your templates, then they obviously won't show up. You should try to specify all the values your templates have in them to make your site appear normal.
