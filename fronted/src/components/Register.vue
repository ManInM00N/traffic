<template>
    <div class="backhome">
        <a href="/" style="color: black;">首页</a>
    </div>
    <div class="container" :class="{ 'sign-up-mode': isSignUp }">
        <div class="form-warp">
            <form class="sign-in-form" v-if="!isSignUp">
                <span class="form-title">管理员注册</span>
                <input placeholder="email" v-model="email"/>
                <input type="password" placeholder="密码" v-model="password" />
              <input type="password" placeholder="确认密码" v-model="password1" />
              <button class="submit-btn" @click.prevent="getyzm">发送验证码</button>
                <input type="text" placeholder="验证码" v-model="yzm" />
                <button class="submit-btn" @click.prevent="login">立即注册</button>
              <a href="/login" style="text-decoration:none;" >已有账号?</a>
            </form>
            <form class="sign-up-form" v-if="isSignUp">
                <span class="form-title">司机注册</span>
                <input placeholder="email" v-model="email" />
                <input type="password" placeholder="密码" v-model="password"/>
              <input type="password" placeholder="确认密码" v-model="password1" />
              <button class="submit-btn" @click.prevent="getyzm">发送验证码</button>
                <input type="text" placeholder="验证码" v-model="yzm" />
                <button class="submit-btn" @click.prevent="register">立即注册</button>
                <a href="/login" >已有账号?</a>
            </form>
        </div>
        <div class="desc-warp">

            <div class="desc-warp-item sign-up-desc" >
                <button @click="isSignUp = true">我是司机</button>
                <img src="../assets/administrators.svg" alt="">
            </div>
            <div class="desc-warp-item sign-in-desc">
                <button @click="isSignUp = false">我是管理员</button>
                <img src="../assets/subway.svg" alt="">
            </div>
        </div>
    </div>
</template>

<script>
import { ref } from 'vue';
import axios, {all} from 'axios';
import router from "@/router/index.js";

export default {
    name: 'AuthForm',
    setup() {
      const email = ref('');
        const password = ref('');
      const yzm= ref('');
        const isSignUp = ref(false);
        const name="langgod"
        const password1=ref('')
        const getyzm=()=>{

            axios.post("goapi/api/email_verification_code",{"email":email.value,"use":"register"})
            .then((response) =>{
              console.log(response.data)
               if(response.data.code==0){
                  alert(response.data.error)
               }else{
                    alert("验证码已发送")
              // localStorage.setItem("Authorization", response.data.token);

               }
           }
            )



        }
        const login = () => {
            // console.log('登录');
            axios.post("goapi/api/register",{"email":email.value,"password":password.value,"name":name.value,"again_password":password1.value,"usertype":'b',"code":yzm.value})
          .then((response) =>{
            console.log(response.data)
            if(response.data.code==0){
              alert(response.data.message)
            }else{
              alert("注册成功")
              localStorage.setItem("Authorization", response.data.token);
              router.push('/login').then(()=>{
              window.location.reload();
              });
            }
         })

        };


        const register = () => {
            axios.post("goapi/api/register",{"email":email.value,"password":password.value,"name":name.value,"again_password":password1.value,"usertype":'a',"code":yzm.value})
          .then((response) =>{
            console.log(response.data)
            if(response.data.code==0){
              alert(response.data.message)
            }else{
              alert("注册成功")
              localStorage.setItem("Authorization", response.data.token);
              router.push('/login').then(()=>{
              window.location.reload();
              });
            }
         })

        };

        return { isSignUp, login, register,email,password,name,password1,yzm,getyzm};
    },
};
</script>

<style scoped>
* {
    padding: 0;
    margin: 0;
    box-sizing: border-box;
    color: #333;
}
a{
    text-decoration: none;
    color: rgb(103, 103, 236);
}
.backhome{
    position:absolute;
    top:0px;
    left:0px;
    z-index: 100;
}
.backhome a{
    font-size: 30px;
    font-family: 楷体;
    text-decoration: none;
    
}
.container {
    position: relative;
    min-height: 100vh;
    width: 100%;
    overflow: hidden;
}

