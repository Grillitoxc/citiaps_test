<!-- app/error.vue -->
<script setup>
const props = defineProps({
  error: Object
})

const handleError = () => clearError({ redirect: '/posts' })
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 text-center px-4">
    <h1 class="text-5xl font-bold text-gray-800">404</h1>
    <p class="text-gray-600 mt-2">Oops... la página que buscas no existe.</p>
    
    <div class="mt-6 flex gap-4">
      <button @click="handleError" class="px-4 py-2 rounded bg-black text-white">
        Ir a Publicaciones
      </button>
      <button @click="$router.back()" class="px-4 py-2 rounded border">
        Volver atrás
      </button>
    </div>
  </div>
</template>
