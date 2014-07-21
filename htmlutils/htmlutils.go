package htmlutils

import (
	"code.google.com/p/go.net/html"
	"strings"
)

type HTMLPred func(*html.Node) bool

// Search will traverse the html tree and return a slice of all nodes that
// match the given predicate
func Search(n *html.Node, pred HTMLPred) []*html.Node {
	matches := make([]*html.Node, 0)

	// Add current node if it matches
	if pred(n) {
		matches = append(matches, n)
	}

	// Iterate over child nodes
	for snaker := n.FirstChild; snaker != nil; snaker = snaker.NextSibling {
		snakerMatches := Search(snaker, pred)
		matches = append(matches, snakerMatches...)
	}

	// Return matches
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

func Attr(n *html.Node, attrKey string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == attrKey {
			return attr.Val, true
		}
	}
	return "", false
}

func SetAttr(n *html.Node, attrKey string, attrValue string) {
	for i, attr := range n.Attr {
		if attr.Key == attrKey {
			attr.Val = attrValue
			n.Attr[i] = attr
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{
		Key: attrKey,
		Val: attrValue,
	})
}

func GetTextContent(n *html.Node) string {
	child := n.FirstChild
	if child.Type == html.TextNode {
		return child.Data
	}
	return ""
}

func SetTextContent(n *html.Node, text string) {
	child := n.FirstChild
	if child.Type == html.TextNode {
		child.Data = text
	}
}

func GetElementById(doc *html.Node, id string) *html.Node {
	matches := Search(doc, func(n *html.Node) bool {
		nodeid, ok := Attr(n, "id")
		return ok && nodeid == id
	})
	if len(matches) == 1 {
		return matches[0]
	}
	return nil
}

func ParseFragment(fragment string) ([]*html.Node, error) {
	context := &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	}
	r := strings.NewReader(fragment)
	ns, err := html.ParseFragment(r, context)
	if err != nil {
		return nil, err
	}

	// Set the chain
	for i, n := range ns {
		if i != 0 {
			n.PrevSibling = ns[i-1]
		}
		if i != len(ns)-1 {
			n.NextSibling = ns[i+1]
		}
	}

	return ns, nil
}

func Replace(old, ns *[]html.Node) {
	parent := old.Parent
	prev := old.PrevSibling
	next := old.NextSibling

	first := ns[0]
	last := ns[len(ns)-1]

	// Set all new nodes' parent
	for n := range ns {
		n.Parent = parent
	}

	// Insert into linked list
	first.PrevSibling = prev
	last.NextSibling = next

	if parent.FirstChild == old {
		parent.FirstChild = n
	}
	if parent.LastChild == old {
		parent.LastChild = n
	}
}
