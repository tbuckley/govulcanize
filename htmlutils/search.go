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

// Closest will find the nearest parent that matches the given predicate
func Closest(n *html.Node, pred HTMLPred) *html.Node {
	for parent := n.Parent; parent != nil; parent = parent.Parent {
		if pred(parent) {
			return parent
		}
	}
	return nil
}

// HasAnyAttrP creates a predicate that checks whether a node has any of the attributes
func HasAnyAttrP(attrKeys ...string) HTMLPred {
	preds := make([]HTMLPred, 0, len(attrKeys))
	for _, key := range attrKeys {
		preds = append(preds, HasAttrP(key))
	}
	return OrP(preds...)
}

// HasAttrP creates a predicate that checks whether a node has the attribute
func HasAttrP(attrKey string) HTMLPred {
	return func(n *html.Node) bool {
		_, ok := Attr(n, attrKey)
		return ok
	}
}

// HasAttrValueP creates a predicate that checks whether a node has the attribute
// with the value
func HasAttrValueP(attrKey, attrValue string) HTMLPred {
	return func(n *html.Node) bool {
		val, ok := Attr(n, attrKey)
		return ok && val == attrValue
	}
}

// HasTagnameP creates a predicate that checks whether a node has the tagname
func HasTagnameP(tagname string) HTMLPred {
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tagname
	}
}

// AndP creates a predicate that checks whether all of the predicates are met
func AndP(preds ...HTMLPred) HTMLPred {
	return func(n *html.Node) bool {
		for _, p := range preds {
			if !p(n) {
				return false
			}
		}
		return true
	}
}

// OrP creates a predicate that checks whether any of the predicates are met
func OrP(preds ...HTMLPred) HTMLPred {
	return func(n *html.Node) bool {
		for _, p := range preds {
			if p(n) {
				return true
			}
		}
		return false
	}
}

// NotP creates a predicate that inverts the given predicate
func NotP(pred HTMLPred) HTMLPred {
	return func(n *html.Node) bool {
		return !pred(n)
	}
}
