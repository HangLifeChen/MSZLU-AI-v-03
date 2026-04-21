# 增长趋势统计接口文档

## 概述

提供用户、知识库、智能体三个维度的月度累计增长趋势数据，用于前端折线图展示。

## 接口列表

### 1. 用户增长趋势

- **URL:** `GET /api/v1/admin/statistics/users/growth`
- **参数:**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| year   | int  | 否   | 年份，默认当前年份 |

- **响应示例:**
```json
{
  "code": 200,
  "data": {
    "months": ["1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"],
    "data": [120, 180, 250, 310, 420, 530, 610, 720, 850, 960, 1100, 1280]
  }
}
```

### 2. 知识库增长趋势

- **URL:** `GET /api/v1/admin/statistics/knowledge-bases/growth`
- **参数:**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| year   | int  | 否   | 年份，默认当前年份 |

- **响应示例:**
```json
{
  "code": 200,
  "data": {
    "months": ["1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"],
    "data": [10, 18, 30, 45, 52, 68, 80, 95, 110, 130, 145, 160]
  }
}
```

### 3. 智能体增长趋势

- **URL:** `GET /api/v1/admin/statistics/agents/growth`
- **参数:**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| year   | int  | 否   | 年份，默认当前年份 |

- **响应示例:**
```json
{
  "code": 200,
  "data": {
    "months": ["1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"],
    "data": [5, 12, 20, 35, 48, 60, 75, 90, 105, 120, 140, 165]
  }
}
```

## 响应字段说明

| 字段    | 类型     | 说明 |
|---------|----------|------|
| months  | string[] | 固定12个月标签 |
| data    | int64[]  | 截至每月末的累计数量，长度固定12。未来月份会沿用上一个有数据的月份的值 |

## 前端对接示例

```js
const initLineChart = async () => {
  if (!lineChartRef.value) return
  lineChart = echarts.init(lineChartRef.value)

  const [usersRes, kbRes, agentsRes] = await Promise.all([
    fetch('/api/v1/admin/statistics/users/growth?year=2026').then(r => r.json()),
    fetch('/api/v1/admin/statistics/knowledge-bases/growth?year=2026').then(r => r.json()),
    fetch('/api/v1/admin/statistics/agents/growth?year=2026').then(r => r.json()),
  ])

  const months = usersRes.data.months

  lineChart.setOption({
    title: { text: '增长趋势', left: 'center', textStyle: { fontSize: 14 } },
    tooltip: { trigger: 'axis' },
    legend: { bottom: 0, data: ['用户', '知识库', '智能体'] },
    xAxis: { type: 'category', data: months, boundaryGap: false },
    yAxis: { type: 'value' },
    grid: { top: 40, bottom: 40, left: 40, right: 20 },
    series: [
      {
        name: '用户',
        type: 'line',
        smooth: true,
        data: usersRes.data.data,
        areaStyle: { opacity: 0.1 },
        itemStyle: { color: '#409EFF' },
      },
      {
        name: '知识库',
        type: 'line',
        smooth: true,
        data: kbRes.data.data,
        areaStyle: { opacity: 0.1 },
        itemStyle: { color: '#67C23A' },
      },
      {
        name: '智能体',
        type: 'line',
        smooth: true,
        data: agentsRes.data.data,
        areaStyle: { opacity: 0.1 },
        itemStyle: { color: '#E6A23C' },
      },
    ],
  })
}
```
