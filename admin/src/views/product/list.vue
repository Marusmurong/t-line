<template>
  <div class="product-list">
    <!-- 顶部标题栏 -->
    <div class="page-header">
      <h2 class="page-title">商品管理</h2>
      <a-button type="primary">
        <template #icon><icon-plus /></template>
        新增商品
      </a-button>
    </div>

    <!-- 筛选栏 -->
    <a-card class="filter-card">
      <a-row :gutter="16" align="center">
        <a-col :span="6">
          <a-select
            v-model="filterCategory"
            placeholder="商品分类"
            allow-clear
            style="width: 100%"
          >
            <a-option value="course">课程</a-option>
            <a-option value="equipment">球具</a-option>
            <a-option value="service">服务</a-option>
          </a-select>
        </a-col>
        <a-col :span="6">
          <a-select
            v-model="filterStatus"
            placeholder="商品状态"
            allow-clear
            style="width: 100%"
          >
            <a-option value="on">上架</a-option>
            <a-option value="off">下架</a-option>
          </a-select>
        </a-col>
        <a-col :span="8">
          <a-input
            v-model="searchKeyword"
            placeholder="搜索商品名称..."
            allow-clear
          >
            <template #prefix><icon-search /></template>
          </a-input>
        </a-col>
      </a-row>
    </a-card>

    <!-- 商品表格 -->
    <a-card style="margin-top: 16px">
      <a-table :data="filteredProducts" :pagination="false" stripe>
        <template #columns>
          <a-table-column title="商品图" :width="80">
            <template #cell>
              <div class="product-thumb" />
            </template>
          </a-table-column>
          <a-table-column title="商品名称" data-index="name">
            <template #cell="{ record }">
              <div class="product-name">
                <span class="name-text">{{ record.name }}</span>
                <span class="name-sub">{{ record.subCategory }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="分类" data-index="category" :width="100">
            <template #cell="{ record }">
              <a-tag :color="categoryColorMap[record.category]">
                {{ record.category }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="价格" data-index="price" :width="150">
            <template #cell="{ record }">
              <div class="price-cell">
                <span v-if="record.originalPrice" class="original-price">
                  ¥{{ record.originalPrice }}
                </span>
                <span class="current-price">¥{{ record.price }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="库存" data-index="stock" :width="80" />
          <a-table-column title="销量" data-index="sales" :width="80" />
          <a-table-column title="状态" data-index="status" :width="100">
            <template #cell="{ record }">
              <a-tag :color="record.status === '上架' ? 'green' : 'gray'">
                {{ record.status }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="操作" :width="160">
            <template #cell="{ record }">
              <a-button type="text" size="small">编辑</a-button>
              <a-button
                type="text"
                size="small"
                :status="record.status === '下架' ? 'success' : 'warning'"
              >
                {{ record.status === '下架' ? '上架' : '下架' }}
              </a-button>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <a-pagination
          v-model:current="currentPage"
          :total="products.length"
          :page-size="10"
          show-total
        />
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { IconPlus, IconSearch } from '@arco-design/web-vue/es/icon'

const filterCategory = ref<string | undefined>(undefined)
const filterStatus = ref<string | undefined>(undefined)
const searchKeyword = ref('')
const currentPage = ref(1)

const categoryColorMap: Record<string, string> = {
  '课程': 'blue',
  '球具': 'orange',
  '服务': 'green',
}

const categoryFilterMap: Record<string, string> = {
  course: '课程',
  equipment: '球具',
  service: '服务',
}

const statusFilterMap: Record<string, string> = {
  on: '上架',
  off: '下架',
}

interface Product {
  readonly id: number
  readonly name: string
  readonly subCategory: string
  readonly category: string
  readonly price: number
  readonly originalPrice: number | null
  readonly stock: number
  readonly sales: number
  readonly status: string
}

const products = ref<ReadonlyArray<Product>>([
  { id: 1, name: '网球私教课', subCategory: '1对1教学', category: '课程', price: 380, originalPrice: 480, stock: 20, sales: 156, status: '上架' },
  { id: 2, name: '网球小组课', subCategory: '4人小班', category: '课程', price: 180, originalPrice: 220, stock: 40, sales: 328, status: '上架' },
  { id: 3, name: 'Wilson Pro Staff', subCategory: '专业球拍', category: '球具', price: 1580, originalPrice: 1880, stock: 15, sales: 42, status: '上架' },
  { id: 4, name: 'HEAD Speed MP', subCategory: '进阶球拍', category: '球具', price: 1280, originalPrice: null, stock: 8, sales: 23, status: '上架' },
  { id: 5, name: '穿线服务', subCategory: '专业穿线', category: '服务', price: 60, originalPrice: null, stock: 999, sales: 512, status: '上架' },
  { id: 6, name: '球拍保养套餐', subCategory: '清洁+换胶', category: '服务', price: 120, originalPrice: 150, stock: 999, sales: 89, status: '下架' },
  { id: 7, name: '网球训练营', subCategory: '暑期集训', category: '课程', price: 2980, originalPrice: 3580, stock: 30, sales: 18, status: '下架' },
  { id: 8, name: 'Babolat Pure Aero', subCategory: '底线利器', category: '球具', price: 1480, originalPrice: 1680, stock: 0, sales: 67, status: '上架' },
])

const filteredProducts = computed(() => {
  return products.value.filter((p) => {
    if (filterCategory.value && p.category !== categoryFilterMap[filterCategory.value]) return false
    if (filterStatus.value && p.status !== statusFilterMap[filterStatus.value]) return false
    if (searchKeyword.value && !p.name.includes(searchKeyword.value)) return false
    return true
  })
})
</script>

<style lang="scss" scoped>
.product-list {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .page-title {
    font-size: 20px;
    font-weight: 600;
    color: #1d2129;
    margin: 0;
  }

  .filter-card {
    :deep(.arco-card-body) {
      padding: 16px;
    }
  }

  .product-thumb {
    width: 48px;
    height: 48px;
    border-radius: 6px;
    background: linear-gradient(135deg, #e8f4f8, #d1ecf1);
  }

  .product-name {
    display: flex;
    flex-direction: column;

    .name-text {
      font-weight: 500;
      color: #1d2129;
    }

    .name-sub {
      font-size: 12px;
      color: #999;
      margin-top: 2px;
    }
  }

  .price-cell {
    display: flex;
    flex-direction: column;

    .original-price {
      font-size: 12px;
      color: #999;
      text-decoration: line-through;
    }

    .current-price {
      font-weight: 600;
      color: #1d2129;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
  }
}
</style>
