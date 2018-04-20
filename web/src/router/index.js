import Vue from 'vue'
import Router from 'vue-router'
import Products from '@/components/products/Products'
import Product from '@/components/products/Product'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/product',
      name: 'product',
      component: Product
    },
    {
      path: '/',
      name: 'products',
      component: Products
    }
  ]
})
