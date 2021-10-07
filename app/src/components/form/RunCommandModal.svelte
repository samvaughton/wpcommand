<script>

    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import Loading from "../../components/Loading.svelte";

    let loading = false;
    let warningMessage = "";
    let mCommandId = 0;

    export let isOpen = false;
    export let onClose = null;
    export let runnableCommands = [];
    export let selector = "";

    const toggle = () => (isOpen = !isOpen);
    const _onClose = function () {
        loading = false;
        mCommandId = 0;

        if (typeof onClose === 'function') {
            onClose();
        }
    };

    let submitModal = function () {
        loading = true;
        warningMessage = "";
        fetch("/api/command/job", {
            method: "POST",
            body: JSON.stringify({
                CommandId: mCommandId,
                Selector: selector,
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

                    if (data['Jobs'].length === 1) {
                        window.location = "/logs/" + data.Jobs[0].Uuid

                    } else {
                        window.location = "/logs";
                    }

                    isOpen = false;
                });
            }
        });
    };
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
                {#if loading}<Loading />{/if}
                Run Command
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>
