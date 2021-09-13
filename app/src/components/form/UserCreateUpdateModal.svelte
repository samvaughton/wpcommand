<script>

    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';
    import { Label, FormGroup, FormText, Input } from 'sveltestrap';
    import Form from "../../form";

    export let item = null;
    export let accountUuid;
    export let fetchData; // to refresh data
    export let isOpen = false;
    export let onClose = null;
    export let formType = 'UPDATE';

    let loading = false;

    let form = new Form(
        formType,
        formType === 'CREATE' ? ['FirstName', 'LastName', 'Email', 'Password'] : ['FirstName', 'LastName', 'Email', 'Password'],
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
        let endpoint = '/api/account/' + accountUuid + '/user';

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

                    if (formType === 'CREATE') {
                        window.location = '/accounts/' + accountUuid + '/users/' + item.Uuid;
                    }

                });

                toggle();
            }
        });
    };
</script>


<Modal isOpen={isOpen} {toggle} on:close={_onClose}>
    <form on:submit|preventDefault={submitModal}>
        <ModalHeader>{formType === 'CREATE' ? 'Add' : 'Update'} User</ModalHeader>
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
                    {#if form.current.Email !== undefined}
                    <FormGroup>
                        <Label>Email</Label>
                        <Input type="text" bind:value={form.current.Email.value} valid={form.isValid(form.current.Email)} invalid={form.isInvalid(form.current.Email)} feedback={form.current.Email.error} />
                    </FormGroup>
                    {/if}

                    {#if form.current.FirstName !== undefined}
                        <FormGroup>
                            <Label>First Name</Label>
                            <Input type="text" bind:value={form.current.FirstName.value} valid={form.isValid(form.current.FirstName)} invalid={form.isInvalid(form.current.FirstName)} feedback={form.current.FirstName.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.LastName !== undefined}
                        <FormGroup>
                            <Label>Last Name</Label>
                            <Input type="text" bind:value={form.current.LastName.value} valid={form.isValid(form.current.LastName)} invalid={form.isInvalid(form.current.LastName)} feedback={form.current.LastName.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.Password !== undefined}
                        <FormGroup>
                            <Label>Password</Label>
                            <Input type="password" bind:value={form.current.Password.value} valid={form.isValid(form.current.Password)} invalid={form.isInvalid(form.current.Password)} feedback={form.current.Password.error} />
                            {#if formType === 'UPDATE'}
                                <div class="help text-muted">
                                    Only enter the password if you wish to change it.
                                </div>
                            {/if}
                        </FormGroup>
                    {/if}
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-primary" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                {formType === 'CREATE' ? 'Add' : 'Update'} User
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>