import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { act } from 'react-dom/test-utils';

import AllowancesDetailForm from './AllowancesDetailForm';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

const initialValues = {
  authorizedWeight: '11000',
  proGearWeight: '2000',
  proGearWeightSpouse: '500',
  requiredMedicalEquipmentWeight: '1000',
  organizationalClothingAndIndividualEquipment: true,
  weightRestriction: '500',
  ubWeightRestriction: '400',
};

const initialValuesOconusAdditions = {
  accompaniedTour: true,
  dependentsTwelveAndOver: '2',
  dependentsUnderTwelve: '4',
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
          {
            touched: false,
          },
          {
            setValue: jest.fn(),
            setTouched: jest.fn(),
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
  weightRestriction: 500,
  ubWeightRestriction: '400',
};

const entitlementOconusAdditions = {
  accompaniedTour: true,
  dependentsTwelveAndOver: 2,
  dependentsUnderTwelve: 4,
};

jest.mock('../../../utils/featureFlags', () => ({
  ...jest.requireActual('../../../utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('AllowancesDetailForm', () => {
  it('renders the form', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} />
      </Formik>,
    );

    expect(await screen.findByTestId('proGearWeightInput')).toHaveDisplayValue('2,000');
    expect(screen.getByTestId('proGearWeightHint')).toHaveTextContent('Max. 2,000 lbs');
    expect(screen.getByTestId('proGearWeightSpouseHint')).toHaveTextContent('Max. 500 lbs');
  });

  it('renders the pro-gear hints', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} />
      </Formik>,
    );

    expect(await screen.findByTestId('proGearWeightHint')).toHaveTextContent('Max. 2,000 lbs');
    expect(screen.getByTestId('proGearWeightSpouseHint')).toHaveTextContent('Max. 500 lbs');
  });

  it('renders the title header', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} header="Test Header" />
      </Formik>,
    );

    expect(await screen.findByRole('heading', { level: 3 })).toHaveTextContent('Test Header');
  });

  it('does not render conditional oconus fields on load', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} header="Test Header" />
      </Formik>,
    );

    expect(screen.queryByText('Accompanied tour')).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).not.toBeInTheDocument();
  });

  it('does render conditional oconus fields when present in entitlement', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));

    await act(async () => {
      render(
        <Formik initialValues={{ ...initialValues, ...initialValuesOconusAdditions }}>
          <AllowancesDetailForm
            entitlements={{ ...entitlements, ...entitlementOconusAdditions }}
            branchOptions={branchOptions}
          />
        </Formik>,
      );
    });
    // Wait for state
    await waitFor(() => expect(screen.queryByLabelText(/Accompanied tour/)).toBeInTheDocument());
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).toBeInTheDocument();
  });

  it('does not render conditional civilian TDY UB allowance field on load', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm
          entitlements={entitlements}
          branchOptions={branchOptions}
          civilianTDYUBMove={false}
          header="Test Header"
        />
      </Formik>,
    );

    expect(screen.queryByTestId('civilianTdyUbAllowance')).not.toBeInTheDocument();
  });

  it('renders conditional civilian TDY UB allowance field on load', async () => {
    isBooleanFlagEnabled.mockResolvedValue(true);
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm
          entitlements={entitlements}
          branchOptions={branchOptions}
          civilianTDYUBMove
          header="Test Header"
        />
      </Formik>,
    );

    const civilianTdyUbAllowance = await screen.findByTestId('civilianTdyUbAllowance');
    expect(civilianTdyUbAllowance).toBeInTheDocument();
  });
});

