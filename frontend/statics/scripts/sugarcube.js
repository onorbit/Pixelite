var SCURLPattern = /** @class */ (function () {
    function SCURLPattern(pattern) {
        this.tokens = [];
        // remove trailing slash in pattern
        if (pattern[pattern.length - 1] == '/') {
            pattern = pattern.substring(0, pattern.length - 1);
        }
        // explode and scan for placeholders.
        var placeholders = {};
        var inputTokens = pattern.split('/');
        for (var i = 0; i < inputTokens.length; i++) {
            var token = inputTokens[i];
            if (token[0] != '<' || token[token.length - 1] != '>') {
                // in case of plain text
                this.tokens.push(token);
            }
            else {
                // in case of substitute pattern
                token = token.substring(1, token.length - 1);
                var parts = token.split(':');
                if (parts.length != 2) {
                    // invalid pattern format.
                    throw new Error();
                }
                else if (parts[1].length <= 0) {
                    // placeholder name is empty.
                    throw new Error();
                }
                else if (placeholders.hasOwnProperty(parts[1]) == true) {
                    // placeholder name already exists.
                    throw new Error();
                }
                else if (parts[0] != 'str' && parts[0] != 'int') {
                    // invalid type.
                    throw new Error();
                }
                placeholders[parts[1]] = '';
                var newToken = { type: parts[0], name: parts[1] };
                this.tokens.push(newToken);
            }
        }
    }
    SCURLPattern.prototype.render = function (params) {
        var ret = '';
        for (var i = 0; i < this.tokens.length; i++) {
            var token = this.tokens[i];
            switch (typeof (token)) {
                case 'string':
                    ret = ret + token;
                    break;
                case 'object':
                    var name_1 = token['name'];
                    var type = token['type'];
                    if (params.hasOwnProperty(name_1) == false) {
                        return null;
                    }
                    var param = params[name_1];
                    if (type == 'str') {
                        if (typeof (param) != 'string') {
                            return null;
                        }
                        ret = ret + param;
                    }
                    else if (type == 'int') {
                        if (typeof (param) != 'number') {
                            return null;
                        }
                        ret = ret + param.toString();
                    }
                    else {
                        return null;
                    }
                    break;
                default:
                    return null;
            }
            if (i != this.tokens.length - 1) {
                ret = ret + '/';
            }
        }
        return ret;
    };
    SCURLPattern.prototype.parse = function (path) {
        var pathTokens = path.split('/');
        if (pathTokens.length != this.tokens.length) {
            return null;
        }
        // process pattern.
        var parsedValues = {};
        for (var i = 0; i < this.tokens.length; i++) {
            var token = this.tokens[i];
            switch (typeof (token)) {
                case 'string':
                    if (token != pathTokens[i]) {
                        return null;
                    }
                    break;
                case 'object':
                    var name_2 = token['name'];
                    var type = token['type'];
                    if (type == 'str') {
                        parsedValues[name_2] = pathTokens[i];
                    }
                    else if (type == 'int') {
                        var value = parseInt(pathTokens[i]);
                        if (value == NaN) {
                            return null;
                        }
                        parsedValues[name_2] = value;
                    }
                    else {
                        return null;
                    }
                    break;
                default:
                    return null;
            }
        }
        return parsedValues;
    };
    SCURLPattern.render = function (pattern, params) {
        var ptn = new SCURLPattern(pattern);
        return ptn.render(params);
    };
    SCURLPattern.parse = function (pattern, path) {
        var ptn = new SCURLPattern(pattern);
        return ptn.parse(path);
    };
    return SCURLPattern;
}());
export { SCURLPattern };
var SCAjax = /** @class */ (function () {
    function SCAjax() {
    }
    SCAjax.get = function (URL, callbackFunc) {
        var request = new XMLHttpRequest();
        if (callbackFunc != null) {
            request.onload = function () {
                callbackFunc(request.status, request.response);
            };
        }
        request.open('GET', URL);
        request.send();
    };
    SCAjax["delete"] = function (URL, callbackFunc) {
        var request = new XMLHttpRequest();
        if (callbackFunc != null) {
            request.onload = function () {
                callbackFunc(request.status, request.response);
            };
        }
        request.open('DELETE', URL);
        request.send();
    };
    SCAjax.post = function (URL, payload, callbackFunc) {
        var request = new XMLHttpRequest();
        if (callbackFunc != null) {
            request.onload = function () {
                callbackFunc(request.status, request.response);
            };
        }
        var toSend = null;
        if (payload != null) {
            switch (typeof (payload)) {
                case 'string':
                    toSend = payload;
                    break;
                case 'object':
                    toSend = new FormData();
                    for (var key in payload) {
                        toSend.append(key, payload[key]);
                    }
                    break;
                default:
                    return;
            }
        }
        request.open('POST', URL);
        request.send(toSend);
    };
    return SCAjax;
}());
export { SCAjax };
var SCLoad = /** @class */ (function () {
    function SCLoad() {
    }
    SCLoad.loadScript = function (scriptURL, onLoadFunc) {
        var scriptElem = document.createElement("script");
        scriptElem.src = scriptURL;
        scriptElem.onload = onLoadFunc;
        document.body.appendChild(scriptElem);
    };
    SCLoad.loadCSS = function (cssURL) {
        var cssElem = document.createElement("link");
        cssElem.rel = "stylesheet";
        cssElem.type = "text/css";
        cssElem.href = cssURL;
        document.head.appendChild(cssElem);
    };
    return SCLoad;
}());
export { SCLoad };
var SCModal = /** @class */ (function () {
    function SCModal() {
    }
    SCModal.popUp = function (width, height, contentElem) {
        if (document.getElementById(SCModal.BACKGROUND_ID) != null) {
            return;
        }
        // background
        var backgroundElem = document.createElement("div");
        backgroundElem.id = SCModal.BACKGROUND_ID;
        backgroundElem.setAttribute('style', SCModal.BACKGROUND_STYLE);
        // close button of dialog
        var closeElem = document.createElement("a");
        closeElem.href = "#";
        closeElem.onclick = SCModal.close;
        closeElem.setAttribute('style', SCModal.CLOSE_STYLE);
        closeElem.innerText = "×";
        // dialog
        var dialogStyle = SCModal.DIALOG_STYLE.concat("width:", width.toString(), "px; height:", height.toString(), "px; margin-top:", (height / -2).toString(), "px; margin-left:", (width / -2).toString(), "px;");
        var dialogElem = document.createElement("div");
        dialogElem.setAttribute('style', dialogStyle);
        dialogElem.appendChild(closeElem);
        if (contentElem != null) {
            dialogElem.appendChild(contentElem);
        }
        backgroundElem.appendChild(dialogElem);
        document.body.appendChild(backgroundElem);
    };
    SCModal.close = function () {
        var backgroundElem = document.getElementById(SCModal.BACKGROUND_ID);
        if (backgroundElem != null) {
            document.body.removeChild(backgroundElem);
        }
    };
    SCModal.BACKGROUND_ID = "sugarCubeModalBackground";
    SCModal.BACKGROUND_STYLE = "position:fixed; padding:0; margin:0; top:0; left:0; width:100%; height:100%; background:rgba(180,180,180,0.5); backdrop-filter: blur(3px); z-index:65535;";
    SCModal.DIALOG_STYLE = "position:fixed; top:50%; left:50%; background-color:#FFF; border-radius:10px; z-index:65536; padding: 40px;";
    SCModal.CLOSE_STYLE = "font-size:24pt; position:absolute; top:0px; right:15px;";
    return SCModal;
}());
export { SCModal };
