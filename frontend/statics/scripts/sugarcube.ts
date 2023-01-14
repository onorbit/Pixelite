export class SCURLPattern {
    tokens: any[]

    public constructor(pattern: string) {
        this.tokens = [];

        // remove trailing slash in pattern
        if (pattern[pattern.length - 1] == '/') {
            pattern = pattern.substring(0, pattern.length - 1);
        }

        // explode and scan for placeholders.
        let placeholders = {};
        let inputTokens = pattern.split('/');
        for (let i = 0; i < inputTokens.length; i++) {
            let token = inputTokens[i];

            if (token[0] != '<' || token[token.length - 1] != '>') {
                // in case of plain text
                this.tokens.push(token);
            } else {
                // in case of substitute pattern
                token = token.substring(1, token.length - 1);
                let parts = token.split(':');

                if (parts.length != 2) {
                    // invalid pattern format.
                    throw new Error();
                } else if (parts[1].length <= 0) {
                    // placeholder name is empty.
                    throw new Error();
                } else if (placeholders.hasOwnProperty(parts[1]) == true) {
                    // placeholder name already exists.
                    throw new Error();
                } else if (parts[0] != 'str' && parts[0] != 'int') {
                    // invalid type.
                    throw new Error();
                }

                placeholders[parts[1]] = ''
                let newToken = { type: parts[0], name: parts[1] }
                this.tokens.push(newToken);
            }
        }
    }

    public render(params: object): string {
        let ret = '';
        for (let i = 0; i < this.tokens.length; i++) {
            let token = this.tokens[i]
            switch (typeof(token)) {
                case 'string':
                    ret = ret + token;
                    break;
                case 'object':
                    let name = token['name'];
                    let type = token['type'];

                    if (params.hasOwnProperty(name) == false) {
                        return null;
                    }

                    let param = params[name];
                    if (type == 'str') {
                        if (typeof(param) != 'string') {
                            return null;
                        }
                        ret = ret + param;
                    } else if (type == 'int') {
                        if (typeof(param) != 'number') {
                            return null;
                        }
                        ret = ret + param.toString();
                    } else {
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
    }

    public parse(path: string): object {
        let pathTokens = path.split('/');
        if (pathTokens.length != this.tokens.length) {
            return null;
        }

        // process pattern.
        let parsedValues = {};
        for (let i = 0; i < this.tokens.length; i++) {
            let token = this.tokens[i];
            switch (typeof(token)) {
                case 'string':
                    if (token != pathTokens[i]) {
                        return null;
                    }
                    break;
                case 'object':
                    let name = token['name'];
                    let type = token['type'];

                    if (type == 'str') {
                        parsedValues[name] = pathTokens[i];
                    } else if (type == 'int') {
                        let value = parseInt(pathTokens[i]);
                        if (value == NaN) {
                            return null;
                        }
                        parsedValues[name] = value;
                    } else {
                        return null;
                    }
                    break;
                default:
                    return null;
            }
        }

        return parsedValues;
    }

    public static render(pattern: string, params: object): string {
        let ptn = new SCURLPattern(pattern);
        return ptn.render(params);
    }

    public static parse(pattern: string, path: string): object {
        let ptn = new SCURLPattern(pattern);
        return ptn.parse(path);
    }
}

export abstract class SCAjax {
    public static get(URL: string, callbackFunc: any) {
        let request = new XMLHttpRequest();

        if (callbackFunc != null) {
            request.onload = function() {
                callbackFunc(request.status, request.response);
            }
        }

        request.open('GET', URL);
        request.send();
    }

    public static delete(URL: string, callbackFunc: any) {
        let request = new XMLHttpRequest();

        if (callbackFunc != null) {
            request.onload = function() {
                callbackFunc(request.status, request.response);
            }
        }

        request.open('DELETE', URL);
        request.send();
    }

    public static post(URL: string, payload: any, callbackFunc: any) {
    let request = new XMLHttpRequest();

    if (callbackFunc != null) {
        request.onload = function() {
            callbackFunc(request.status, request.response);
        }
    }

    let toSend = null;
    if (payload != null) {
        switch (typeof(payload)) {
            case 'string':
                toSend = payload;
                break;
            case 'object':
                toSend = new FormData();
                for (const key in payload) {
                    toSend.append(key, payload[key])
                }
                break;
            default:
                return;
        }
    }

    request.open('POST', URL);
    request.send(toSend);
    }
}

export abstract class SCLoad {
    public static loadScript(scriptURL: string, onLoadFunc: any) {
        let scriptElem = document.createElement("script");
        scriptElem.src = scriptURL;
        scriptElem.onload = onLoadFunc;

        document.body.appendChild(scriptElem);
    }

    public static loadCSS(cssURL: string) {
        let cssElem = document.createElement("link");
        cssElem.rel = "stylesheet";
        cssElem.type = "text/css";
        cssElem.href = cssURL;

        document.head.appendChild(cssElem);
    }
}

export abstract class SCModal {
    static BACKGROUND_ID = "sugarCubeModalBackground";
    static BACKGROUND_STYLE = "position:fixed; padding:0; margin:0; top:0; left:0; width:100%; height:100%; background:rgba(180,180,180,0.5); backdrop-filter: blur(3px); z-index:65535;"
    static DIALOG_STYLE = "position:fixed; top:50%; left:50%; background-color:#FFF; border-radius:10px; z-index:65536; padding: 40px;"
    static CLOSE_STYLE = "font-size:24pt; position:absolute; top:0px; right:15px;"

    public static popUp(width: number, height: number, contentElem: HTMLElement) {
        if (document.getElementById(SCModal.BACKGROUND_ID) != null) {
            return;
        }

        // background
        let backgroundElem = document.createElement("div");
        backgroundElem.id = SCModal.BACKGROUND_ID;
        backgroundElem.setAttribute('style', SCModal.BACKGROUND_STYLE);

        // close button of dialog
        let closeElem = document.createElement("a");
        closeElem.href = "#";
        closeElem.onclick = SCModal.close;
        closeElem.setAttribute('style', SCModal.CLOSE_STYLE);
        closeElem.innerText = "×";

        // dialog
        let dialogStyle = SCModal.DIALOG_STYLE.concat("width:", width.toString(), "px; height:", height.toString(), "px; margin-top:", (height / -2).toString(), "px; margin-left:", (width / -2).toString(), "px;");
        let dialogElem = document.createElement("div");
        dialogElem.setAttribute('style', dialogStyle);
        dialogElem.appendChild(closeElem);

        if (contentElem != null) {
            dialogElem.appendChild(contentElem)
        }

        backgroundElem.appendChild(dialogElem);
        document.body.appendChild(backgroundElem)
    }

    public static close() {
        let backgroundElem = document.getElementById(SCModal.BACKGROUND_ID);
        if (backgroundElem != null) {
            document.body.removeChild(backgroundElem);
        }
    }
}
