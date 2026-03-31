<template>
    <div class="flex flex-col w-48 justify-center items-center relative bg-black z-50">

        <div class="flex items-center w-full cursor-pointer py-1" @click="dropped = !dropped">

            <div class="w-6 opacity-0 pointer-events-none"></div>

            <span class="flex-1 text-center whitespace-nowrap px-4">
                {{ selected }}
            </span>

            <div class="w-6 flex justify-center items-center transform transition-transform duration-200"
                :class="dropped ? 'rotate-180' : 'rotate-0'">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M15 12l-5-5-5 5" />
                </svg>
            </div>

        </div>

        <div :class="dropped ? 'hidden' : 'flex'"
            class="absolute top-full left-0 w-full flex-col m-y-1 items-center bg-black truncate">
            <div v-for="option of options" :key="option" @click="updateSelected(option)"
                class="cursor-pointer w-full text-center  hover:bg-gray-900">
                {{ option }}
            </div>
        </div>

    </div>
</template>
<script setup lang="ts">
import { GroupBy } from '@/Dashboard.vue';
import { ref } from 'vue';
const options = Object.values(GroupBy)
const selected = ref(GroupBy.Languages)
const dropped = ref(false)

const updateSelected = (val: GroupBy): void => {
    selected.value = val;
    emit('displayUpdated', val)
}

const emit = defineEmits<{
    (e: 'displayUpdated', value: GroupBy): void
}>()
</script>
