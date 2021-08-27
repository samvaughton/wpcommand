<script>

    import {Router, Link} from "svelte-routing";
    import Enabled from "../components/Enabled.svelte";
    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import {AuthEnum, hasAccess} from "../store/user";
    import BlueprintObjectType from "../components/BlueprintObjectType.svelte";
    import DeleteModal from "../components/DeleteModal.svelte";
    import Loading from "../components/Loading.svelte";

    /*
     * Fetch site details
     */

    export let blueprintUuid;
    export let objectUuid;
    export let revisionId;

    let item = null;
    let revisions = [];
    let isDeleteModalOpen = false;

    const fetchBlueprintData = function () {
        fetch("/api/blueprint/" + blueprintUuid + "/object/" + objectUuid + "/revision/" + revisionId).then(resp => resp.json()).then(data => {
            item = data;
        });

        fetch("/api/blueprint/" + blueprintUuid + "/object/" + objectUuid + "/" + revisionId + "/revision").then(resp => resp.json()).then(data => {
            revisions = data;
        });
    };

    let isOpen = false;
    let loading = false;
    let warningMessage = "";

    const pristineObj = {
        Type: {
            value: '',
            error: ''
        },
        Name: {
            value: '',
            error: ''
        },
        ExactName: {
            value: '',
            error: ''
        },
        Version: {
            value: '',
            error: ''
        },
        Url: {
            value: '',
            error: ''
        },
    };

    const newObj = function () {
        return JSON.parse(JSON.stringify(pristineObj));
    }

    const getValuesFromObj = function () {
        let values = {};

        for (let key in currentObj) {
            values[key] = currentObj[key].value;
        }

        return values;
    }

    let currentObj = newObj();

    const toggle = () => (isOpen = !isOpen);
    const onClose = function () {
        warningMessage = "";
        loading = false;
        currentObj = newObj();
    };

    let submitModal = function () {
        loading = true;
        warningMessage = "";
        fetch("/api/blueprint/" + item.Uuid + "/object", {
            method: "POST",
            body: JSON.stringify(getValuesFromObj())
        }).then(resp => {
            loading = false;

            if (resp.status !== 200) {
                resp.json().then(data => {
                    warningMessage = data.Message;
                });
            } else {
                fetchBlueprintData()
            }
        });
    };

    fetchBlueprintData();

</script>


<Modal isOpen={isOpen} {toggle} on:close={onClose}>
    <form on:submit|preventDefault={submitModal}>
        <ModalHeader>Add Object</ModalHeader>
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
            <div class="row mb-3">
                <div class="col-12">
                    <label for="type" class="form-label">Type</label>
                    <select bind:value={currentObj.Type.value} required id="type" class="form-control" aria-describedby="typeHelp">
                        <option>Select type</option>
                        <option value="PLUGIN">Plugin</option>
                        <option value="THEME">Theme</option>
                    </select>
                </div>
            </div>
            <div class="row mb-3">
                <div class="col-12">
                    <label for="name" class="form-label">Name</label>
                    <input bind:value={currentObj.Name.value} required type="text" id="name" class="form-control" aria-describedby="nameHelp">
                </div>
            </div>
            <div class="row mb-3">
                <div class="col-12">
                    <label for="exactName" class="form-label">Exact Name</label>
                    <input bind:value={currentObj.ExactName.value} required type="text" id="exactName" class="form-control" aria-describedby="exactNameHelp">
                    <div id="exactNameHelp" class="form-text">
                        What the theme or plugin calls itself, ie Advanced Custom Fields might be <code>acf</code>
                    </div>
                </div>
            </div>
            <div class="row mb-3">
                <div class="col-12">
                    <label for="version" class="form-label">Version</label>
                    <input bind:value={currentObj.Version.value} required type="text" id="version" class="form-control" aria-describedby="versionHelp">
                    <div id="versionHelp" class="form-text">
                        The version of the theme or plugin eg `3.1.2`
                    </div>
                </div>
            </div>
            <div class="row mb-3">
                <div class="col-12">
                    <label for="url" class="form-label">Zip File Url</label>
                    <input bind:value={currentObj.Url.value} required type="text" id="url" class="form-control" aria-describedby="urlHelp">
                    <div id="urlHelp" class="form-text">
                        URL of the zip file
                    </div>
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                Add Object
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>


<Router>
    {#if item}
        <div class="row">
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Blueprint #{item.Id}</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">
                        {#await hasAccess(AuthEnum.ObjectBlueprintObject, AuthEnum.ActionWrite)}
                            <Loading />
                        {:then result}
                            <button on:click={toggle} class="btn btn-primary">Add Object</button>
                        {/await}
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

                <div class="d-flex bd-highlight">
                    <div class="ms-auto p-2 bd-highlight">
                        {#await hasAccess(AuthEnum.ObjectBlueprint, AuthEnum.ActionDelete)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isDeleteModalOpen = !isDeleteModalOpen} class="btn btn-sm btn-danger">Delete Blueprint</button>
                            <DeleteModal isOpen={isDeleteModalOpen} onClose="{() => isDeleteModalOpen = false}" name="Blueprint" endpoint={"/api/blueprint/" + item.Uuid} redirectTo="/blueprints" />
                        {/await}
                    </div>
                </div>
            </div>
            {#await hasAccess(AuthEnum.ObjectBlueprintObject, AuthEnum.ActionRead)}
            {:then result}
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Objects</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">

                    </div>
                </div>
                <div class="d-flex bd-highlight">
                    <p class="p-2">Set order runs from lowest to highest.</p>
                </div>
                <table class="table table-borderless table-striped">
                    <thead>
                    <tr>
                        <th scope="col">Set Order</th>
                        <th scope="col">Type</th>
                        <th scope="col">Name</th>
                        <th scope="col">Version</th>
                        <th scope="col">Rev. ID</th>
                        <th scope="col">Enabled</th>
                        <th scope="col"></th>
                    </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td colspan="7"><span class="badge bg-success">First</span></td>
                        </tr>
                        {#each objects as oItem, index}
                            <tr>
                                <td>{oItem.SetOrder}</td>
                                <td><BlueprintObjectType value={oItem.Type} /></td>
                                <td>{oItem.Name}</td>
                                <td>{oItem.Version}</td>
                                <td>{oItem.RevisionId}</td>
                                <td><Enabled value={oItem.Enabled} /></td>
                                <td>
                                    <Link to="/blueprints/{item.Uuid}/object/{oItem.Uuid}/{oItem.RevisionId}">Details</Link>
                                </td>
                            </tr>
                        {/each}
                        <tr>
                            <td colspan="7"><span class="badge bg-danger">Last</span></td>
                        </tr>
                    </tbody>
                </table>
            </div>
            {/await}
        </div>
    {:else}
        <div class="spinner-border m-5" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    {/if}

</Router>