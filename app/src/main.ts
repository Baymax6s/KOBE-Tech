import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'

import 'unfonts.css'
import 'highlight.js/styles/github.css'
import './styles/tailwind.css'
import './main.css'
import App from './App.vue'

async function main() {
  const useMSW = import.meta.env.VITE_USE_MSW === 'true'

  if (import.meta.env.DEV && useMSW) {
    const { worker } = await import('./mocks/browser')

    await worker.start({
      onUnhandledRequest: 'bypass',
    })
  }

  const app = createApp(App)
  registerPlugins(app)
  app.mount('#app')
}

main()
