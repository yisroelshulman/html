package html

import (
	"testing"
)

func TestPropsToHTML(t *testing.T) {
	tests := []struct {
		name      string
		props     Props
		wantProps string
		wantErr   bool
	}{
		{
			name:      "no props",
			props:     Props{},
			wantProps: "",
			wantErr:   false,
		},
		{
			name: "single prop",
			props: Props{
				props: map[string]string{"href": "https://google.com"},
			},
			wantProps: " href=\"https://google.com\"",
			wantErr:   false,
		},
		{
			name: "multiple props",
			props: Props{
				props: map[string]string{"class": "first", "href": "https://google.com"},
			},
			wantProps: " class=\"first\" href=\"https://google.com\"",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProps, gotErr := tt.props.PropsToHTML()
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("PropsToHTML(): error =  %v, wantErr %v", gotErr, tt.wantErr)
			}
			if gotProps != tt.wantProps {
				t.Errorf("PropsToHTML(): Props = %v, want %v", gotProps, tt.wantProps)
			}
		})
	}
}

func TestNewLeafNode(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		value   string
		wantErr bool
	}{
		{
			name:    "no tags",
			tag:     "",
			value:   "no tags",
			wantErr: false,
		},
		{
			name:    "no value",
			tag:     "p",
			value:   "",
			wantErr: true,
		},
		{
			name:    "tag and value",
			tag:     "p",
			value:   "tag and value",
			wantErr: false,
		},
		{
			name:    "self closing tag no value",
			tag:     "img",
			value:   "",
			wantErr: false,
		},
		{
			name:    "self closing tag with value",
			tag:     "img",
			value:   "image",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := NewLeafNode(tt.tag, tt.value, nil)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("NewLeafNode(): error = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}

}

func TestLeafNode(t *testing.T) {
	no_tag, _ := NewLeafNode("", "i am a leaf node", nil)
	no_value, _ := NewLeafNode("p", "", nil)
	no_props, _ := NewLeafNode("p", "i am a leaf node", nil)
	one_prop, _ := NewLeafNode("p", "i am a leaf node", map[string]string{"id": "leaf"})
	multi_props, _ := NewLeafNode("p", "i am a leaf node", map[string]string{"class": "node", "id": "leaf"})
	self_closing_tag, _ := NewLeafNode("img", "", map[string]string{"alt": "image", "src": "https://google.com"})
	tests := []struct {
		name     string
		node     LeafNode
		wantHtml string
		wantStr  string
		wantErr  bool
	}{
		{
			name:     "unset leaf node",
			node:     LeafNode{},
			wantHtml: "",
			wantStr:  "LeafNode(, , map[])",
			wantErr:  true,
		},
		{
			name:     "leaf with no tag",
			node:     no_tag,
			wantHtml: "i am a leaf node",
			wantStr:  "LeafNode(, i am a leaf node, map[])",
			wantErr:  false,
		},
		{
			name:     "leaf with no value",
			node:     no_value,
			wantHtml: "",
			wantStr:  "LeafNode(, , map[])",
			wantErr:  true,
		},
		{
			name:     "leaf with no props",
			node:     no_props,
			wantHtml: "<p>i am a leaf node</p>",
			wantStr:  "LeafNode(p, i am a leaf node, map[])",
			wantErr:  false,
		},
		{
			name:     "leaf with one prop",
			node:     one_prop,
			wantHtml: "<p id=\"leaf\">i am a leaf node</p>",
			wantStr:  "LeafNode(p, i am a leaf node, map[id:leaf])",
			wantErr:  false,
		},
		{
			name:     "leaf with multiple props",
			node:     multi_props,
			wantHtml: "<p class=\"node\" id=\"leaf\">i am a leaf node</p>",
			wantStr:  "LeafNode(p, i am a leaf node, map[class:node id:leaf])",
			wantErr:  false,
		},
		{
			name:     "leaf with self closing tag",
			node:     self_closing_tag,
			wantHtml: "<img alt=\"image\" src=\"https://google.com\">",
			wantStr:  "LeafNode(img, , map[alt:image src:https://google.com])",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHtml, gotErr := ToHtml(tt.node)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("LeafNode html: got error = %v, wantErr %v", gotErr, tt.wantErr)
			}
			if gotHtml != tt.wantHtml {
				t.Errorf("LeafNode html: got html = %v, want %v", gotHtml, tt.wantHtml)
			}

			gotStr := tt.node.String()
			if gotStr != tt.wantStr {
				t.Errorf("LeafNode str: got str = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}

func TestNewParentNode(t *testing.T) {
	child, _ := NewLeafNode("p", "i am leaf", nil)
	tests := []struct {
		name     string
		tag      string
		children []HTMLNode
		wantErr  bool
	}{
		{
			name:     "no tag",
			tag:      "",
			children: []HTMLNode{child},
			wantErr:  true,
		},
		{
			name:     "self closing tag",
			tag:      "img",
			children: []HTMLNode{child},
			wantErr:  true,
		},
		{
			name:     "no children",
			tag:      "div",
			children: nil,
			wantErr:  true,
		},
		{
			name:     "valid tag with childten",
			tag:      "div",
			children: []HTMLNode{child},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := NewParentNode(tt.tag, tt.children, nil)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("NewParentNode(): error = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestParentLeaf(t *testing.T) {
	no_children, _ := NewParentNode("div", nil, nil)
	grandchild, _ := NewLeafNode("p", "i am a leaf node", nil)
	child, _ := NewParentNode("div", []HTMLNode{grandchild}, nil)
	parent, _ := NewParentNode("div", []HTMLNode{child}, nil)
	child2, _ := NewLeafNode("a", "i am a link", map[string]string{"href": "https://google.com"})
	child3, _ := NewLeafNode("", "no tags", nil)
	multi_children, _ := NewParentNode("div", []HTMLNode{grandchild, child2, child3}, nil)
	no_tag, _ := NewParentNode("", []HTMLNode{child2}, nil)
	void_tag, _ := NewParentNode("img", []HTMLNode{child2}, nil)
	no_props, _ := NewParentNode("div", []HTMLNode{child3}, nil)
	one_prop, _ := NewParentNode("div", []HTMLNode{child3}, map[string]string{"id": "prop"})
	multi_props, _ := NewParentNode("div", []HTMLNode{child3}, map[string]string{"id": "prop", "class": "parent"})

	tests := []struct {
		name     string
		node     ParentNode
		wantHtml string
		wantStr  string
		wantErr  bool
	}{
		{
			name:     "unset parent",
			node:     ParentNode{},
			wantHtml: "",
			wantStr:  "ParentNode(, [], map[])",
			wantErr:  true,
		},
		{
			name:     "parent with no children",
			node:     no_children,
			wantHtml: "",
			wantStr:  "ParentNode(, [], map[])",
			wantErr:  true,
		},
		{
			name:     "parent with child",
			node:     child,
			wantHtml: "<div><p>i am a leaf node</p></div>",
			wantStr:  "ParentNode(div, [LeafNode(p, i am a leaf node, map[])], map[])",
			wantErr:  false,
		},
		{
			name:     "parent with grandchild",
			node:     parent,
			wantHtml: "<div><div><p>i am a leaf node</p></div></div>",
			wantStr:  "ParentNode(div, [ParentNode(div, [LeafNode(p, i am a leaf node, map[])], map[])], map[])",
			wantErr:  false,
		},
		{
			name:     "parent with children",
			node:     multi_children,
			wantHtml: "<div><p>i am a leaf node</p><a href=\"https://google.com\">i am a link</a>no tags</div>",
			wantStr:  "ParentNode(div, [LeafNode(p, i am a leaf node, map[]) LeafNode(a, i am a link, map[href:https://google.com]) LeafNode(, no tags, map[])], map[])",
			wantErr:  false,
		},
		{
			name:     "parent with no tag",
			node:     no_tag,
			wantHtml: "",
			wantStr:  "ParentNode(, [], map[])",
			wantErr:  true,
		},
		{
			name:     "parent with void tag",
			node:     void_tag,
			wantHtml: "",
			wantStr:  "ParentNode(, [], map[])",
			wantErr:  true,
		},
		{
			name:     "parent with no props",
			node:     no_props,
			wantHtml: "<div>no tags</div>",
			wantStr:  "ParentNode(div, [LeafNode(, no tags, map[])], map[])",
			wantErr:  false,
		},
		{
			name:     "parent with one prop",
			node:     one_prop,
			wantHtml: "<div id=\"prop\">no tags</div>",
			wantStr:  "ParentNode(div, [LeafNode(, no tags, map[])], map[id:prop])",
			wantErr:  false,
		},
		{
			name:     "parent with multiple props",
			node:     multi_props,
			wantHtml: "<div class=\"parent\" id=\"prop\">no tags</div>",
			wantStr:  "ParentNode(div, [LeafNode(, no tags, map[])], map[class:parent id:prop])",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHtml, gotErr := ToHtml(tt.node)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("Parentnode html: error = %v, want %v", gotErr, tt.wantErr)
			}
			if gotHtml != tt.wantHtml {
				t.Errorf("ParentNode html: html = %v, want %v", gotHtml, tt.wantHtml)
			}

			gotStr := tt.node.String()
			if gotStr != tt.wantStr {
				t.Errorf("ParenNode str: str = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}
