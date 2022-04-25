<script>
    export let label;
    export let fn;
</script>

<button class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600"
        type="button"
        on:click={fn}>
    <slot>
        {label}
    </slot>
</button>