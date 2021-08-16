<script>

    import { Router } from "svelte-routing";
    import {getSites} from "../store/site";
    import { Modal, ModalHeader, ModalBody, ModalFooter } from 'sveltestrap';

    /*
     * Fetch site details
     */

    export let key;
    let site = null;
    let runnableCommands = [];

    fetch("/api/site/" + key).then(resp => resp.json()).then(data => {
        site = data;

        fetch("/api/site/" + key + "/command").then(resp => resp.json()).then(data => {
            runnableCommands = data;
        })
    });

    /*
     * Run command modal
     */

    let isOpen = false;
    let loading = false;
    let warningMessage = "";
    let mCommandId = 0;

    const toggle = () => (isOpen = !isOpen);
    const onClose = function() {
        warningMessage = "";
        loading = false;
        mCommandId = 0;
    };

    getSites();

    let submitModal = function() {
        loading = true;
        warningMessage = "";
        fetch("/api/command/job", {
            method: "POST",
            body: JSON.stringify({
                CommandId: mCommandId,
                Selector: site.Key,
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
                    //window.location = "/commands/" + data.CommandId + "/jobs/" + data.Id
                    isOpen = false;
                });
            }
        });
    };


</script>

<Modal isOpen={isOpen} {toggle} on:close={onClose}>
    <form on:submit|preventDefault={submitModal}>
        <ModalHeader>Run Command</ModalHeader>
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
                    <label for="command" class="form-label">Command</label>
                    <select bind:value={mCommandId} required id="command" class="form-control" aria-describedby="commandHelp">
                        <option>Select command</option>
                        {#each runnableCommands as command}
                            <option value={command.Id}>
                                {command.Description}
                            </option>
                        {/each}
                    </select>
                    <div id="commandHelp" class="form-text">
                        The command to run.
                    </div>
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                Run Command
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>

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
                        <button type="button" class="btn btn-primary" on:click={toggle}>Run Command</button>
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