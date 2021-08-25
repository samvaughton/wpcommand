import {writable} from 'svelte/store';
import {register, TYPE_REQUEST} from "../fetch/fetchm";

const apiURL = "/auth";
const storedUser = JSON.parse(localStorage.getItem("user"));
let accessCache = JSON.parse(localStorage.getItem("accessCache"))

const setAndSaveAccessCache = function(key, value) {
    accessCache[key] = value;
    localStorage.setItem("accessCache", JSON.stringify(accessCache));
}

export let enforcer = null;
export const userStore = writable(storedUser);

export const hasAccess = async function(obj, action) {
    // We need to cache these results into application storage
    if (accessCache === undefined || accessCache === null) {
        accessCache = {};
    }

    const key = obj + "," + action;
    return new Promise((resolve, reject) => {
        if (accessCache[key] === undefined) {
            fetch("/api/access?params=" + key, {method: "POST"}).then(resp => {
                if (resp.status === 200) {
                    setAndSaveAccessCache(key, true);
                    resolve();
                } else {
                    setAndSaveAccessCache(key, false);
                    reject();
                }
            });
        } else {
            if (accessCache[key]) {
                resolve();
            } else {
                reject();
            }
        }
    });
};

userStore.subscribe(value => {
    localStorage.setItem("user", JSON.stringify(value))

    // configure fetchm
    if (value !== null) {
        register(TYPE_REQUEST, 'auth', (url, reqOpts) => {
            reqOpts.headers = {...(reqOpts.headers ?? {}), ...{"Token": value.Token}};
            reqOpts.headers = {...(reqOpts.headers ?? {}), ...{"Content-Type": "application/json"}};
        })
    }
});

export function logout() {
    userStore.set(null);
}

export function authenticate(account, email, password) {
    return fetch(apiURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            Account: account,
            Email: email,
            Password: password
        })
    });
}