package htmlutils

import (
	"code.google.com/p/go.net/html"
)

type HTMLPred func(*html.Node) bool

type NodeFn func(*html.Node)

// dfs runs a depth-first search over the HTML tree
func DFS(n *html.Node, pre NodeFn, post NodeFn) {
	if pre != nil {
		pre(n)
	}
	for snaker := n.FirstChild; snaker != nil; snaker = snaker.NextSibling {
		DFS(snaker, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// Search will traverse the html tree and return a slice of all nodes that
// match the given predicate
func Search(n *html.Node, pred HTMLPred) []*html.Node {
	matches := make([]*html.Node, 0)
	DFS(n, func(n *html.Node) {
		if pred(n) {
			matches = append(matches, n)
		}
	}, nil)
	return matches
}

// Return a predicate that checks to see if the given node has any of the
// attributes
func HasAnyAttr(attrKeys ...string) HTMLPred {
	attrMap := make(map[string]bool)
	for _, attrKey := range attrKeys {
		attrMap[attrKey] = true
	}

	return func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if _, ok := attrMap[attr.Key]; ok {
				return true
			}
		}
		return false
	}
}

func HasAttrValue(attrKey, attrValue string) HTMLPred {
	return func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Key == attrKey && attr.Val == attrValue {
				return true
			}
		}
		return false
	}
}

func HasTagname(tagname string) HTMLPred {
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tagname
	}
}

func And(preds ...HTMLPred) HTMLPred {
	return func(n *html.Node) bool {
		for _, p := range preds {
			if !p(n) {
				return false
			}
		}
		return true
	}
}

func Or(preds ...HTMLPred) HTMLPred {
	return func(n *html.Node) bool {
		for _, p := range preds {
			if p(n) {
				return true
			}
		}
		return false
	}
}
