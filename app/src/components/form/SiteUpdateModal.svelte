<script>

    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import { Label, FormGroup, FormText, Input } from 'sveltestrap';
    import Form from "../../form";

    export let item;
    export let fetchData; // to refresh data
    export let isOpen = false;
    export let onClose = null;
    export let formType = 'UPDATE';

    let loading = false;

    const fields = ['Description', 'LabelSelector', 'Namespace', 'Enabled', 'SiteConfig', 'SiteEmail', 'SiteUsername', 'SitePassword', 'WpDomain', 'DockerRegistryName'];

    let form = new Form(
        formType,
        formType === 'UPDATE' ? fields : fields,
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
        let endpoint = '/api/site';

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
        <ModalHeader>{formType === 'CREATE' ? 'Add' : 'Update'} Site</ModalHeader>
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
                    {#if form.current.Description !== undefined}
                    <FormGroup>
                        <Label>Description</Label>
                        <Input type="text" bind:value={form.current.Description.value} valid={form.isValid(form.current.Description)} invalid={form.isInvalid(form.current.Description)} feedback={form.current.Description.error} />
                    </FormGroup>
                    {/if}

                    {#if form.current.Namespace !== undefined}
                        <FormGroup>
                            <Label>Namespace</Label>
                            <Input type="text" bind:value={form.current.Namespace.value} valid={form.isValid(form.current.Namespace)} invalid={form.isInvalid(form.current.Namespace)} feedback={form.current.Namespace.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.LabelSelector !== undefined}
                        <FormGroup>
                            <Label>LabelSelector</Label>
                            <Input type="text" bind:value={form.current.LabelSelector.value} valid={form.isValid(form.current.LabelSelector)} invalid={form.isInvalid(form.current.LabelSelector)} feedback={form.current.LabelSelector.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.SiteConfig !== undefined}
                        <FormGroup>
                            <Label>Site Config</Label>
                            <Input type="textarea" bind:value={form.current.SiteConfig.value} valid={form.isValid(form.current.SiteConfig)} invalid={form.isInvalid(form.current.SiteConfig)} feedback={form.current.SiteConfig.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.SiteUsername !== undefined}
                        <FormGroup>
                            <Label>Site Username</Label>
                            <Input type="text" bind:value={form.current.SiteUsername.value} valid={form.isValid(form.current.SiteUsername)} invalid={form.isInvalid(form.current.SiteUsername)} feedback={form.current.SiteUsername.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.SiteEmail !== undefined}
                        <FormGroup>
                            <Label>Site Email</Label>
                            <Input type="text" bind:value={form.current.SiteEmail.value} valid={form.isValid(form.current.SiteEmail)} invalid={form.isInvalid(form.current.SiteEmail)} feedback={form.current.SiteEmail.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.SitePassword !== undefined}
                        <FormGroup>
                            <Label>Site Config</Label>
                            <Input type="text" bind:value={form.current.SitePassword.value} valid={form.isValid(form.current.SitePassword)} invalid={form.isInvalid(form.current.SitePassword)} feedback={form.current.SitePassword.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.WpDomain !== undefined}
                        <FormGroup>
                            <Label>Wordpress Domain</Label>
                            <Input type="text" bind:value={form.current.WpDomain.value} valid={form.isValid(form.current.WpDomain)} invalid={form.isInvalid(form.current.WpDomain)} feedback={form.current.WpDomain.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.DockerRegistryName !== undefined}
                        <FormGroup>
                            <Label>Docker Registry Name</Label>
                            <Input type="text" bind:value={form.current.DockerRegistryName.value} valid={form.isValid(form.current.DockerRegistryName)} invalid={form.isInvalid(form.current.DockerRegistryName)} feedback={form.current.DockerRegistryName.error} />
                        </FormGroup>
                    {/if}
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                {formType === 'CREATE' ? 'Add' : 'Update'} Site
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>