<script>

    import {Router} from "svelte-routing";
    import CommandJobStatus from "../components/CommandJobStatus.svelte";
    import CommandJobEventLogStatus from "../components/CommandJobEventLogStatus.svelte";
    import CommandJobEventLogType from "../components/CommandJobEventLogType.svelte";

    /*
     * Fetch site details
     */

    export let uuid;
    let item = null;
    let events = [];

    let timer = null;
    let currentTimeout = 0;
    let maxTimeout = 300 * 3; // 15min

    let fetchEvents = () => {
        fetch("/api/command/job/" + uuid + "/event").then(resp => resp.json()).then(data => {
            events = data

            // also check if the last item is JOB_FINISHED, if so clear timer
            if (data.length > 0 && data[data.length - 1].Type === "JOB_FINISHED" && timer !== null) {
                clearInterval(timer);
            }

            data.forEach(item => {
                if (item.Type === 'DATA') {

                }
            });
        });

        currentTimeout++;

        if (currentTimeout > maxTimeout && timer !== null) {
            clearInterval(timer);
        }
    }

    fetch("/api/command/job/" + uuid).then(resp => resp.json()).then(data => {
        item = data;

        fetchEvents();

        if (data.status !== "SUCCESS" && data.status !== "FAILURE") {
            timer = setInterval(fetchEvents, 5000); // every second
        }
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
                            <CommandJobStatus value={item.Status} />
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
                    {#each events as eItem, index}
                        <tr>
                            <td>{eItem.ExecutedAt}</td>
                            <td>
                                <CommandJobEventLogType value={eItem.Type} />
                            </td>
                            <td>
                                <CommandJobEventLogStatus value={eItem.Status} />
                            </td>
                            <td><code>{eItem.Command}</code></td>
                            <td>
                                <a target="_blank" href="/public/command/job/{item.Uuid}/event/{eItem.Uuid}">Details</a>
                                {#if eItem.MetaData !== "{}"}
                                <a target="_blank" href="/public/command/job/{item.Uuid}/event/{eItem.Uuid}?metadata=yes">Metadata</a>
                                {/if}
                            </td>
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