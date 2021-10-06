<script>

    import {Router, Link} from "svelte-routing";
    import {hasAccess, AuthEnum} from "../store/user";
    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import Enabled from "../components/Enabled.svelte";
    import Loading from "../components/Loading.svelte";
    import SiteUpdateModal from "../components/form/SiteUpdateModal.svelte";
    import CommandCreateUpdateModal from "../components/form/CommandCreateUpdateModal.svelte";
    import WpUserCreateUpdateModal from "../components/form/WpUserCreateUpdateModal.svelte";

    /*
     * Fetch site details
     */

    export let key;
    let item = null;
    let runnableCommands = [];
    let blueprintSets = [];
    let itemSpecificCommands = [];
    let wpUsers = [];

    function fetchData() {
        fetch("/api/site/" + key).then(resp => resp.json()).then(data => {
            item = data;

            fetch("/api/site/" + key + "/command?type=runnable").then(resp => resp.json()).then(data => {
                runnableCommands = data;
            })

            fetch("/api/site/" + key + "/command?type=attached").then(resp => resp.json()).then(data => {
                itemSpecificCommands = data;

                itemSpecificCommands.forEach(cmd => {
                    isUpdateCommandModalOpen[cmd.Uuid] = false;
                });
            })

            hasAccess(AuthEnum.ObjectBlueprint, AuthEnum.ActionRead).then(() => {
                fetch("/api/site/" + key + "/blueprint").then(resp => resp.json()).then(data => {
                    blueprintSets = data;
                })
            });

            hasAccess(AuthEnum.ObjectWordpressUser, AuthEnum.ActionRead).then(() => {
                fetch("/api/site/" + item.Uuid + "/wp/user").then(resp => resp.json()).then(data => {
                    wpUsers = data;
                })
            });
        });
    }

    /*
     * Site update modal
     */
    let isUpdateSiteModalOpen = false;

    /*
     * Command update modal map
     */
    let isUpdateCommandModalOpen = {};

    /*
    * Command create modal
    */
    let isAddCommandModalOpen = false;

    /*
    * WpUser update modal map
    */
    let isUpdateWpUserModalOpen = {};

    /*
     * WpUser create modal
     */
    let isAddWpUserModalOpen = false;

    /*
     * Run command modal
     */

    let isOpen = false;
    let loading = false;
    let warningMessage = "";
    let mCommandId = 0;

    const toggle = () => (isOpen = !isOpen);
    const onClose = function () {
        warningMessage = "";
        loading = false;
        mCommandId = 0;
    };

    let submitModal = function () {
        loading = true;
        warningMessage = "";
        fetch("/api/command/job", {
            method: "POST",
            body: JSON.stringify({
                CommandId: mCommandId,
                Selector: item.Key,
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
                    window.location = "/logs/" + data.Jobs[0].Uuid
                    isOpen = false;
                });
            }
        });
    };

    fetchData();

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
    {#if item}

        <div class="row">
            <div class="col-12">
                <div class="d-flex bd-highlight mb-3">
                    <div class="p-2 bd-highlight">
                        <h1 class="float-start">{item.Description}</h1>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">
                        <div class="btn-group" role="group" aria-label="Site Actions">
                            <button type="button" class="btn btn-primary" on:click={toggle}>Run Command</button>
                        </div>
                        {#await hasAccess(AuthEnum.ObjectSite, AuthEnum.ActionWrite)}
                            <Loading />
                        {:then result}
                            <button on:click={() => isUpdateSiteModalOpen = !isUpdateSiteModalOpen} class="btn btn-info">Edit</button>
                            <SiteUpdateModal bind:isOpen={isUpdateSiteModalOpen} bind:item={item} fetchData={fetchData} formType={"UPDATE"} />
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
                        <th scope="col"></th>
                        <th scope="col"></th>
                    </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <th>Status</th>
                            <td><Enabled value="{item.Enabled}"/></td>
                        </tr>
                        <tr>
                            <th>Namespace</th>
                            <td>{item.Namespace}</td>
                        </tr>
                        <tr>
                            <th>Label Selector</th>
                            <td>{item.LabelSelector}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        {#await hasAccess(AuthEnum.ObjectBlueprint, AuthEnum.ActionRead)}
        {:then result}
            <div class="row mt-5">
                <div class="col-12">
                    <div class="d-flex bd-highlight">
                        <div class="p-2 bd-highlight">
                            <h3 class="float-start">Blueprints</h3>
                        </div>
                        <div class="ms-auto p-2 bd-highlight">

                        </div>
                    </div>
                    <table class="table table-borderless table-striped">
                        <thead>
                        <tr>
                            <th scope="col">Blueprint Set</th>
                            <th scope="col">Status</th>
                            <th scope="col"></th>
                        </tr>
                        </thead>
                        <tbody>
                        {#each blueprintSets as item, index}
                        <tr>
                            <td>{item.Name}</td>
                            <td><Enabled value={item.Enabled} /></td>
                            <td><Link to="/blueprints/{item.Uuid}">Details</Link></td>
                        </tr>
                        {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        {/await}

        {#await hasAccess(AuthEnum.ObjectCommand, AuthEnum.ActionRead)}
        {:then result}
            <div class="row mt-5">
                <div class="col-12">
                    <div class="d-flex bd-highlight">
                        <div class="p-2 bd-highlight">
                            <h3 class="float-start">{item.Description} specific commands</h3>
                        </div>
                        <div class="ms-auto p-2 bd-highlight">
                            {#await hasAccess(AuthEnum.ObjectCommand, AuthEnum.ActionWrite)}
                                <Loading />
                            {:then result}
                                <button on:click={() => isAddCommandModalOpen = !isAddCommandModalOpen} class="btn btn-sm btn-primary">Create Command</button>
                                <CommandCreateUpdateModal bind:isOpen={isAddCommandModalOpen} bind:site={item} fetchData={fetchData} formType={"CREATE"} />
                            {/await}
                        </div>
                    </div>
                    <table class="table table-borderless table-striped">
                        <thead>
                        <tr>
                            <th scope="col">Description</th>
                            <th scope="col">Type</th>
                            <th scope="col"></th>
                        </tr>
                        </thead>
                        <tbody>
                        {#each itemSpecificCommands as cmd, index}
                            <tr>
                                <td>{cmd.Description}</td>
                                <td>{cmd.Type}</td>
                                <td>
                                    {#await hasAccess(AuthEnum.ObjectCommand, AuthEnum.ActionWrite)}
                                        <Loading />
                                    {:then result}
                                        <a href="javascript:;" on:click={() => isUpdateCommandModalOpen[cmd.Uuid] = !isUpdateCommandModalOpen[cmd.Uuid]} style="cursor: pointer;">Update</a>
                                        <CommandCreateUpdateModal bind:isOpen={isUpdateCommandModalOpen[cmd.Uuid]} bind:site={item} bind:item={cmd} fetchData={fetchData} formType={"UPDATE"} />
                                    {/await}
                                </td>
                            </tr>
                        {:else}
                            <tr><td colspan="3">No commands found</td></tr>
                        {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        {/await}

        {#await hasAccess(AuthEnum.ObjectWordpressUser, AuthEnum.ActionRead)}
        {:then result}
            <div class="row mt-5">
                <div class="col-12">
                    <div class="d-flex bd-highlight">
                        <div class="p-2 bd-highlight">
                            <h3 class="float-start">Wordpress Users</h3>
                        </div>
                        <div class="ms-auto p-2 bd-highlight">
                            {#await hasAccess(AuthEnum.ObjectWordpressUser, AuthEnum.ActionWrite)}
                                <Loading />
                            {:then result}
                                <button on:click={() => isAddWpUserModalOpen = !isAddWpUserModalOpen} class="btn btn-sm btn-primary">Create User</button>
                                <WpUserCreateUpdateModal bind:isOpen={isAddWpUserModalOpen} bind:site={item} fetchData={fetchData} formType={"CREATE"} />
                            {/await}
                        </div>
                    </div>
                    <table class="table table-borderless table-striped">
                        <thead>
                        <tr>
                            <th scope="col">Username</th>
                            <th scope="col">Email</th>
                            <th scope="col">Roles</th>
                            <th scope="col"></th>
                        </tr>
                        </thead>
                        <tbody>
                        {#each wpUsers as wpUser, index}
                            <tr>
                                <td>{wpUser['user_login']}</td>
                                <td>{wpUser['user_email']}</td>
                                <td>{wpUser['roles']}</td>
                                <td>
                                    {#await hasAccess(AuthEnum.ObjectWordpressUser, AuthEnum.ActionWrite)}
                                        <Loading />
                                    {:then result}
                                        <a href="javascript:;" on:click={() => isUpdateWpUserModalOpen[wpUser['ID']] = !isUpdateWpUserModalOpen[wpUser['ID']]} style="cursor: pointer;">Update</a>
                                        <WpUserCreateUpdateModal bind:isOpen={isUpdateWpUserModalOpen[wpUser['ID']]} bind:site={item} bind:item={wpUser} fetchData={fetchData} formType={"UPDATE"} />
                                    {/await}
                                </td>
                            </tr>
                        {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        {/await}


    {:else}
        <div class="spinner-border m-5" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    {/if}

</Router>