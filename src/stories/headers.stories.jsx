import React from 'react';

import FOUOHeader from '../components/FOUOHeader';
import MilMoveHeader from '../components/MilMoveHeader';
import CustomerHeader from '../components/CustomerHeader';

export default {
  title: 'Components|Headers',
};

export const all = () => (
  <div>
    <FOUOHeader />
    <MilMoveHeader />
    <CustomerHeader />
  </div>
);
