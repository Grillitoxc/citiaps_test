<!-- components/Modal.vue -->
<script setup>
const props = defineProps({
  open: { type: Boolean, default: false },
  title: { type: String, default: '' },
})
const emit = defineEmits(['close'])
</script>

<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div class="absolute inset-0 bg-black/50" @click="emit('close')" />
    <div class="relative z-10 w-full max-w-2xl bg-white rounded-lg shadow-lg">
      <header class="flex items-center justify-between p-4 border-b">
        <h3 class="text-lg font-semibold">{{ title }}</h3>
        <button class="text-gray-500 hover:text-black" @click="emit('close')">✕</button>
      </header>
      <div class="p-4 max-h-[80vh] overflow-y-auto">
        <slot />
      </div>
    </div>
  </div>
</template>
