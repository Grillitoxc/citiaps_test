<!-- components/PostForm.vue -->
<script setup>
const props = defineProps({
    initial: { type: Object, default: null },
    submitText: { type: String, default: 'Crear Publicación' },
})
const emit = defineEmits(['submit', 'cancel'])

const form = reactive({
    title: props.initial?.title || '',
    author: props.initial?.author || '',
    content: props.initial?.content || '',
    tags: Array.isArray(props.initial?.tags) ? [...props.initial.tags] : [],
    published: !!props.initial?.published,
})

const isEdit = computed(() => !!props.initial)
const isAlreadyPublished = computed(() => !!props.initial?.published)

const newTag = ref('')
const pending = ref(false)
const errors = reactive({ title: '', author: '', content: '' })

const validate = () => {
    errors.title = ''
    errors.author = ''
    errors.content = ''
    if (!form.title.trim()) errors.title = 'El título es requerido'
    else if (form.title.length < 5 || form.title.length > 140) errors.title = 'Entre 5 y 140 caracteres'
    if (!form.author.trim()) errors.author = 'El autor es requerido'
    if (!form.content.trim()) errors.content = 'El contenido es requerido'
    return !errors.title && !errors.author && !errors.content
}

const addTag = () => {
    const t = newTag.value.trim()
    if (t && !form.tags.includes(t)) form.tags.push(t)
    newTag.value = ''
}
const removeTag = (t) => { form.tags = form.tags.filter(x => x !== t) }
const onKey = (e) => { if (e.key === 'Enter') { e.preventDefault(); addTag() } }

const submit = async () => {
    if (!validate()) return
    pending.value = true
    try {
        emit('submit', {
            title: form.title,
            author: form.author,
            content: form.content,
            tags: form.tags,
            published: isAlreadyPublished.value ? true : form.published,
        })
    } finally {
        pending.value = false
    }
}
</script>

<template>
    <form class="space-y-4" @submit.prevent="submit">
        <p class="text-sm text-gray-600">
            {{ isEdit ? 'Modifica los datos de la publicación' : 'Completa los datos para crear una nueva publicación'
            }}
        </p>

        <div>
            <label class="block text-sm mb-1">Título *</label>
            <input v-model="form.title" placeholder="Título de la publicación (5-140 caracteres)"
                class="w-full border rounded px-3 py-2" :class="errors.title && 'border-red-500'" />
            <p v-if="errors.title" class="text-xs text-red-600 mt-1">{{ errors.title }}</p>
        </div>

        <div>
            <label class="block text-sm mb-1">Autor *</label>
            <input v-model="form.author" placeholder="Nombre del autor" class="w-full border rounded px-3 py-2"
                :class="errors.author && 'border-red-500'" />
            <p v-if="errors.author" class="text-xs text-red-600 mt-1">{{ errors.author }}</p>
        </div>

        <div>
            <label class="block text-sm mb-1">Contenido *</label>
            <textarea v-model="form.content" rows="7" placeholder="Contenido de la publicación"
                class="w-full border rounded px-3 py-2" :class="errors.content && 'border-red-500'"></textarea>
            <p v-if="errors.content" class="text-xs text-red-600 mt-1">{{ errors.content }}</p>
        </div>

        <div>
            <label class="block text-sm mb-1">Tags</label>
            <div class="flex gap-2">
                <input v-model="newTag" @keydown="onKey" placeholder="Agregar tag"
                    class="flex-1 border rounded px-3 py-2" />
                <button type="button" class="border rounded px-3 py-2" @click="addTag">Agregar</button>
            </div>
            <div class="flex flex-wrap gap-2 mt-2">
                <span v-for="t in form.tags" :key="t"
                    class="inline-flex items-center gap-1 text-xs bg-gray-100 border px-2 py-1 rounded">
                    {{ t }}
                    <button type="button" class="text-gray-500 hover:text-red-600" @click="removeTag(t)">✕</button>
                </span>
            </div>
        </div>

        <!-- Checkbox sólo para creación o para borradores en edición -->
        <template v-if="!isAlreadyPublished">
            <label class="inline-flex items-center gap-2">
                <input type="checkbox" v-model="form.published" />
                <span>Publicar ahora</span>
            </label>
        </template>
        <template v-else>
            <div class="inline-flex items-center gap-2 text-sm">
                <span class="px-2 py-0.5 rounded bg-black text-white">Ya publicado</span>
                <span class="text-gray-500">(no editable)</span>
            </div>
        </template>

        <div class="flex gap-2 pt-2">
            <button :disabled="pending" class="bg-black text-white rounded px-4 py-2">
                {{ pending ? 'Guardando…' : submitText }}
            </button>
            <button type="button" class="border rounded px-4 py-2" @click="$emit('cancel')">Cancelar</button>
        </div>
    </form>
</template>
