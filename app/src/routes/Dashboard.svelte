<script>

    import { Router, Link, Route } from "svelte-routing";
    import { getSites, siteStore } from '../store/site.js';
    import { Modal, ModalHeader, ModalBody, ModalFooter } from 'sveltestrap';

    let isOpen = false;
    let loading = false;
    let warningMessage = "";
    let mNamespace = "";
    let mLabelSelector = "";

    const toggle = () => (isOpen = !isOpen);
    const onClose = function() {
        warningMessage = "";
        loading = false;
        mNamespace = "";
        mLabelSelector = "";
    };

    getSites();

    let submitAddSite = function() {
        loading = true;
        warningMessage = "";
        fetch("/api/site", {
            method: "POST",
            body: JSON.stringify({
                Namespace: mNamespace,
                LabelSelector: mLabelSelector
            })
        }).then(resp => {
            loading = false;

            if (resp.status !== 200) {
                resp.json().then(data => {
                    warningMessage = data.Message;
                });
            } else {
                // redirect to website
                resp.json().then(data => {
                    window.location = "/sites/" + data.Key;
                });
            }
        });
    };

</script>

<Modal isOpen={isOpen} {toggle} on:close={onClose}>
    <form on:submit|preventDefault={submitAddSite}>
    <ModalHeader>Add Site</ModalHeader>
    <ModalBody>
        {#if warningMessage !== ""}
        <div class="row">
            <div class="col-12">
                <div class="alert alert-warning" role="alert">
                    {warningMessage}
                </div>
            </div>
        </div>
        {/if}
        <div class="row">
            <div class="col-12">
                <label for="namespace" class="form-label">Namespace</label>
                <input bind:value={mNamespace} required type="text" id="namespace" class="form-control" aria-describedby="namespaceHelp">
                <div id="namespaceHelp" class="form-text">
                    The kubernetes namespace that the site resides in.
                </div>
            </div>
        </div>
        <div class="row mt-2">
            <div class="col-12">
                <label for="labelSelector" class="form-label">Label Selector</label>
                <input bind:value={mLabelSelector} required type="text" id="labelSelector" class="form-control" aria-describedby="labelSelectorHelp">
                <div id="labelSelectorHelp" class="form-text">
                    The label associated with the pod. This is combined with the base selector defined in the config file to find the correct pod. Eg <code>app.kubernetes.io/instance=demowp</code>
                </div>
            </div>
        </div>
    </ModalBody>
    <ModalFooter>
        <button type="submit" class="btn btn-primary" disabled={loading}>
            {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
            Add Site
        </button>
        <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
    </ModalFooter>
    </form>
</Modal>

<Router>
    <div class="row">
        <div class="col-12">
            <div class="d-flex bd-highlight mb-3">
                <div class="p-2 bd-highlight">
                    <h1>Sites</h1>
                    <p class="lead">Currently managed sites.</p>
                </div>
                <div class="ms-auto p-2 bd-highlight">
                    <button on:click={toggle} class="btn btn-primary">Add Site</button>
                </div>
            </div>
        </div>
    </div>
    <div class="row">
        <div class="col-12">
            <table class="table table-borderless table-striped">
                <thead>
                <tr>
                    <th scope="col">#</th>
                    <th scope="col">Description</th>
                    <th scope="col">Enabled</th>
                    <th scope="col">Namespace</th>
                    <th scope="col">Instance Selector</th>
                    <th scope="col"></th>
                </tr>
                </thead>
                <tbody>
                {#each $siteStore as site, index}
                <tr>
                    <th scope="row">{index + 1}</th>
                    <td>{site.Description}</td>
                    <td>{#if site.Enabled}<span class="badge bg-success">Enabled</span>{:else}<span class="badge bg-danger">Disabled</span>{/if}</td>
                    <td>{site.Namespace}</td>
                    <td><code>{site.LabelSelector}</code></td>
                    <td><Link to="/sites/{site.Key}">Details</Link></td>
                </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </div>

</Router>