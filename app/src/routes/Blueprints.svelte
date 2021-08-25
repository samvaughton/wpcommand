<script>

    import {Router, Link, Route} from "svelte-routing";
    import {getBlueprints, blueprintStore} from '../store/blueprints.js';
    import Enabled from "../components/Enabled.svelte";
    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import {AuthEnum, hasAccess} from "../store/user";

    getBlueprints();

    let isOpen = false;
    let loading = false;
    let warningMessage = "";
    let mName = "";

    const toggle = () => (isOpen = !isOpen);
    const onClose = function () {
        warningMessage = "";
        loading = false;
        mName = "";
    };

    let submitAddBlueprint = function () {
        loading = true;
        warningMessage = "";
        fetch("/api/blueprint", {
            method: "POST",
            body: JSON.stringify({
                Name: mName,
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
                    window.location = "/blueprints/" + data.Uuid;
                });
            }
        });
    };

</script>

<Modal isOpen={isOpen} {toggle} on:close={onClose}>
    <form on:submit|preventDefault={submitAddBlueprint}>
        <ModalHeader>Add Blueprint</ModalHeader>
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
                    <label for="name" class="form-label">Name</label>
                    <input bind:value={mName} required type="text" id="name" class="form-control" aria-describedby="nameHelp">
                    <div id="nameHelp" class="form-text">
                        Blueprint name
                    </div>
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                Add Blueprint
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
                    <h1>Blueprints</h1>
                    <p class="lead">Overview of all account blueprints.</p>
                </div>
                <div class="ms-auto p-2 bd-highlight">
                    {#await hasAccess(AuthEnum.ObjectBlueprint, AuthEnum.ActionWrite)}
                    {:then result}
                        <button on:click={toggle} class="btn btn-primary">Add Blueprint</button>
                    {/await}
                </div>
            </div>
        </div>
    </div>
    <div class="row">
        <div class="col-12">
            <table class="table table-borderless table-striped">
                <thead>
                <tr>
                    <th scope="col">Blueprint</th>
                    <th scope="col">Status</th>
                    <th scope="col">Created</th>
                    <th scope="col"></th>
                </tr>
                </thead>
                <tbody>
                {#each $blueprintStore as item, index}
                    <tr>
                        <td>{item.Name}</td>
                        <td><Enabled value={item.Enabled} /></td>
                        <td>{item.CreatedAt}</td>
                        <td>
                            <Link to="/blueprints/{item.Uuid}">Details</Link>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </div>

</Router>