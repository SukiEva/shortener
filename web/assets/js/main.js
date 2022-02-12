function doPost() {
    $.ajax({
        type: "POST",
        url: "/api/v1/generate",
        data: JSON.stringify({
            url: $("#url").val(),
            // exp: parseInt($("#exp").val())
        }),
        dataType: "json",
        contentType: "application/json; charset=utf-8",
        success: function (data) {
            if (data.code === "0") {
                alert("错误：" + data.error);
                return;
            }
            const shorter = window.location.href + data.token
            const result = document.getElementById("result");
            result.innerHTML = "<p>shorten to: </p><a href=\"" + shorter + "\" target=\"_blank\">" + shorter + "</a>"
            document.getElementById("form").appendChild(result)
        }
    })
}