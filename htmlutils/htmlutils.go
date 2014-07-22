package htmlutils

import (
	"code.google.com/p/go.net/html"
)

// Attr returns the value of the attribute. ok indicates if the attribute
// exists for the given node.
func Attr(n *html.Node, attrKey string) (val string, ok bool) {
	for _, attr := range n.Attr {
		if attr.Key == attrKey {
			return attr.Val, true
		}
	}
	return "", false
}

// SetAttr sets the value of the attribute for the given node
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

// GetTextContent returns the text within the given node
func GetTextContent(n *html.Node) string {
	child := n.FirstChild
	if child.Type == html.TextNode {
		return child.Data
	}
	return ""
}

// SetTextContent sets the text within the given node
func SetTextContent(n *html.Node, text string) {
	child := n.FirstChild
	if child.Type == html.TextNode {
		child.Data = text
	}
}

// GetElementByID returns the element with the given id, if one exists
func GetElementByID(doc *html.Node, id string) *html.Node {
	matches := Search(doc, func(n *html.Node) bool {
		nodeid, ok := Attr(n, "id")
		return ok && nodeid == id
	})
	if len(matches) == 1 {
		return matches[0]
	}
	return nil
}

func RemoveNode(doc *Fragment, n *html.Node) {

}

func ReplaceNodeWithFragment(doc *Fragment, node *html.Node, fragment *Fragment) {
	// Set all new nodes' parent
	fragment.eachNode(func(n *html.Node) {
		n.Parent = node.Parent
	})

	// Insert into linked list
	fragment.FirstNode.PrevSibling = node.PrevSibling
	fragment.LastNode.NextSibling = node.NextSibling

	// Update parent
	if node.Parent != nil {
		if node.Parent.FirstChild == node {
			node.Parent.FirstChild = fragment.FirstNode
		}
		if node.Parent.LastChild == node {
			node.Parent.LastChild = fragment.LastNode
		}
	}

	// Update doc
	if doc.FirstNode == node {
		doc.FirstNode = fragment.FirstNode
	}
	if doc.LastNode == node {
		doc.LastNode = fragment.LastNode
	}
}
