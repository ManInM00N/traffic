<template>

    <body>
        
    </body>
</template>


<script setup>
import { ref, onMounted } from 'vue';
let isSliding = ref(false);
onMounted(() => {
    setTimeout(() => {
        isSliding.value = true;
    }, 1000); // 延迟 1 秒后开始动画
});
</script>


<style scoped>
.headbar {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 250px;
    overflow: hidden;
}

.headbar>img {
    position: absolute;
    top: -40px;
    left: 0;
    width: 100%;
    height: 400px;
    object-fit: fill;
    object-position: center center;
    z-index: -2000;
}

.headbar .mask {
    position: absolute;
    left: 0;
    bottom: 0;
    width: 100%;
    height: 200px;
    background-color: black;
    background: linear-gradient(to top, rgba(13, 13, 13, 0.6) 25%, transparent 100%);
    z-index: -1000;
}

.headbar .navbar {
    position: absolute;
    bottom: 0px;
    right: 100%;
    transform: translateY(-50%);
    width: 100%;
    height: 90px;
    background-color: transparent;
    transition: all 1s ease-in-out;
    z-index: 1000;
}

.headbar .navbar.slide {
    right: 0;
    /* 滑动到左侧 */
}


.headbar .navbar svg {
    position: absolute;
    right: 5px;
}



/* nav-content */
.headbar .nav-content {
    position: absolute;
    bottom: 0px;
    width: 100%;
    height: 85px;
    align-content: center;
    display: flex;
    z-index: 2000;
}


.headbar .nav-content img {
    width: 50px;
    height: 50px;
    border-radius: 50%;
}

.headbar .nav-content-left {
    position: relative;
    left: 0;
    width: 200px;
    height: 70px;
    display: flex;
    align-content: center;
    z-index: 4000;
    background-color: black;
}


.headbar .nav-content-left img {
    position: relative;
    top: 17px;
    left: 10px;
}


.headbar .nav-content-left .nav-content-rl {
    color: #fff;
    position: relative;
    left: 50px;
    top: 30px;

}



.headbar .nav-content-center {
    background-color: black;
    position: relative;
    left: 950px;
    width: 200px;
    align-content: center;
}


.headbar .nav-content-center span {
    color: #fff;
}



.headbar .nav-content-right {
    background-color: black;
    position: relative;
    left: 1680px;
    width: 220px;
    display: flex;
    align-content: center;
}


.headbar .nav-content-right span {
    color: #fff;
    position: relative;
    left: 100px;
    top: 30px;
}


.headbar .nav-content-right img {
    position: relative;
    left: 130px;
    top: 17px;
}
</style>