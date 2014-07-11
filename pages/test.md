Headings are surrounded by Hashes
=================================

Markdown paragraphs are just
consecutive lines of text separated
by whitespace. You can use a number of
inline formattings, *emphasis is surrounded by asterisks*
or _underscores_, strong text by **double asterisks**
or __double underscores__.

##The number of hashes determine the heading level##

You can use “forbidden” HTML characters
like &, >, <, " and ', they are escaped
automatically. If you need additional HTML
formatting you can <span class="mySpan">just embed</span>
it into the Markdown source. You can include links
[easily](http://example.com "with a title") or
[without a title](http://foo.example.com). If you use links often,
[define them once][someid] and reference them by any id.
Just add the link definition anywhere, it will be removed
from the output:

[someid]: http://example.com "You can add a title"

You can add a horizontal ruler easily:

*******************************************

##Trailing hashes in a heading are optional

<div class="special"> If necessary you can even include
whole HTML blocks. Note however that Markdown code is *not*
evaluated inside HTML blocks, so if you want emphasis, you
have to <em>add it yourself</em>
</div>

You can embed verbatim code examples by indenting them by
four spaces or a tab:

    //This is verbatim code. Markdown syntax is *not*
    //processed here. However, HTML special chars are
    //escaped: < & > " '
    def foo() = <span>Hello World</span>

If you want verbatim code inline in your text
you can surround it with `def backticks():String`
or ``def doubleBackticks() = "To add ` to your code"``

* Unordered
* List items
* are started
* by asterisks

> To quote something, add ">" in front of
> the quoted text

1. Numbered Lists
2. start with a number and a "."
234. the numbering is
1. ignored in the output
250880. however and replaced
9. by consecutive numbers.

> To round things off
>
> * you can
> * nest lists
>
> #headings#
> > and quotes
> >
> >     //and code
> >
> > ##as much as you want##
> >
> > 1. also
> >     * a list
> >     * in a list
> > 2. isn't that cool?

``` java
//Actuarius now supports fenced code blocks!
System.out.println("Hello World!");
```