<script>
    import {onMount} from 'svelte'

    let videos = []
    let config = {}

    onMount(async () => {
        const resCfg = await fetch('/config')
        config = await resCfg.json()

        if (config.Dir === '') {
            console.log('set dir')
        } else {
            const resVid = await fetch('/videos')
            videos = await resVid.json()
        }
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