declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<object, object, unknown>
  export default component
}

declare module '@/stores/auth' {
  import type { StoreDefinition } from 'pinia'
  export const useAuthStore: StoreDefinition
}

declare module '@/services/api' {
  import type { AxiosInstance } from 'axios'
  export const api: AxiosInstance
}
