import React from 'react';

import FOUOHeader from '../components/FOUOHeader';
import MilMoveHeader from '../components/MilMoveHeader';
import CustomerHeader from '../components/CustomerHeader';

export default {
  title: 'Components|Headers',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d9ad20e6-944c-48a2-bbd2-1c7ed8bc1315?mode=design',
    },
  },
};

export const all = () => (
  <div>
    <FOUOHeader />
    <MilMoveHeader />
    <CustomerHeader />
  </div>
);
