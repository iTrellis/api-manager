<template>
  <div class="panel-body">
    <vue-form-generator :schema="schema" :model="item" :options="formOptions"></vue-form-generator>
  </div>
</template>
<script>
import VueFormGenerator from "vue-form-generator";
import "vue-form-generator/dist/vfg.css";

export default {
  name: 'Product',

  components: {
    "vue-form-generator": VueFormGenerator.component
  },
  data() {
    return {
      item: null,
      selected: '',
      statusList: [],
      schema: {},
      formOptions: {
        validateAfterLoad: true,
        validateAfterChanged: true
      }
    }
  },
  created() {
    this.fetchData();
    this.id = this.$route.query.id
    if (this.id = 0) {
      return
    }
  },

  methods: {
    fetchData() {
      try {
        this.$http.get('http://localhost:8001/h/v1/project/status').then((response) => {
          console.log(response)
          if (response.data.code === 0) {
            // this.statusList = response.data.data
            this.statusList = new Array(response.data.data.length)
            for (var i = response.data.data.length - 1; i >= 0; i--) {
              this.statusList[i] = {
                id: response.data.data[i].id,
                name: response.data.data[i].description
              }
            }
            console.log(this.statusList);

            this.schema = {
              fields: [{
                type: "input",
                inputType: "text",
                label: "ID (不可编辑)",
                model: "id",
                readonly: true,
                featured: false
                // disabled: true
              }, {
                type: "input",
                inputType: "text",
                label: "项目名称",
                model: "name",
                readonly: false,
                featured: true,
                required: true,
                disabled: false,
                placeholder: "项目名称",
                validator: VueFormGenerator.validators.string
              }, {
                type: "input",
                inputType: "text",
                label: "联系人",
                model: "contact_name",
                readonly: false,
                featured: true,
                required: true,
                disabled: false,
                required: true,
                placeholder: "联系人",
                validator: VueFormGenerator.validators.string
              }, {
                type: "input",
                inputType: "text",
                label: "联系电话",
                model: "contact_cellphone",
                readonly: false,
                featured: true,
                required: true,
                disabled: false,
                required: true,
                placeholder: "请填写联系电话",
                validator: VueFormGenerator.validators.string
              }, {
                type: "input",
                inputType: "text",
                label: "联系地址",
                model: "address",
                readonly: false,
                featured: true,
                required: true,
                disabled: false,
                required: false,
                placeholder: "联系地址",
                validator: VueFormGenerator.validators.string
              }, {
                type: "input",
                inputType: "text",
                label: "违约金",
                model: "brokerage",
                readonly: false,
                featured: true,
                required: true,
                disabled: false,
                required: false,
                placeholder: "0",
                validator: VueFormGenerator.validators.double
              }, {
                type: "input",
                inputType: "text",
                label: "订金",
                model: "deposit",
                readonly: false,
                featured: true,
                required: true,
                disabled: false,
                required: false,
                placeholder: "0",
                validator: VueFormGenerator.validators.double
              }, {
                type: "input",
                inputType: "text",
                label: "退款",
                model: "refund",
                readonly: false,
                featured: true,
                required: true,
                disabled: false,
                required: false,
                placeholder: "0",
                validator: VueFormGenerator.validators.double
              }, {
                type: "select",
                label: "状态",
                model: "status",
                required: true,
                values: this.statusList,
                default: this.selected,
                validator: VueFormGenerator.validators.required
              }]
            }
          }
        });
        this.$http.get('http://localhost:8001/h/v1/project?id=' + this.$route.query.id).then(
          (response) => {
            if (response.data.code === 0) {
              this.item = response.data.data
              this.selected = this.item.status
            }
          });
      } catch (e) {
        console.log(e);
      }
    }
  }
}

</script>
<style scoped>


</style>
