<script>

    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import { Label, FormGroup, FormText, Input } from 'sveltestrap';
    import Form from "../../form";

    export let item;

    export let blueprint;
    export let fetchData; // to refresh data
    export let isOpen = false;
    export let onClose = null;
    export let formType = 'UPDATE';

    let loading = false;

    let form = new Form(
        formType,
        formType === 'UPDATE' ? ['Name'] : ['Type', 'Name', 'ExactName', 'Version', 'Url'],
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
        let endpoint = "/api/blueprint/" + blueprint.Uuid + "/object";

        if (formType === 'UPDATE') {
            method = 'PUT';
            endpoint += '/' + item.Uuid;
        }

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
                if (typeof fetchData === 'function') {
                    fetchData();
                }

                resp.json().then(data => {
                    item = data; // trigger update

                    form.initWithNewData(data);

                    form = form; // trigger update
                });

                toggle();
            }
        });
    };
</script>


<Modal isOpen={isOpen} {toggle} on:close={_onClose}>
    <form on:submit|preventDefault={submitModal}>
        <ModalHeader>{formType === 'CREATE' ? 'Add' : 'Update'} Object</ModalHeader>
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

                    {#if form.current.Type !== undefined}
                    <FormGroup>
                        <Label>Type</Label>
                        <Input type="select" bind:value={form.current.Type.value} valid={form.isValid(form.current.Type)} invalid={form.isInvalid(form.current.Type)} feedback={form.current.Type.error}>
                            <option>Select type</option>
                            <option value="PLUGIN">Plugin</option>
                            <option value="THEME">Theme</option>
                        </Input>
                    </FormGroup>
                    {/if}

                    {#if form.current.Name !== undefined}
                    <FormGroup>
                        <Label>Name</Label>
                        <Input type="text" bind:value={form.current.Name.value} valid={form.isValid(form.current.Name)} invalid={form.isInvalid(form.current.Name)} feedback={form.current.Name.error} />
                    </FormGroup>
                    {/if}

                    {#if form.current.ExactName !== undefined}
                    <FormGroup>
                        <Label>Exact Name</Label>
                        <Input type="text" bind:value={form.current.ExactName.value} valid={form.isValid(form.current.ExactName)} invalid={form.isInvalid(form.current.ExactName)} feedback={form.current.ExactName.error} />
                        <FormText color="muted">
                            What the theme or plugin calls itself, ie Advanced Custom Fields might be <code>acf</code>
                        </FormText>
                    </FormGroup>
                    {/if}

                    {#if form.current.Version !== undefined}
                    <FormGroup>
                        <Label>Version</Label>
                        <Input type="text" bind:value={form.current.Version.value} valid={form.isValid(form.current.Version)} invalid={form.isInvalid(form.current.Version)} feedback={form.current.Version.error} />
                        <FormText color="muted">
                            The version of the theme or plugin eg <code>3.1.2</code>
                        </FormText>
                    </FormGroup>
                    {/if}

                    {#if form.current.Url !== undefined}
                    <FormGroup>
                        <Label>Zip File Url</Label>
                        <Input type="text" bind:value={form.current.Url.value} valid={form.isValid(form.current.Url)} invalid={form.isInvalid(form.current.Url)} feedback={form.current.Url.error} />
                        <FormText color="muted">
                            URL of the zip file
                        </FormText>
                    </FormGroup>
                    {/if}
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                {formType === 'CREATE' ? 'Add' : 'Update'} Object
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>