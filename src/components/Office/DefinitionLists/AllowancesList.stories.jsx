import React from 'react';
import { object } from '@storybook/addon-knobs';

import AllowancesList from './AllowancesList';

export default {
  title: 'Office Components/AllowancesList',
  component: AllowancesList,
};

const info = {
  branch: 'NAVY',
  rank: 'E_6',
  weightAllowance: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
  requiredMedicalEquipmentWeight: 1000,
  organizationalClothingAndIndividualEquipment: true,
};

export const Basic = () => <AllowancesList info={object('info', info)} />;
