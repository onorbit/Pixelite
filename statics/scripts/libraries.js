function makeLibraryElem(id, desc) {
    let libUrl = sC.renderPath('/library/<str:id>', {id: id});

    let divElem = document.createElement("div");
    divElem.innerText = desc;

    let linkElem = document.createElement("a");
    linkElem.href = libUrl;

    linkElem.appendChild(divElem);

    return linkElem;
}

function populateLibraries(libList) {
    let wrapperDiv = document.getElementById('wrapper');

    for (let i = 0; i < libList.length; i++) {
        let lib = libList[i];
        let libElem = makeLibraryElem(lib.id, lib.desc);
        wrapperDiv.appendChild(libElem);
    }
}

function bootStrap() {
    // request library list.
    let apiUrl = '/apis/libraries';
    sC.ajaxGet(apiUrl, function(status, response) {
        if (status == 200) {
            let libList = JSON.parse(response)
            populateLibraries(libList);
        }
    });
}

bootStrap();
