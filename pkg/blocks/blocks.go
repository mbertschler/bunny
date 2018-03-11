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

var (
	focusNowIcon   = "yellow star"
	focusLaterIcon = "red wait"
	focusWatchIcon = "blue unhide"
)

func LayoutBlock(content html.Block) html.Block {
	return html.Blocks{
		html.Doctype("html"),
		html.Head(nil,
			html.Meta(html.Charset("utf-8")),
			html.Meta(html.Attr{{Key: "http-equiv", Value: "X-UA-Compatible"}}.Content("IE=edge,chome=1")),
			html.Meta(html.Name("viewport").Content("width=device-width, initial-scale=1.0, maximum-scale=1.0")),
			html.Meta(html.Name("apple-mobile-web-app-capable").Content("yes")),
			html.Title(nil,
				html.Text("Bunny"),
			),
			html.Link(html.Rel("stylesheet").Href("/static/semantic-ui-css/semantic.min.css")),
		),
		html.Body(nil,
			html.H1(html.Class("ui center aligned header").Styles("padding:32px 0 16px"),
				html.Text("Bunny Work Management Tool")),
			html.Div(html.Id("container"),
				content,
			),
			html.Script(html.Src("/static/jquery/dist/jquery.min.js")),
			html.Script(html.Src("/static/semantic-ui-css/semantic.min.js")),
			html.Script(html.Src("/static/sortablejs/Sortable.min.js")),
			html.Script(html.Src("/js/app.js")),
		),
	}
}

func EditItemBlock(data data.Item, new bool) html.Block {
	var cancelFunc string
	if new {
		cancelFunc = "listView()"
	} else {
		cancelFunc = fmt.Sprintf("itemView(%d)", data.ID)
	}
	return html.Div(html.Class("ui text container"),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				html.Button(append(html.Class("ui right floated positive button"),
					html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemSave(%d, %t)", data.ID, new)}),
					html.Text("Save"),
				),
				html.Button(append(html.Class("ui right floated button"),
					html.AttrPair{Key: "onclick", Value: cancelFunc}),
					html.Text("Cancel"),
				),
			),
		),
		html.Div(html.Class("ui form"),
			html.Div(html.Class("ui big input fluid").Styles("padding-top:15px"),
				html.Input(append(html.Class("itemForm").Name("Title").Type("text").Value(data.Title),
					html.AttrPair{Key: "placeholder", Value: "Item title"})),
			),
			html.Div(html.Class("ui divider")),
			html.Div(html.Class("field"),
				html.Textarea(append(html.Class("itemForm").Name("Body").Styles("font:inherit;"),
					html.AttrPair{Key: "placeholder", Value: "Item description"},
					html.AttrPair{Key: "rows", Value: "8"}),
					html.Text(data.Body),
				),
			),
		),
	)
}

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

func ViewItemBlock(d data.Item) html.Block {
	var status, statusButton html.Block
	var archiveButton, archiveLabel html.Block

	switch d.State {
	case data.ItemComplete:
		status = html.I(html.Class("checkmark icon green").Styles("display:inline-block"))
		archiveButton = html.Button(append(html.Class("ui right floated button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemState(%d, 'archived')", d.ID)}),
			html.Text("Archive item"),
		)
		statusButton = html.Button(append(html.Class("ui right floated yellow button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemState(%d, 'open')", d.ID)}),
			html.Text("Reopen item"),
		)
	case data.ItemArchived:
		status = html.I(html.Class("checkmark icon green").Styles("display:inline-block"))
		archiveButton = html.Button(append(html.Class("ui right floated button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemState(%d, 'complete')", d.ID)}),
			html.Text("Unarchive item"),
		)
		archiveLabel = html.Div(html.Class("ui horizontal label").
			Styles("top: -4px; position: relative; margin-right: 8px;"), html.Text("archived"))
		statusButton = html.Button(append(html.Class("ui right floated red button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemDelete(%d)", d.ID)}),
			html.Text("Delete item"),
		)
	case data.ItemOpen:
		status = html.I(html.Class("radio icon grey").Styles("display:inline-block"))
		statusButton = html.Button(append(html.Class("ui right floated positive button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemState(%d, 'complete')", d.ID)}),
			html.Text("Complete item"),
		)
	}

	var laterClass, focusClass, watchClass string
	var focusIcon = "star"
	switch d.Focus {
	case data.FocusLater:
		laterClass = " red"
	case data.FocusNow:
		focusClass = " yellow"
	case data.FocusWatch:
		watchClass = " blue"
	case data.FocusNone:
		laterClass = " red"
		focusClass = " yellow"
		watchClass = " blue"
	}

	return html.Div(html.Class("ui text container"),
		menuBlock(),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column").Styles("text-align:center"),
				html.Button(append(html.Class("ui left floated button"),
					html.AttrPair{Key: "onclick", Value: "listView()"}),
					html.I(html.Class("chevron left icon")),
					html.Text("List"),
				),
				html.Button(append(html.Class("ui right floated button"),
					html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemEdit(%d)", d.ID)}),
					html.Text("Edit"),
				),
			),
		),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column").Styles("text-align:center"),
				html.Div(html.Class("ui buttons"),
					html.Button(append(html.Class("ui compact button"+laterClass),
						html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemFocus(%d, 'later')", d.ID)}),
						html.I(html.Class("wait icon")),
						html.Text("Later"),
					),
					html.Button(append(html.Class("ui compact button"+focusClass),
						html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemFocus(%d, 'focus')", d.ID)}),
						html.I(html.Class(focusIcon+" icon")),
						html.Text("Focus"),
					),
					html.Button(append(html.Class("ui compact button"+watchClass),
						html.AttrPair{Key: "onclick", Value: fmt.Sprintf("itemFocus(%d, 'watch')", d.ID)}),
						html.I(html.Class("unhide icon")),
						html.Text("Watch"),
					),
				),
			),
		),
		html.H2(nil,
			status,
			archiveLabel,
			html.Text(d.Title),
		),
		html.Div(html.Class("ui divider")),
		html.P(nil, html.Text(d.Body)),
		html.Div(html.Class("ui divider")),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				archiveButton,
				statusButton,
			),
		),
	)
}

