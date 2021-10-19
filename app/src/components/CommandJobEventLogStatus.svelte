<script>
    export let value;
</script>

{#if value === "SKIPPED"}
    <span class="badge bg-warning">Skipped</span>
{:else if value === "SUCCESS"}
    <span class="badge bg-success">Success</span>
{:else if value === "FAILURE"}
    <span class="badge bg-danger">Failure</span>
{:else if value === "TERMINATED"}
    <span class="badge bg-danger">Terminated</span>
{/if}
