$(window).ready(function () {
    InitBot();

    $.getJSON(location.href, 'cols=', function (res) {
        var $water = $('#water');
        var $pump = $("#pump");
        var $remote = $("#remote");
        for (var k in res.water) {
            var temp = k.split(',');
            var name = temp[0];
            var id = temp[1];
            var count = temp[2];
            var $div = createVector(name, true);
            $water.append($div);
            $div.append(createData(id, res.water[k], count, 8, $div));
        }

        for (var k in res.site) {
            var temp = k.split(',');
            var name = temp[0];
            var id = temp[1];
            var count = temp[2];
            var $div = createVector(name, true);
            $pump.append($div);
            $div.append(createData(id, res.site[k], count, 8, $div));
        }

        var $div = createVector("测压点");
        $div.append(createRemote(res.remote));
        $div.css('width', 525);

        $pump.append($div);
        $("#remote").parent().css('width', '605px')
        Update();

        setInterval(Update, 30000);
    });
});

function createVector(t, b) {
    if (b == null) {
        b = false;
    }
    var $div = $(document.createElement("div"));
    $div.addClass("table-vector");

    var $title = $(document.createElement("div"));
    $title.addClass("table-title");
    $title.append("&nbsp;&nbsp;" + t);

    if (b) {
        $div.append("<div class='table-vector-time'>数据时间:--</div>");
    }

    $div.append($title);
    return $div;
}

function createData(id, list, c, s, vec) {
    var $table = $(document.createElement('table'));
    $table.attr('name', id);
    $table.attr('id', id);
    var col = c % s == 0 ? c / s : Math.floor(c / s) + 1;
    var $tr = null;
    if (id == "0118" || id == "0112") {
        col = 2;
    }
    if (col == 1) {
        vec.css("width", "295px");
    }

    var index = 0;
    for (var k in list) {
        if (k != 100) {
            if (index == 0) {
                $tr = $(document.createElement("tr"));
                $table.append($tr);
            }
            var temp = list[k].split(',');
            var id2 = temp[0];
            var name = temp[1];
            if (col == 1) {
                $tr.append("<td style='width:70%'>" + name + ":</td><td style='width:70%; text-align:center;' id='" + id + "_" + id2 + "'>--</td>");
            } else {
                $tr.append("<td >" + name + ":</td><td id='" + id + "_" + id2 + "'>--</td>");
            }
            //$tr.append("<td >" + list[k] + ":</td><td id='" + id + "_" + k + "'>--</td>");

            index++;
            if (index == col) {
                index = 0;
            }
        }
    }

    return $table;
}

function createRemote(res) {
    var $table = $(document.createElement('table'));
    $table.attr('name', 'remote');
    $table.attr('id', 'remote');
    $table.css('width', '75%');
    var index = 0;
    var $tr = null;
    for (var k in res) {
        if (index == 0) {
            $tr = $(document.createElement('tr'));
            $table.append($tr);
        }
        $tr.append("<td id='" + res[k].ID + "_t'>" + res[k].InstallationPositionDesc + "</td><td id='" + res[k].ID + "'>--</td>");
        index++;

        if (index == 2) {
            index = 0;
        }
    }

    return $table;
}

//更新数据
function Update() {
    $.getJSON(location.href, 'data=' + Date(), function (res) {
        for (var k in res.data) {
            var $table = $("#" + k);
            if (res.data[k].HasError != '1') {
                for (var k2 in res.data[k]) {
                    var val = res.data[k][k2];
                    var cell = $("#" + k + "_" + k2);
                    cell.html(val);
                    var $parent = $($table.parent());
                    $($parent.find('.table-vector-time')).html('数据时间:' + res.data[k]['AcquireTime']);
                    $parent.removeClass('table-vector-error');
                }
            } else {
                $($table.parent()).addClass('table-vector-error');
            }
        }

        for (var k in res.remote) {
            var $remote = $("#" + k);
            if (res.remote[k] == null) {
                $remote.addClass('remote-error');
                $("#" + k + "_t").addClass('remote-error');
            } else {
                $remote.removeClass('remote-error');
                $("#" + k + "_t").removeClass('remote-error');
                $remote.html(res.remote[k]);
            }
        }
    });
}