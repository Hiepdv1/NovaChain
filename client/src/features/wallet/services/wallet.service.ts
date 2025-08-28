import { BaseResponse } from '@/shared/types/api';
import { Wallet, WalletConnectPayload } from '../types/wallet';
import { http } from '@/lib/axios/http';
import { handleApiError } from '@/lib/axios/handleErrorApi';
class WalletService {
  public async CreateWallet(payload: WalletConnectPayload) {
    try {
      const res = await http.post<BaseResponse<null>>(
        '/wallet/__pub/new',
        payload,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async ImportWallet(payload: WalletConnectPayload) {
    try {
      const res = await http.post<BaseResponse<null>>(
        '/wallet/__pub/import',
        payload,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetWallet() {
    try {
      const res = await http.get<BaseResponse<Wallet>>('/wallet/__pri/me');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async Disconnect() {
    try {
      const res = await http.post<BaseResponse<null>>(
        '/wallet/__pri/disconnect',
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const walletService = new WalletService();

export default walletService;
