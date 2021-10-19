<script>
    export let value;
</script>

{#if value === "INFO"}
    <span class="badge bg-info">Info</span>
{:else if value === "DATA"}
    <span class="badge bg-primary">Data</span>
{:else if value === "JOB_STARTED"}
    <span class="badge bg-secondary">Job Started</span>
{:else if value === "JOB_FINISHED"}
    <span class="badge bg-secondary">Job Finished</span>
{/if}
