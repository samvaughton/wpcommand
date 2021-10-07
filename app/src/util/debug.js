import {clearCache} from "../store/user";

class DebugUtils {
    help() {
        return "Available commands: reloadCasbin(), clearAuthCache()";
    }

    clearAuthCache() {
        clearCache();
    }

    reloadCasbin() {
        return fetch("/api/casbin/reload").then(resp => resp.json()).then(data => {
            console.log(data);
        });
    }
}

export const Debug = new DebugUtils();