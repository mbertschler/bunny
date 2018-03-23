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
	"github.com/mbertschler/blocks/html"
)

func LayoutBlock(content html.Block) html.Block {
	return html.Blocks{
		html.Doctype("html"),
		layoutHead(),
		layoutBody(content),
	}
}

func layoutHead() html.Block {
	return html.Head(nil,
		html.Meta(html.Charset("utf-8")),
		html.Meta(html.Attr{{Key: "http-equiv", Value: "X-UA-Compatible"}}.Content("IE=edge,chome=1")),
		html.Meta(html.Name("viewport").Content("width=device-width, initial-scale=1.0, maximum-scale=1.0")),
		html.Meta(html.Name("apple-mobile-web-app-capable").Content("yes")),
		html.Title(nil,
			html.Text("Bunny"),
		),
		html.Link(html.Rel("stylesheet").Href("/static/semantic-ui-css/semantic.min.css")),
	)
}

func layoutBody(content html.Block) html.Block {
	return html.Body(nil,
		html.H1(html.Class("ui center aligned header").Styles("padding:32px 0 16px"),
			html.Text("Bunny Work Management Tool")),
		html.Div(html.Id("container"),
			content,
		),
		html.Script(html.Src("/static/jquery/dist/jquery.min.js")),
		html.Script(html.Src("/static/semantic-ui-css/semantic.min.js")),
		html.Script(html.Src("/static/sortablejs/Sortable.min.js")),
		html.Script(html.Src("/js/app.js")),
	)
}
