// app/pages/posts/index.vue
<script setup>
// Meta
definePageMeta({ layout: 'blog' })

// Composición API
const { get, del } = useApi()

// Estados de filtro
const q = ref('')
const tag = ref('')
const published = ref('all')
const sort = ref('-publishedAt')
const page = ref(1)
const limit = ref(6)
const items = ref([])
const total = ref(0)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit.value)))
const loading = ref(false)
const err = ref(null)

// Opciones para "por página" (máx 100)
const limitOptions = [6, 9, 12, 18, 24, 36, 48, 60, 100]

// Tags (para filtro)
const allTags = ref([])

// Modal
const showDetail = ref(false)
const postSelected = ref(null)

// Helpers y load()
const toBoolOrNull = (v) => v === 'all' ? null : v === 'published'
const load = async () => {
    loading.value = true
    err.value = null
    try {
        const publishedParam = toBoolOrNull(published.value)
        const res = await get('/posts', {
            query: {
                q: q.value || undefined,
                tag: tag.value || undefined,
                published: publishedParam === null ? undefined : String(publishedParam),
                page: page.value,
                limit: Math.min(limit.value, 100),
                sort: sort.value,
            },
        })
        items.value = res.items || []
        total.value = res.total || 0
        allTags.value = Array.from(new Set(items.value.flatMap(p => p.tags || []))).sort()
    } catch (e) {
        err.value = e
    } finally {
        loading.value = false
    }
}

// Reseteo de filtros
watch([q, tag, published, sort, limit], () => {
    page.value = 1
    load()
})

watch(page, () => {
    load()
})

// Primera carga (SSR friendly)
await load()

// Acciones
const openDetail = (p) => { postSelected.value = p; showDetail.value = true }
const closeDetail = () => { showDetail.value = false; postSelected.value = null }

const removePost = async (id) => {
    if (!confirm('¿Eliminar esta publicación?')) return
    await del(`/posts/${id}`)
    await load()
}

// Mensaje de feedback
const message = ref(null)
const clearMessage = () => { message.value = null }

// --- Crear vía modal ---
const { showCreate, closeCreate } = usePostUI()
const { post: apiPost } = useApi()

const createPost = async (payload) => {
    try {
        await apiPost('/posts', payload)
        message.value = { type: 'success', text: '✅ Publicación creada correctamente.' }
        await load()
    } catch (e) {
        message.value = { type: 'error', text: '❌ Error al crear la publicación.' }
        throw e
    } finally {
        closeCreate()
        setTimeout(clearMessage, 4000)
    }
}

// --- Edición via modal ---
const showEdit = ref(false)
const editingPost = ref(null)

const { put: apiPut } = useApi()

const openEdit = (p) => {
    editingPost.value = { ...p }
    showEdit.value = true
}

const updatePost = async (payload) => {
    try {
        await apiPut(`/posts/${editingPost.value._id}`, payload)
        message.value = { type: 'success', text: '✅ Publicación actualizada.' }
        await load()
    } catch (e) {
        message.value = { type: 'error', text: '❌ Error al actualizar.' }
        throw e
    } finally {
        showEdit.value = false
        setTimeout(clearMessage, 4000)
    }
}
</script>

<template>
    <div class="min-h-screen bg-gray-50">
        <!-- Mensaje feedback -->
        <div v-if="message" class="mx-auto max-w-6xl pb-4 mt-4">
            <div class="px-4 py-3 rounded text-sm font-medium" :class="message.type === 'success'
                ? 'bg-green-100 text-green-800 border border-green-300'
                : 'bg-red-100 text-red-800 border border-red-300'">
                {{ message.text }}
            </div>
        </div>

        <div class="mx-auto max-w-6xl">
            <!-- Filtros -->
            <div class="bg-white border rounded-lg p-4 mb-6">
                <h2 class="font-medium mb-3">Filtros y Búsqueda</h2>
                <!-- Aumentamos a 5 columnas para incluir "Por página" -->
                <div class="grid gap-3 md:grid-cols-5">
                    <div>
                        <label class="block text-sm mb-1">Buscar</label>
                        <input v-model="q" placeholder="Título o contenido..."
                            class="w-full border rounded px-3 py-2" />
                    </div>
                    <div>
                        <label class="block text-sm mb-1">Tag</label>
                        <select v-model="tag" class="w-full border rounded px-3 py-2">
                            <option value="">Todos los tags</option>
                            <option v-for="t in allTags" :key="t" :value="t">{{ t }}</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm mb-1">Estado</label>
                        <select v-model="published" class="w-full border rounded px-3 py-2">
                            <option value="all">Todos</option>
                            <option value="published">Publicados</option>
                            <option value="draft">Borradores</option>
                        </select>
                    </div>
                    <div>
                        <label class="block text-sm mb-1">Ordenar por</label>
                        <select v-model="sort" class="w-full border rounded px-3 py-2">
                            <option value="-publishedAt">Más recientes</option>
                            <option value="publishedAt">Más antiguos</option>
                        </select>
                    </div>
                    <!-- NUEVO: Por página -->
                    <div>
                        <label class="block text-sm mb-1">Blogs por página</label>
                        <!-- v-model.number para asegurar número -->
                        <select v-model.number="limit" class="w-full border rounded px-3 py-2">
                            <option v-for="n in limitOptions" :key="n" :value="n">{{ n }}</option>
                        </select>
                        <p class="text-[11px] text-gray-500 mt-1">Máximo 100</p>
                    </div>
                </div>
            </div>

            <!-- Lista -->
            <ClientOnly>
                <div v-if="err" class="text-red-600 mb-4">Error: {{ err?.message || err?.data?.message || err }}</div>
                <div v-if="loading">Cargando…</div>

                <div v-else class="grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                    <PostCard v-for="p in items" :key="p._id" :post="p" @view="openDetail" @edit="openEdit"
                        @delete="removePost" />
                </div>

                <!-- Paginación -->
                <div v-if="totalPages > 1" class="mt-6 flex justify-center gap-2">
                    <button class="border rounded px-3 py-1.5" :disabled="page <= 1" @click="page--">Anterior</button>
                    <button v-for="n in totalPages" :key="n" class="border rounded px-3 py-1.5"
                        :class="n === page ? 'bg-black text-white' : ''" @click="page = n">
                        {{ n }}
                    </button>
                    <button class="border rounded px-3 py-1.5" :disabled="page >= totalPages"
                        @click="page++">Siguiente</button>
                </div>
            </ClientOnly>
        </div>

        <!-- Modal detalle -->
        <PostDetail :post="postSelected" :open="showDetail" @close="closeDetail" />

        <!-- Modal Crear -->
        <Modal :open="showCreate" title="Nueva Publicación" @close="closeCreate">
            <PostForm submit-text="Crear Publicación" @submit="createPost" @cancel="closeCreate" />
        </Modal>

        <!-- Modal Editar -->
        <Modal :open="showEdit" title="Editar Publicación" @close="showEdit = false">
            <PostForm :initial="editingPost" submit-text="Actualizar" @submit="updatePost" @cancel="showEdit = false" />
        </Modal>
    </div>
</template>

<style scoped>
main {
    padding: 0.5rem;
}
</style>