function makeAlbumElem(albumSubPath) {
    let libUrl = sC.renderPath('/thumbnails/<str:path>', {path: encodeURIComponent(albumSubPath)});

    let divElem = document.createElement("div");
    divElem.innerText = albumSubPath;

    let linkElem = document.createElement("a");
    linkElem.href = libUrl;

    linkElem.appendChild(divElem);

    return linkElem;
}

function populateLibrary(libDesc) {
    let wrapperDiv = document.getElementById('wrapper');

    for (let i = 0; i < libDesc.albums.length; i++) {
        let albumSubPath = libDesc.albums[i];
        let albumElem = makeAlbumElem(albumSubPath);
        wrapperDiv.appendChild(albumElem);
    }
}

function bootStrap() {
    let pathParam = sC.parsePath('/library/<str:libId>', sC.getPath());
    let apiUrl = sC.renderPath('/apis/library/<str:libId>', pathParam);

    sC.ajaxGet(apiUrl, function(status, response) {
        if (status == 200) {
            let libDesc = JSON.parse(response)
            populateLibrary(libDesc);
        }
    });
}

bootStrap();
