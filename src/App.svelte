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

    const updateEntry = async (v) => {
        await fetch(`/videos`, {
            method: "POST",
            body: JSON.stringify(v),
        })
        await getVideos()
    }
</script>

<main>
    <div class="mb-12 ml-8">
        <input class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 shadow-sm sm:text-sm border-transparent rounded-md bg-slate-500 text-slate-200 py-2 px-4"
               bind:value={config.dir}/>
        <button class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600"
                type="button" on:click={updateDir}>Send
        </button>
    </div>
    <div class="grid gap-4 grid-cols-3 mx-8">
        {#each videos as video}
            <div class="rounded-lg bg-slate-700 shadow-2xl px-4 py-8 text-slate-200 grid grid-flow-row auto-rows-auto content-between gap-4">
                <div class="text-2xl">{video.date}</div>
                <div class="mt-1">
                    <textarea id="description" name="about" rows="3"
                              class="text-mono shadow-sm focus:ring-indigo-500 focus:border-indigo-500 mt-1 block w-full sm:text-sm border border-0 rounded-md text-slate-200 bg-slate-500 p-2"
                              bind:value={video.description}></textarea>
                </div>
                <div class="grid grid-cols-2 justify-between">
                    <div class="text-xl font-semibold underline text-slate-400">{video.name}.mkv</div>
                    <button class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600"
                            type="button"
                            on:click={() => {updateEntry(video)}}>
                        Save
                    </button>
                </div>
            </div>
        {/each}
    </div>
</main>

<style global lang="postcss">
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
</style>