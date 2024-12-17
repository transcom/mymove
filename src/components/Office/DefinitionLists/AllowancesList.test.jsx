import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { act } from 'react-dom/test-utils';

import AllowancesList from './AllowancesList';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

const info = {
  branch: 'NAVY',
  grade: 'E_6',
  weightAllowance: 12000,
  authorizedWeight: 11000,
  totalWeight: 12000,
  progear: 2000,
  spouseProgear: 500,
  storageInTransit: 90,
  dependents: true,
  requiredMedicalEquipmentWeight: 1000,
  organizationalClothingAndIndividualEquipment: true,
  ubAllowance: 400,
};

const initialValuesOconusAdditions = {
  accompaniedTour: true,
  dependentsTwelveAndOver: '2',
  dependentsUnderTwelve: '4',
  ubAllowance: 400,
};

const oconusInfo = {
  accompaniedTour: true,
  dependentsTwelveAndOver: 2,
  dependentsUnderTwelve: 4,
  ubAllowance: 400,
};

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
  useField: (field) => {
    const initialValues = {
      accompaniedTour: true,
      dependentsTwelveAndOver: '2',
      dependentsUnderTwelve: '4',
      ubAllowance: '400',
    };

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

jest.mock('../../../utils/featureFlags', () => ({
  ...jest.requireActual('../../../utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('AllowancesList', () => {
  it('renders formatted branch', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('Navy')).toBeInTheDocument();
  });

  it('renders formatted weight allowance', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('12,000 lbs')).toBeInTheDocument();
  });

  it('renders storage in transit', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('90 days')).toBeInTheDocument();
  });

  it('renders authorized dependents', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByTestId('dependents').textContent).toEqual('Authorized');
  });

  it('renders unauthorized dependents', () => {
    const withUnauthorizedDependents = { ...info, dependents: false };
    render(<AllowancesList info={withUnauthorizedDependents} />);
    expect(screen.getByTestId('dependents').textContent).toEqual('Unauthorized');
  });

  it('renders formatted pro-gear', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('2,000 lbs')).toBeInTheDocument();
  });

  it('renders formatted spouse pro-gear', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('500 lbs')).toBeInTheDocument();
  });

  it('renders formatted Required medical equipment', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByText('1,000 lbs')).toBeInTheDocument();
  });

  it('renders authorized ocie', () => {
    render(<AllowancesList info={info} />);
    expect(screen.getByTestId('ocie').textContent).toEqual('Authorized');
  });

  it('renders unauthorized ocie', () => {
    const withUnauthorizedOcie = { ...info, organizationalClothingAndIndividualEquipment: false };
    render(<AllowancesList info={withUnauthorizedOcie} />);
    expect(screen.getByTestId('ocie').textContent).toEqual('Unauthorized');
  });

  it('renders visual cues classname', () => {
    render(<AllowancesList info={info} showVisualCues />);
    expect(screen.getByText('Pro-gear').parentElement.className).toContain('rowWithVisualCue');
    expect(screen.getByText('Spouse pro-gear').parentElement.className).toContain('rowWithVisualCue');
    expect(screen.getByText('Required medical equipment').parentElement.className).toContain('rowWithVisualCue');
    expect(screen.getByText('OCIE').parentElement.className).toContain('rowWithVisualCue');
  });

  it('does not render oconus fields by default', async () => {
    render(<AllowancesList info={info} showVisualCues />);
    expect(screen.queryByText('Accompanied tour')).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents under the age of 12/)).not.toBeInTheDocument();
    expect(screen.queryByLabelText(/Number of dependents of the age 12 or over/)).not.toBeInTheDocument();
    expect(screen.queryByText('Unaccompanied baggage allowance')).not.toBeInTheDocument();
  });

  it('does render oconus fields when present', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    await act(async () => {
      render(
        <Formik initialValues={initialValuesOconusAdditions}>
          <AllowancesList info={{ ...oconusInfo }} showVisualCues />
        </Formik>,
      );
    });
    // Wait for state
    await waitFor(() => expect(screen.getByTestId('ordersAccompaniedTour')).toBeInTheDocument());
    expect(screen.getByTestId('ordersDependentsUnderTwelve')).toBeInTheDocument();
    expect(screen.getByTestId('ordersDependentsTwelveAndOver')).toBeInTheDocument();
    expect(screen.getByTestId('unaccompaniedBaggageAllowance')).toBeInTheDocument();
    expect(screen.getByTestId('unaccompaniedBaggageAllowance').textContent).toEqual('400 lbs');
  });
});
