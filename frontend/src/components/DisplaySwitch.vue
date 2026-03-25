<template>
    <div class="flex justify-center items-center flex-col w-1/6 text-wrap relative">
        <div class="flex justify-center items-center w-full">
            <span class="flex-5/6 text-center">
                {{ selected }}
            </span>
            <div :class="dropped ? 'rotate-180' : 'rotate-0'" @click="dropped = !dropped;"
                class="cursor-pointer transform transition-transform duration-200 absolute right-0">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M15 12l-5-5-5 5" />
                </svg>
            </div>
        </div>
        <div :class="dropped ? 'hidden' : ''" class="w-full flex items-center justify-center flex-col self-start">
            <div v-for="option of options" @click="updateSelected(option)" class="cursor-pointer">
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
