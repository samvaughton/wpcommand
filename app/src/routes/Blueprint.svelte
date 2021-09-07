<script>

    import {Router, Link} from "svelte-routing";
    import Enabled from "../components/Enabled.svelte";
    import {AuthEnum, hasAccess} from "../store/user";
    import BlueprintObjectType from "../components/BlueprintObjectType.svelte";
    import DeleteModal from "../components/DeleteModal.svelte";
    import Loading from "../components/Loading.svelte";
    import BlueprintObjectCreateUpdateModal from "../components/form/BlueprintObjectCreateUpdateModal.svelte";
    import Active from "../components/Active.svelte";

    /*
     * Fetch site details
     */

    export let uuid;

    let item = null;
    let objects = [];
    let isDeleteModalOpen = false;
    let isAddObjectModalOpen = false;

    const fetchBlueprintData = function () {
        fetch("/api/blueprint/" + uuid).then(resp => resp.json()).then(data => {
            item = data;
        });

        fetch("/api/blueprint/" + uuid + "/object").then(resp => resp.json()).then(data => {
            objects = data;
        });
    };

    fetchBlueprintData();

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
                        {#await hasAccess(AuthEnum.ObjectBlueprintObject, AuthEnum.ActionWrite)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isAddObjectModalOpen = !isAddObjectModalOpen} class="btn btn-primary">Add Object</button>
                            <BlueprintObjectCreateUpdateModal bind:isOpen={isAddObjectModalOpen} blueprint={item} fetchData={fetchBlueprintData} formType={"CREATE"} />
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
                            <DeleteModal bind:isOpen={isDeleteModalOpen} name="Blueprint" endpoint={"/api/blueprint/" + item.Uuid} redirectTo="/blueprints" />
                        {/await}
                    </div>
                </div>
            </div>
            {#await hasAccess(AuthEnum.ObjectBlueprintObject, AuthEnum.ActionRead)}
                <Loading />
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
                        <th scope="col"></th>
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
                                <td><Active value={oItem.Active} /></td>
                                <td>
                                    <Link to="/blueprints/{item.Uuid}/object/{oItem.Uuid}/revision/{oItem.RevisionId}">Details</Link>
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