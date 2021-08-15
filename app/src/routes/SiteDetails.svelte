<script>

    import { Router, Link, Route } from "svelte-routing";
    import {configStore, fetchConfig, sitesStore} from '../store/config.js';

    export let name;

    fetchConfig();

    let site = null;

    sitesStore.subscribe(sites => {
        for (let i = 0; i < sites.length; i++) {
            if (sites[i].name === name) {
                site = sites[i];
                break;
            }
        }
    });


</script>

<Router>
    <div class="row">
        <div class="col-12">
            <h1 class=" float-start">{name}</h1>
        </div>
    </div>
    {#if site}
    <div class="row">
        <div class="col-12">
            <table class="table table-borderless table-striped">
                <thead>
                <tr>
                    <th scope="col"></th>
                    <th scope="col"></th>
                </tr>
                </thead>
                <tbody>
                    <tr>
                        <th>Status</th>
                        <td>{#if site.enabled}<span class="badge bg-success">Enabled</span>{:else}<span class="badge bg-danger">Disabled</span>{/if}</td>
                    </tr>
                    <tr>
                        <th>Namespace</th>
                        <td>{site.namespace}</td>
                    </tr>
                    <tr>
                        <th>Label Selector</th>
                        <td>{site.labelSelector}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    {/if}

</Router>