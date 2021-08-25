import {derived, writable} from 'svelte/store';

export const logStore = writable([])

const apiURL = "/api/command/job";

export function getLogs() {
    return fetch(apiURL, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }).then(response => response.json()).then(json => {
        logStore.set(json);

        return json;
    });
}