<script>

    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import { Label, FormGroup, FormText, Input } from 'sveltestrap';
    import Form from "../../form";

    export let item;
    export let site;
    export let fetchData; // to refresh data
    export let isOpen = false;
    export let onClose = null;
    export let formType = 'UPDATE';

    let loading = false;

    const fields = ['Type', 'Description', 'HttpMethod', 'HttpUrl', 'HttpHeaders', 'HttpBody', 'BuildPreviewRef', 'Public', 'GithubActionName'];

    let form = new Form(
        formType,
        formType === 'UPDATE' ? fields : fields,
        item,
        {
            "Public": "bool"
        },
        {
            onHydrate: data => {
                if (data === undefined) {
                    return data;
                }

                const cfg = JSON.parse(data['Config']);

                data['BuildPreviewRef'] = cfg['BuildPreviewRef'];

                return data;
            }
        }
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
        let endpoint = '/api/site/' + site.Uuid + '/command';

        if (formType === 'UPDATE') {
            method = 'PUT';
            endpoint += '/' + item.Uuid;
        }

        loading = true;
        form.clearErrors();

        let values = form.getValuesFromObj();

        if (values['Public'] === '' || values['Public'] === undefined) {
            values['Public'] = false;
        }

        fetch(endpoint, {
            method: method,
            body: JSON.stringify(values)
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
        <ModalHeader>{formType === 'CREATE' ? 'Add' : 'Update'} Command</ModalHeader>
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
                                <option value="BUILD_RELEASE">Build & Deploy</option>
                                <option value="PREVIEW_BUILD">Preview Build</option>
                                <option value="HTTP_CALL">Http Call</option>
                            </Input>
                        </FormGroup>
                    {/if}

                    {#if form.current.Description !== undefined}
                        <FormGroup>
                            <Label>Description</Label>
                            <Input type="text" bind:value={form.current.Description.value} valid={form.isValid(form.current.Description)} invalid={form.isInvalid(form.current.Description)} feedback={form.current.Description.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.Type.value === "PREVIEW_BUILD" }
                        {#if form.current.BuildPreviewRef !== undefined}
                            <FormGroup>
                                <Label>Github Ref (Branch/Commit/Tag)</Label>
                                <Input type="text" bind:value={form.current.BuildPreviewRef.value} valid={form.isValid(form.current.BuildPreviewRef)} invalid={form.isInvalid(form.current.BuildPreviewRef)} feedback={form.current.BuildPreviewRef.error} />
                            </FormGroup>
                        {/if}
                    {/if}

                    {#if form.current.Type.value === "BUILD_RELEASE" }
                        {#if form.current.GithubActionName !== undefined}
                            <FormGroup>
                                <Label>Workflow Action Name</Label>
                                <Input type="text" bind:value={form.current.GithubActionName.value} valid={form.isValid(form.current.GithubActionName)} invalid={form.isInvalid(form.current.GithubActionName)} feedback={form.current.GithubActionName.error} />
                            </FormGroup>
                        {/if}
                    {/if}

                    {#if form.current.Type.value === "HTTP_CALL" }
                        {#if form.current.HttpMethod !== undefined}
                        <FormGroup>
                            <Label>Http Method</Label>
                            <Input type="select" bind:value={form.current.HttpMethod.value} valid={form.isValid(form.current.HttpMethod)} invalid={form.isInvalid(form.current.HttpMethod)} feedback={form.current.HttpMethod.error}>
                                <option>Select method</option>
                                <option value="GET">GET</option>
                                <option value="POST">POST</option>
                            </Input>
                        </FormGroup>
                        {/if}

                        {#if form.current.HttpUrl !== undefined}
                            <FormGroup>
                                <Label>Http Url</Label>
                                <Input type="text" bind:value={form.current.HttpUrl.value} valid={form.isValid(form.current.HttpUrl)} invalid={form.isInvalid(form.current.HttpUrl)} feedback={form.current.HttpUrl.error} />
                            </FormGroup>
                        {/if}

                        {#if form.current.HttpHeaders !== undefined}
                            <FormGroup>
                                <Label>Http Headers</Label>
                                <Input type="textarea" bind:value={form.current.HttpHeaders.value} valid={form.isValid(form.current.HttpHeaders)} invalid={form.isInvalid(form.current.HttpHeaders)} feedback={form.current.HttpHeaders.error} />
                            </FormGroup>
                        {/if}

                        {#if form.current.HttpBody !== undefined}
                            <FormGroup>
                                <Label>Http Body</Label>
                                <Input type="textarea" bind:value={form.current.HttpBody.value} valid={form.isValid(form.current.HttpBody)} invalid={form.isInvalid(form.current.HttpBody)} feedback={form.current.HttpBody.error} />
                            </FormGroup>
                        {/if}
                    {/if}

                    {#if form.current.Public !== undefined}
                        <FormGroup>
                            <Input id="wp_cmd" label="Public" type="checkbox" bind:checked={form.current.Public.value} valid={form.isValid(form.current.Public)} invalid={form.isInvalid(form.current.Public)} feedback={form.current.Public.error} />
                        </FormGroup>
                    {/if}
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                {formType === 'CREATE' ? 'Add' : 'Update'} Command
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>