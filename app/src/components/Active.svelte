<script>
    export let value;
</script>

{#if value}<span class="badge bg-success">Active</span>{:else}<span class="badge bg-warning">Inactive</span>{/if}
