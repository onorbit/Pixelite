// @ts-ignore
import { SCURLPattern, SCAjax } from "/statics/scripts/sugarcube.js";
function makeThumbnailElem(libId, albumId, fileName) {
    var thumbnailUrl = SCURLPattern.render('/apis/album/thumbnail/<str:libId>/<str:albumId>/<str:fileName>', { libId: libId, albumId: albumId, fileName: encodeURIComponent(fileName) });
    var imgElem = document.createElement("img");
    imgElem.src = thumbnailUrl;
    imgElem.className = 'thumbnail';
    var divElem = document.createElement("div");
    divElem.className = 'thumbnail';
    var rawImageUrl = SCURLPattern.render('/apis/album/image/<str:libId>/<str:albumId>/<str:fileName>', { libId: libId, albumId: albumId, fileName: encodeURIComponent(fileName) });
    var linkElem = document.createElement("a");
    linkElem.href = rawImageUrl;
    linkElem.className = 'thumbnail';
    divElem.appendChild(imgElem);
    linkElem.appendChild(divElem);
    return linkElem;
}
function populateThumbnails(libId, albumId, fileList) {
    var wrapperDiv = document.getElementById('wrapper');
    for (var i = 0; i < fileList.length; i++) {
        var fileEntry = fileList[i];
        if (fileEntry.type == 1) {
            var thumbnailElem = makeThumbnailElem(libId, albumId, fileEntry.name);
            wrapperDiv.appendChild(thumbnailElem);
        }
    }
}
function bootStrap() {
    // extract subpath.
    var pathParam = SCURLPattern.parse('/thumbnails/<str:libId>/<str:albumId>', window.location.pathname);
    // request file list.
    var apiUrl = SCURLPattern.render('/apis/album/list/<str:libId>/<str:albumId>', pathParam);
    SCAjax.get(apiUrl, function (status, response) {
        if (status == 200) {
            var resultObj = JSON.parse(response);
            populateThumbnails(pathParam.libId, pathParam.albumId, resultObj);
        }
    });
}
bootStrap();
