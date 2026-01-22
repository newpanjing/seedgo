import {createApp} from 'vue'
import {createPinia} from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'
import {hasPerm, permission} from './directives/permission'
import {superDirective} from './directives/super'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.directive('permission', permission)
app.directive('super', superDirective)
app.directive('has', hasPerm)

app.mount('#app')
