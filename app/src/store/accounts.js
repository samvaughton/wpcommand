import {writable} from 'svelte/store';

export const accountStore = writable([])

const apiURL = "/api/account";

export function getAccounts() {
    return fetch(apiURL, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }).then(response => response.json()).then(json => {
        accountStore.set(json);

        return json;
    });
}