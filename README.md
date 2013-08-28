goblog
======

A static blog generator implemented in Go. This all started when
looking for a static blog tool. I asked around in my local Linux user
group and one of them joked you could write your own over a
weekend. Mission accomplished. 

Installation
============

You'll want go installed on your system, then you just need to 

    go get github.com/icub3d/goblog
	
This will install the binary *goblog* into *$GOPATH/bin*. In this case, you have to [set](http://golang.org/doc/code.html#GOPATH) properly GOPATH and GOROOT. Alternatively, you can download the source and do it yourself.

Usage
=====

All directory locations are configurable, but it is generally considered
wise to have a single place for your blog.

The default blog directory structure is this:

    $ ls
    blogs  public  static  templates

  * The *blogs* directory contains all of your blog entries. Each file
within that directory or any sub-directory that ends in *.md* will be
processed as a blog entry. Goblog uses markdown (like github), so feel
free to mark down your blog.

  * The *public* directory is where your generated code will go. You'll want
to point your web server to that location.

  * The *static* directory contains static assets like CSS, JavaScript,
images, etc that your blog needs to function.

  * The *templates* directory contains a list of html templates to use
when generating the site. Each of the templates use go's templating
system to display specific values. You can see
[my own blog](https://github.com/icub3d/joshua.themarshians.com) for
an example.

A *site.html* template is the template for every page. A
*archive.html* template is used for printing a list of all your blog
entries. A *about.html* template is used for displaying information
about yourself. A *entries.html* template is used to display multiple
blog entries on the *index.html* page. A *entry.html* template renders
a single blog entry. A *tags.html* template renders all of the blog
tags into a page.
