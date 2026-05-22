import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/idc-dashboard'
  },
  {
    path: '/login',
    name: 'LoginPage',
    component: () => import('@/views/LoginPage.vue')
  },
  {
    path: '/idc-dashboard',
    name: 'IDCDashboard',
    component: () => import('@/views/IDCDashboard.vue')
  },
  {
    path: '/teacher-console',
    name: 'TeacherConsole',
    component: () => import('@/views/TeacherConsole.vue')
  },
  {
    path: '/student/login',
    name: 'StudentLogin',
    component: () => import('@/views/StudentLogin.vue')
  },
  {
    path: '/student/console/:studentId',
    name: 'StudentConsole',
    component: () => import('@/views/StudentTerminal.vue')
  },
  {
    path: '/student/workspace',
    name: 'StudentWorkspace',
    component: () => import('@/views/StudentWorkspace.vue')
  },
  {
    path: '/teacher/audit',
    name: 'TeacherAudit',
    component: () => import('@/views/TeacherAudit.vue')
  },
  {
    path: '/teacher/vm-management',
    name: 'TeacherVMManagement',
    component: () => import('@/views/TeacherVMManagement.vue')
  },
  {
    path: '/enterprise/sandbox',
    name: 'EnterpriseSandbox',
    component: () => import('@/views/EnterpriseSandbox.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
