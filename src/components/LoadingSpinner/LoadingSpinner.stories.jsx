import React from 'react';

import LoadingSpinner from './LoadingSpinner';

export default {
  title: 'Components/Loading Spinner',
};

export const LoadingSpinnerComponent = () => <LoadingSpinner />;

export const LoadingSpinnerComponentWithMessage = () => <LoadingSpinner message="custom message" />;
