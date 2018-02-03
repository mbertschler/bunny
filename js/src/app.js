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

function hello(name) {
	callGuiAPI("hello",{
		name: name,
	})
}

function editItem(id) {
	callGuiAPI("editItem", id)
}

function saveItem(id) {
	var data = {}
	$(".itemForm").each(function(i, el){
		console.log(el)
		data[el.name] = el.value
	})
	callGuiAPI("saveItem", data)
}

function viewItem(id) {
	callGuiAPI("viewItem", id)
}

function editItemClosed(closed) {
	callGuiAPI("editItemClosed", closed)
}

function editItemArchived(archived) {
	callGuiAPI("editItemArchived", archived)
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

function handleResponse(resp) {
	for (var i=0; i< resp.Results.length; i++) {
		var r = resp.Results[i]
		if (r.HTML) {
			for (var j=0; j< r.HTML.length; j++) {
				var update = r.HTML[j]
				console.log(update)
				if (update.Operation == 1) {
					$(update.Selector).html(update.Content)
				} else {
					console.warn("update type not implemented :(")
				}
			}
		}
	}
}