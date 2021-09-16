<script>

    import {Router, Link} from "svelte-routing";
    import Active from "../components/Active.svelte";
    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import {AuthEnum, hasAccess} from "../store/user";
    import BlueprintObjectType from "../components/BlueprintObjectType.svelte";
    import DeleteModal from "../components/DeleteModal.svelte";
    import Loading from "../components/Loading.svelte";
    import BlueprintObjectCreateUpdateModal from "../components/form/BlueprintObjectCreateUpdateModal.svelte";
    import BlueprintObjectUpdateVersionModal from "../components/form/BlueprintObjectUpdateVersionModal.svelte";

    export let blueprintUuid;
    export let objectUuid;
    export let revisionId;

    console.log(`loaded bp ${blueprintUuid} obj ${objectUuid} rev ${revisionId}`);

    let item = null;
    let revisions = [];
    let isDeleteRevisionModalOpen = false;
    let isDeleteObjectModalOpen = false;
    let isObjectModalOpen = false;
    let isObjectVersionModalOpen = false;

    const fetchBlueprintObjectData = function () {
        fetch("/api/blueprint/" + blueprintUuid + "/object/" + objectUuid + "/revision/" + revisionId).then(resp => resp.json()).then(data => {
            item = data;
        });

        fetch("/api/blueprint/" + blueprintUuid + "/object/" + objectUuid + "/revision").then(resp => resp.json()).then(data => {
            revisions = data;
        });
    };

    fetchBlueprintObjectData();

    const viewRevision = function(uuid, revId) {
        window.location = `/blueprints/${item.BlueprintSet.Uuid}/object/${uuid}/revision/${revId}`
    }

</script>

<Router>
    {#if item}
        <div class="row">
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Blueprint Object #{item.Id} - <small class="text-muted">for set</small> {item.BlueprintSet.Name}</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">
                        {#await hasAccess(AuthEnum.ObjectBlueprintObject, AuthEnum.ActionWrite)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isObjectModalOpen = !isObjectModalOpen} class="btn btn-info">Edit</button>
                            <BlueprintObjectCreateUpdateModal bind:isOpen={isObjectModalOpen} bind:item={item} blueprint={item.BlueprintSet} fetchData={fetchBlueprintObjectData} formType={"UPDATE"} />
                        {/await}

                        {#await hasAccess(AuthEnum.ObjectBlueprintObject, AuthEnum.ActionWriteSpecial)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isObjectVersionModalOpen = !isObjectVersionModalOpen} class="btn btn-primary">Update Version</button>
                            <BlueprintObjectUpdateVersionModal bind:isOpen={isObjectVersionModalOpen} bind:item={item} blueprint={item.BlueprintSet} fetchData={fetchBlueprintObjectData} />
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
                        <th>Revision</th>
                        <td>{item.RevisionId}</td>
                    </tr>
                    <tr>
                        <th>Name</th>
                        <td>{item.Name}</td>
                    </tr>
                    <tr>
                        <th>Exact Name</th>
                        <td><code>{item.ExactName}</code></td>
                    </tr>
                    <tr>
                        <th>Version</th>
                        <td>{item.Version}</td>
                    </tr>
                    <tr>
                        <th>Active</th>
                        <td><Active value={item.Active} /></td>
                    </tr>
                    <tr>
                        <th>Created</th>
                        <td>{item.CreatedAt}</td>
                    </tr>
                    <tr>
                        <th>Set Order</th>
                        <td>{item.SetOrder}</td>
                    </tr>
                    <tr>
                        <th>URL</th>
                        <td><a href={item.OriginalObjectUrl}>Original Object URL</a> - <span class="text-muted">Depending on how you store plugins, this url may not reflect the version that is stored on this object.</span></td>
                    </tr>
                    </tbody>
                </table>

                <div class="d-flex bd-highlight">
                    <div class="ms-auto p-2 bd-highlight">
                        {#await hasAccess(AuthEnum.ObjectBlueprint, AuthEnum.ActionDelete)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isDeleteRevisionModalOpen = !isDeleteRevisionModalOpen} class="btn btn-sm btn-warning">Delete Revision</button>
                            <DeleteModal isOpen={isDeleteRevisionModalOpen} onClose="{() => isDeleteRevisionModalOpen = false}" name="Revision" endpoint={"/api/blueprint/" + item.BlueprintSet.Uuid + "/object/" + item.Uuid + "/revision/" + item.RevisionId} redirectTo={"/blueprints/" + item.BlueprintSet.Uuid} />
                            <button on:click={() => isDeleteObjectModalOpen = !isDeleteObjectModalOpen} class="btn btn-sm btn-danger">Delete Object</button>
                            <DeleteModal isOpen={isDeleteObjectModalOpen} onClose="{() => isDeleteObjectModalOpen = false}" name="Object" endpoint={"/api/blueprint/" + item.BlueprintSet.Uuid + "/object/" + item.Uuid} redirectTo={"/blueprints/" + item.BlueprintSet.Uuid} />
                        {/await}
                    </div>
                </div>
            </div>
            {#await hasAccess(AuthEnum.ObjectBlueprintObject, AuthEnum.ActionRead)}
            {:then result}
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">All Revisions</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">

                    </div>
                </div>
                <div class="d-flex bd-highlight">
                    <p class="p-2">All revisions including the one currently being viewed.</p>
                </div>
                <table class="table table-borderless table-striped">
                    <thead>
                    <tr>
                        <th scope="col">Rev. ID</th>
                        <th scope="col">Set Order</th>
                        <th scope="col">Type</th>
                        <th scope="col">Name</th>
                        <th scope="col">Version</th>
                        <th scope="col">Status</th>
                        <th scope="col"></th>
                    </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td colspan="7"><span class="badge bg-success">Newest</span></td>
                        </tr>
                        {#each revisions as oItem, index}
                            <tr>
                                <td>{oItem.RevisionId}</td>
                                <td>{oItem.SetOrder}</td>
                                <td><BlueprintObjectType value={oItem.Type} /></td>
                                <td>{oItem.Name}</td>
                                <td>{oItem.Version}</td>
                                <td>
                                    <Active value={oItem.Active} />
                                </td>
                                <td>
                                    {#if oItem.RevisionId === item.RevisionId}
                                        <span class="badge bg-info">Viewing</span>
                                    {:else}
                                        <a href="#" on:click={viewRevision(oItem.Uuid, oItem.RevisionId)}>Details</a>
                                    {/if}
                                </td>
                            </tr>
                        {/each}
                        <tr>
                            <td colspan="7"><span class="badge bg-danger">Oldest</span></td>
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