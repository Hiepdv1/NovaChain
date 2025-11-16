import { BaseErrorResponse, BaseResponse } from '@/shared/types/api';
import { useQuery } from '@tanstack/react-query';
import { DownloadInfo } from '../types/docs';
import docService from '../services/docs.service';

export const useDownloadInfoQuery = () => {
  return useQuery<BaseResponse<DownloadInfo>, BaseErrorResponse>({
    queryKey: ['txPending'],
    queryFn: docService.InfoDowloadNovaChain,
  });
};
