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
	cancelFunc := fmt.Sprintf("itemView(%d)", data.ID)
	if new {
		cancelFunc = "listView()"
	}
	saveFunc := fmt.Sprintf("itemSave(%d, %t)", data.ID, new)
	return html.Div(html.Class("ui text container"),
		gridColumnBlock(
			floatedButton("positive right", saveFunc, "Save"),
			floatedButton("right", cancelFunc, "Cancel"),
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

func ViewItemPage(d data.Item) html.Block {
	var status, statusButton html.Block
	var archiveButton, archiveLabel html.Block
	switch d.State {
	case data.ItemComplete:
		status = completeItemElement
		archiveButton = floatedButton("right",
			fmt.Sprintf("itemState(%d, 'archived')", d.ID), "Archive item")
		statusButton = floatedButton("right yellow",
			fmt.Sprintf("itemState(%d, 'open')", d.ID), "Reopen item")
	case data.ItemArchived:
		status = completeItemElement
		archiveButton = floatedButton("right",
			fmt.Sprintf("itemState(%d, 'complete')", d.ID), "Unarchive item")
		statusButton = floatedButton("right red",
			fmt.Sprintf("itemDelete(%d)", d.ID), "Delete item")
		archiveLabel = html.Div(html.Class("ui horizontal label").
			Styles("top: -4px; position: relative; margin-right: 8px;"), html.Text("archived"))
	case data.ItemOpen:
		status = openItemElement
		statusButton = floatedButton("right positive",
			fmt.Sprintf("itemState(%d, 'complete')", d.ID), "Complete item")
	}

	var laterClass, focusClass, watchClass string
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
		gridColumnBlock(
			floatedIconButton("left", "listView()", "chevron left", "List"),
			floatedButton("right", fmt.Sprintf("itemEdit(%d)", d.ID), "Edit"),
		),
		buttonGroupBlock(
			compactIconButton(laterClass,
				fmt.Sprintf("itemFocus(%d, 'later')", d.ID), "wait", "List"),
			floatedIconButton(focusClass,
				fmt.Sprintf("itemFocus(%d, 'focus')", d.ID), "star", "Focus"),
			floatedIconButton(watchClass,
				fmt.Sprintf("itemFocus(%d, 'watch')", d.ID), "unhide", "Watch"),
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

func ViewAreaPage(d []data.Thing) html.Block {
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
		gridColumnBlock(
			floatedButton("positive right", "itemNew()", "New item"),
			floatedButton("purple right", "listNew()", "New list"),
		),
		html.Div(html.Id("item-list").Class("ui relaxed selection list"),
			list,
		),
		html.Div(html.Id("archive-list").Class("ui relaxed selection list"),
			archived,
		),
	)
}

func ViewListPage(d []data.Item) html.Block {
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
		gridColumnBlock(
			floatedIconButton("left", "areaView()", "chevron left", "Area"),
			floatedButton("positive right", "itemNew()", "New item"),
		),
		html.Div(html.Id("item-list").Class("ui relaxed selection list"),
			list,
		),
		html.Div(html.Id("archive-list").Class("ui relaxed selection list"),
			archived,
		),
	)
}

func ViewFocusPage(focus data.FocusData) html.Block {
	var list html.Blocks
	if len(focus.Focus) > 0 {
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
