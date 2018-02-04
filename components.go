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
			html.Script(html.Src("/static/jquery/dist/jquery.min.js")),
			html.Script(html.Src("/static/semantic-ui-css/semantic.min.js")),
		),
		html.Body(nil,
			html.H1(html.Class("ui center aligned header").Styles("padding:30px"),
				html.Text("Bunny Work Management Tool")),
			html.Div(html.Id("container"),
				content,
			),
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
	if data.Closed {
		status = html.I(html.Class("remove circle outline icon red").Styles("display:inline-block"))
		statusButton = html.Button(append(html.Class("ui right floated positive button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemClosed(%d, false)", data.ID)}),
			html.Text("Reopen item"),
		)
		if data.Archived {
			archiveButton = html.Button(append(html.Class("ui right floated  button"),
				html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemArchived(%d, false)", data.ID)}),
				html.Text("Unarchive item"),
			)
			archiveLabel = html.Div(html.Class("ui horizontal label").
				Styles("top: -4px; position: relative; margin-right: 8px;"), html.Text("archived"))
		} else {
			archiveButton = html.Button(append(html.Class("ui right floated  button"),
				html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemArchived(%d, true)", data.ID)}),
				html.Text("Archive item"),
			)
		}
	} else {
		status = html.I(html.Class("selected radio icon green").Styles("display:inline-block"))
		statusButton = html.Button(append(html.Class("ui right floated negative button"),
			html.AttrPair{Key: "onclick", Value: fmt.Sprintf("editItemClosed(%d, true)", data.ID)}),
			html.Text("Close item"),
		)
	}

	return html.Div(html.Class("ui text container"),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
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
	var list html.Blocks
	for _, item := range data {
		var iconClass string
		if item.Closed {
			iconClass = "remove circle outline red"
		} else {
			iconClass = "selected radio green"
		}
		block := html.A(html.Class("item").Href(fmt.Sprint("/item/", item.ID)),
			html.I(html.Class("large middle aligned icon "+iconClass)),
			html.Div(html.Class("middle aligned content").Styles("color:rgba(0,0,0,0.87)"),
				html.Text(item.Title),
			),
		)
		list.Add(block)
	}
	return html.Div(html.Class("ui text container"),
		html.Div(html.Class("ui relaxed selection list"),
			list,
		),
	)
}
