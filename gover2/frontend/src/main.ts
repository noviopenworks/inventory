import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from '@/router'
import App from '@/App.vue'
import { useUiStore } from '@/stores/ui'
import '@/style.css'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

// Apply the persisted theme/density before the first paint so a dark-mode user
// doesn't see a light-themed flash on startup. loadFromConfig swallows bridge
// errors (returns defaults), so mount always proceeds.
const ui = useUiStore(pinia)
ui.loadFromConfig().finally(() => app.mount('#app'))
