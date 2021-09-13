<script>

    import {Router, Link, Route} from "svelte-routing";
    import {getAccounts, accountStore} from '../store/accounts.js';
    import Enabled from "../components/Enabled.svelte";
    import {AuthEnum, hasAccess} from "../store/user";
    import Loading from "../components/Loading.svelte";
    import AccountCreateUpdateModal from "../components/form/AccountCreateUpdateModal.svelte";

    getAccounts();

    let isAccountModalOpen = false;

</script>

<Router>
    <div class="row">
        <div class="col-12">
            <div class="d-flex bd-highlight mb-3">
                <div class="p-2 bd-highlight">
                    <h1>Accounts</h1>
                    <p class="lead">Overview of all accounts.</p>
                </div>
                <div class="ms-auto p-2 bd-highlight">
                    {#await hasAccess(AuthEnum.ObjectAccount, AuthEnum.ActionWriteSpecial)}
                        <Loading />
                    {:then result}
                        <button on:click={() => isAccountModalOpen = !isAccountModalOpen} class="btn btn-primary">Create</button>
                        <AccountCreateUpdateModal bind:isOpen={isAccountModalOpen} fetchData={getAccounts} formType={"CREATE"} />
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
                    <th scope="col">Name</th>
                    <th scope="col">Status</th>
                    <th scope="col">Created</th>
                    <th scope="col"></th>
                </tr>
                </thead>
                <tbody>

                {#each $accountStore as item, index}
                    <tr>
                        <td>{item.Name}</td>
                        <td><Enabled value={item.Enabled} /></td>
                        <td>{item.CreatedAt}</td>
                        <td>
                            <Link to="/accounts/{item.Uuid}">Details</Link>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </div>

</Router>