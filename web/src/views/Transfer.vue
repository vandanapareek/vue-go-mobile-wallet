<template>
    <div id="overlay" :class="loading == true ? 'overlay_active' : 'overlay_inactive'"></div>
    <v-card class="text-center">
      <v-card-text>
        <v-row>
          <v-col cols="10" class="header-title">
            <h1 class="display-1">            
              <img src="../assets/wallet_icon.png" width="40" height="40" class="header-img"/>
              Send Money
            </h1>
          </v-col>
          <v-col cols="2" class="header-btn">
            <v-btn class="btn-warning" color="warning" type="warning" @click="logout()">
              Log Out
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

      <div class="wrapper" >
        <div class="card-form">
          <div class="card-list">
            <div class="card-item" >
              <div class="card-item__side -front">
                <div class="card-item__focus"  ref="focusElement"></div>
                <div class="card-item__cover">
                  <img
                  v-bind:src="'https://raw.githubusercontent.com/muhammederdem/credit-card-form/master/src/assets/images/' + currentCardBackground + '.jpeg'" class="card-item__bg">
                </div>
                
                <div class="card-item__wrapper">
                  <div class="card-item__top">
                    <img class="avatar_img" :src="getPic()" alt="Avatar" style="width:80px;height: 80px;">
                      <label for="cardName" class="card-item__info" ref="cardName">
                          <div class="card-item__name"  key="2"><span class="card-item__sub">Hi,</span>{{ username }}</div>
                      </label>

                    <div class="card-item__type">
                      <img src="../assets/chip.png" class="card-item__chip"> 
                    </div>
                  </div>
                  <label for="cardNumber" class="card-item__balance" ref="cardNumber">
                    <div class="card-item__holder balance_title">Total Balance</div>
                    <div class="card-item__name" >{{ currency }} {{ balance }}</div>
                  </label>

                </div>
              </div>
        
            </div>
          </div>
          <div class="card-form__inner">
            <div class="card-input">
              <label for="to_account_id" class="card-input__label">Send To</label>
                  <select class="card-input__input -select" id="to_account_id" v-model="to_account_id"  data-ref="cardDate">
                    <option value="" disabled selected>Choose Beneficiary</option>
                    <option v-for="item in beneficiaries" :key='item.account_id' v-bind:value="item.account_id">
                      {{item.username}}
                    </option>
                  </select>
            </div>

            <div class="card-input">
              <label for="amount" class="card-input__label">Amount (SGD)</label>
              <input type="text" id="amount" class="card-input__input" v-model="amount" autocomplete="off">
            </div>
            <div v-if="loading" id="loading">
              <img id="loading-image" src="../assets/send.gif" alt="Loading..." />
            </div>
            <button class="card-form__button" @click="sendMoney">
              Send
            </button>
          </div>
        </div> 
      </div>
    <v-snackbar v-model="snackbar"  :timeout="timeout" :color="color" location="top">
          {{ snackbarText }}
    </v-snackbar>
  </template>

<script>
import axios from 'axios';
export default {
   
  el: "#app",
  data() {
    return {
      username:localStorage.getItem('user_name'),
      balance:localStorage.getItem('user_balance'),
      currency:localStorage.getItem('user_currency'),
      beneficiaries:[],
      currentCardBackground: Math.floor(Math.random()* 25 + 1),
      to_account_id: "",
      amount:0,
      loading: false,
      snackbar: false,
      snackbarText: "",
      timeout: 3000,
      color: "",
    };
  },
  mounted() {
    const config = {
      headers: { Authorization: `Bearer `+localStorage.getItem("thunes_token") }
    };
    axios.get("/beneficiaries",config).then(response => {
           this.beneficiaries = response.data;
           console.log(this.beneficiaries)
        }).catch((error) => {
          console.log(error)
          if((error.code == "ERR_NETWORK") || (error.response && error.response.status == 401)) {
            this.logout()
          }
        })
  },
  methods: {
    sendMoney() {
      if (!isNaN(this.amount) && this.amount.toString().indexOf('.') != -1)
      {
        this.snackbar = true;
        this.color = 'red';
        this.snackbarText = "Decimal Places are not permitted!"
      } else {
        var formData = {
            to_account_id: parseInt(this.to_account_id),
            amount: parseInt(this.amount),
            currency: "SGD"
        }
        const config = {
        headers: { Authorization: `Bearer `+localStorage.getItem("thunes_token"), 'Content-Type': 'application/json'  }
        };
        axios.post("/transfer", JSON.stringify(formData) ,config).then(response => {
              this.isActive = !this.isActive;
              this.loading = true;
              this.timer = setInterval(() => {
                //update balance in localstorage
                this.loading=false;
                this.balance = response.data.Details.balance;
                localStorage.setItem("user_balance",this.balance)
                this.$router.go(0);
              }, 4000) 
                    
        }).catch((error) => {
            console.log(error)  
            if(error.response && error.response.data.Code == 422) {
              this.snackbar = true;
              this.color = 'red';
              this.snackbarText = error.response.data.Message
            } 
            if((error.code == "ERR_NETWORK") || (error.response &&  error.response.data.Code == 401)) {
              this.logout()
            } 
          })  
      }
    },
    getPic() {
       return 'src/assets/' + this.username + '.jpeg';
    },
    logout() {
      localStorage.clear();
      this.$router.go(0);
    }
  }
};
</script>
