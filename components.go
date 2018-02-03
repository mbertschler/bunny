package main

import "github.com/mbertschler/blocks/html"

var pageBlock = html.Blocks{
	html.Blocks{html.Doctype("html"),
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
			html.Script(html.Src("/static/semantic-ui-css/semantic.js")),
		),
	},
	html.Body(nil,
		html.H1(html.Class("ui center aligned header").Styles("padding:30px"),
			html.Text("Bunny Work Management Tool")),
		html.Div(html.Id("container"),
			displayBlock(getItemData()),
		),
		html.Script(html.Src("/js/app.js")),
	),
}

func editBlock(data itemData) html.Block {
	return html.Div(html.Class("ui text container"),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				html.Button(append(html.Class("ui right floated positive button"),
					html.AttrPair{Key: "onclick", Value: "saveItem('id')"}),
					html.Text("Save"),
				),
				html.Button(append(html.Class("ui right floated button"),
					html.AttrPair{Key: "onclick", Value: "viewItem('id')"}),
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

func displayBlock(data itemData) html.Block {
	var status, statusButton html.Block
	var archiveButton, archiveLabel html.Block
	if data.Closed {
		status = html.I(html.Class("remove circle outline icon red").Styles("display:inline-block"))
		statusButton = html.Button(append(html.Class("ui right floated positive button"),
			html.AttrPair{Key: "onclick", Value: "editItemClosed(false)"}),
			html.Text("Reopen item"),
		)
		if data.Archived {
			archiveButton = html.Button(append(html.Class("ui right floated  button"),
				html.AttrPair{Key: "onclick", Value: "editItemArchived(false)"}),
				html.Text("Unarchive item"),
			)
			archiveLabel = html.Div(html.Class("ui horizontal label").
				Styles("top: -4px; position: relative; margin-right: 8px;"), html.Text("archived"))
		} else {
			archiveButton = html.Button(append(html.Class("ui right floated  button"),
				html.AttrPair{Key: "onclick", Value: "editItemArchived(true)"}),
				html.Text("Archive item"),
			)
		}
	} else {
		status = html.I(html.Class("selected radio icon green").Styles("display:inline-block"))
		statusButton = html.Button(append(html.Class("ui right floated negative button"),
			html.AttrPair{Key: "onclick", Value: "editItemClosed(true)"}),
			html.Text("Close item"),
		)
	}

	return html.Div(html.Class("ui text container"),
		html.Div(html.Class("ui grid"),
			html.Div(html.Class("column"),
				html.Button(append(html.Class("ui right floated button"),
					html.AttrPair{Key: "onclick", Value: "editItem('id')"}),
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
