<script>

    import { Router, Link, Route } from "svelte-routing";
    import { fetchConfig, configStore, sitesStore } from '../store/config.js';

    let k8LabelSelector = '';
    configStore.subscribe((data) => {
        if (data['kubernetes']) {
            k8LabelSelector = data.kubernetes.labelSelector;
        }
    });

    fetchConfig();

</script>

<Router>
    <div class="row">
        <div class="col-12">
            <h1 class="">Sites</h1>
            <p class="lead ">Currently managed sites.</p>
            <p class="">Base selector: <code>{k8LabelSelector}</code></p>
        </div>
    </div>
    <div class="row">
        <div class="col-12">
            <table class="table table-borderless table-striped">
                <thead>
                <tr>
                    <th scope="col">#</th>
                    <th scope="col">Name</th>
                    <th scope="col">Enabled</th>
                    <th scope="col">Namespace</th>
                    <th scope="col">App Selector</th>
                    <th scope="col"></th>
                </tr>
                </thead>
                <tbody>
                {#each $sitesStore as site, index}
                <tr>
                    <th scope="row">{index + 1}</th>
                    <td>{site.name}</td>
                    <td>{#if site.enabled}<span class="badge bg-success">Enabled</span>{:else}<span class="badge bg-danger">Disabled</span>{/if}</td>
                    <td>{site.namespace}</td>
                    <td><code>{site.labelSelector}</code></td>
                    <td><Link to="/sites/{site.name}">Details</Link></td>
                </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </div>

</Router>