func ViewThingsBlock(d []data.Thing) html.Block {
	var list, archived html.Blocks
	for _, t := range d {
		block := listItemBlock(t)
		if t.Archived() {
			if len(archived) == 0 {
				archived.Add(html.H4(html.Styles("padding-left:48px"),
					html.Text("Archived"),
				))
			}
			archived.Add(block)
		} else {
			list.Add(block)
		}
	}
	return html.Div(html.Class("ui text container"),
		menuBlock(),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				html.Button(append(html.Class("ui right floated positive button"),
					html.AttrPair{Key: "onclick", Value: "itemNew()"}),
					html.Text("New item"),
				),
				html.Button(append(html.Class("ui right floated purple button"),
					html.AttrPair{Key: "onclick", Value: "listNew()"}),
					html.Text("New list"),
				),
			),
		),
		html.Div(html.Id("item-list").Class("ui relaxed selection list"),
			list,
		),
		html.Div(html.Id("archive-list").Class("ui relaxed selection list"),
			archived,
		),
	)
}

func ViewListBlock(d []data.Item) html.Block {
	var list, archived html.Blocks
	for _, t := range d {
		block := listItemBlock(t)
		if t.Archived() {
			if len(archived) == 0 {
				archived.Add(html.H4(html.Styles("padding-left:48px"),
					html.Text("Archived"),
				))
			}
			archived.Add(block)
		} else {
			list.Add(block)
		}
	}
	return html.Div(html.Class("ui text container"),
		menuBlock(),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				html.Button(append(html.Class("ui left floated button"),
					html.AttrPair{Key: "onclick", Value: "areaView()"}),
					html.I(html.Class("chevron left icon")),
					html.Text("Area"),
				),
				html.Button(append(html.Class("ui right floated positive button"),
					html.AttrPair{Key: "onclick", Value: "itemNew()"}),
					html.Text("New item"),
				),
			),
		),
		html.Div(html.Id("item-list").Class("ui relaxed selection list"),
			list,
		),
		html.Div(html.Id("archive-list").Class("ui relaxed selection list"),
			archived,
		),
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

func ViewFocusBlock(focus data.FocusData) html.Block {
	var list html.Blocks
	if len(focus.Later) > 0 {
		list.Add(html.H4(html.Styles("padding-left:10px; margin: 32px 0 0;"),
			html.I(html.Class("large middle aligned icon "+focusNowIcon).Styles("padding-right:12px")),
			html.Text("Focus"),
		))
	}
	for _, item := range focus.Focus {
		item.Focus = data.FocusNone
		list.Add(listItemBlock(item))
	}
	if len(focus.Later) > 0 {
		list.Add(html.H4(html.Styles("padding-left:10px; margin: 32px 0 0;"),
			html.I(html.Class("large middle aligned icon "+focusLaterIcon).Styles("padding-right:12px")),
			html.Text("Later"),
		))
	}
	for _, item := range focus.Later {
		item.Focus = data.FocusNone
		list.Add(listItemBlock(item))
	}
	if len(focus.Watch) > 0 {
		list.Add(html.H4(html.Styles("padding-left:10px; margin: 32px 0 0;"),
			html.I(html.Class("large middle aligned icon "+focusWatchIcon).Styles("padding-right:12px")),
			html.Text("Watched"),
		))
	}
	for _, item := range focus.Watch {
		item.Focus = data.FocusNone
		list.Add(listItemBlock(item))
	}
	return html.Div(html.Class("ui text container"),
		menuBlock(),
		html.Div(html.Id("focus-list").Class("ui relaxed selection list"),
			list,
		),
	)
}
