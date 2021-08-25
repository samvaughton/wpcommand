<script>

    import {Router} from "svelte-routing";
    import Enabled from "../components/Enabled.svelte";

    /*
     * Fetch site details
     */

    export let uuid;
    let item = null;
    let objects = [];

    fetch("/api/blueprint/" + uuid).then(resp => resp.json()).then(data => {
        item = data;
    });

</script>

<Router>
    {#if item}
        <div class="row">
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Blueprint #{item.Id}</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">

                    </div>
                </div>

                <table class="table table-borderless table-striped">
                    <thead>
                    <tr>
                        <th scope="col"></th>
                        <th scope="col"></th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr>
                        <th>UUID</th>
                        <td>{item.Uuid}</td>
                    </tr>
                    <tr>
                        <th>Name</th>
                        <td>{item.Name}</td>
                    </tr>
                    <tr>
                        <th>Enabled</th>
                        <td><Enabled value={item.Enabled} /></td>
                    </tr>
                    <tr>
                        <th>Created</th>
                        <td>{item.CreatedAt}</td>
                    </tr>
                    </tbody>
                </table>
            </div>
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Objects</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">

                    </div>
                </div>
                <table class="table table-borderless table-striped">
                    <thead>
                    <tr>
                        <th scope="col">Set Order</th>
                        <th scope="col">Revision ID</th>
                        <th scope="col">Type</th>
                        <th scope="col">Name</th>
                        <th scope="col">Version</th>
                        <th scope="col">Enabled</th>
                        <th scope="col"></th>
                    </tr>
                    </thead>
                    <tbody>
                    {#each objects as item, index}
                        <tr>
                            <td>object</td>
                            <td></td>
                        </tr>
                    {/each}
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