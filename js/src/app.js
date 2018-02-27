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

enableSorting()

function enableSorting() {
	activateList("item-list", sortUpdate)
	activateList("focus-pause-list", sortFocusUpdate)
}

function activateList(id, cb) {
	var el = document.getElementById(id);
	if (el) {
		var options = {
			animation: 150,
			onUpdate: cb,
		}
		Sortable.create(el, options);
	}
}

function sortUpdate(event) {
	callGuiAPI("listSort",{
		Item: parseInt(event.item.dataset.itemId, 10),
		Pos: event.newIndex+1,
	})
}

function sortFocusUpdate(event) {
	callGuiAPI("focusSort",{
		Old: event.oldIndex,
		New: event.newIndex,
	})
}

function itemFocus(id, status) {
	callGuiAPI("itemFocus",{
		ID: id,
		Focus: status,
	})
}

function listView() {
	callGuiAPI("listView", null)
}

function itemEdit(id) {
	callGuiAPI("itemEdit", id)
}

function itemNew() {
	callGuiAPI("itemNew", null)
}

function itemSave(id, isNew) {
	var data = {
		ID: id,
		New: isNew,
	}
	$(".itemForm").each(function(i, el){
		data[el.name] = el.value
	})
	callGuiAPI("itemSave", data)
}

function itemDelete(id) {
	callGuiAPI("itemDelete", id)
}

function itemView(id) {
	callGuiAPI("itemView", id)
}

function itemState(id, state) {
	callGuiAPI("itemState", {
		ID: id,
		State: state,
	})
}

function focusView(id) {
	callGuiAPI("focusView", id)
}

function callGuiAPI(name, args) {
	var req = {
		Actions: [{
			Name: name,
			Args: args,
		}]
	}
	$.ajax({
		method: "POST",
		url: "/gui/",
		data: JSON.stringify(req),
		success: function (data) {
			var ret = JSON.parse(data)
			handleResponse(ret)
		},
		error: function (error) {
			console.error("error:", error)
		},
	})
}

var callableFunctions = {
	"setURL": setURL,
	"enableSorting": enableSorting,
}

function setURL(args) {
	history.pushState(args[0], args[1], args[2])
}

function handleResponse(resp) {
	for (var i=0; i< resp.Results.length; i++) {
		var r = resp.Results[i]
		if (r.HTML) {
			for (var j=0; j< r.HTML.length; j++) {
				var update = r.HTML[j]
				if (update.Operation == 1) {
					$(update.Selector).html(update.Content)
				} else {
					console.warn("update type not implemented :(", update)
				}
			}
		}
		if (r.JS) {
			for (var j=0; j< r.JS.length; j++) {
				var call = r.JS[j]
				var func = callableFunctions[call.Name]
				if (func) {
					func(call.Arguments)
				} else {
					console.warn("function call not implemented :(", call)
				}
			}
		}
	}
}