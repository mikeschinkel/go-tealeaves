module github.com/mikeschinkel/go-tealeaves/site/examples/syntax-highlighting

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teahilite v0.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teahilite => ../../../teahilite
)
