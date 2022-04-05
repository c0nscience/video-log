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
            <li>
                <strong>{video.name}</strong><br>
                <i>{video.description}</i>
            </li>
        {/each}
    </ul>
</main>

<style global lang="postcss">
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
</style>