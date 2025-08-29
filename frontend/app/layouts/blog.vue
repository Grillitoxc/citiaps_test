<script setup>
const route = useRoute()
const isMetrics = computed(() => route.path.startsWith('/posts/metrics'))
const isPosts   = computed(() => route.path.startsWith('/posts') && !isMetrics.value)

const { openCreate } = usePostUI()
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <header class="border-b bg-white">
      <div class="mx-auto max-w-6xl px-4 py-4">
        <div class="flex items-center justify-between">
          <h1 class="text-2xl font-bold">Gestión de Blogs - Prueba Técnica - Citiaps</h1>

          <!-- Botón siempre presente, pero se deshabilita en métricas -->
          <button
            class="rounded px-3 py-2 text-sm transition"
            :class="isMetrics 
              ? 'bg-gray-300 text-gray-500 cursor-not-allowed' 
              : 'bg-black text-white hover:bg-gray-800'"
            :disabled="isMetrics"
            @click="!isMetrics && openCreate()">
            + Nueva Publicación
          </button>
        </div>

        <!-- Tabs -->
        <div class="mt-4 inline-flex rounded border bg-white p-1">
          <NuxtLink to="/posts" class="px-3 py-1.5 text-sm rounded"
            :class="isPosts ? 'bg-black text-white' : 'text-gray-600'">
            Publicaciones
          </NuxtLink>
          <NuxtLink to="/posts/metrics" class="px-3 py-1.5 text-sm rounded"
            :class="isMetrics ? 'bg-black text-white' : 'text-gray-600'">
            Métricas
          </NuxtLink>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-6xl px-4 py-6">
      <slot />
    </main>
  </div>
</template>
