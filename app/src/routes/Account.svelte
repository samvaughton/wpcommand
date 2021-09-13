<script>

    import {Router, Link} from "svelte-routing";
    import Enabled from "../components/Enabled.svelte";
    import {AuthEnum, hasAccess} from "../store/user";
    import Loading from "../components/Loading.svelte";
    import AccountCreateUpdateModal from "../components/form/AccountCreateUpdateModal.svelte";
    import UserCreateUpdateModal from "../components/form/UserCreateUpdateModal.svelte";

    /*
     * Fetch site details
     */

    export let uuid;

    let item = null;

    let users = [];
    let isUserModalOpen = false;
    let isAccountModalOpen = false;

    const fetchData = function () {
        fetch("/api/account/" + uuid).then(resp => resp.json()).then(data => {
            item = data;
        });

        fetch("/api/account/" + uuid + "/user").then(resp => resp.json()).then(data => {
            users = data;
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
                        <h3 class="float-start">Account #{item.Id}</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">
                        {#await hasAccess(AuthEnum.ObjectAccount, AuthEnum.ActionWriteSpecial)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isAccountModalOpen = !isAccountModalOpen} class="btn btn-info">Update</button>
                            <AccountCreateUpdateModal bind:isOpen={isAccountModalOpen} item={item} fetchData={fetchData} formType={"UPDATE"} />
                        {/await}
                        {#await hasAccess(AuthEnum.ObjectUser, AuthEnum.ActionWriteSpecial)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isUserModalOpen = !isUserModalOpen} class="btn btn-primary">Add User</button>
                            <UserCreateUpdateModal bind:isOpen={isUserModalOpen} accountUuid={item.Uuid} fetchData={fetchData} formType={"CREATE"} />
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
            </div>
            {#await hasAccess(AuthEnum.ObjectUser, AuthEnum.ActionRead)}
                <Loading />
            {:then result}
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Users</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">

                    </div>
                </div>
                <div class="d-flex bd-highlight">
                    <p class="p-2">Users associated with this account.</p>
                </div>
                <table class="table table-borderless table-striped">
                    <thead>
                    <tr>
                        <th scope="col">First Name</th>
                        <th scope="col">Last Name</th>
                        <th scope="col">Email</th>
                        <th scope="col">Roles</th>
                        <th scope="col"></th>
                    </tr>
                    </thead>
                    <tbody>
                        {#each users as oItem, index}
                            <tr>
                                <td>{oItem.FirstName}</td>
                                <td>{oItem.LastName}</td>
                                <td>{oItem.Email}</td>
                                <td>{oItem.Roles}</td>
                                <td>
                                    <Link to="/accounts/{uuid}/users/{oItem.Uuid}">Details</Link>
                                </td>
                            </tr>
                        {:else}
                            <tr><td colspan="5">No users found</td></tr>
                        {/each}
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