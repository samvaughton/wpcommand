import {writable} from 'svelte/store';
import {register, TYPE_REQUEST} from "../fetch/fetchm";

const apiURL = "/auth";
const storedUser = JSON.parse(localStorage.getItem("user"));

export const userStore = writable(storedUser);

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
    }).then(response => response.json());
}