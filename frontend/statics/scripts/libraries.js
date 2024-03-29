// @ts-ignore
import { SCURLPattern, SCAjax } from "/statics/scripts/sugarcube.js";
function sendRescanRequest(id) {
    var scanPattern = new SCURLPattern('apis/library/<str:id>/rescan');
    var scanPath = scanPattern.render({ 'id': id });
    SCAjax.post(scanPath, null, function (status, response) {
        // TODO : notification?
    });
}
function sendUnmountRequest(id) {
    var unmountPattern = new SCURLPattern('apis/library/<str:id>');
    var unmountPath = unmountPattern.render({ 'id': id });
    SCAjax["delete"](unmountPath, function (status, response) {
        location.reload();
    });
}
function makeLibraryElem(id, title) {
    var wrapperElem = document.createElement("div");
    // make library link
    var spanElem = document.createElement("span");
    spanElem.innerText = title;
    var linkElem = document.createElement("a");
    var libPattern = new SCURLPattern('/library/<str:id>');
    linkElem.href = libPattern.render({ 'id': id });
    linkElem.appendChild(spanElem);
    // append library link to wrapper
    wrapperElem.appendChild(linkElem);
    // make rescan button
    var rescanButtonElem = document.createElement("span");
    rescanButtonElem.innerText = "rescan";
    rescanButtonElem.addEventListener('click', function () { sendRescanRequest(id); });
    // append rescan button to wrapper
    wrapperElem.appendChild(rescanButtonElem);
    // make unmount button
    var unmountButtonElem = document.createElement("span");
    unmountButtonElem.innerText = "unmount";
    unmountButtonElem.addEventListener('click', function () { sendUnmountRequest(id); });
    // append rescan button to wrapper
    wrapperElem.appendChild(unmountButtonElem);
    return wrapperElem;
}
function populateLibraries(libList) {
    var wrapperDiv = document.getElementById('wrapper');
    for (var i = 0; i < libList.length; i++) {
        var lib = libList[i];
        var libElem = makeLibraryElem(lib.id, lib.title);
        wrapperDiv.appendChild(libElem);
    }
}
function submitNewLibrary() {
    var inputElem = document.getElementById('rootPath');
    var rootPath = inputElem.value;
    if (rootPath.length == 0) {
        return;
    }
    var apiUrl = '/apis/library';
    SCAjax.post(apiUrl, { 'rootPath': rootPath }, function (status, response) {
        if (status != 200) {
            return;
        }
        // TODO : the API should return appropriate response.
        location.reload();
    });
}
function bootStrap() {
    // append new library submit button.
    var submitButton = document.createElement('input');
    submitButton.type = 'button';
    submitButton.value = 'new library';
    submitButton.onclick = submitNewLibrary;
    document.getElementById('newLibraryDiv').append(submitButton);
    // request library list.
    var apiUrl = '/apis/libraries';
    SCAjax.get(apiUrl, function (status, response) {
        if (status == 200) {
            var libList = JSON.parse(response);
            populateLibraries(libList);
        }
    });
}
bootStrap();
