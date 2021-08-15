import {derived, writable} from 'svelte/store';

export const siteStore = writable([])

const apiURL = "/api/site";

export function getSites() {
    return fetch(apiURL, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }).then(response => response.json()).then(json => {
        siteStore.set(json);

        return json;
    });
}