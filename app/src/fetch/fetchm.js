/*
 * Fetchm (fetch-em) is a wrapper for fetch() enables support for middleware such as authentication, headers etc
 */

const windowImpl = window.fetch;

let middlewares = [];

export const TYPE_REQUEST = 'REQUEST';
export const TYPE_RESPONSE = 'RESPONSE';

function preRequestFilter(item) {
   return item.type === TYPE_REQUEST;
}

function postResponseFilter(item) {
    return item.type === TYPE_RESPONSE;
}

/**
 * This will replace the "fetch" function in the browser with fetchm's - removing the need for an import in every file
 * that performs http requests
 */
export function takeover() {
    window.fetch = fetchm;
}

export function restore() {
    window.fetch = windowImpl;
}

export function register(type, key, fn) {
    middlewares.push({
        type: type,
        key: key,
        fn: fn
    });
}

export function deregister(key) {
    middlewares = middlewares.filter(item => item.key !== key);
}

export function fetchm(url, reqOpts) {
    middlewares.filter(preRequestFilter).forEach(item => {
        reqOpts = reqOpts ?? {};

        item.fn(url, reqOpts)
    });

    let prom = windowImpl(url, reqOpts);

    middlewares.filter(postResponseFilter).forEach(item => {
       item.fn(prom);
    });

    return prom;
}