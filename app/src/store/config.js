import {derived, writable} from 'svelte/store';

export const configStore = writable({})

const apiURL = "/api/config";

export async function fetchConfig() {
    const response = await fetch(apiURL);
    let data = (await response.json());
    configStore.set(data)
}

export const sitesStore = derived(configStore, ($configData) => {
    if ($configData.wordpress && $configData.wordpress.sites) {
        return $configData.wordpress.sites.map(site => site);
    }

    return [];
});