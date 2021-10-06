<script>

    import {hasAccess, AuthEnum} from "../../store/user";
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

    let form = new Form(
        formType,
        formType === 'UPDATE' ? ['Email', 'Role', 'Password'] : ['Username', 'Email', 'Role', 'Password'],
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
        let endpoint = '/api/site/' + site.Uuid + '/wp/user';

        if (formType === 'UPDATE') {
            method = 'PUT';
            endpoint += '/' + item['Id'];
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
                    // do not update from response since it does not return the new object

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
                    {#if form.current.Username !== undefined}
                        <FormGroup>
                            <Label>Username</Label>
                            <Input type="text" bind:value={form.current.Username.value} valid={form.isValid(form.current.Username)} invalid={form.isInvalid(form.current.Username)} feedback={form.current.Username.error} />
                        </FormGroup>
                    {/if}

                    {#if form.current.Email !== undefined}
                    <FormGroup>
                        <Label>Email</Label>
                        <Input type="text" bind:value={form.current.Email.value} valid={form.isValid(form.current.Email)} invalid={form.isInvalid(form.current.Email)} feedback={form.current.Email.error} />
                    </FormGroup>
                    {/if}

                    {#if form.current.Role !== undefined}
                        <FormGroup>
                            <Label>Role</Label>
                            <Input type="select" bind:value={form.current.Role.value} valid={form.isValid(form.current.Role)} invalid={form.isInvalid(form.current.Role)} feedback={form.current.Role.error}>
                                <option>Select type</option>
                                <option value="owner">Site Owner</option>
                                {#await hasAccess(AuthEnum.ObjectWordpressUser, AuthEnum.ActionWriteSpecial)}
                                {:then result}
                                    <option value="administrator">Administrator</option>
                                {/await}

                            </Input>
                        </FormGroup>
                    {/if}

                    {#if form.current.Password !== undefined}
                        <FormGroup>
                            <Label>Password</Label>
                            <Input type="password" bind:value={form.current.Password.value} valid={form.isValid(form.current.Password)} invalid={form.isInvalid(form.current.Password)} feedback={form.current.Password.error} />
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