<script>
    import {Link} from "svelte-routing";

    import { userStore, hasAccess, logout } from '../store/user'

    function doLogout() {
        logout();
        window.location = '/login';
    }

    function getProps({location, href, isPartiallyCurrent, isCurrent}) {
        const isActive = href === "/" ? isCurrent : isPartiallyCurrent || isCurrent;

        if (isActive) {
            return {class: "nav-link px-3 text-warning"};
        }

        return {class: "text-white nav-link px-3 "};
    }

    function getButtonProps({location, href, isPartiallyCurrent, isCurrent}) {
        const isActive = href === "/" ? isCurrent : isPartiallyCurrent || isCurrent;

        if (isActive) {
            return {class: "btn btn-warning me-2"};
        }

        return {class: "btn btn-outline-light me-2"};
    }
</script>

<header class="p-3 mb-4 bg-dark text-white">
    <div class="container">
        <div class="d-flex flex-wrap align-items-center justify-content-center justify-content-lg-start">
            <a href="/" class="d-flex align-items-center mb-2 mb-lg-0 text-decoration-none text-white me-3">
                WP_CMD
            </a>

            <ul class="nav col-12 col-lg-auto me-lg-auto mb-2 justify-content-center mb-md-0">
                {#if $userStore}
                    {#await hasAccess("site", "read")}
                    {:then result}
                        <li>
                            <Link to="/" getProps="{getProps}">Sites</Link>
                        </li>
                    {/await}

                    {#await hasAccess("command_job", "read")}
                    {:then result}
                        <li>
                            <Link to="/logs" getProps="{getProps}">Command Log</Link>
                        </li>
                    {/await}

                    {#await hasAccess("blueprint", "read")}
                    {:then result}
                        <li>
                            <Link to="/blueprints" getProps="{getProps}">Blueprint Sets</Link>
                        </li>
                    {/await}

                    {#await hasAccess("config", "read")}
                    {:then result}
                        <li>
                            <Link to="/config" getProps="{getProps}">Config</Link>
                        </li>
                    {/await}
                {/if}
            </ul>

            {#if $userStore}
                <div class="text-end">

                    <ul class="nav col-12 col-lg-auto me-lg-auto mb-2 justify-content-center mb-md-0">
                        <li class="text-white me-2">
                            {$userStore.Email} - {$userStore.Account.Name}
                        </li>
                        <li class="ms-2 me-3">
                            |
                        </li>
                        <li>
                            <a href="#" on:click={doLogout} class="d-flex align-items-center mb-2 mb-lg-0 text-decoration-none text-white">
                                Logout
                            </a>
                        </li>
                    </ul>

                </div>
            {:else}
                <div class="text-end">
                    <Link to="/login" getProps="{getButtonProps}">Login</Link>
                </div>
            {/if}

        </div>
    </div>
</header>

