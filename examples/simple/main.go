package main

import (
	"fmt"

	"github.com/if1live/maya"
	"github.com/op/go-logging"
)

func main() {
	logging.SetLevel(logging.CRITICAL, "maya")

	intext := `gist sample begin
~~~maya:gist
id=b23494b9e42ae89e6f28
file=factorial.sh
~~~
gist sample end`
	article := maya.NewArticle(intext, "empty")
	outtext := article.OutputString()
	fmt.Println(outtext)
}
