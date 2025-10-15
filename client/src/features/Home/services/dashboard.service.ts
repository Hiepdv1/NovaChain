import { handleApiError } from '@/lib/axios/handleErrorApi';
import { http } from '@/lib/axios/http';
import { BaseResponse } from '@/shared/types/api';
import { NetworkOverview, RecentActivityResponse } from '../types/dashboard';

class DashBoardService {
  public async GetNetWorkOverview() {
    try {
      const res = await http.get<BaseResponse<NetworkOverview>>('/dashboard');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetRecentActivity() {
    try {
      const res = await http.get<BaseResponse<RecentActivityResponse>>(
        '/dashboard/recent-activity',
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const dashboardService = new DashBoardService();

export default dashboardService;
