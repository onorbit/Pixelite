// @ts-ignore
import { SCURLPattern, SCAjax } from "/statics/scripts/sugarcube.js";

function sendRescanRequest(id) {
    let scanPattern = new SCURLPattern('apis/library/<str:id>/rescan');
    let scanPath = scanPattern.render({ 'id': id});
    SCAjax.post(scanPath, null, function(status, response) {
        // TODO : notification?
    });
}

function sendUnmountRequest(id) {
    let unmountPattern = new SCURLPattern('apis/library/<str:id>');
    let unmountPath = unmountPattern.render({ 'id': id });
    SCAjax.delete(unmountPath, function(status, response) {
        location.reload();
    });
}

function makeLibraryElem(id, title) {
    let wrapperElem = document.createElement("div");

    // make library link
    let spanElem = document.createElement("span");
    spanElem.innerText = title;

    let linkElem = document.createElement("a");
    let libPattern = new SCURLPattern('/library/<str:id>');
    linkElem.href = libPattern.render({ 'id': id });
    linkElem.appendChild(spanElem);

    // append library link to wrapper
    wrapperElem.appendChild(linkElem);

    // make rescan button
    let rescanButtonElem = document.createElement("span");
    rescanButtonElem.innerText = "rescan";
    rescanButtonElem.addEventListener('click', function(){ sendRescanRequest(id); });

    // append rescan button to wrapper
    wrapperElem.appendChild(rescanButtonElem);

    // make unmount button
    let unmountButtonElem = document.createElement("span");
    unmountButtonElem.innerText = "unmount";
    unmountButtonElem.addEventListener('click', function() { sendUnmountRequest(id); })

    // append rescan button to wrapper
    wrapperElem.appendChild(unmountButtonElem);

    return wrapperElem;
}

function populateLibraries(libList) {
    let wrapperDiv = document.getElementById('wrapper');

    for (let i = 0; i < libList.length; i++) {
        let lib = libList[i];
        let libElem = makeLibraryElem(lib.id, lib.title);
        wrapperDiv.appendChild(libElem);
    }
}

function submitNewLibrary() {
    let inputElem = document.getElementById('rootPath') as HTMLInputElement;

    let rootPath = inputElem.value;
    if (rootPath.length == 0) {
        return;
    }

    let apiUrl = '/apis/library';
    SCAjax.post(apiUrl, {'rootPath': rootPath}, function(status, response) {
        if (status != 200) {
            return;
        }

        // TODO : the API should return appropriate response.
        location.reload();
    });
}

function bootStrap() {
    // append new library submit button.
    let submitButton = document.createElement('input');
    submitButton.type = 'button';
    submitButton.value = 'new library';
    submitButton.onclick = submitNewLibrary;

    document.getElementById('newLibraryDiv').append(submitButton);

    // request library list.
    let apiUrl = '/apis/libraries';
    SCAjax.get(apiUrl, function(status, response) {
        if (status == 200) {
            let libList = JSON.parse(response);
            populateLibraries(libList);
        }
    });
}

bootStrap();
