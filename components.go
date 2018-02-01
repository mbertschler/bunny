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
		displayBlock,
		html.H1(html.Class("ui center aligned header").Styles("padding:30px"),
			html.Text("Edit Mode")),
		editBlock,
	),
}

var editBlock = html.Div(html.Class("ui text container"),
	html.Div(html.Class("ui grid"),
		html.Div(html.Class("column"),
			html.Button(html.Class("ui right floated positive button"),
				html.Text("Save"),
			),
			html.Button(html.Class("ui right floated basic negative button"),
				html.Text("Cancel"),
			),
		),
	),
	html.Div(html.Class("ui form"),
		html.Div(html.Class("ui big input fluid").Styles("padding-top:15px"),
			html.Input(append(html.Class("text"),
				html.AttrPair{Key: "placeholder", Value: "Item title"})),
		),
		html.Div(html.Class("ui divider")),
		html.Div(html.Class("field"),
			html.Textarea(append(html.Styles("font:inherit;"),
				html.AttrPair{Key: "placeholder", Value: "Item description"},
				html.AttrPair{Key: "rows", Value: "8"})),
		),
		html.Div(html.Class("ui divider")),
	),
	html.Div(html.Class("ui grid"),
		html.Div(html.Class("column"),
			html.Button(html.Class("ui right floated button"),
				html.Text("Archive item"),
			),
			html.Button(html.Class("ui right floated negative button"),
				html.Text("Close item"),
			),
			html.Button(html.Class("ui right floated positive button"),
				html.Text("Reopen item"),
			),
		),
	),
)

var displayBlock = html.Div(html.Class("ui text container"),
	html.Div(html.Class("ui grid"),
		html.Div(html.Class("column"),
			html.Button(html.Class("ui right floated button"),
				html.Text("Edit"),
			),
		),
	),
	html.H2(nil,
		html.I(html.Class("remove circle outline icon red").Styles("display:inline-block")),
		html.I(html.Class("selected radio icon green").Styles("display:inline-block")),
		html.Text("Item title"),
	),
	html.Div(html.Class("ui divider")),
	html.P(nil, html.Text("Some body text for this item. It should really be worked on!")),
	html.Div(html.Class("ui divider")),
	html.Div(html.Class("ui grid"),
		html.Div(html.Class("column"),
			html.Button(html.Class("ui right floated button"),
				html.Text("Archive item"),
			),
			html.Button(html.Class("ui right floated negative button"),
				html.Text("Close item"),
			),
			html.Button(html.Class("ui right floated positive button"),
				html.Text("Reopen item"),
			),
		),
	),
)