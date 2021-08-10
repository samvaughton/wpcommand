<script>
    import {Link} from "svelte-routing";

    function getProps({location, href, isPartiallyCurrent, isCurrent}) {
        const isActive = href === "/" ? isCurrent : isPartiallyCurrent || isCurrent;

        if (isActive) {
            return {class: "nav-link px-3 text-warning"};
        }

        return {class: "nav-link px-3 text-white"};
    }
</script>

<header class="p-3 bg-dark text-white mb-4" style="border-bottom: 2px solid #36393e; background-color: #2d3034 !important;">
    <div class="container">
        <div class="d-flex flex-wrap align-items-center justify-content-center justify-content-lg-start">
            <ul class="nav col-12 col-lg-auto me-lg-auto mb-2 justify-content-center mb-md-0">
                <li>
                    <Link to="/" getProps="{getProps}">Sites</Link>
                </li>
                <li>
                    <Link to="/config" getProps="{getProps}">Config</Link>
                </li>
            </ul>

            <div class="text-end">
                <a href="/" class="d-flex align-items-center mb-2 mb-lg-0 text-white text-decoration-none">
                    WPCOMMAND
                </a>
            </div>
        </div>
    </div>
</header>

