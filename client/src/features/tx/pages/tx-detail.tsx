'use server';

import NotFoundState from '@/components/notfound';
import transactionService from '../services/transactions.service';
import { ApiError } from 'next/dist/server/api-utils';
import ErrorState from '@/components/errorState';
import TransactionCard from '../components/tx-detail/transaction-card';

interface PageParams {
  tx_hash: string;
}

const TransactionDetailPage = async ({
  params,
}: {
  params: Promise<PageParams>;
}) => {
  const { tx_hash } = await params;

  try {
    const res = await transactionService.GetDetailTransaction(tx_hash);
    if (!res) {
      return <ErrorState message={'Failed to fetch data'} />;
    }

    const transaction = res.data;

    return <TransactionCard transaction={transaction} />;
  } catch (err) {
    if (err instanceof ApiError) {
      if (err.statusCode === 404) return <NotFoundState />;
      return <ErrorState message={err.message} />;
    }
    return <ErrorState />;
  }
};

export default TransactionDetailPage;
