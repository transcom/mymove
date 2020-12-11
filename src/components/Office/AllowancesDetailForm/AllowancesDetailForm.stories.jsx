import React from 'react';
import { withKnobs, object } from '@storybook/addon-knobs';
import { Formik } from 'formik';

import AllowancesDetailForm from './AllowancesDetailForm';

export default {
  title: 'TOO/TIO Components|AllowancesDetailForm',
  component: AllowancesDetailForm,
  decorators: [
    withKnobs,
    (Story) => (
      <div style={{ padding: `20px`, background: `#f0f0f0` }}>
        <Story />
      </div>
    ),
  ],
};

const entitlement = {
  authorizedWeight: 1950,
  dependentsAuthorized: true,
  nonTemporaryStorage: true,
  privatelyOwnedVehicle: false,
  proGearWeight: 1500,
  proGearWeightSpouse: 1000,
  storageInTransit: 90,
  totalWeight: 12875,
  totalDependents: 2,
};

export const Basic = () => {
  return (
    <Formik
      initialValues={{
        authorizedWeight: '8000',
      }}
      onSubmit={() => {}}
    >
      <form>
        <AllowancesDetailForm entitlements={object('entitlement', entitlement)} />
      </form>
    </Formik>
  );
};
