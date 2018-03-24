// Copyright 2018 Martin Bertschler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blocks

import (
	"fmt"

	"github.com/mbertschler/blocks/html"

	"github.com/mbertschler/bunny/pkg/data"
)

func menuBlock() html.Block {
	return html.Div(html.Class("ui two item menu"),
		// html.A(append(html.Class("item"),
		// 	html.AttrPair{Key: "onclick", Value: "listView()"}),
		// 	html.I(html.Class("comments purple icon")),
		// 	html.Text("Updates")),
		html.A(append(html.Class("item"),
			html.AttrPair{Key: "onclick", Value: "focusView()"}),
			html.I(html.Class(focusNowIcon+" icon")),
			html.Text("Focus")),
		html.A(append(html.Class("item"),
			html.AttrPair{Key: "onclick", Value: "listView()"}),
			html.I(html.Class("clone violet icon")),
			html.Text("Workspace")),
	)
}

func listItemBlock(item data.Thing) html.Block {
	switch it := item.(type) {
	case data.Item:
		return itemBlock(it)
	case data.List:
		return listBlock(it)
	}
	return nil
}

func itemBlock(item data.Item) html.Block {
	var iconClass string
	switch item.State {
	case data.ItemOpen:
		iconClass = "radio grey"
	default:
		iconClass = "checkmark green"
	}

	var focusIcon html.Block
	switch item.Focus {
	case data.FocusLater:
		focusIcon = html.I(html.Class("large middle aligned icon " + focusLaterIcon).Styles("padding-left:10px"))
	case data.FocusNow:
		focusIcon = html.I(html.Class("large middle aligned icon " + focusNowIcon).Styles("padding-left:10px"))
	case data.FocusWatch:
		focusIcon = html.I(html.Class("large middle aligned icon " + focusWatchIcon).Styles("padding-left:10px"))
	}

	return html.Div(append(html.Class("item").Data("item-id", item.ID),
		html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemView(%d)", item.ID)}),
		html.I(html.Class("large middle aligned icon "+iconClass)),
		focusIcon,
		html.Div(html.Class("middle aligned content").Styles("color:rgba(0,0,0,0.87)"),
			html.Text(item.Title),
		),
	)
}

func listBlock(list data.List) html.Block {
	var iconClass string
	switch list.State {
	case data.ItemOpen:
		iconClass = "violet square"
	default:
		iconClass = "purple square check"
	}

	return html.Div(append(html.Class("item").Data("list-id", list.ID),
		html.AttrPair{Key: "onclick", Value: fmt.Sprintf("listView(%d)", list.ID)}),
		html.I(html.Class("large middle aligned icon "+iconClass)),
		html.Div(html.Class("middle aligned content").Styles("color:rgba(0,0,0,0.87)"),
			html.Text(list.Title),
		),
	)
}

func gridColumnBlock(children ...html.Block) html.Block {
	return html.Div(html.Class("ui grid"),
		html.Div(html.Class("column"),
			children...,
		),
	)
}

func floatedButton(class, action, text string) html.Block {
	return html.Button(append(
		html.Class("ui floated button "+class),
		html.AttrPair{Key: "onclick", Value: action}),
		html.Text(text),
	)
}
