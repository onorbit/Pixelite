function makeThumbnailElem(url) {
    let imgElem = document.createElement("img");
    imgElem.src = url;
    imgElem.className = 'thumbnail';

    let divElem = document.createElement("div");
    divElem.className = 'thumbnail';

    let linkElem = document.createElement("a");
    linkElem.href = '#';
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
            let thumbnailUrl = sC.renderPath('/apis/thumbnail/<str:libId>/<str:albumId>/<str:fileName>',
                { libId: libId, albumId: albumId, fileName: encodeURIComponent(fileEntry.name) });
            let thumbnailElem = makeThumbnailElem(thumbnailUrl);

            wrapperDiv.appendChild(thumbnailElem);
        }
    }
}

function bootStrap() {
    // extract subpath.
    let pathParam = sC.parsePath('/thumbnails/<str:libId>/<str:albumId>', sC.getPath());

    // request file list.
    let apiUrl = sC.renderPath('/apis/list/<str:libId>/<str:albumId>', pathParam);
    sC.ajaxGet(apiUrl, function(status, response) {
        if (status == 200) {
            let resultObj = JSON.parse(response);
            populateThumbnails(pathParam.libId, pathParam.albumId, resultObj);
        }
    });
}

bootStrap();
