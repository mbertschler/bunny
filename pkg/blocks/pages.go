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

func EditItemPage(data data.Item, new bool) html.Block {
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
