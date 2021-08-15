<script>

    import { Router } from "svelte-routing";

    export let key;

    let site = null;

    fetch("/api/site/" + key).then(resp => resp.json()).then(data => {
        site = data;
    });

</script>

<Router>
    {#if site}
    <div class="row">
        <div class="col-12">
            <div class="d-flex bd-highlight mb-3">
                <div class="p-2 bd-highlight">
                    <h1 class="float-start">{site.Description}</h1>
                </div>
                <div class="ms-auto p-2 bd-highlight">
                    <div class="btn-group" role="group" aria-label="Site Actions">
                        <button type="button" class="btn btn-outline-primary">Build & Deploy</button>
                        <button type="button" class="btn btn-outline-primary">Sync Plugins</button>
                    </div>
                </div>
            </div>

        </div>
    </div>
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
                        <td>{#if site.Enabled}<span class="badge bg-success">Enabled</span>{:else}<span class="badge bg-danger">Disabled</span>{/if}</td>
                    </tr>
                    <tr>
                        <th>Namespace</th>
                        <td>{site.Namespace}</td>
                    </tr>
                    <tr>
                        <th>Label Selector</th>
                        <td>{site.LabelSelector}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    {:else}
        <div class="spinner-border m-5" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    {/if}

</Router>