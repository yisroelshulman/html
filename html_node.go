package html

import (
	"fmt"
	"sort"
	"strings"
)

type Props struct {
	props map[string]string
}

func (p Props) PropsToHTML() (string, error) {

	if len(p.props) == 0 {
		return "", nil
	}

	var keys []string
	for k := range p.props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var props strings.Builder
	for _, k := range keys {
		_, err := props.WriteString(fmt.Sprintf(" %s=\"%s\"", k, p.props[k]))
		if err != nil {
			return "", err
		}
	}
	return props.String(), nil
}

var VOID_TAGS = [16]string{"area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr"}

func voidTag(tag string) bool {
	for _, t := range VOID_TAGS {
		if tag == t {
			return true
		}
	}
	return false
}

type HTMLNode interface {
	toHtml() (string, error)
	String() string
}

func ToHtml(h HTMLNode) (string, error) {
	return h.toHtml()
}

type LeafNode struct {
	tag   string
	value string
	props Props
}

func NewLeafNode(tag, value string, props map[string]string) (LeafNode, error) {
	if value == "" && !voidTag(tag) {
		return LeafNode{}, fmt.Errorf("leaf node needs a value unless its self closing")
	}
	if voidTag(tag) && value != "" {
		return LeafNode{}, fmt.Errorf("self closing tags shouldn't contain content")
	}
	node := LeafNode{
		tag:   tag,
		value: value,
		props: Props{
			props: props,
		},
	}
	return node, nil
}

func (l LeafNode) toHtml() (string, error) {

	value := l.value

	if value == "" && !voidTag(l.tag) {
		return "", fmt.Errorf("leaf node missing value")
	}

	tag := l.tag

	if tag == "" {
		return value, nil
	}

	props, err := l.props.PropsToHTML()
	if err != nil {
		return "", err
	}

	if voidTag(tag) {
		return fmt.Sprintf("<%s%s>%s", tag, props, value), nil
	}

	return fmt.Sprintf("<%s%s>%s</%s>", tag, props, value, tag), nil
}

func (l LeafNode) String() string {
	return fmt.Sprintf("LeafNode(%s, %s, %s)", l.tag, l.value, l.props.props)
}

type ParentNode struct {
	tag      string
	children []HTMLNode
	props    Props
}

func NewParentNode(tag string, children []HTMLNode, props map[string]string) (ParentNode, error) {
	if tag == "" || voidTag(tag) {
		return ParentNode{}, fmt.Errorf("parent must containg a tag that is not self closing")
	}

	if len(children) == 0 {
		return ParentNode{}, fmt.Errorf("parent must contain children")
	}

	node := ParentNode{
		tag:      tag,
		children: children,
		props: Props{
			props: props,
		},
	}
	return node, nil
}

func (p ParentNode) toHtml() (string, error) {

	tag := p.tag
	if tag == "" {
		return "", fmt.Errorf("parent node missing tag")
	}

	if voidTag(tag) {
		return "", fmt.Errorf("parent node can't have a void tag")
	}

	children := p.children
	if len(children) == 0 {
		return "", fmt.Errorf("parent node missing children")
	}

	props, err := p.props.PropsToHTML()
	if err != nil {
		return "", err
	}

	var html strings.Builder
	_, err = html.WriteString(fmt.Sprintf("<%s%s>", tag, props))
	if err != nil {
		return "", err
	}
	for _, c := range children {
		cHtml, err := ToHtml(c)
		if err != nil {
			return "", err
		}

		_, err = html.WriteString(cHtml)
		if err != nil {
			return "", err
		}
	}
	_, err = html.WriteString(fmt.Sprintf("</%s>", tag))
	if err != nil {
		return "", err
	}

	return html.String(), nil
}

func (p ParentNode) String() string {
	return fmt.Sprintf("ParentNode(%s, %s, %s)", p.tag, p.children, p.props.props)
}
