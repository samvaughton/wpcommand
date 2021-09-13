<script>

    import {Router, Link} from "svelte-routing";
    import Enabled from "../components/Enabled.svelte";
    import {AuthEnum, hasAccess} from "../store/user";
    import BlueprintObjectType from "../components/BlueprintObjectType.svelte";
    import DeleteModal from "../components/DeleteModal.svelte";
    import Loading from "../components/Loading.svelte";
    import BlueprintObjectCreateUpdateModal from "../components/form/BlueprintObjectCreateUpdateModal.svelte";
    import Active from "../components/Active.svelte";
    import AccountCreateUpdateModal from "../components/form/AccountCreateUpdateModal.svelte";
    import UserCreateUpdateModal from "../components/form/UserCreateUpdateModal.svelte";

    /*
     * Fetch site details
     */

    export let accUuid;
    export let userUuid;
    export let isUserModalOpen;

    let item = null;

    const fetchData = function () {
        fetch("/api/account/" + accUuid + "/user/" + userUuid).then(resp => resp.json()).then(data => {
            item = data;
        });
    };

    fetchData();

</script>


<Router>
    {#if item}
        <div class="row">
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">User {item.Email} - <small class="text-muted">for account</small> {item.Account}</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">
                        {#await hasAccess(AuthEnum.ObjectUser, AuthEnum.ActionWriteSpecial)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isUserModalOpen = !isUserModalOpen} class="btn btn-primary">Update</button>
                            <UserCreateUpdateModal bind:isOpen={isUserModalOpen} accountUuid={accUuid} item={item} fetchData={fetchData} formType={"UPDATE"} />
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
                        <th>First Name</th>
                        <td>{item.FirstName}</td>
                    </tr>
                    <tr>
                        <th>Last Name</th>
                        <td>{item.LastName}</td>
                    </tr>
                    <tr>
                        <th>Email</th>
                        <td>{item.LastName}</td>
                    </tr>
                    <tr>
                        <th>Roles</th>
                        <td>{item.Roles}</td>
                    </tr>
                    <tr>
                        <th>Created</th>
                        <td>{item.CreatedAt}</td>
                    </tr>
                    </tbody>
                </table>

                <div class="d-flex bd-highlight">
                    <div class="ms-auto p-2 bd-highlight">

                    </div>
                </div>
            </div>
        </div>
    {:else}
        <div class="spinner-border m-5" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    {/if}

</Router>