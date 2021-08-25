import {derived, writable} from 'svelte/store';

export const blueprintStore = writable([])

const apiURL = "/api/blueprint";

export function getBlueprints() {
    return fetch(apiURL, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }).then(response => response.json()).then(json => {
        blueprintStore.set(json);

        return json;
    });
}