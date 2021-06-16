import React from 'react';

import { ReactComponent as MmLogo } from '../shared/images/milmove-logo.svg';
import { ReactComponent as TcmLogo } from '../shared/images/transcom-emblem.svg';

export default {
  title: 'Global/Brand',
};

export const MilmoveLogo = () => (
  <div style={{ backgroundColor: '#71767a', margin: '-15px', padding: '30px' }}>
    <MmLogo />
  </div>
);

export const TranscomLogo = () => (
  <div style={{ maxWidth: '300px' }}>
    <TcmLogo />
  </div>
);
