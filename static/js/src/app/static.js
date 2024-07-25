var dirlist = {}
var base_path = '/'

function change_dir(dir_name) {
    if (dir_name == '..') {
	path_array = base_path.split('/')
	path_array.pop()
	base_path = '/' + path_array.join('/')
    } else {
	base_path += dir_name + '/'
    }
    load()
}

function load() {
    $("#dirTable").hide()
    $("#emptyMessage").hide()
    $("#loading").show()
    api.static_content.get(base_path)
	.success(function (response) {
            $("#loading").hide()
	    if ((Array.isArray(response.files) || Array.isArray(response.dirs)) && response.files.length + response.dirs.length > 0) {
                files = response.files
		dirs = response.dirs
		$("#emptyMessage").hide()
		$("#dirTable").show()
		var dirTable = $('#dirTable').DataTable({
		    columns: [
			{data: 'type'},
			{data: 'name'}
		    ],
		    destroy: true, 
	            columnDefs: [{
			orderable: false,
			targets: "no-sort"
                    }]
		});
		dirTable.clear();
		dirRows = []
		if (base_path != '/') {
		    dirrows.push({
			    type: escapeHtml("Dir"),
			    name: escapeHtml("..")
		    })
		}
		$.each(dirs, function (i, dir) {
		    dirRows.push({
			    type: escapeHtml("Dir"),
			    name: escapeHtml(dir)
		    })
		})
		$.each(files, function (i, file) {
	            dirRows.push({
			    type: escapeHtml("File"), 
			    name: escapeHtml(file)
		    })
		})
		dirTable.rows.add(dirRows).draw()
	    } else {
		$("#emptyMessage").show()
	    }
	})
	.error(function () {
	    errorFlash("Error listing directory")
	})
}

$(document).ready(function () {
    load()
});
