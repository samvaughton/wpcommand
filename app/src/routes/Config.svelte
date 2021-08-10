<script>
    import { Router, Link, Route } from "svelte-routing";

    import { fetchConfig, configStore, sitesStore } from '../store/config.js';

    let configJson = '';
    configStore.subscribe((data) => {
        configJson = JSON.stringify(data, null, 2)
    });

    fetchConfig();

</script>

<Router>
    <div class="row">
        <div class="col-12">
            <h1 class="text-white">Config</h1>
            <p class="lead text-white">WPCOMMAND's current config</p>
        </div>
    </div>
    <div class="row">
        <div class="col-12">
            <textarea readonly rows="25" style="width: 100%;" class="text-white bg-dark">{configJson}</textarea>
        </div>
    </div>
</Router>