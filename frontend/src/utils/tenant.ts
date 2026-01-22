import type {Tenant} from '@/api/tenant'

export const getAdminUsername = (tenant: Tenant) => {
    if (tenant.users && tenant.users.length > 0) {
        return tenant.users[0].username
    }
    return '-'
}
