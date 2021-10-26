<script>
    import {Modal, ModalHeader, ModalBody, ModalFooter} from 'sveltestrap';

    export let name = null;
    export let notice = "";
    export let endpoint = null;
    export let redirectTo = null;
    export let isOpen = false;
    export let onClose = null;

    let loading = false;
    let warningMessage = '';

    const toggle = () => (isOpen = !isOpen);

    const onCloseFn = function () {
        loading = false;
        warningMessage = '';

        if (typeof onClose === 'function') {
            onClose();
        }
    };

    let performDelete = function () {
        loading = true;
        warningMessage = '';

        fetch(endpoint, {
            method: "DELETE",
        }).then(resp => {
            loading = false;

            if (resp.status === 200) {
                resp.json().then(data => {
                    if (redirectTo !== null) {
                        window.location = redirectTo;
                    } else {
                        onCloseFn();
                    }
                });
            } else {
                warningMessage = `Something went wrong when deleting the ${name}.`;
            }
        });
    };
</script>

<Modal isOpen={isOpen} {toggle} on:close={onCloseFn}>
    <form on:submit|preventDefault={performDelete}>
        <ModalHeader>Delete {name}</ModalHeader>
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
            {#if notice !== ""}
                <div class="row">
                    <div class="col-12">
                        <div class="alert alert-danger" role="alert">
                            {notice}
                        </div>
                    </div>
                </div>
            {/if}
            <div class="row mb-3">
                <div class="col-12">
                    <p>Are you sure you want to delete this {name}?</p>
                </div>
            </div>
        </ModalBody>
        <ModalFooter>
            <button type="submit" class="btn btn-danger" disabled={loading}>
                {#if loading}<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>{/if}
                Delete {name}
            </button>
            <button type="button" class="btn btn-secondary" on:click={toggle}>Cancel</button>
        </ModalFooter>
    </form>
</Modal>