<template>
  <div>
    <router-view></router-view>
    <v-table is-horizontal-resize style="width:90%" column-width-drag :columns="columns" :table-data="projects" row-hover-color="#eee" row-click-color="#edf7ff" @on-custom-comp="customCompFunc"></v-table>
  </div>
</template>
<script>
export default {

  name: 'Products',

  data() {
    return {
      projects: null,
      psLength: 0,
      columns: [{
          field: 'id',
          title: '编号',
          width: 50,
          titleAlign: 'center',
          columnAlign: 'center',
          formatter: function(rowData, rowIndex, pagingIndex, field) {
            return '<span style="color:red;font-weight: bold;">' + rowData.id + '</span>'
          },
        },
        { field: 'name', title: '项目名称', width: 200, titleAlign: 'center', columnAlign: 'center' },
        { field: 'host', title: '域名', width: 200, titleAlign: 'center', columnAlign: 'center' },
        { field: 'contact_name', title: '联系人', width: 100, titleAlign: 'center', columnAlign: 'center' },
        { field: 'contact_cellphone', title: '联系电话', width: 150, titleAlign: 'center', columnAlign: 'center' },
        { field: 'deposit', title: '订金', width: 100, titleAlign: 'center', columnAlign: 'center' },
        { field: 'brokerage', title: '违约金', width: 100, titleAlign: 'center', columnAlign: 'center' },
        { field: 'refund', title: '偿还佣金', width: 100, titleAlign: 'center', columnAlign: 'center' },
        { field: 'custome-adv', title: '操作', width: 200, titleAlign: 'center', columnAlign: 'center', componentName: 'table-operation', isResize: true }
      ]
    }
  },
  methods: {
    async fetchData() {
      this.isLoading = true;
      try {
        var response = await this.$http.get("http://localhost:8001/h/v1/projects")
        if (response.data.code === 0) {
          this.projects = response.data.data
          this.psLength = this.projects.length
        }
      } catch (e) {
        console.log(e);
      }
      this.isLoading = false;
    },
    customCompFunc(params) {}
  },
  created() {
    this.fetchData();
  }
}
import Vue from 'vue'

// 自定义列组件
Vue.component('table-operation', {
  // <router-link :to="{path:'/product/'+rowData.id}">编辑</router-link>
  template: `<span>
  <router-link :to="{path:'/product', query: {id: rowData.id}}">编辑</router-link>
        </span>`,
  props: {
    rowData: {
      type: Object
    },
    field: {
      type: String
    },
    index: {
      type: Number
    }
  },
  data() {
    return {

    }
  },
  methods: {
    update() {

      // 参数根据业务场景随意构造
      let params = { type: 'edit', index: this.index, rowData: this.rowData };
      this.$emit('on-custom-comp', params);
    }
  }
})

</script>
<style lang="css" scoped>


</style>
