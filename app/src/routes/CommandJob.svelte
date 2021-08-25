<script>

    import { Router } from "svelte-routing";

    /*
     * Fetch site details
     */

    export let uuid;
    let item = null;
    let events = [];

    let timer = null;
    let currentTimeout = 0;
    let maxTimeout = 300;

    let fetchEvents = () => {
        fetch("/api/command/job/" + uuid + "/events").then(resp => resp.json()).then(data => {
            events = data

            // also check if the last item is JOB_FINISHED, if so clear timer
            if (data.length > 0 && data[data.length - 1].Type === "JOB_FINISHED") {
                clearInterval(timer);
            }
        });

        currentTimeout++;

        if (currentTimeout > maxTimeout) {
            clearInterval(timer);
        }
    }

    fetch("/api/command/job/" + uuid).then(resp => resp.json()).then(data => {
        item = data;

        fetchEvents();
        timer = setInterval(fetchEvents, 1000); // every second
    });

</script>

<Router>
    {#if item}
        <div class="row">
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Log #{item.Id}</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">

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
                        <th>Command</th>
                        <td>{item.Command.Description}</td>
                    </tr>
                    <tr>
                        <th>Site</th>
                        <td>{item.Site.Description}</td>
                    </tr>
                    <tr>
                        <th>Created</th>
                        <td>{item.CreatedAt}</td>
                    </tr>
                    <tr>
                        <th>Status</th>
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
                    </tr>
                    <tr>
                        <th>Run By</th>
                        <td>{item.RunByUser ? item.RunByUser.Email : '-'}</td>
                    </tr>

                    </tbody>
                </table>
            </div>
            <div class="col-12">
                <div class="d-flex bd-highlight">
                    <div class="p-2 bd-highlight">
                        <h3 class="float-start">Events</h3>
                    </div>
                    <div class="ms-auto p-2 bd-highlight">

                    </div>
                </div>
                <table class="table table-borderless table-striped">
                    <thead>
                    <tr>
                        <th scope="col">Executed</th>
                        <th scope="col">Type</th>
                        <th scope="col">Status</th>
                        <th scope="col">Command</th>
                        <th scope="col"></th>
                    </tr>
                    </thead>
                    <tbody>
                    {#each events as item, index}
                        <tr>
                            <td>{item.ExecutedAt}</td>
                            <td>
                                {#if item.Type === "INFO"}
                                    <span class="badge bg-info">Info</span>
                                {:else if item.Type === "DATA"}
                                    <span class="badge bg-primary">Data</span>
                                {:else if item.Type === "JOB_STARTED"}
                                    <span class="badge bg-secondary">Job Started</span>
                                {:else if item.Type === "JOB_FINISHED"}
                                    <span class="badge bg-secondary">Job Finished</span>
                                {/if}
                            </td>
                            <td>
                                {#if item.Status === "SKIPPED"}
                                    <span class="badge bg-warning">Skipped</span>
                                {:else if item.Status === "SUCCESS"}
                                    <span class="badge bg-success">Success</span>
                                {:else if item.Status === "FAILURE"}
                                    <span class="badge bg-danger">Failure</span>
                                {/if}
                            </td>
                            <td><code>{item.Command}</code></td>
                            <td></td>
                        </tr>
                    {/each}
                    </tbody>
                </table>
            </div>
        </div>
    {:else}
        <div class="spinner-border m-5" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    {/if}

</Router>