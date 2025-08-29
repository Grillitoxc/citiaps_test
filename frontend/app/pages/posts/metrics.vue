<script setup>
definePageMeta({ layout: 'blog' })
const { get } = useApi()

// Estado
const loading = ref(true)
const err = ref(null)

const posts = ref([])                   // todos los posts (para métricas locales)
const topTags = ref([])                 // [{ tag, count }]
const totals = reactive({               // tarjetas
    total: 0,
    published: 0,
    drafts: 0,
    uniqueTags: 0
})
const authors = ref([])                 // [{ name, count }]
const allTags = ref([])                 // [{ tag, count }]

// Helpers
const fetchAllPosts = async () => {
    // traemos todo paginando en bloques de 100
    const perPage = 100
    let page = 1
    let acc = []
    while (true) {
        const res = await get('/posts', {
            query: {
                page,
                limit: perPage,
                sort: '-publishedAt'
            }
        })
        acc = acc.concat(res.items || [])
        const total = res.total || acc.length
        if (acc.length >= total) break
        page++
    }
    return acc
}

const computeFromPosts = (list) => {
    totals.total = list.length
    const pub = list.filter(p => !!p.published)
    totals.published = pub.length
    totals.drafts = list.length - pub.length

    // Unique tags
    const tagSet = new Set()
    list.forEach(p => (p.tags || []).forEach(t => tagSet.add(t)))
    totals.uniqueTags = tagSet.size

    // Posts por autor
    const byAuthor = new Map()
    list.forEach(p => {
        const name = p.author || 'Desconocido'
        byAuthor.set(name, (byAuthor.get(name) || 0) + 1)
    })
    authors.value = [...byAuthor.entries()]
        .map(([name, count]) => ({ name, count }))
        .sort((a, b) => b.count - a.count)

    // Todos los tags (con conteo)
    const tagCount = new Map()
    list.forEach(p => (p.tags || []).forEach(t => tagCount.set(t, (tagCount.get(t) || 0) + 1)))
    allTags.value = [...tagCount.entries()]
        .map(([tag, count]) => ({ tag, count }))
        .sort((a, b) => b.count - a.count)
}

const load = async () => {
    loading.value = true
    err.value = null
    try {
        // 1) Top 10 tags SOLO publicados (backend)
        const tagsRes = await get('/posts/metrics/by-tag', {
            query: { limit: 10, published: 'true' }
        })
        // normaliza posibles nombres: TagMetric => { tag, count }
        topTags.value = (tagsRes || []).map(r => ({
            tag: r.tag ?? r._id ?? r.Tag ?? r.id ?? '—',
            count: r.count ?? r.Count ?? 0
        }))

        // 2) Todos los posts (front metrics)
        const full = await fetchAllPosts()
        posts.value = full
        computeFromPosts(full)
    } catch (e) {
        err.value = e
    } finally {
        loading.value = false
    }
}

// Primera carga
await load()

// Derivados para gráficos simples
const publishedPct = computed(() => {
    if (!totals.total) return 0
    return Math.round((totals.published / totals.total) * 100)
})
const draftsPct = computed(() => 100 - publishedPct.value)
</script>

