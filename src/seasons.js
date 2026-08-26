// ============================================================
// 季节限定:当季活动期间完成指定动作,点亮限定成就(过期不候)
// ============================================================
export const seasons = [
  {
    id: 'newyear',
    name: '新年季',
    months: [1, 2],
    icon: '🧧',
    achName: '开年第一爪',
    desc: '新年季(1-2月)期间盖爪签到',
  },
  {
    id: 'sakura',
    name: '樱花季',
    months: [3, 4],
    icon: '🌸',
    achName: '樱花巡礼',
    desc: '樱花季(3-4月)期间盖爪签到',
  },
  {
    id: 'summer',
    name: '夏日祭',
    months: [7, 8],
    icon: '🎆',
    achName: '夏日花火',
    desc: '夏日祭(7-8月)期间盖爪签到',
  },
  {
    id: 'autumn',
    name: '秋收祭',
    months: [9, 10, 11],
    icon: '🍂',
    achName: '秋收满仓',
    desc: '秋收祭(9-11月)期间盖爪签到',
  },
  {
    id: 'xmas',
    name: '圣诞季',
    months: [12],
    icon: '🎄',
    achName: '圣诞袜底',
    desc: '圣诞季(12月)期间盖爪签到',
  },
]

// 当前季节(没有匹配则返回 null)
export function currentSeason(date = new Date()) {
  const m = date.getMonth() + 1
  return seasons.find((s) => s.months.includes(m)) || null
}

// 点亮某季节的限定成就
export function unlockSeasonAch(season) {
  if (season) localStorage.setItem('catcafe_ach_season_' + season.id, '1')
}
