class DebugUtils {
    help() {
        return "Available commands: reloadCasbin()";
    }

    reloadCasbin() {
        return fetch("api/casbin/reload").then(resp => resp.json()).then(data => {
            console.log(data);
        });
    }
}

export const Debug = new DebugUtils();