function makeAlbumElem(libraryId, albumSubPath) {
    let albumUrl = sC.renderPath('/thumbnails/<str:libid>/<str:path>',
        {libid: libraryId, path: encodeURIComponent(albumSubPath)});

    let divElem = document.createElement("div");
    divElem.innerText = albumSubPath;

    let linkElem = document.createElement("a");
    linkElem.href = albumUrl;

    linkElem.appendChild(divElem);

    return linkElem;
}

function populateLibrary(libDesc) {
    let pathParam = sC.parsePath('/library/<str:libId>', sC.getPath());
    let libraryId = pathParam.libId;

    let wrapperDiv = document.getElementById('wrapper');

    libDesc.albums.sort();

    for (let i = 0; i < libDesc.albums.length; i++) {
        let albumSubPath = libDesc.albums[i];
        let albumElem = makeAlbumElem(libraryId, albumSubPath);
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
