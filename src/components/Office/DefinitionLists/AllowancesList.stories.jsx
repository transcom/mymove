import React from 'react';
import { object } from '@storybook/addon-knobs';

import AllowancesList from './AllowancesList';

export default {
  title: 'Office Components/AllowancesList',
  component: AllowancesList,
  argTypes: {
    showVisualCues: {
      defaultValue: true,
    },
  },
};

const info = {
  branch: 'NAVY',
  grade: 'E_6',
  totalWeight: 11000,
  authorizedWeight: 11000,
  progear: 2000,
  spouseProgear: 500,
  gunSafeWeight: 500,
  storageInTransit: 90,
  requiredMedicalEquipmentWeight: 1000,
  organizationalClothingAndIndividualEquipment: true,
  ubAllowance: 400,
};

export const Basic = () => <AllowancesList info={object('info', info)} isOconusMove={false} />;

export const VisualCues = (argTypes) => (
  <AllowancesList info={object('info', info)} showVisualCues={argTypes.showVisualCues} isOconusMove={false} />
);
