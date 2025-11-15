import { throwApiError } from '@/lib/axios/handleErrorApi';
import { http } from '@/lib/axios/http';
import https from 'https';

class DocumentService {
  public async InfoDowloadNovaChain() {
    try {
      const agent = new https.Agent({
        rejectUnauthorized: false,
      });

      const res = await http.head('/download/novachain.rar', {
        httpsAgent: agent,
      });

      return res;
    } catch (err) {
      throwApiError(err);
    }
  }
}

const docService = new DocumentService();

export default docService;
