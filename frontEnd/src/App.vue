<template>
  <v-app id="inspire">
    <v-navigation-drawer
    clipped
    fixed
    v-model="drawer"
    app
    >
    <v-list two-line class="grey darken-2">
      <v-subheader class="grey darken-3 white--text">Direct Reports</v-subheader>
      <template v-for="(emp,index) in emps">
        
      <v-list-tile avatar ripple class="grey darken-1" @click="GetEmployee(emp.id)" >
        <v-list-tile-content>
          <v-list-tile-title class="indigo--text text--darken-4" v-html="emp.name"></v-list-tile-title>
          <v-list-tile-sub-title class="yellow--text" v-html="emp.title"></v-list-tile-sub-title>
          <v-list-tile-sub-title class="white--text" v-html="emp.email"></v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
      <v-divider v-if="index + 1 < emps.length"></v-divider>
      </template>
    </v-list>
  </v-navigation-drawer>
  <v-toolbar app fixed clipped-left>
    <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
    <img src="static/imgs/4lc.svg" style="width:2%;">
    <v-toolbar-title>Forest</v-toolbar-title>
  </v-toolbar>
  <v-content>
    <v-container fluid>
      <v-layout justify-center align-center>
        <v-flex xs12>
          <!-- dont-fill-mask-blanks -->
          <!-- label="Search Employees" -->
          <v-select
          :filter="customFilter"
          placeholder="First or Family name"
          autocomplete
          :async-loading="loading"
          :items="items"
          solo
          :search-input.sync="search"
          v-model="select"
          color="light-green accent-3"
          dense
          @input="OnChange"
          :rules="[() => select.length > 0 || 'You must choose at least one']"
          ></v-select>
          <!-- <vheatmap v-if="yearcounts.year!=null" -->
            <!-- :CalendarData="yearcounts" -->
            <!-- > </vheatmap>  -->
            <!-- {{m.title}} -->
            <teamHeatMap v-if="!!TeamLeaves" :CalendarData="TeamLeaves"  :LeaveTypes="leaveTypeNames"></teamHeatMap>


            <vheatmap v-for="m in yearsLeave" :key="m.id" v-bind:data="m" :CalendarData="m"  :LeaveTypes="leaveTypeNames">        </vheatmap>

          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
    <v-footer app fixed>
      <span>&copy; ATOP Workforce Analytics</span>
    </v-footer>
  </v-app>
</template>

<script>
import axios from 'axios';

import vheatmap from './components/vheatmap.vue';
import teamHeatMap from './components/teamHeatMap.vue';

export default {
  data: () => ({
    drawer: false,
    search: null,
    items: [],
    select: [],
    loading: false,
    yearcounts : {days: {}, year: null, title:"Not Data Available"},
    yearsLeave : [],
    emps : [],
    leaveTypeNames : [],
    TeamLeaves : null,
    customFilter (item, queryText, itemText) {
      return true;
    }
  }),
  watch: {
    search (val) {
      val && this.querySelections(val)
    }
  },
  methods: {
    getLeaveTypes(){
      this.leaveTypeNames = [];
      axios.get("api/leaves")
        .then((response)  =>  {
            this.leaveTypeNames = response.data;
            console.log(this.leaveTypeNames);
          }, (error)  =>  {
            console.log(error);
          });
    },
    getReports(id){
      this.emps = [];
      axios.get("api/emp/"+ id)
        .then((response)  =>  {
            this.search = "";
            this.emps = response.data;
            if (this.emps.length == 0){
              this.drawer = false;
            }else{
              this.drawer = true;
            }
            console.log(this.emps);
          }, (error)  =>  {
            console.log(error);
          });
    },
    getLeaves(id){
      this.loading = true;
      this.yearsLeave = [];
      axios.get("api/list/"+ id + "/5")
        .then((response)  =>  {
          this.items =[];
          this.search = "";
          this.loading = false;
          this.yearsLeave = response.data;
        }, (error)  =>  {
             console.log(error);
        });
    },
    getTeam(id){
      this.TeamLeaves = null;
        axios.get("api/team/"+ id + "/5")
        .then((response)  =>  {
          this.TeamLeaves = response.data;
          console.log("==----===>>", this.TeamLeaves.year);
        }, (error)  =>  {
             console.log(error);
        });
    },
    GetEmployee(id){
        this.getLeaves(id);
        this.getTeam(id);
        this.getReports(id);
    },
    OnChange(){
      var id = this.select;
      this.select = "";
      // this.loading = true;
      this.search = "";
      // this.yearsLeave = [];
      this.GetEmployee(id);
      
    },
    querySelections (v) {
      if (!v || v.length < 3 ){
        this.items = [];
        return;
      }
      axios.get("api/search?query="+v)
      .then((response)  =>  {
        this.items = response.data;
      }, (error)  =>  {
        console.log(error);
      })

    },
  },
  created: function(){
      this.getLeaveTypes();
  },
  components: {
    'vheatmap': vheatmap,
    'teamHeatMap': teamHeatMap,
  },
}
</script>
