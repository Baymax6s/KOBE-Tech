import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'

import 'unfonts.css'
import './styles/tailwind.css'
import './main.css'
import App from './App.vue'

const app = createApp(App)
registerPlugins(app)
app.mount('#app')
