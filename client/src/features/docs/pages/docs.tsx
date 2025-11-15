'use server';

import { ApiError } from 'next/dist/server/api-utils';
import NovaChain from '../components/nova-chain-page';
import docService from '../services/docs.service';
import ErrorState from '@/components/errorState';
import { FormatFloat } from '@/shared/utils/format';

export const revalidate = 0;

const DocumentPage = async () => {
  try {
    const res = await docService.InfoDowloadNovaChain();

    if (!res) {
      return <ErrorState message={'Failed to fetch data'} />;
    }

    const size = res.headers['content-length']
      ? Number(res.headers['content-length'])
      : 0;

    const sizeFile = FormatFloat(size / 1024 / 1024, 2);

    return <NovaChain SizeFile={sizeFile} />;
  } catch (err) {
    if (err instanceof ApiError) {
      return <ErrorState message={err.message} />;
    }

    return <ErrorState message={'Unknown error'} />;
  }
};

export default DocumentPage;
