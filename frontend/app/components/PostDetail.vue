<!-- components/PostDetail.vue -->
<script setup>
const props = defineProps({
    post: { type: Object, default: null },
    open: { type: Boolean, default: false },
})
const emit = defineEmits(['close'])

const formatDate = (d) => {
    if (!d) return '-'
    const date = typeof d === 'string' ? new Date(d) : d
    return date.toLocaleDateString()
}
</script>

<template>
    <div v-if="open && post" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- backdrop -->
        <div class="absolute inset-0 bg-black/50" @click="emit('close')" />
        <!-- modal -->
        <div class="relative z-10 w-full max-w-3xl bg-white rounded-lg shadow-lg p-6 max-h-[90vh] overflow-y-auto">
            <div class="flex items-start justify-between gap-4">
                <h2 class="text-xl font-semibold leading-tight">{{ post.title }}</h2>
                <span class="shrink-0 inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
                    :class="post.published ? 'bg-black text-white' : 'bg-gray-200 text-gray-700'">
                    {{ post.published ? 'Publicado' : 'Borrador' }}
                </span>
            </div>

            <div class="mt-4 grid gap-2 text-sm text-gray-600">
                <div>Autor: <span class="font-medium text-gray-900">{{ post.author }}</span></div>
                <div>
                    Fecha: <span class="font-medium text-gray-900">
                        {{ post.published && post.publishedAt ? `Publicado el ${formatDate(post.publishedAt)}` : `Creado
                        el
                        ${formatDate(post.createdAt)}` }}
                    </span>
                </div>
            </div>

            <div v-if="post.tags?.length" class="mt-3 flex flex-wrap gap-2">
                <span v-for="t in post.tags" :key="t"
                    class="text-xs bg-gray-100 text-gray-700 border border-gray-200 px-2 py-1 rounded">
                    {{ t }}
                </span>
            </div>

            <div class="mt-6 whitespace-pre-wrap leading-relaxed">
                {{ post.content }}
            </div>

            <div class="mt-6 flex justify-end">
                <button class="border rounded px-3 py-1.5 hover:bg-gray-50" @click="emit('close')">Cerrar</button>
            </div>
        </div>
    </div>
</template>
