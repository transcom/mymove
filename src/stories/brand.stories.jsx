import MmLogo from '../../shared/images/milmove-logo.svg';
import TcmLogo from '../../shared/images/transcom-emblem.svg';
import React from 'react';

export default {
  title: 'Global/Brand',
  parameters: {},
};

export const MilmoveLogo = () => (
  <div>
    <MmLogo />
  </div>
);

export const TranscomLogo = () => (
  <div>
    <TcmLogo />
  </div>
);
