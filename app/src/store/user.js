import {writable} from 'svelte/store';
import {register, TYPE_REQUEST, TYPE_RESPONSE} from "../fetch/fetchm";

const apiURL = "/auth";
const storedUser = JSON.parse(localStorage.getItem("user"));
let accessCache = JSON.parse(localStorage.getItem("accessCache"))

const setAndSaveAccessCache = function (key, value) {
    accessCache[key] = value;
    localStorage.setItem("accessCache", JSON.stringify(accessCache));
}

export let enforcer = null;
export const userStore = writable(storedUser);

export const hasAccess = async function (obj, action) {
    // We need to cache these results into application storage
    if (accessCache === undefined || accessCache === null) {
        accessCache = {};
    }

    // only start caching if storedUser is not null...
    if (storedUser !== null) {
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
    }
};

userStore.subscribe(value => {
    localStorage.setItem("user", JSON.stringify(value))

    // configure fetchm
    if (value !== null) {
        register(TYPE_REQUEST, 'auth', (url, reqOpts) => {
            reqOpts.headers = {...(reqOpts.headers ?? {}), ...{"Token": value.Token}};
            reqOpts.headers = {...(reqOpts.headers ?? {}), ...{"Content-Type": "application/json"}};
        })

        register(TYPE_RESPONSE, 'session_expired', (promise) => {
            promise.then(resp => {
                if (resp.status === 401) {
                    logout();
                    window.location = '/login';
                }
            });
        });
    }
});

export function logout() {
    // clear access cache
    localStorage.setItem("accessCache", JSON.stringify({}));
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

export const AuthEnum = {
    ObjectSite: "site",

    ObjectCommand: "command",
    ObjectCommandJob: "command_job",
    ObjectCommandJobEvent: "command_job_event",

    ObjectBlueprint: "blueprint",
    ObjectBlueprintObject: "blueprint_object",

    ObjectUser: "user",
    ObjectAccount: "account",
    ObjectConfig: "config",

    ObjectWordpressUser: "wp_user",

    ActionRead: "read",
    ActionReadSpecial: "read_special",
    ActionWrite: "write",
    ActionWriteSpecial: "write_special",
    ActionDelete: "delete",
    ActionRun: "run",      // things like deploy site etc
    ActionRunSpecial: "run_special", // things like setting up plugins/themes etc
    ActionConfigure: "configure",
};
