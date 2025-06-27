import React from 'react';

import LoadingButton from './LoadingButton';

export default {
  title: 'Components/LoadingButton',
  component: LoadingButton,
};

export const BasicLoadingButton = () => {
  return (
    <div style={{ padding: 20, fontFamily: 'sans-serif' }}>
      <LoadingButton
        type="button"
        onClick={() => {}}
        isLoading={false}
        labelText="Click to Load"
        loadingText="Loading"
      />
    </div>
  );
};

export const LoadingButtonInLoadingState = () => {
  return (
    <div style={{ padding: 20, fontFamily: 'sans-serif' }}>
      <LoadingButton type="button" onClick={() => {}} isLoading labelText="Click to Load" loadingText="Loading" />
    </div>
  );
};
