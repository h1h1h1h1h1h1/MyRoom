/**
 * 格式化工具函数
 */

/**
 * 格式化货币金额
 * @param amount 金额（单位：元）
 * @returns 格式化后的货币字符串
 */
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount);
}

/**
 * 格式化日期
 * @param date 日期字符串或Date对象
 * @returns 格式化后的日期字符串
 */
export function formatDate(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }).format(d);
}

/**
 * 格式化日期（仅日期部分）
 * @param date 日期字符串或Date对象
 * @returns 格式化后的日期字符串
 */
export function formatDateOnly(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  }).format(d);
}

/**
 * 格式化用电量（千瓦时）
 * @param usage 用电量（千瓦时）
 * @returns 格式化后的用电量字符串
 */
export function formatElectricityUsage(usage: number): string {
  if (usage >= 1000) {
    return `${(usage / 1000).toFixed(2)} 千千瓦时`;
  }
  return `${usage.toFixed(2)} 千瓦时`;
}

/**
 * 格式化电费（元）
 * @param amount 电费金额（元）
 * @returns 格式化后的电费字符串
 */
export function formatElectricityBill(amount: number): string {
  return formatCurrency(amount);
}

/**
 * 格式化手机号码
 * @param phone 手机号码
 * @returns 格式化后的手机号码
 */
export function formatPhone(phone: string): string {
  if (phone.length === 11) {
    return `${phone.slice(0, 3)}-${phone.slice(3, 7)}-${phone.slice(7)}`;
  }
  return phone;
}

/**
 * 格式化客户编号
 * @param customerNumber 客户编号
 * @returns 格式化后的客户编号
 */
export function formatCustomerNumber(customerNumber: string): string {
  if (customerNumber.length >= 8) {
    return `${customerNumber.slice(0, 4)}-${customerNumber.slice(4, 8)}${customerNumber.length > 8 ? '-' + customerNumber.slice(8) : ''}`;
  }
  return customerNumber;
}

/**
 * 格式化状态文本
 * @param status 状态代码
 * @returns 状态文本
 */
export function formatStatus(status: string): string {
  const statusMap: Record<string, string> = {
    'pending': '待处理',
    'processing': '处理中',
    'completed': '已完成',
    'cancelled': '已取消',
    'paid': '已支付',
    'unpaid': '未支付',
    'active': '正常',
    'inactive': '停用',
    'warning': '余额不足',
    'critical': '即将断电',
    'normal': '正常'
  };

  return statusMap[status] || status;
}

/**
 * 格式化状态颜色
 * @param status 状态代码
 * @returns 状态对应的CSS颜色类名
 */
export function formatStatusColor(status: string): string {
  const colorMap: Record<string, string> = {
    'pending': 'text-yellow-600',
    'processing': 'text-blue-600',
    'completed': 'text-green-600',
    'cancelled': 'text-red-600',
    'paid': 'text-green-600',
    'unpaid': 'text-red-600',
    'active': 'text-green-600',
    'inactive': 'text-gray-600',
    'warning': 'text-orange-600',
    'critical': 'text-red-600',
    'normal': 'text-green-600'
  };

  return colorMap[status] || 'text-gray-600';
}

/**
 * 格式化支付方式
 * @param paymentMethod 支付方式代码
 * @returns 支付方式文本
 */
export function formatPaymentMethod(paymentMethod: string): string {
  const methodMap: Record<string, string> = {
    'alipay': '支付宝',
    'wechat': '微信支付',
    'bank_card': '银行卡支付',
    'balance': '余额支付',
    'cash': '现金支付'
  };

  return methodMap[paymentMethod] || paymentMethod;
}

/**
 * 格式化申请类型
 * @param applicationType 申请类型代码
 * @returns 申请类型文本
 */
export function formatApplicationType(applicationType: string): string {
  const typeMap: Record<string, string> = {
    'new_installation': '新装用户申请',
    'name_change': '更名过户申请',
    'meter_check': '电能表校验',
    'address_change': '地址变更',
    'tariff_change': '电价变更'
  };

  return typeMap[applicationType] || applicationType;
}
