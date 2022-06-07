function makeAlbumElem(libraryId, albumSubPath) {
    // container
    let divElem = document.createElement("div");
    divElem.className = 'albumEntry';

    // cover image
    let coverUrl = sC.renderPath('/apis/album/cover/<str:libid>/<str:path>',
        {libid: libraryId, path: encodeURIComponent(albumSubPath)});
    let coverElem = document.createElement("img");
    coverElem.src = coverUrl;
    divElem.appendChild(coverElem);

    // link text
    let albumUrl = sC.renderPath('/thumbnails/<str:libid>/<str:path>',
        {libid: libraryId, path: encodeURIComponent(albumSubPath)});
    let linkElem = document.createElement("div");
    linkElem.innerHTML = "<a href=\"" + albumUrl + "\">" + albumSubPath + "</a>";
    divElem.appendChild(linkElem);

    return divElem;
}

function populateLibrary(libDesc) {
    let pathParam = sC.parsePath('/library/<str:libId>', sC.getPath());
    let libraryId = pathParam.libId;

    let wrapperDiv = document.getElementById('wrapper');

    libDesc.albums.sort();
    for (let i = 0; i < libDesc.albums.length; i++) {
        let albumId = libDesc.albums[i];
        let albumElem = makeAlbumElem(libraryId, albumId);
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
