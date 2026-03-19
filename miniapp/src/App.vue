<script setup lang="ts">
import { onLaunch, onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/stores/user'

onLaunch(() => {
  console.log('T-Line App Launch')
  const userStore = useUserStore()
  userStore.autoLogin()
})

onShow(() => {
  console.log('T-Line App Show')
})
</script>

<style lang="scss">
@import './uni.scss';

page {
  background-color: #f5f6f8;
  font-family: -apple-system, 'PingFang SC', 'Microsoft YaHei', sans-serif;
  font-size: 28rpx;
  color: #333;
}
</style>
