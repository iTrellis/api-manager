// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

Vue.config.productionTip = false


// 导入 table 组件 和分页组件
import {VTable,VPagination} from 'vue-easytable'
Vue.component(VTable.name, VTable)
Vue.component(VPagination.name, VPagination)

import axios from 'axios'

Vue.prototype.$http = axios

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
