function copy_tooltip(action) {
    let copy_btn = document.querySelector('#clip_copy');

    if (action === "on") {
        copy_btn.classList.add('tooltip', 'tooltip-bottom');
    } else if (action === "off") {
        copy_btn.classList.remove('tooltip', 'tooltip-bottom');
    }
}

document.addEventListener('DOMContentLoaded', function () {
    let clipboard = new ClipboardJS('#clip_copy');

    clipboard.on('success', function (e) {
        let render_body = document.getElementById('render_results');
        if (render_body.innerHTML) {
            copy_tooltip("on");
        }

        e.clearSelection();
    });
}, false);

function request_render(event) {
    event.preventDefault();

    let j2_template = $('#j2_template').val();
    let j2_data = $('#j2_data').val();

    let render_options = {
        strict: $('#opt_strict').prop('checked') === true,
        trim: $('#opt_trim').prop('checked') === true,
        lstrip: $('#opt_lstrip').prop('checked') === true,
    };

    let additional_filters = [];
    $('.add_filters').each(function (index, obj) {
        if ($(obj).prop('checked') === true) {
            additional_filters.push(obj.id);
        }
    });

    //let render_mode = $("input[name='render_mode']:checked").val();

    let render_data = {
        template: j2_template,
        data: j2_data,
        /*options: render_options,
        filters: additional_filters,
        render_mode: render_mode*/
    };

    $.ajax({
        type: 'POST',
        url: 'http://127.0.0.1:8080/device-translate',
        data: JSON.stringify(render_data),
        contentType: 'application/json',
     //   dataType: 'yaml',
    }).done(function (reply) {
        //let rendered_template = reply['render_result'];
        let rendered_template = reply;
        let rendered_html = classify_whitespaces(rendered_template);
        let render_body = document.getElementById('render_results');
        if (render_body.firstChild) {
            render_body.replaceChild(rendered_html, render_body.firstChild);
        } else {
            render_body.appendChild(rendered_html);
        }
        toggle_whitespaces();
        copy_tooltip("off");
    }).fail(function (e) {
        //alert("I'm so sorry, did not get reply from server.");
        alert(e.responseText)
        alert(e.readyState)
        alert(e.status)
        alert(e.statusText)
    });
}

function reset_render(event) {
    let render_body = document.getElementById('render_results');
    if (render_body.innerHTML) {
        render_body.innerHTML = '';
        copy_tooltip("off");
    }
}

function toggle_whitespaces() {
    let show_ws_on = $('#toggle_whitespaces').prop('checked') === true;

    if (show_ws_on === true) {
        let hidden_ws = document.querySelectorAll('.ws_space, .ws_tab, .ws_newline');
        hidden_ws.forEach(el => el.classList.add('ws_vis'));
    } else {
        let visible_ws = document.querySelectorAll('.ws_vis');
        visible_ws.forEach(el => el.classList.remove('ws_vis'))
    }
}


function classify_whitespaces(text) {
    let html_w_ws = document.createElement('span');

    let ws_space = document.createElement('span');
    ws_space.classList.add('ws_space');
    ws_space.innerHTML = ' ';

    let ws_tab = document.createElement('span');
    ws_tab.classList.add('ws_tab');
    ws_tab.innerHTML = '\t';

    let ws_newline = document.createElement('span');
    ws_newline.classList.add('ws_newline');
    ws_newline.innerHTML = ' ';

    let normal_text = [];
    for (let i = 0; i < text.length; i++) {
        if ([' ', '\t', '\n'].includes(text[i])) {
            if (normal_text.length > 0) {
                html_w_ws.append(document.createTextNode(normal_text.join('')));
                normal_text = []
            }
            switch (text[i]) {
                case ' ':
                    html_w_ws.append(ws_space.cloneNode(true));
                    break;
                case '\t':
                    html_w_ws.append(ws_tab.cloneNode(true));
                    break;
                case '\n':
                    html_w_ws.append(ws_newline.cloneNode(true));
                    normal_text.push('\n');
                    break;
            }
        } else {
            normal_text.push(text[i])
        }
    }
    if (normal_text.length > 0) {
        html_w_ws.append(document.createTextNode(normal_text.join('')));
    }

    return html_w_ws
}

$(function () {
    $('#request_render').on('click', request_render);
});

$(function () {
    $('#reset_render').on('click', reset_render);
});

$(function () {
    $('#toggle_whitespaces').change(toggle_whitespaces);
});