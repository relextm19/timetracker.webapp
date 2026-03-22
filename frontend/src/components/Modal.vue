<template>
    <div
        class="absolute top-1/2 left-1/2 -translate-1/2 flex justify-center items-center flex-col border-2 border-white rounded-lg gap-4 p-4 z-50 bg-black">
        <h2 class="text-2xl tracking-wider">{{ props.title }}</h2>
        <slot>
        </slot>
    </div>
</template>

<script setup lang="ts">
const props = defineProps<{
    title: string
}>()
</script>
