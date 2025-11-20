// Простая реализация toast-уведомлений
export const toast = {
  success(message) {
    console.log(`✅ ${message}`)
    // В реальном приложении здесь будет отображение уведомления
  },
  
  error(message) {
    console.log(`❌ ${message}`)
    // В реальном приложении здесь будет отображение уведомления
  },
  
  info(message) {
    console.log(`ℹ️ ${message}`)
    // В реальном приложении здесь будет отображение уведомления
  }
}