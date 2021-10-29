import React from 'react';

import FinancialReviewModal from './FinancialReviewModal';

export default {
  title: 'Office Components/FinancialReviewModal',
  component: FinancialReviewModal,
};

export const Basic = () => <FinancialReviewModal Submit={() => {}} onClose={() => {}} />;
