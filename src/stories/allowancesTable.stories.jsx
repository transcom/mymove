import React from 'react';
import { withKnobs, object } from '@storybook/addon-knobs';

import AllowancesTable from '../components/Office/AllowancesTable';

const info = {
  branch: 'Navy',
  rank: 'E-6',
  weightAllowance: '11,000 lbs',
  authorizedWeight: '11,000 lbs',
  progear: '2,000 lbs',
  spouseProgear: '500 lbs',
  storageInTransit: '90 days',
  dependents: 'Authorized',
};

export default {
  title: 'TOO&#47;TIO Components|AllowancesTable',
  decorator: withKnobs,
};

export const Default = () => <AllowancesTable info={object('info', info)} />;
