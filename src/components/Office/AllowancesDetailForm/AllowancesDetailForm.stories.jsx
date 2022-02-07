import React from 'react';
import { withKnobs, object } from '@storybook/addon-knobs';
import { Formik } from 'formik';

import AllowancesDetailForm from './AllowancesDetailForm';

export default {
  title: 'Office Components/AllowancesDetailForm',
  component: AllowancesDetailForm,
  decorators: [
    withKnobs,
    (Story) => (
      <div className="officeApp" style={{ padding: `20px`, background: `#f0f0f0` }}>
        <Story />
      </div>
    ),
  ],
  argTypes: {
    editableAuthorizedWeight: {
      defaultValue: false,
      control: {
        type: 'select',
        options: [true, false],
      },
    },
    header: {
      defaultValue: null,
      control: {
        type: 'text',
      },
    },
  },
};

const entitlement = {
  authorizedWeight: 1950,
  dependentsAuthorized: true,
  nonTemporaryStorage: true,
  privatelyOwnedVehicle: false,
  proGearWeight: 2000,
  proGearWeightSpouse: 500,
  requiredMedicalEquipmentWeight: 1000,
  storageInTransit: 90,
  organizationalClothingAndIndividualEquipment: true,
  totalWeight: 12875,
  totalDependents: 2,
};

export const Basic = (data) => {
  return (
    <Formik
      initialValues={{
        authorizedWeight: '8000',
        storageInTransit: '90',
      }}
      onSubmit={() => {}}
    >
      <form>
        <AllowancesDetailForm
          entitlements={object('entitlement', entitlement)}
          editableAuthorizedWeight={data.editableAuthorizedWeight}
          header={data.header}
        />
      </form>
    </Formik>
  );
};
