console.log("hello from app.js")
hello("app.js")

function hello(name) {
	callGuiAPI("hello",{
		name: name,
	})
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