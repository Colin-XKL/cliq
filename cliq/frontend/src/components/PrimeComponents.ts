import { App } from 'vue';
import Button from 'primevue/button';
import InputText from 'primevue/inputtext';
import Card from 'primevue/card';
import Divider from 'primevue/divider';
import FileUpload from 'primevue/fileupload';
import Dropdown from 'primevue/dropdown';
import Checkbox from 'primevue/checkbox';
import Toast from 'primevue/toast';
import ToastService from 'primevue/toastservice';

export default {
  install: (app: App) => {
    app.use(ToastService);
    app.component('PButton', Button);
    app.component('PInputText', InputText);
    app.component('PCard', Card);
    app.component('PDivider', Divider);
    app.component('PFileUpload', FileUpload);
    app.component('PDropdown', Dropdown);
    app.component('PCheckbox', Checkbox);
    app.component('Toast', Toast);
  }
};