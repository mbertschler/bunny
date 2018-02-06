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

package main

import (
	"fmt"

	"github.com/mbertschler/blocks/html"
)

func pageBlock(content html.Block) html.Block {
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
			html.H1(html.Class("ui center aligned header").Styles("padding:30px"),
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

func editItemBlock(data itemData) html.Block {
	return html.Div(html.Class("ui text container"),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				html.Button(append(html.Class("ui right floated positive button"),
					html.AttrPair{Key: "onclick", Value: fmt.Sprintf("saveItem(%d)", data.ID)}),
					html.Text("Save"),
				),
				html.Button(append(html.Class("ui right floated button"),
					html.AttrPair{Key: "onclick", Value: fmt.Sprintf("viewItem(%d)", data.ID)}),
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

func displayItemBlock(data itemData) html.Block {
	var status, statusButton html.Block
	var archiveButton, archiveLabel html.Block
	if data.Complete {
		status = html.I(html.Class("checkmark icon green").Styles("display:inline-block"))
		if data.Archived {
			archiveButton = html.Button(append(html.Class("ui right floated  button"),
				html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemArchived(%d, false)", data.ID)}),
				html.Text("Unarchive item"),
			)
			archiveLabel = html.Div(html.Class("ui horizontal label").
				Styles("top: -4px; position: relative; margin-right: 8px;"), html.Text("archived"))
		} else {
			archiveButton = html.Button(append(html.Class("ui right floated button"),
				html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemArchived(%d, true)", data.ID)}),
				html.Text("Archive item"),
			)
			statusButton = html.Button(append(html.Class("ui right floated red button"),
				html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemComplete(%d, false)", data.ID)}),
				html.Text("Reopen item"),
			)
		}
	} else {
		status = html.I(html.Class("radio icon grey").Styles("display:inline-block"))
		statusButton = html.Button(append(html.Class("ui right floated positive button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemComplete(%d, true)", data.ID)}),
			html.Text("Complete item"),
		)
	}

	var laterClass, focusClass, watchClass string
	if data.Later || data.Focus || data.Watch {
		laterClass, focusClass, watchClass = "", "", ""
		if data.Later {
			laterClass = " red"
		}
		if data.Focus {
			focusClass = " yellow"
		}
		if data.Watch {
			watchClass = " blue"
		}
	} else {
		laterClass = " red"
		focusClass = " yellow"
		watchClass = " blue"
	}

	return html.Div(html.Class("ui text container"),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column").Styles("text-align:center"),
				html.Button(append(html.Class("ui left floated button"),
					html.AttrPair{Key: "onclick", Value: "viewList()"}),
					html.I(html.Class("chevron left icon")),
					html.Text("List"),
				),
				html.Button(append(html.Class("ui right floated button"),
					html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItem(%d)", data.ID)}),
					html.Text("Edit"),
				),
			),
		),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column").Styles("text-align:center"),
				html.Div(html.Class("ui buttons"),
					html.Button(append(html.Class("ui compact button"+laterClass),
						html.AttrPair{Key: "onclick", Value: fmt.Sprintf("focusItem(%d, 'later')", data.ID)}),
						html.I(html.Class("wait icon")),
						html.Text("Later"),
					),
					html.Button(append(html.Class("ui compact button"+focusClass),
						html.AttrPair{Key: "onclick", Value: fmt.Sprintf("focusItem(%d, 'focus')", data.ID)}),
						html.I(html.Class("star icon")),
						html.Text("Focus"),
					),
					html.Button(append(html.Class("ui compact button"+watchClass),
						html.AttrPair{Key: "onclick", Value: fmt.Sprintf("focusItem(%d, 'watch')", data.ID)}),
						html.I(html.Class("unhide icon")),
						html.Text("Watch"),
					),
				),
			),
		),
		html.H2(nil,
			status,
			archiveLabel,
			html.Text(data.Title),
		),
		html.Div(html.Class("ui divider")),
		html.P(nil, html.Text(data.Body)),
		html.Div(html.Class("ui divider")),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				archiveButton,
				statusButton,
			),
		),
	)
}

func displayListBlock(data []itemData) html.Block {
	var list, archived html.Blocks
	for _, item := range data {
		var iconClass string
		if item.Complete {
			iconClass = "checkmark green"
		} else {
			iconClass = "radio grey"
		}
		var focusIcon html.Block
		switch {
		case item.Later:
			focusIcon = html.I(html.Class("large middle aligned icon red wait").Styles("padding-left:10px"))
		case item.Focus:
			focusIcon = html.I(html.Class("large middle aligned icon yellow star").Styles("padding-left:10px"))
		case item.Watch:
			focusIcon = html.I(html.Class("large middle aligned icon blue unhide").Styles("padding-left:10px"))
		}
		block := html.Div(append(html.Class("item"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("viewItem(%d)", item.ID)}),
			html.I(html.Class("large middle aligned icon "+iconClass)),
			focusIcon,
			html.Div(html.Class("middle aligned content").Styles("color:rgba(0,0,0,0.87)"),
				html.Text(item.Title),
			),
		)
		if item.Archived {
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
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				html.Button(append(html.Class("ui right floated positive button"),
					html.AttrPair{Key: "onclick", Value: "newItem()"}),
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
