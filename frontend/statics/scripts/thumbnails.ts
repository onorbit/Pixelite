// @ts-ignore
import { SCURLPattern, SCAjax } from "/statics/scripts/sugarcube.js";

function makeThumbnailElem(libId, albumId, fileName) {
    let thumbnailUrl = SCURLPattern.render('/apis/album/thumbnail/<str:libId>/<str:albumId>/<str:fileName>',
                { libId: libId, albumId: albumId, fileName: encodeURIComponent(fileName) });

    let imgElem = document.createElement("img");
    imgElem.src = thumbnailUrl;
    imgElem.className = 'thumbnail';

    let divElem = document.createElement("div");
    divElem.className = 'thumbnail';

    let rawImageUrl = SCURLPattern.render('/apis/album/image/<str:libId>/<str:albumId>/<str:fileName>',
                { libId: libId, albumId: albumId, fileName: encodeURIComponent(fileName) });

    let linkElem = document.createElement("a");
    linkElem.href = rawImageUrl;
    linkElem.className = 'thumbnail';

    divElem.appendChild(imgElem);
    linkElem.appendChild(divElem);

    return linkElem;
}

function populateThumbnails(libId, albumId, fileList) {
    let wrapperDiv = document.getElementById('wrapper');

    for (let i = 0; i < fileList.length; i++) {
        let fileEntry = fileList[i];

        if (fileEntry.type == 1) {
            let thumbnailElem = makeThumbnailElem(libId, albumId, fileEntry.name);
            wrapperDiv.appendChild(thumbnailElem);
        }
    }
}

function bootStrap() {
    // extract subpath.
    let pathParam = SCURLPattern.parse('/thumbnails/<str:libId>/<str:albumId>', window.location.pathname);

    // request file list.
    let apiUrl = SCURLPattern.render('/apis/album/list/<str:libId>/<str:albumId>', pathParam);
    SCAjax.get(apiUrl, function(status, response) {
        if (status == 200) {
            let resultObj = JSON.parse(response);
            populateThumbnails(pathParam.libId, pathParam.albumId, resultObj);
        }
    });
}

bootStrap();
