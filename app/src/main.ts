import { createApp } from 'vue'
import { config } from 'md-editor-v3'
import { registerPlugins } from '@/plugins'

import 'unfonts.css'
import './styles/tailwind.css'
import './main.css'
import App from './App.vue'

// XSS攻撃を防止するため、md-editor-v3 の HTML タグの解釈を無効化する。
config({
  markdownItConfig: (md) => {
    md.options.html = false
  },
})

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
