import { BaseErrorResponse, BaseResponse } from '@/shared/types/api';
import { useQuery } from '@tanstack/react-query';
import { NetworkOverview, RecentActivityResponse } from '../types/dashboard';
import dashboardService from '../services/dashboard.service';
import { QueryOptions } from '@/shared/types/query';

export const useNetworkOverview = (
  opts?: QueryOptions<BaseResponse<NetworkOverview>>,
) => {
  return useQuery<BaseResponse<NetworkOverview>, BaseErrorResponse>({
    queryKey: ['overview'],
    queryFn: dashboardService.GetNetWorkOverview,
    ...opts,
  });
};

export const useRecentActivity = (
  opts?: QueryOptions<BaseResponse<RecentActivityResponse>>,
) => {
  return useQuery<BaseResponse<RecentActivityResponse>, BaseErrorResponse>({
    queryKey: ['activity'],
    queryFn: dashboardService.GetRecentActivity,
    ...opts,
  });
};
