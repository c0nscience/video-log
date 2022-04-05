<script>
    import {onMount} from 'svelte'

    let videos = []

    onMount(async () => {
        const res = await fetch(`/videos`)
        videos = await res.json()
    })
</script>

<main>
    <ul>
        {#each videos as video}
            <li>{video}</li>
        {/each}
    </ul>
</main>

<style>

</style>