
import axios from './axios'
import qs from 'qs';

export function stringify(params: any) {
  return qs.stringify(params, { skipNulls: true, addQueryPrefix: true })
}
/**上网策略列表 */
export function getInternetList (params: any) {
  return axios.get('/api/v1/policy' + stringify(params))
};
/**新增网络策略 */
export function addInternet(params: any) {
  return axios.post('/api/v1/policy', params)
};
export function getFirewalls() {
  return axios.get(`/api/v1/firewall`)
};