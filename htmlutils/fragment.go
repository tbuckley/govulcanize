package htmlutils

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"os"
)

type Fragment struct {
	FirstNode, LastNode *html.Node
}

// FromFile loads an Fragment from a file
func FromFile(filename string) (*Fragment, error) {
	context := &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	ns, err := html.ParseFragment(f, context)
	if err != nil {
		return nil, err
	}

	// Set the sibling chain
	for i, n := range ns {
		if i != 0 {
			n.PrevSibling = ns[i-1]
		}
		if i != len(ns)-1 {
			n.NextSibling = ns[i+1]
		}
	}

	return &Fragment{
		FirstNode: ns[0],
		LastNode:  ns[len(ns)-1],
	}, nil
}

func (f *Fragment) Search(pred HTMLPred) []*html.Node {
	matches := make([]*html.Node, 0)
	f.eachNode(func(n *html.Node) {
		submatches := Search(n, pred)
		matches = append(matches, submatches...)
	})
	return matches
}

func (f *Fragment) String() string {
	contents := ""
	f.eachNode(func(n *html.Node) {
		buf := new(bytes.Buffer)
		html.Render(buf, n)
		contents += buf.String()
	})
	return contents
}

func (f Fragment) eachNode(fn NodeFn) {
	for snaker := f.FirstNode; snaker != nil; snaker = snaker.NextSibling {
		fn(snaker)
	}
}
