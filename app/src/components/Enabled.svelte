<script>
    export let value;
</script>

{#if value}<span class="badge bg-success">Enabled</span>{:else}<span class="badge bg-danger">Disabled</span>{/if}
