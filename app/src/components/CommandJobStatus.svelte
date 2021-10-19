<script>
    export let value;
</script>

{#if value === "CREATED"}
    <span class="badge bg-secondary">Created</span>
{:else if value === "PENDING"}
    <span class="badge bg-info">Pending</span>
{:else if value === "RUNNING"}
    <span class="badge bg-primary">Running</span>
{:else if value === "SUCCESS"}
    <span class="badge bg-success">Success</span>
{:else if value === "FAILURE"}
    <span class="badge bg-danger">Failure</span>
{:else if value === "TERMINATED"}
    <span class="badge bg-danger">Terminated</span>
{/if}
