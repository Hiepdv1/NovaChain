import { handleApiError } from '@/lib/axios/handleErrorApi';
import { http } from '@/lib/axios/http';
import { BaseResponseList } from '@/shared/types/api';
import { PaginationParam } from '@/shared/types/query';
import { Miner } from '../types/miner';

class MinerService {
  public async GetMiners(params: PaginationParam) {
    try {
      const res = await http.get<BaseResponseList<Miner[]>>('/chain/miners', {
        params,
      });

      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const minerService = new MinerService();

export default minerService;
