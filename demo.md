title: this is title
subtitle: this is subtitle
tags: python, demo
author: if1live
slug: sample-article

## write article

attach text file.

~~~maya:view
file=demo.py
lang=python
~~~

attach text file with line number.

~~~maya:view
file=demo.py
start_line=0
end_line=1
lang=python
~~~

print stdout/stderr as markdown code format.

~~~maya:execute
cmd=python demo.py
~~~

print stdout/stderr as markdown blockquote format.

~~~maya:execute
cmd=python demo.py
format=blockquote
~~~
