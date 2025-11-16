import { BaseResponse } from './../../../shared/types/api';
import { handleApiError } from '@/lib/axios/handleErrorApi';
import { http } from '@/lib/axios/http';
import { DownloadInfo } from '../types/docs';

class DocumentService {
  public async InfoDowloadNovaChain() {
    try {
      const res = await http.get<BaseResponse<DownloadInfo>>(
        '/download/novachain.rar?query=info',
      );

      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const docService = new DocumentService();

export default docService;