.container::before {
    content: " ";
    position: absolute;
    width: 2000px;
    height: 2000px;
    border-radius: 50%;
    background-image: linear-gradient(-45deg, #6266f5 0%, #04befe 100%);
    transition: 1.8s ease-in-out;
    z-index: 6;
    top: -10%;
    right: 48%;
    transform: translateY(-50%);
}

.container.sign-up-mode::before {
    transform: translate(100%, -50%);
}

.form-warp {
    width: 50%;
    position: absolute;
    z-index: 5;
    left: 60%;
    top: 50%;
    z-index: 5;
    transform: translate(-40%, -40%);
    display: grid;
    grid-template-columns: 1fr;
    transition: 1s 0.7s ease-in-out;
}

.form-warp form {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    gap: 20px;
    grid-row: 1 / 2;
    grid-column: 1 / 2;
    transition: all 0.2s 0.7s;
    opacity: 1;
    z-index: 4;
}

.form-title {
    color: #6266f5;
}

.form-warp .sign-up-form {
    opacity: 0;
    z-index: 3;
}

.container.sign-up-mode .form-warp {
    left: 25%;
}

.container.sign-up-mode .sign-in-form {
    opacity: 0;
    z-index: 3;
}

.container.sign-up-mode .sign-up-form {
    opacity: 1;
    z-index: 4;
}

.container .form-warp span{
    font-size:50px;
}




input,
.submit-btn {
    min-width: 300px;
    outline: none;
    padding: 12px 30px;
    line-height: 1;
    font-size: 16px;
    border-radius: 60px;
    color: #333;
    background-color: #6267f513;
    border: none;
}

input::placeholder {
    color: #cccc;
}


.container .form-warp input{
    width:500px;
    height: 60px;
    font-size:40px;
}

.submit-btn {
    background-color: #6266f5;
    color: #FFF;
    text-align: center;
    min-width: 150px;
    font-size: initial;
    font-weight: bold;
    letter-spacing: 1.5px;
    cursor: pointer;
}

.container .form-warp button{
    width:180px;
    height:60px;
    font-size:20px;
}

.desc-warp {
    width: 100%;
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
    display: grid;
    grid-template-columns: repeat(2, 1fr);
}

.desc-warp-item {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: space-around;
    text-align: center;
    text-align: center;
    padding: 3rem 17% 2rem 12%;
    z-index: 6;
}

.desc-warp-item button {

    font-size: 18px;

    margin-bottom: 20px;

}

.desc-warp-item button{
    width:200px;
    height: 60px;
    font-size:30px;
}



.sign-in-desc {
    pointer-events: none;
}

.sign-up-mode .sign-in-desc {
    pointer-events: all;
}

.sign-up-mode .sign-up-desc {
    pointer-events: none;
}

.content {
    width: 100%;
    transition: transform 0.9s ease-in-out;
    transition-delay: .6s;
}

.sign-in-desc img,
.sign-in-desc .content {
    transform: translateX(800px);
}

.sign-up-mode .sign-in-desc img,
.sign-up-mode .sign-in-desc .content {
    transform: translateX(0);
}

.sign-up-mode .sign-up-desc img,
.sign-up-mode .sign-up-desc .content {
    transform: translateX(-800px);
}

button {
    outline: none;
    padding: 6px 8px;
    min-width: 100px;
    text-align: center;
    border-radius: 30px;
    border: 2px solid #FFF;
    background: none;
    color: #FFF;
    cursor: pointer;
    transition: all .3s ease;
}

button:active {
    background: rgba(255, 255, 255, .1);
}

img {
    width: 100%;
    display: block;
    transition: transform 0.9s ease-in-out;
    transition-delay: .5s;
}


@media screen and (max-width: 870px) {
    .container::before {
        width: 1500px;
        height: 1500px;
        transform: translateX(-50%);
        left: 30%;
        bottom: 68%;
        right: initial;
        top: initial;
        transition: 2s ease-in-out;
    }

    .container.sign-up-mode::before {
        transform: translate(-50%, 100%);
        bottom: 32%;
        right: initial;
    }

    .form-warp {
        width: 100%;
        top: 75%;
        left: 50%;
        transform: translate(-50%, -100%);
        transition: 1s 0.8s ease-in-out;
    }

    .container.sign-up-mode .form-warp {
        top: 25%;
        left: 50%;
        transform: translate(-50%, 0);
    }

    img {
        width: 200px;
        transition: transform 0.9s ease-in-out;
        transition-delay: 0.7s;
    }

    .desc-warp {
        grid-template-columns: 1fr;
        grid-template-rows: 1fr 2fr 1fr;
    }

    .desc-warp-item {
        flex-direction: row;
        justify-content: space-around;
        align-items: center;
        padding: 2.5rem 8%;
        grid-column: 1 / 2;
    }

    .sign-in-desc {
        grid-row: 3 / 4;
    }

    .sign-in-desc img,
    .sign-in-desc .content {
        transform: translateY(800px);
    }

    .sign-up-mode .sign-in-desc img,
    .sign-up-mode .sign-in-desc .content {
        transform: translateY(0);
    }

    .sign-up-mode .sign-up-desc img,
    .sign-up-mode .sign-up-desc .content {
        transform: translateY(-800px);
    }
}

@media screen and (max-width: 570px) {
    .container::before {
        bottom: 72%;
        left: 50%;
    }

    img {
        display: none;
    }
}
</style>