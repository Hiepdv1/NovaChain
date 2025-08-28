/* eslint-disable @typescript-eslint/no-explicit-any */
export interface BaseResponse<T> {
  statusCode: number;
  message?: string;
  traceId?: string;
  data: T;
}

export interface BaseErrorResponse extends Error {
  success?: false;
  message: string;
  code?: string | number;
  statusCode?: number;
  errors?: Record<string, any>;
}

export interface NullTime {
  Time: string;
  Valid: boolean;
}
