<script>

    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import { Label, FormGroup, FormText, Input } from 'sveltestrap';
    import Form from "../../form";

    export let item;
    export let blueprint;
    export let fetchData; // to refresh data
    export let isOpen = false;
    export let onClose = null;

    let loading = false;

    let form = new Form(
        'UPDATE',
        ['Version', 'Url'],
        item
    );

    const toggle = () => (isOpen = !isOpen);
    const _onClose = function () {
        loading = false;
        form.reset();

        if (typeof onClose === 'function') {
            onClose();
        }
    };

    let submitModal = function () {
        let method = 'POST';
        let endpoint = "/api/blueprint/" + blueprint.Uuid + "/object/" + item.Uuid + "/revision/" + item.RevisionId + "/version";

        loading = true;
        form.clearErrors();
        fetch(endpoint, {
            method: method,
            body: JSON.stringify(form.getValuesFromObj())
        }).then(resp => {
            loading = false;

            if (resp.status !== 200) {
                resp.json().then(data => {
                    if (data.Status === "VALIDATION_ERRORS") {
                        form.hydrateErrorsFromValidationErrors(data.Errors);
                    } else {
                        form.errorMessage = data.Message;
                    }

                    form = form; // trigger update
                });
            } else {
                resp.json().then(data => {
                    form = form; // trigger update
                    window.location = '/blueprints/' + blueprint.Uuid + '/object/' + data.Uuid + '/revision/' + data.RevisionId;
                });

                if (typeof fetchData === 'function') {
                    fetchData();
                }

                toggle();
            }
        });
    };
</script>


<Modal isOpen={isOpen} {toggle} on:close={_onClose}>
    <form on:submit|preventDefault={submitModal}>
        <ModalHeader>Update Version</ModalHeader>
        <ModalBody>
            {#if form.errorMessage !== ""}
                <div class="row">
                    <div class="col-12">
                        <div class="alert alert-warning" role="alert">
                            {form.errorMessage}
                        </div>
                    </div>
                </div>
            {/if}
            <div class="row">
                <div class="col-12">

                    <FormGroup>
                        <Label>Version</Label>
                        <Input type="text" bind:value={form.current.Version.value} valid={form.isValid(form.current.Version)} invalid={form.isInvalid(form.current.Version)} feedback={form.current.Version.error} />
                        <FormText color="muted">
                            The version of the theme or plugin eg <code>3.1.2</code>
                        </FormText>
                    </FormGroup>

                    <FormGroup>
                        <Label>Zip File Url</Label>
                        <Input type="text" bind:value={form.current.Url.value} valid={form.isValid(form.current.Url)} invalid={form.isInvalid(form.current.Url)} feedback={form.current.Url.error} />
                        <FormText color="muted">
                            Leave blank to use the stored object URL
                        </FormText>
                    </FormGroup>

                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                Update Version
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>