<script>
    export let value;
</script>

{#if value}<span class="badge bg-primary">Plugin</span>{:else}<span class="badge bg-warning">Theme</span>{/if}
