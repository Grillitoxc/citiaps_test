
// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE,
    },
  },
  modules: [
    '@nuxtjs/tailwindcss'
  ],
  app: {
    pageTransition: { name: 'page', mode: 'out-in' }
  }
})