<template>
    <div class=" bg-gray-50 min-h-screen">
        <div class="mx-auto max-w-6xl">
            <div v-if="err" class="text-red-600 mb-4">
                Error: {{ err?.message || err?.data?.message || err }}
            </div>

            <!-- Tarjetas -->
            <div class="grid gap-4 grid-cols-1 md:grid-cols-2 xl:grid-cols-4">
                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-sm font-medium text-gray-700 flex items-center justify-between">
                        Total Posts
                        <svg class="w-4 h-4 text-gray-400" viewBox="0 0 24 24" fill="none">
                            <path d="M6 4h9a2 2 0 0 1 2 2v1h1a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2V4Z"
                                stroke="currentColor" stroke-width="1.5" />
                        </svg>
                    </div>
                    <div class="mt-3 text-2xl font-bold">{{ totals.total }}</div>
                    <div class="text-xs text-gray-500 mt-1">Todas las publicaciones</div>
                </div>

                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-sm font-medium text-gray-700 flex items-center justify-between">
                        Publicados
                        <svg class="w-4 h-4 text-gray-400" viewBox="0 0 24 24" fill="none">
                            <path d="M12 5c3.866 0 7 3.582 7 8s-3.134 8-7 8-7-3.582-7-8 3.134-8 7-8Z"
                                stroke="currentColor" stroke-width="1.5" />
                            <circle cx="12" cy="13" r="3" stroke="currentColor" stroke-width="1.5" />
                        </svg>
                    </div>
                    <div class="mt-3 text-2xl font-bold">{{ totals.published }}</div>
                    <div class="text-xs text-gray-500 mt-1">{{ publishedPct }}% del total</div>
                </div>

                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-sm font-medium text-gray-700 flex items-center justify-between">
                        Borradores
                        <svg class="w-4 h-4 text-gray-400" viewBox="0 0 24 24" fill="none">
                            <path d="M7 3h10v4H7zM5 7h14v14H5z" stroke="currentColor" stroke-width="1.5" />
                        </svg>
                    </div>
                    <div class="mt-3 text-2xl font-bold">{{ totals.drafts }}</div>
                    <div class="text-xs text-gray-500 mt-1">Pendientes de publicar</div>
                </div>

                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-sm font-medium text-gray-700 flex items-center justify-between">
                        Tags Únicos
                        <svg class="w-4 h-4 text-gray-400" viewBox="0 0 24 24" fill="none">
                            <path d="M7 7h.01M3 12l7.586-7.586A2 2 0 0 1 12 4h8v8a2 2 0 0 1-.586 1.414L12 21 3 12Z"
                                stroke="currentColor" stroke-width="1.5" stroke-linecap="round"
                                stroke-linejoin="round" />
                        </svg>
                    </div>
                    <div class="mt-3 text-2xl font-bold">{{ totals.uniqueTags }}</div>
                    <div class="text-xs text-gray-500 mt-1">Categorías diferentes</div>
                </div>
            </div>

            <!-- Gráficos principales -->
            <div class="grid gap-4 grid-cols-1 xl:grid-cols-2 mt-6">
                <!-- Top tags (Top 10 publicados) -->
                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-base font-semibold">Posts por Tag (Top 10)</div>
                    <div class="text-sm text-gray-500 mb-4">Los tags más utilizados en las publicaciones</div>

                    <div v-if="loading" class="text-sm text-gray-500">Cargando…</div>
                    <div v-else>
                        <div v-if="!topTags.length" class="text-sm text-gray-500">No hay datos</div>
                        <ul v-else class="space-y-2">
                            <li v-for="(t, i) in topTags" :key="t.tag" class="flex items-center gap-3">
                                <div class="w-6 text-right text-xs text-gray-500">{{ i + 1 }}</div>
                                <div class="flex-1">
                                    <div class="flex items-center justify-between text-sm">
                                        <span class="font-medium text-gray-800">{{ t.tag }}</span>
                                        <span class="text-gray-500">{{ t.count }}</span>
                                    </div>
                                    <!-- bar SVG simple -->
                                    <svg :width="'100%'" height="8" class="mt-1">
                                        <rect x="0" y="0" :width="(t.count / (topTags[0]?.count || 1)) * 100 + '%'"
                                            height="8" rx="4" class="fill-gray-800" />
                                        <rect :x="(t.count / (topTags[0]?.count || 1)) * 100 + '%'" y="0"
                                            :width="(100 - (t.count / (topTags[0]?.count || 1)) * 100) + '%'" height="8"
                                            rx="4" class="fill-gray-200" />
                                    </svg>
                                </div>
                            </li>
                        </ul>
                    </div>
                </div>

                <!-- Posts por autor -->
                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-base font-semibold">Posts por Autor</div>
                    <div class="text-sm text-gray-500 mb-4">Productividad de cada autor</div>

                    <ul class="divide-y">
                        <li v-for="(a, idx) in authors" :key="a.name" class="flex items-center justify-between py-3">
                            <div class="flex items-center gap-3">
                                <span
                                    class="w-8 h-8 rounded-full bg-gray-100 text-gray-700 flex items-center justify-center text-sm font-medium">
                                    {{ idx + 1 }}
                                </span>
                                <span class="font-medium text-gray-800">{{ a.name }}</span>
                            </div>
                            <span class="text-xs bg-gray-700 text-white rounded px-2 py-1">{{ a.count }} posts</span>
                        </li>
                    </ul>
                </div>

                
            </div>

            <!-- Secundarias -->
            <div class="grid gap-4 grid-cols-1 xl:grid-cols-2 mt-6 ">
                

                <!-- Todos los Tags -->
                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-base font-semibold">Todos los Tags</div>
                    <div class="text-sm text-gray-500 mb-4">Lista completa de tags utilizados</div>

                    <div class="flex flex-wrap gap-2">
                        <span v-for="t in allTags" :key="t.tag"
                            class="inline-flex items-center gap-2 text-xs bg-gray-100 text-gray-700 border border-gray-200 px-3 py-1.5 rounded">
                            <svg class="h-3.5 w-3.5 text-gray-500" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                                <path
                                    d="M7 7h.01M3 12l7.586-7.586a2 2 0 0 1 1.414-.586H20a1 1 0 0 1 1 1v8a2 2 0 0 1-.586 1.414L12 21 3 12z"
                                    stroke="currentColor" stroke-width="1.5" stroke-linecap="round"
                                    stroke-linejoin="round" />
                            </svg>
                            <span>{{ t.tag }}</span>
                            <span class="ml-1 inline-block rounded bg-white border px-1.5 py-0.5">{{ t.count }}</span>
                        </span>
                    </div>
                </div>

                <!-- Distribución publicados/borradores -->
                <div class="rounded-lg border bg-white p-4 hover:shadow">
                    <div class="text-base font-semibold">Estado de Publicaciones</div>
                    <div class="text-sm text-gray-500 mb-4">Distribución entre publicados y borradores</div>

                    <div class="mt-2">
                        <div class="w-full h-6 bg-gray-200 rounded overflow-hidden">
                            <div class="h-6 bg-black" :style="{ width: publishedPct + '%' }"></div>
                        </div>
                        <div class="mt-2 flex justify-between text-sm text-gray-700">
                            <div><span class="inline-block w-3 h-3 bg-black mr-2"></span>Publicados ({{ totals.published
                                }})</div>
                            <div><span class="inline-block w-3 h-3 bg-gray-300 mr-2"></span>Borradores ({{ totals.drafts
                                }})</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
