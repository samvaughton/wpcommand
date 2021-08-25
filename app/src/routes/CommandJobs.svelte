<script>

    import {Router, Link, Route} from "svelte-routing";
    import {getLogs, logStore} from '../store/logs.js';

    getLogs();

</script>

<Router>
    <div class="row">
        <div class="col-12">
            <div class="d-flex bd-highlight mb-3">
                <div class="p-2 bd-highlight">
                    <h1>Command Logs</h1>
                    <p class="lead">Overview of all commands run.</p>
                </div>
                <div class="ms-auto p-2 bd-highlight">

                </div>
            </div>
        </div>
    </div>
    <div class="row">
        <div class="col-12">
            <table class="table table-borderless table-striped">
                <thead>
                <tr>
                    <th scope="col">Created</th>
                    <th scope="col">Status</th>
                    <th scope="col">Command</th>
                    <th scope="col">Site</th>
                    <th scope="col">Run By</th>
                    <th scope="col"></th>
                </tr>
                </thead>
                <tbody>
                {#each $logStore as item, index}
                    <tr>
                        <td>{item.CreatedAt}</td>
                        <td>
                            {#if item.Status === "CREATED"}
                                <span class="badge bg-secondary">Created</span>
                            {:else if item.Status === "PENDING"}
                                <span class="badge bg-info">Pending</span>
                            {:else if item.Status === "RUNNING"}
                                <span class="badge bg-primary">Running</span>
                            {:else if item.Status === "SUCCESS"}
                                <span class="badge bg-success">Success</span>
                            {:else if item.Status === "FAILURE"}
                                <span class="badge bg-danger">Failure</span>
                            {/if}
                        </td>
                        <td>{item.Command.Description}</td>
                        <td>{item.Site.Description}</td>
                        <td>{item.RunByUser ? item.RunByUser.Email : '-'}</td>
                        <td>
                            <Link to="/logs/{item.Uuid}">Details</Link>
                        </td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </div>

</Router>