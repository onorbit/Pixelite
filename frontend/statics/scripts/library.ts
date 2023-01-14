// @ts-ignore
import { SCURLPattern, SCAjax } from "/statics/scripts/sugarcube.js";

function makeAlbumElem(libraryId, albumSubPath) {
    // container
    let divElem = document.createElement("div");
    divElem.className = 'albumEntry widthLimiter';

    // cover image
    let coverPattern = new SCURLPattern('/apis/album/cover/<str:libid>/<str:path>');
    let coverPath = coverPattern.render({ 'libid': libraryId, 'path': encodeURIComponent(albumSubPath) });
    let coverElem = document.createElement("img");
    coverElem.src = coverPath;
    divElem.appendChild(coverElem);

    // link text
    let albumPattern = new SCURLPattern('/thumbnails/<str:libid>/<str:path>');
    let albumPath = albumPattern.render({ 'libid': libraryId, 'path': encodeURIComponent(albumSubPath) });
    let linkElem = document.createElement("div");
    linkElem.innerHTML = "<a href=\"" + albumPath + "\">" + albumSubPath + "</a>";
    divElem.appendChild(linkElem);

    return divElem;
}

function populateLibrary(libDesc) {
    let libPattern = new SCURLPattern('/library/<str:libId>');
    let params = libPattern.parse(window.location.pathname);
    let libraryId = params.libId;

    let wrapperDiv = document.getElementById('wrapper');

    libDesc.albums.sort();
    for (let i = 0; i < libDesc.albums.length; i++) {
        let albumId = libDesc.albums[i];
        let albumElem = makeAlbumElem(libraryId, albumId);
        wrapperDiv.appendChild(albumElem);
    }
}

function bootStrap() {
    // fetch albums.
    let pathParam = SCURLPattern.parse('/library/<str:libId>', window.location.pathname);
    let apiUrl = SCURLPattern.render('/apis/library/<str:libId>', pathParam);
    SCAjax.get(apiUrl, function(status, response) {
        if (status == 200) {
            let libDesc = JSON.parse(response);
            populateLibrary(libDesc);
        }
    });

    // get configuration and set max-width of cover.
    SCAjax.get('/apis/configs', function(status, response) {
        if (status == 200) {
            let conf = JSON.parse(response);
            let coverMaxWidth = conf.coverSize;
            let classWidthLimiter = '.widthLimiter { max-width: ' + coverMaxWidth + 'px }'

            document.styleSheets[0].insertRule(classWidthLimiter, 0)
        }
    });
}

bootStrap();
