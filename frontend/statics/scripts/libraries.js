function sendRescanRequest(id) {
    let rescanUrl = sC.renderPath('apis/library/<str:id>/rescan', {id: id})
    sC.ajaxPost(rescanUrl, null, function(status, response) {
        // TODO : notification?
    })
}

function makeLibraryElem(id, title) {
    let wrapperElem = document.createElement("div");

    // make library link
    let spanElem = document.createElement("span");
    spanElem.innerText = title;

    let linkElem = document.createElement("a");
    linkElem.href = sC.renderPath('/library/<str:id>', {id: id});
    linkElem.appendChild(spanElem);

    // append library link to wrapper
    wrapperElem.appendChild(linkElem);

    // make rescan button
    let rescanButtonElem = document.createElement("span");
    rescanButtonElem.innerText = "rescan";
    rescanButtonElem.addEventListener('click', function(){ sendRescanRequest(id); });

    // append rescan button to wrapper
    wrapperElem.appendChild(rescanButtonElem);

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
    let inputElem = document.getElementById('rootPath');

    let rootPath = inputElem.value;
    if (rootPath.length == 0) {
        return;
    }

    let apiUrl = '/apis/library';
    sC.ajaxPost(apiUrl, {'rootPath': rootPath}, function(status, response) {
        if (status != 200) {
            return;
        }

        // TODO : the API should return appropriate response.
        location.reload();
    });
}

function bootStrap() {
    // request library list.
    let apiUrl = '/apis/libraries';
    sC.ajaxGet(apiUrl, function(status, response) {
        if (status == 200) {
            let libList = JSON.parse(response);
            populateLibraries(libList);
        }
    });
}

bootStrap();
