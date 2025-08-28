import { BaseErrorResponse } from '@/shared/types/api';
import axios, { AxiosError } from 'axios';

export function handleApiError(error: unknown): BaseErrorResponse {
  if (axios.isAxiosError(error)) {
    const err = error as AxiosError<BaseErrorResponse>;
    return (
      err.response?.data || {
        success: false,
        message: err.message || 'Unknown error',
        name: error.name,
        statusCode: err.response?.data?.statusCode || 500,
      }
    );
  }
  return {
    success: false,
    name: (error as Error).name,
    message: 'Unexpected error occurred',
    statusCode: 500,
  };
}
