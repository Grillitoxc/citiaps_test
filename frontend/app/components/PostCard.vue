<script setup>
const props = defineProps({
  post: { type: Object, required: true }
})
const emit = defineEmits(['view', 'edit', 'delete'])

const formatDate = (d) => {
  if (!d) return '-'
  const date = typeof d === 'string' ? new Date(d) : d
  return date.toLocaleDateString('es-CL', { timeZone: 'UTC' })
}
</script>

<template>
  <article class="card bg-white border rounded-lg p-4 hover:shadow transition-shadow flex flex-col">
    <header class="flex items-start justify-between gap-3">
      <h3 class="font-semibold text-lg line-clamp-2">{{ post.title }}</h3>
      <span
        class="shrink-0 inline-flex items-center rounded px-2 py-1 text-xs font-medium"
        :class="post.published ? 'bg-black text-white' : 'bg-gray-200 text-gray-700'"
      >
        {{ post.published ? 'Publicado' : 'Borrador' }}
      </span>
    </header>

    <div class="card-body flex-1 flex flex-col">
      <p class="card-desc text-sm text-gray-600 mt-1 line-clamp-2">
        {{ post.content }}
      </p>

      <div class="mt-3 text-sm text-gray-600 flex items-center gap-2">
        <span class="text-gray-500">Autor:</span>
        <span class="font-medium text-gray-900">{{ post.author }}</span>
      </div>

      <div class="text-sm text-gray-600 flex items-center gap-2">
        <span class="text-gray-500">Fecha:</span>
        <span class="font-medium text-gray-900">
          {{ post.published && post.publishedAt ? formatDate(post.publishedAt) : formatDate(post.createdAt) }}
        </span>
      </div>

      <!-- Tags con icono -->
      <div class="card-tags mt-3 flex flex-wrap gap-1">
        <span
          v-for="t in (post.tags || [])"
          :key="t"
          class="inline-flex items-center gap-1 text-xs bg-gray-100 text-gray-700 border border-gray-200 px-2 py-1 rounded"
        >
          <!-- Tag icon -->
          <svg class="h-3.5 w-3.5 text-gray-500" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M7 7h.01M3 12l7.586-7.586a2 2 0 0 1 1.414-.586H20a1 1 0 0 1 1 1v7.999a2 2 0 0 1-.586 1.414L12 22 3 13z"
                  stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t }}
        </span>
      </div>

      <!-- Acciones con iconos -->
      <footer class="mt-auto pt-4 flex gap-2">
        <button
          class="inline-flex items-center gap-2 border rounded px-3 py-1.5 hover:bg-gray-50"
          @click="$emit('view', post)"
        >
          <!-- Eye icon -->
          <svg class="h-4 w-4 text-gray-700" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M2.25 12s3.75-6.75 9.75-6.75S21.75 12 21.75 12s-3.75 6.75-9.75 6.75S2.25 12 2.25 12Z"
                  stroke="currentColor" stroke-width="1.5"/>
            <circle cx="12" cy="12" r="3.25" stroke="currentColor" stroke-width="1.5"/>
          </svg>
          Ver
        </button>

        <button
          class="inline-flex items-center gap-2 border rounded px-3 py-1.5 hover:bg-gray-50"
          @click="$emit('edit', post)"
        >
          <!-- Pencil icon -->
          <svg class="h-4 w-4 text-gray-700" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="m16.862 3.487 3.651 3.651M4 20h4.5l10.5-10.5a2.58 2.58 0 0 0 0-3.651l-.849-.849a2.58 2.58 0 0 0-3.651 0L4 15.5V20Z"
                  stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          Editar
        </button>

        <button
          class="inline-flex items-center gap-2 border rounded px-3 py-1.5 text-red-600 hover:bg-red-50"
          @click="$emit('delete', post._id)"
        >
          <!-- Trash icon -->
          <svg class="h-4 w-4 text-red-600" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M4 7h16M9.5 7V5.5A1.5 1.5 0 0 1 11 4h2a1.5 1.5 0 0 1 1.5 1.5V7M6 7l1 12a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-12"
                  stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          Eliminar
        </button>
      </footer>
    </div>
  </article>
</template>