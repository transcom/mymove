import React from 'react';
import { object } from '@storybook/addon-knobs';

import AllowancesTable from './AllowancesTable';

const info = {
  branch: 'NAVY',
  rank: 'E_6',
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
};

export default {
  title: 'TOO/TIO Components/AllowancesTable',
};

export const Default = () => <AllowancesTable info={object('info', info)} />;
