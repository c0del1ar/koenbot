package libs

var XNXXSearchCache = map[string][]XNXXResult{}

type XNXXResult struct {
	Title string
	URL   string
}
