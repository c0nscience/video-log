<script>
    import {onMount} from 'svelte'
    import Button from "./components/Button.svelte";
    import Modal from "./components/Modal.svelte";

    let videos = []
    let config = {}
    let tools = {}
    let showDeleteModal = false
    let showTmspModal = false
    let videoToDelete = undefined
    let tmsp = ""

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

    const deleteEntry = async (v) => {
        await fetch(`/videos`, {
            method: "DELETE",
            body: JSON.stringify(v),
        })
        await getVideos()
    }

    const transformTmsp = async () => {
        const res = await fetch(`/tools/timestamps`, {
            method: "POST",
            body: JSON.stringify(tools),
        })
        const t = await res.json()
        tmsp = t.timestamp
        showTmspModal = true
    }
</script>

<main>
    <div class="mb-12 ml-8 grid grid-cols-2 gap-4">
        <div>
            <input
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 shadow-sm sm:text-sm border-transparent rounded-md bg-slate-500 text-slate-200 py-2 px-4"
                    bind:value={config.dir}/>
            <Button label="Send" fn={updateDir}/>
        </div>
        <div>
            <input
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 shadow-sm sm:text-sm border-transparent rounded-md bg-slate-500 text-slate-200 py-2 px-4"
                    bind:value={tools.path}/>
            <Button label="Send" fn={transformTmsp}/>
        </div>
    </div>
    <div class="mb-8 ml-8">
        <Button fn={getVideos}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                 stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
        </Button>
    </div>
    <div class="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3 mx-16 md:mx-8">
        {#each videos as video}
            <div class="rounded-lg bg-slate-700 shadow-2xl px-4 py-8 text-slate-200 grid grid-flow-row auto-rows-auto content-between gap-4">
                <div class="text-2xl">{video.date}</div>
                <div class="mt-1">
                    <textarea id="description" name="about" rows="3"
                              class="text-mono shadow-sm focus:ring-indigo-500 focus:border-indigo-500 mt-1 block w-full sm:text-sm border border-0 rounded-md text-slate-200 bg-slate-500 p-2"
                              bind:value={video.description}></textarea>
                </div>
                <div class="grid gap-1 grid-rows-2 2xl:grid-rows-none 2xl:grid-cols-2 2xl:justify-between">
                    <div class="text-xl font-semibold underline text-slate-400">{video.name}.mkv</div>
                    <Button label="Save" fn={() => {updateEntry(video)}}/>
                    <Button label="Delete" fn={() => {
                        showDeleteModal = true
                        videoToDelete = video
                    }}/>
                </div>
            </div>
        {/each}
    </div>
</main>

{#if showDeleteModal}
    <Modal on:close="{() => showDeleteModal = false}">
        <p>Do you really want to delete this?</p>

        <Button label="Delete" fn={() => {
            deleteEntry(videoToDelete)
            videoToDelete = undefined
            showDeleteModal = false
        }}/>
    </Modal>
{/if}

{#if showTmspModal}
    <Modal on:close="{() => showTmspModal = false}">
        <p>Timestamp:</p>
        <textarea rows="20"
                  class="text-mono shadow-sm focus:ring-indigo-500 focus:border-indigo-500 mt-1 block w-full sm:text-sm border border-0 rounded-md text-slate-200 bg-slate-500 p-2">{tmsp}</textarea>
    </Modal>
{/if}


<style global lang="postcss">
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
</style>