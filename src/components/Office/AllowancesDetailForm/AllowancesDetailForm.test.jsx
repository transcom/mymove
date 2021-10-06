import React from 'react';
import { render, screen } from '@testing-library/react';

import AllowancesDetailForm from './AllowancesDetailForm';

const initialValues = {
  authorizedWeight: '11000',
  proGearWeight: '2000',
  proGearWeightSpouse: '500',
  requiredMedicalEquipmentWeight: '1000',
  organizationalClothingAndIndividualEquipment: true,
};

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
  useField: (field) => {
    switch (field.type) {
      case 'checkbox': {
        return [
          {
            name: field.name,
            value: !!initialValues[field.name],
            checked: !!initialValues[field.name],
            onChange: jest.fn(),
            onBlur: jest.fn(),
          },
        ];
      }

      default: {
        return [
          {
            value: initialValues[field.name],
          },
          {
            touched: false,
          },
          {
            setValue: jest.fn(),
            setTouched: jest.fn(),
          },
        ];
      }
    }
  },
}));

const { Formik } = jest.requireActual('formik');

const rankOptions = [
  { key: 'E_1', value: 'E-1' },
  { key: 'E_2', value: 'E-2' },
  { key: 'E_3', value: 'E-3' },
  { key: 'E_4', value: 'E-4' },
  { key: 'E_5', value: 'E-5' },
  { key: 'E_6', value: 'E-6' },
  { key: 'E_7', value: 'E-7' },
  { key: 'E_8', value: 'E-8' },
  { key: 'E_9', value: 'E-9' },
];

const branchOptions = [
  { key: 'Army', value: 'Army' },
  { key: 'Navy', value: 'Navy' },
  { key: 'Marines', value: 'Marines' },
];

const entitlements = {
  authorizedWeight: 11000,
  dependentsAuthorized: true,
  nonTemporaryStorage: false,
  privatelyOwnedVehicle: false,
  proGearWeight: 2000,
  proGearWeightSpouse: 500,
  requiredMedicalEquipmentWeight: 1000,
  organizationalClothingAndIndividualEquipment: true,
  storageInTransit: 90,
  totalWeight: 11000,
  totalDependents: 2,
};

describe('AllowancesDetailForm', () => {
  it('renders the form', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} rankOptions={rankOptions} branchOptions={branchOptions} />
      </Formik>,
    );

    expect(await screen.findByTestId('proGearWeightInput')).toHaveDisplayValue('2,000');
    expect(screen.getByTestId('proGearWeightHint')).toHaveTextContent('Max. 2,000 lbs');
    expect(screen.getByTestId('proGearWeightSpouseHint')).toHaveTextContent('Max. 500 lbs');
  });

  it('renders the pro-gear hints', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} rankOptions={rankOptions} branchOptions={branchOptions} />
      </Formik>,
    );

    expect(await screen.findByTestId('proGearWeightHint')).toHaveTextContent('Max. 2,000 lbs');
    expect(screen.getByTestId('proGearWeightSpouseHint')).toHaveTextContent('Max. 500 lbs');
  });

  it('renders the title header', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm
          entitlements={entitlements}
          rankOptions={rankOptions}
          branchOptions={branchOptions}
          header="Test Header"
        />
      </Formik>,
    );

    expect(await screen.findByRole('heading', { level: 3 })).toHaveTextContent('Test Header');
  });
});
