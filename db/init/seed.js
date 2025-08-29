// db/init/seed.js
// Ejecutado por /docker-entrypoint-initdb.d de la imagen oficial de Mongo
const db = connect("mongodb://root:rootpass@localhost:27017/blog?authSource=admin");

const posts = db.getCollection("posts");

// Limpia por si existían datos residuales en escenarios locales
posts.deleteMany({});

// Autores y tags base
const AUTHORS = ["Ana García", "Carlos López", "María Rodríguez", "Christopher Jackson"];
const TAGS = ["go", "mongodb", "api", "frontend", "vue", "nuxt", "performance"];

// Genera fechas (UTC)
const now = new Date();
function daysAgo(n) {
  const d = new Date(now);
  d.setUTCDate(d.getUTCDate() - n);
  return d;
}

// 12 documentos mixtos
const docs = [
  {
    title: "Introducción a MongoDB con Go",
    author: AUTHORS[0],
    content: "Cómo integrar MongoDB con aplicaciones Go…",
    tags: ["mongodb", "go", "backend"],
    published: true,
    publishedAt: daysAgo(220),
    createdAt: daysAgo(225),
  },
  {
    title: "Construyendo APIs REST con Echo",
    author: AUTHORS[1],
    content: "Crea APIs robustas con Echo y buenas prácticas.",
    tags: ["go", "api", "rest"],
    published: true,
    publishedAt: daysAgo(200),
    createdAt: daysAgo(205),
  },
  {
    title: "Frontend Moderno con Nuxt 3",
    author: AUTHORS[2],
    content: "Novedades de Nuxt 3 y cómo aprovecharlas.",
    tags: ["nuxt", "vue", "frontend"],
    published: false,
    createdAt: daysAgo(180),
  },
  {
    title: "Optimización de Consultas en MongoDB",
    author: AUTHORS[0],
    content: "Índices, proyecciones y agregaciones eficientes.",
    tags: ["mongodb", "performance", "database"],
    published: true,
    publishedAt: daysAgo(160),
    createdAt: daysAgo(165),
  },
  {
    title: "Buenas prácticas en Go",
    author: AUTHORS[1],
    content: "Patrones de diseño y organización de proyectos Go.",
    tags: ["go"],
    published: true,
    publishedAt: daysAgo(120),
    createdAt: daysAgo(125),
  },
  {
    title: "Arquitectura de una API limpia",
    author: AUTHORS[3],
    content: "Separación por capas, DTOs, servicios y controladores.",
    tags: ["api", "backend"],
    published: false,
    createdAt: daysAgo(110),
  },
  {
    title: "Vue 3: Composition API",
    author: AUTHORS[2],
    content: "Cómo migrar de Options API a Composition.",
    tags: ["vue", "frontend"],
    published: true,
    publishedAt: daysAgo(95),
    createdAt: daysAgo(100),
  },
  {
    title: "Nuxt 3 + Tailwind: UI productiva",
    author: AUTHORS[3],
    content: "Plantillas, layouts, components y slots.",
    tags: ["nuxt", "frontend"],
    published: false,
    createdAt: daysAgo(70),
  },
  {
    title: "Paginar y filtrar en MongoDB",
    author: AUTHORS[0],
    content: "Text search, sort y paginación con índices.",
    tags: ["mongodb", "api"],
    published: true,
    publishedAt: daysAgo(60),
    createdAt: daysAgo(65),
  },
  {
    title: "Optimización de rendimiento en Nuxt",
    author: AUTHORS[2],
    content: "Hydration, lazy routes, image optimization.",
    tags: ["nuxt", "performance", "frontend"],
    published: false,
    createdAt: daysAgo(40),
  },
  {
    title: "Mejores prácticas con Docker en Go",
    author: AUTHORS[1],
    content: "Multi-stage builds, minimal images y healthchecks.",
    tags: ["go", "api"],
    published: true,
    publishedAt: daysAgo(20),
    createdAt: daysAgo(25),
  },
  {
    title: "Agregaciones en MongoDB",
    author: AUTHORS[0],
    content: "Pipelines, $match, $group, $sort, $limit.",
    tags: ["mongodb", "performance"],
    published: true,
    publishedAt: daysAgo(5),
    createdAt: daysAgo(10),
  },
];

posts.insertMany(docs);

// Índices (tu backend también los crea, esto es por si consultas fuera)
posts.createIndex({ title: "text", content: "text" }, { name: "text_title_content" });
posts.createIndex({ published: 1, publishedAt: -1 }, { name: "idx_published_publishedAt" });

print(`✅ Seed listo: ${posts.countDocuments()} posts`);
