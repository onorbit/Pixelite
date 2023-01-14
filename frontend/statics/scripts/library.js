// @ts-ignore
import { SCURLPattern, SCAjax } from "/statics/scripts/sugarcube.js";
function makeAlbumElem(libraryId, albumSubPath) {
    // container
    var divElem = document.createElement("div");
    divElem.className = 'albumEntry widthLimiter';
    // cover image
    var coverPattern = new SCURLPattern('/apis/album/cover/<str:libid>/<str:path>');
    var coverPath = coverPattern.render({ 'libid': libraryId, 'path': encodeURIComponent(albumSubPath) });
    var coverElem = document.createElement("img");
    coverElem.src = coverPath;
    divElem.appendChild(coverElem);
    // link text
    var albumPattern = new SCURLPattern('/thumbnails/<str:libid>/<str:path>');
    var albumPath = albumPattern.render({ 'libid': libraryId, 'path': encodeURIComponent(albumSubPath) });
    var linkElem = document.createElement("div");
    linkElem.innerHTML = "<a href=\"" + albumPath + "\">" + albumSubPath + "</a>";
    divElem.appendChild(linkElem);
    return divElem;
}
function populateLibrary(libDesc) {
    var libPattern = new SCURLPattern('/library/<str:libId>');
    var params = libPattern.parse(window.location.pathname);
    var libraryId = params.libId;
    var wrapperDiv = document.getElementById('wrapper');
    libDesc.albums.sort();
    for (var i = 0; i < libDesc.albums.length; i++) {
        var albumId = libDesc.albums[i];
        var albumElem = makeAlbumElem(libraryId, albumId);
        wrapperDiv.appendChild(albumElem);
    }
}
function bootStrap() {
    // fetch albums.
    var pathParam = SCURLPattern.parse('/library/<str:libId>', window.location.pathname);
    var apiUrl = SCURLPattern.render('/apis/library/<str:libId>', pathParam);
    SCAjax.get(apiUrl, function (status, response) {
        if (status == 200) {
            var libDesc = JSON.parse(response);
            populateLibrary(libDesc);
        }
    });
    // get configuration and set max-width of cover.
    SCAjax.get('/apis/configs', function (status, response) {
        if (status == 200) {
            var conf = JSON.parse(response);
            var coverMaxWidth = conf.coverSize;
            var classWidthLimiter = '.widthLimiter { max-width: ' + coverMaxWidth + 'px }';
            document.styleSheets[0].insertRule(classWidthLimiter, 0);
        }
    });
}
bootStrap();
