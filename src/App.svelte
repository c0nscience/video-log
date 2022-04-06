<script>
    import {onMount} from 'svelte'

    let videos = []
    let config = {}

    const getVideos = async () => {
        const res = await fetch('/videos')
        videos = await res.json()
    }

    onMount(async () => {
        const res = await fetch('/config')
        config = await res.json()

        if (config.dir !== '') {
            await getVideos()
        }
    })

    const updateDir = async () => {
        await fetch('/config', {
            method: "POST",
            body: JSON.stringify(config),
        })
        await getVideos()
    }
</script>

<main>
    <input bind:value={config.dir}/>
    <button type="button" on:click={updateDir}>Send</button>
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