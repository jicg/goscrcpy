<script lang="ts" setup>
import {reactive, ref, watch} from 'vue'
import * as runtime from '../wailsjs/runtime/runtime.js';
import {ConnPhone, InstallApk, Msg, UnInstallApk} from "../wailsjs/go/main/App";
import {Fold, Open, TurnOff, Close} from "@element-plus/icons-vue";

const LocalStorageWBit = "LocalStorageWBit"
const LocalStorageHBit = "LocalStorageHBit"
const LocalStorageIsTop = "LocalStorageIsTop"
const LocalStorageIsWalk = "LocalStorageIsWalk"
const alwaysTop = ref(localStorage.getItem(LocalStorageIsTop) == "true")
const alwaysWalk = ref(localStorage.getItem(LocalStorageIsWalk) == "true")
let canEdit = ref(true)
let showDetail = ref(true)

let isDoing = ref(false)

const connPhone = async () => {
  try {
    canEdit.value = false
    await ConnPhone(alwaysTop.value, alwaysWalk.value)
  } catch (e) {
    await Msg("操作失败", `${e}`)
  } finally {
    canEdit.value = true
  }
}

watch(alwaysTop, (c) => {
  localStorage.setItem(LocalStorageIsTop, "" + c)
})
watch(alwaysWalk, (c) => {
  localStorage.setItem(LocalStorageIsWalk, "" + c)
})

const minView = async () => {
  showDetail.value = false;
  await changeSize(140, 40)
}
let wBit = parseFloat(localStorage.getItem(LocalStorageWBit) ?? "0"),
    hBit = parseFloat(localStorage.getItem(LocalStorageHBit) ?? "0");
const changeSize = async (width: number, height: number) => {
  let size = await runtime.WindowGetSize()
  let point = await runtime.WindowGetPosition()
  let _ = await runtime.WindowSetSize(width, height)
  if (wBit == 0 || hBit == 0) {
    setTimeout(async () => {
      let newTrueSize = await runtime.WindowGetSize()
      if (wBit == 0) {
        wBit = newTrueSize.w / width;
        localStorage.setItem(LocalStorageWBit, `${wBit}`);
      }
      if (hBit == 0) {
        hBit = newTrueSize.h / height;
        localStorage.setItem(LocalStorageHBit, `${hBit}`);
      }
      await runtime.WindowSetPosition(point.x + size.w - width * wBit, point.y + size.h - height * hBit)
    }, 10)
  } else {
    console.log(wBit, hBit, point, {
      x: point.x + size.w - width * wBit, y: point.y + size.h - height * hBit
    })
    await runtime.WindowSetPosition(point.x + size.w - width * wBit, point.y + size.h - height * hBit)
  }


}

const bigView = async () => {
  showDetail.value = true;
  await changeSize(500, 120)
}

const closeWin = () => {
  runtime.LogInfo("closeWin")
  runtime.Quit()
}

const installApk = async () => {
  try {
    isDoing.value=true
    await InstallApk()
  } catch (e) {
    await Msg("操作失败", `${e}`)
  }finally {
    isDoing.value=false
  }
}

const unInstallApk =async () => {
  try {
    isDoing.value=true
    await UnInstallApk()
  } catch (e) {
    await Msg("操作失败", `${e}`)
  }finally {
    isDoing.value=false
  }
}


(async () => {
  await bigView();
  // await connPhone();
})()
</script>
<template>

  <el-container v-if="!showDetail" style="align-items: center;justify-content: center;height: 100%;">
    <el-space v-if="canEdit" fill>
      <el-button type="primary" text :icon="Open" style="flex: 1;" @click="connPhone()">
      </el-button>
    </el-space>
    <span v-if="canEdit" style="width: 20px"/>
    <el-space fill style="align-items: center;justify-content: center;height: 100%;">
      <el-button type="danger" text :icon="Fold" v-if="!showDetail" @click="bigView()">
      </el-button>
    </el-space>
  </el-container>

  <el-container class="detailContainer" v-if="showDetail" direction="vertical">
    <el-button :icon="Close" style="position: absolute;top: 20px;right: 30px;" type="danger" size="small"
               @click="closeWin()" round></el-button>
    <el-main>
      <el-row>
        <el-switch
            v-model="alwaysTop"
            :disabled="!canEdit"
            active-text="总是置顶"
        />
        <span style="width: 32px"/>
        <el-switch
            :disabled="!canEdit"
            v-model="alwaysWalk"
            active-text="总是唤醒"
        />
      </el-row>
    </el-main>
    <el-footer style="display: flex;">
      <el-space style="flex: 1" alignment="start">
        <el-button :disabled="isDoing" size="small" @click="installApk()" round>安装apk</el-button>
        <el-button :disabled="isDoing" size="small" @click="unInstallApk()" round>卸载apk</el-button>
      </el-space>
      <el-button type="primary" size="small" :disabled="!canEdit" @click="connPhone()" round>连接手机</el-button>
      <el-button type="danger" size="small" @click="minView()" round>隐藏</el-button>

    </el-footer>
  </el-container>
</template>

<style>
.detailContainer {
  overflow: hidden;
  padding: 10px;
}
</style>
