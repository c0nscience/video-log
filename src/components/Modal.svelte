<script>
    import {createEventDispatcher, onDestroy} from 'svelte';
    import Button from "./Button.svelte";

    const dispatch = createEventDispatcher();
    const close = () => dispatch('close');

    let modal;

    const handle_keydown = e => {
        if (e.key === 'Escape') {
            close();
            return;
        }

        if (e.key === 'Tab') {
            // trap focus
            const nodes = modal.querySelectorAll('*');
            const tabbable = Array.from(nodes).filter(n => n.tabIndex >= 0);

            let index = tabbable.indexOf(document.activeElement);
            if (index === -1 && e.shiftKey) index = 0;

            index += tabbable.length + (e.shiftKey ? -1 : 1);
            index %= tabbable.length;

            tabbable[index].focus();
            e.preventDefault();
        }
    };

    const previously_focused = typeof document !== 'undefined' && document.activeElement;

    if (previously_focused) {
        onDestroy(() => {
            previously_focused.focus();
        });
    }
</script>

<svelte:window on:keydown={handle_keydown}/>

<div class="fixed top-0 left-0 w-full h-full bg-slate-900/40" on:click={close}></div>

<div class="absolute left-1/2 top-1/2 overflow-auto p-1 rounded-md text-slate-900 bg-slate-200 modal" role="dialog"
     aria-modal="true" bind:this={modal}>
    <slot name="header"></slot>
    <slot></slot>
    <!-- svelte-ignore a11y-autofocus -->
    <Button label="Close" autofocus fn={close}/>
</div>

<style>
    .modal {
        width: calc(100vw - 4em);
        max-width: 32em;
        max-height: calc(100vh - 4em);
        transform: translate(-50%, -50%);
    }
</style>