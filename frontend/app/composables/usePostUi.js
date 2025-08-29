// app/composables/usePostUI.js
export const usePostUI = () => {
  const showCreate = useState('post.showCreate', () => false)

  const openCreate  = () => { showCreate.value = true }
  const closeCreate = () => { showCreate.value = false }

  return { showCreate, openCreate, closeCreate }
}