describe('AllowancesDetailForm additional tests', () => {
  it('renders gun safe checkbox field', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} />
      </Formik>,
    );

    expect(await screen.findByTestId('gunSafeInput')).toBeInTheDocument();
  });

  it('renders admin weight location section with conditional weight restriction field', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} />
      </Formik>,
    );

    const adminWeightCheckbox = await screen.findByTestId('adminWeightLocation');
    expect(adminWeightCheckbox).toBeInTheDocument();
    expect(screen.getByLabelText('Admin restricted weight location')).toBeChecked();

    const weightRestrictionInput = screen.getByTestId('weightRestrictionInput');
    expect(weightRestrictionInput).toBeInTheDocument();
    expect(weightRestrictionInput).toHaveValue('500');
  });

  it('does not render the admin weight location section when the weightRestriction entitlement is null', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm
          entitlements={{ ...entitlements, weightRestriction: null }}
          branchOptions={branchOptions}
        />
      </Formik>,
    );

    const adminWeightCheckbox = await screen.findByTestId('adminWeightLocation');
    expect(adminWeightCheckbox).toBeInTheDocument();
    expect(screen.queryByTestId('weightRestrictionInput')).not.toBeInTheDocument();
  });

  it('renders admin UB weight location section with conditional UB weight restriction field', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} />
      </Formik>,
    );

    const adminUBWeightCheckbox = await screen.findByTestId('adminUBWeightLocation');
    expect(adminUBWeightCheckbox).toBeInTheDocument();
    expect(screen.getByLabelText('Admin restricted UB weight location')).toBeChecked();

    const ubWeightRestrictionInput = screen.getByTestId('ubWeightRestrictionInput');
    expect(ubWeightRestrictionInput).toBeInTheDocument();
    expect(ubWeightRestrictionInput).toHaveValue('400');
  });

  it('does not render the admin UB weight location section when the ubWeightRestriction entitlement is null', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm
          entitlements={{ ...entitlements, ubWeightRestriction: null }}
          branchOptions={branchOptions}
        />
      </Formik>,
    );

    const adminUBWeightCheckbox = await screen.findByTestId('adminUBWeightLocation');
    expect(adminUBWeightCheckbox).toBeInTheDocument();
    expect(screen.queryByTestId('ubWeightRestrictionInput')).not.toBeInTheDocument();
  });

  it('displays the total weight allowance correctly', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} />
      </Formik>,
    );

    expect(await screen.findByTestId('weightAllowance')).toHaveTextContent('11,000');
  });

  it('renders the form disabled with all information if flag is passed', async () => {
    render(
      <Formik initialValues={initialValues}>
        <AllowancesDetailForm entitlements={entitlements} branchOptions={branchOptions} formIsDisabled />
      </Formik>,
    );

    expect(await screen.findByTestId('proGearWeightInput')).toHaveDisplayValue('2,000');
    expect(await screen.findByTestId('proGearWeightInput')).toBeDisabled();
    expect(screen.getByTestId('proGearWeightHint')).toHaveTextContent('Max. 2,000 lbs');
    expect(await screen.findByTestId('proGearWeightSpouseInput')).toHaveDisplayValue('500');
    expect(await screen.findByTestId('proGearWeightSpouseInput')).toBeDisabled();
    expect(screen.getByTestId('proGearWeightSpouseHint')).toHaveTextContent('Max. 500 lbs');
    expect(await screen.findByTestId('rmeInput')).toHaveDisplayValue('1,000');
    expect(await screen.findByTestId('rmeInput')).toBeDisabled();
    expect(await screen.findByTestId('branchInput')).toHaveDisplayValue('Army');
    expect(await screen.findByTestId('branchInput')).toBeDisabled();
    expect(await screen.findByTestId('sitInput')).toHaveDisplayValue('0');
    expect(await screen.findByTestId('sitInput')).toBeDisabled();
    expect(await screen.findByTestId('ocieInput')).toBeInTheDocument();
    expect(await screen.findByTestId('ocieInput')).toBeDisabled();
    expect(await screen.findByTestId('gunSafeInput')).toBeInTheDocument();
    expect(await screen.findByTestId('gunSafeInput')).toBeDisabled();
    expect(await screen.findByTestId('adminWeightLocation')).toBeInTheDocument();
    expect(await screen.findByTestId('adminWeightLocation')).toBeDisabled();
    expect(await screen.findByTestId('adminUBWeightLocation')).toBeInTheDocument();
    expect(await screen.findByTestId('adminUBWeightLocation')).toBeDisabled();
    expect(await screen.findByTestId('dependentsAuthorizedInput')).toBeInTheDocument();
    expect(await screen.findByTestId('dependentsAuthorizedInput')).toBeDisabled();
  });
});
