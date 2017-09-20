package template

func genHTML_example(applications []map[string]interface{}, xs []map[string]interface{}, student_pinyin string, content string, A string) (_gen string) {
	_gen = ``
	_gen += `<!DOCTYPE html>\n<html lang="en">\n<head>\n    <meta charset="UTF-8">\n    <title>Student: `
	_gen += student_pinyin
	_gen += ` | System of Application Manager</title>\n    <script src="/static/jq.js"></script>\n    <link rel="stylesheet" href="/static/main.css" type="text/css">\n    <script>\n    function UploadFile(app_id) {\n        var form = new FormData($("#upload_"+app_id).get(0));\n        form.append('csrfmiddlewaretoken',$('#csrf').children()[0].value);\n        $("#upload-prompt").show();\n        $.ajax({\n            type:'POST',\n            url:'/upload/'+app_id,\n            data:form,\n            processData:false,\n            contentType:false,\n            success: function () {\n                alert("Data Uploaded");\n                $("#upload-prompt").hide();\n                window.location.reload(false);\n            }\n        })\n    }\n    </script>\n</head>\n<body>\n<div id="big">\n    <h1 class="title-text"> List of Current Application Essays: </h1>\n    <ul>\n        <li>\n            <span class="span1 demo">File Name</span>\n            <span class="span2 demo">Student PDF</span>\n            <span class="span3 demo">Teacher PDF</span>\n            <span class="span4 demo">Teacher checked status</span>\n        </li>\n        `
	for _, app := range applications {
		app__teacher_checked := app[`teacher_checked`].(string)
		app__student_upload := app[`student_upload`].(string)
		app__filename := app[`filename`].(string)
		app__id := app[`id`].(string)
		app__teacher_checked := app[`teacher_checked`].(string)
		app__teacher_upload := app[`teacher_upload`].(string)
		app__teacher_upload := app[`teacher_upload`].(string)
		app__id := app[`id`].(string)
		app__student_upload := app[`student_upload`].(string)
		app__filename := app[`filename`].(string)
		_gen += `\n            <li>\n                <div class="cur_app">\n                    <span class="span1">`
		_gen += app__filename
		_gen += `</span>\n                    <form style="display: none;" id="upload_`
		_gen += app__id
		_gen += `" action="/upload/`
		_gen += app__id
		_gen += `" enctype="multipart/form-data" method="POST">\n                        <input type="file" name="upload" accept="application/pdf" onchange="UploadFile(`
		_gen += app__id
		_gen += `)"/>\n                    </form>\n                    <span class="span2">\n                        `
		if app__student_upload == "" {
			_gen += `\n                            <a href='#' onclick="$('#upload_`
			_gen += app__id
			_gen += `').children()[0].click()">upload</a>\n                        `
		} else {
			_gen += `\n                            <a href='#' onclick="$('#upload_`
			_gen += app__id
			_gen += `').children()[0].click()">update</a>\n                            <a href="/download/student/`
			_gen += app__id
			_gen += `">download</a>\n                        `
		}
		_gen += `\n                    </span>\n                    <span class="span3">\n                        `
		if app__teacher_upload == "" {
			_gen += `\n                            Not Available\n                        `
		} else {
			_gen += `\n                            <a href="/download/teacher/`
			_gen += app__id
			_gen += `">download</a>\n                        `
		}
		_gen += `\n                    </span>\n                    <span class="span4">`
		if app__teacher_checked {
			_gen += ` &Sqrt; `
		} else {
			_gen += ` &times; `
		}
		_gen += `</span>\n                </div>\n            </li>\n        `
	}
	_gen += `\n        <li><div id="add_app"><a href="#" style="font-size: 14px;" onclick="showIframe();">Add New</a></div></li>\n        <li id="upload-prompt" style="display: none;"><div style="color: red;">Do not close the window until the uploaded prompt is shown</div></li>\n    </ul>\n</div>\n\n<div class="overlay"></div>\n    <div id="add-iframe" style="display: none;z-index:999;">\n        <iframe src="/add/application" frameborder="0" width="400px" height="600px;"></iframe>\n    </div>\n    <form id="csrf">`
	_gen += `</form>\n    <script>\n        function showIframe(){\n            $("#add-iframe").show();\n            $(".overlay").show().css('background', 'rgba(0,0,0,0.7)');\n        }\n        window.addEventListener('message', function(e){\n            $("#add-iframe").hide();\n            location.reload();\n        });\n  </script>\n</body>\n</html>\n`
	_gen += `<body>\n    `
	_gen += content
	_gen += `\n    `
	for _, x := range xs {
		x__id := x[`id`].(map[string]interface{})
		x__prim := x[`prim`].(map[string]interface{})
		x__id__first := x__id[`first`].(string)
		x__prim__first := x__prim[`first`].(string)
		x__id__second := x__id[`second`].(string)
		_gen += `\n        `
		_gen += x__id__first
		_gen += `\n        `
		_gen += x__prim__first
		_gen += `\n        `
		_gen += x__id__second
		_gen += `\n    `
	}
	_gen += `\n    `
	if A == "" {
		_gen += ` STH `
	}
	_gen += `\n</body>\n`

	return
}

func Construct_example(m map[string]interface{}) string {
	applications := m[`applications`].([]map[string]interface{})
	student_pinyin := m[`student_pinyin`].(string)
	applications := m[`applications`].([]map[string]interface{})
	A := m[`A`].(string)
	content := m[`content`].(string)
	student_pinyin := m[`student_pinyin`].(string)
	xs := m[`xs`].([]map[string]interface{})

	return genHTML_example(applications, xs, student_pinyin, content, A)
}
