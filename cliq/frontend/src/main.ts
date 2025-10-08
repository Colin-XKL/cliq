import {createApp} from 'vue'
import App from './App.vue'
import './style.css';

// PrimeVue
import PrimeVue from 'primevue/config';
import 'primevue/resources/themes/lara-light-blue/theme.css';
import 'primeicons/primeicons.css';
import PrimeComponents from './components/PrimeComponents';

const app = createApp(App);
app.use(PrimeVue);
app.use(PrimeComponents);
app.mount('#app');
