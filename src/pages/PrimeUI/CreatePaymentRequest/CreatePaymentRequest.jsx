import React from 'react';
import { useParams } from 'react-router-dom';

import LoadingPlaceholder from '../../../shared/LoadingPlaceholder';
import SomethingWentWrong from '../../../shared/SomethingWentWrong';

import { usePrimeSimulatorGetMove } from 'hooks/queries';

const CreatePaymentRequest = () => {
  const { moveCode } = useParams();

  const { data, isLoading, isError } = usePrimeSimulatorGetMove(moveCode);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return <div style={{ 'white-space': 'pre' }}>{JSON.stringify(data, null, '\t')}</div>;
};

export default CreatePaymentRequest;
