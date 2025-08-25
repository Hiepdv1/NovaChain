import axios, { AxiosError } from 'axios';
import { handleUnauthorized } from './handleAuth';

export const http = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
});

http.interceptors.response.use(
  (res) => {
    return res;
  },
  async (err: AxiosError) => {
    if (err?.status === 401) {
      await handleUnauthorized();
    }
    return Promise.reject(err);
  },
);
