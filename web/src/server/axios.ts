import axios from "axios";
import { message } from 'antd';

const Axios = axios.create({
  baseURL: '/api',
});

Axios.interceptors.request.use(
  function (config) {
    return config;
  },
  function (error) {
    return Promise.reject(error);
  }
);

Axios.interceptors.response.use(
  function (res) {
    if(res.data.code === undefined){
      return res
    }else if(res.data.code === 0){
      return res
    }else {
      message.error(res.data.msg)
      return res
    }
  },
  function (error) {
    message.error(error.message)
    return Promise.reject(error);
  }
);

export default Axios;
