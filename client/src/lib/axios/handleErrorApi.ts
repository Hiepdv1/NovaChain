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

export class ApiError extends Error implements BaseErrorResponse {
  code?: string | number;
  statusCode?: number;
  errors?: Record<string, unknown>;

  constructor(
    message: string,
    statusCode?: number,
    code?: string | number,
    errors?: Record<string, unknown>,
  ) {
    super(message);
    this.name = 'ApiError';
    this.statusCode = statusCode;
    this.code = code;
    this.errors = errors;

    Object.setPrototypeOf(this, ApiError.prototype);
  }
}

export function throwApiError(error: unknown): never {
  if (error instanceof ApiError) throw error;

  if (typeof error === 'object' && error !== null && 'message' in error) {
    const e = error as BaseErrorResponse;
    throw new ApiError(e.message, e.statusCode, e.code, e.errors);
  }

  if (error instanceof Error) {
    throw new ApiError(error.message, 500, undefined, undefined);
  }

  throw new ApiError('Unknown error', 500);
}
