<script>

    import {Router, Link, Route} from "svelte-routing";
    import {getLogs, logStore} from '../store/logs.js';
    import CommandJobStatus from "../components/CommandJobStatus.svelte";

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
                            <CommandJobStatus value={item.Status} />
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