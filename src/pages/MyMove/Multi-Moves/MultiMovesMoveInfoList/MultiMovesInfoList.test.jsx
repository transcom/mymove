import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; // For expect assertions

import MultiMovesMoveInfoList from './MultiMovesMoveInfoList';

describe('MultiMovesMoveInfoList', () => {
  const mockMoveSeparation = {
    status: 'DRAFT',
    orders: {
      date_issued: '2022-01-01',
      ordersType: 'SEPARATION',
      reportByDate: '2022-02-01',
      originDutyLocation: {
        name: 'Fort Bragg North Station',
        address: {
          streetAddress1: '123 Main Ave',
          streetAddress2: 'Apartment 9000',
          streetAddress3: '',
          city: 'Anytown',
          state: 'AL',
          postalCode: '90210',
          country: 'USA',
        },
      },
      destinationDutyLocation: {
        name: 'Fort Bragg North Station',
        address: {
          streetAddress1: '123 Main Ave',
          streetAddress2: 'Apartment 9000',
          streetAddress3: '',
          city: 'Anytown',
          state: 'AL',
          postalCode: '90210',
          country: 'USA',
        },
      },
    },
  };

  const mockMoveRetirement = {
    status: 'DRAFT',
    orders: {
      date_issued: '2022-01-01',
      ordersType: 'RETIREMENT',
      reportByDate: '2022-02-01',
      originDutyLocation: {
        name: 'Fort Bragg North Station',
        address: {
          streetAddress1: '123 Main Ave',
          streetAddress2: 'Apartment 9000',
          streetAddress3: '',
          city: 'Anytown',
          state: 'AL',
          postalCode: '90210',
          country: 'USA',
        },
      },
      destinationDutyLocation: {
        name: 'Fort Bragg North Station',
        address: {
          streetAddress1: '123 Main Ave',
          streetAddress2: 'Apartment 9000',
          streetAddress3: '',
          city: 'Anytown',
          state: 'AL',
          postalCode: '90210',
          country: 'USA',
        },
      },
    },
  };

  const mockMovePCS = {
    status: 'DRAFT',
    orders: {
      date_issued: '2022-01-01',
      ordersType: 'PERMANENT_CHANGE_OF_DUTY_STATION',
      reportByDate: '2022-02-01',
      originDutyLocation: {
        name: 'Fort Bragg North Station',
        address: {
          streetAddress1: '123 Main Ave',
          streetAddress2: 'Apartment 9000',
          streetAddress3: '',
          city: 'Anytown',
          state: 'AL',
          postalCode: '90210',
          country: 'USA',
        },
      },
      destinationDutyLocation: {
        name: 'Fort Bragg North Station',
        address: {
          streetAddress1: '123 Main Ave',
          streetAddress2: 'Apartment 9000',
          streetAddress3: '',
          city: 'Anytown',
          state: 'AL',
          postalCode: '90210',
          country: 'USA',
        },
      },
    },
  };

  it('renders move information correctly', () => {
    const { getByText } = render(<MultiMovesMoveInfoList move={mockMoveSeparation} />);

    expect(getByText('Move Status')).toBeInTheDocument();
    expect(getByText('DRAFT')).toBeInTheDocument();

    expect(getByText('Orders Issue Date')).toBeInTheDocument();
    expect(getByText('2022-01-01')).toBeInTheDocument();

    expect(getByText('Orders Type')).toBeInTheDocument();
    expect(getByText('SEPARATION')).toBeInTheDocument();

    expect(getByText('Separation Date')).toBeInTheDocument();
    expect(getByText('2022-02-01')).toBeInTheDocument();

    expect(getByText('Current Duty Location')).toBeInTheDocument();
    expect(getByText('HOR or PLEAD')).toBeInTheDocument();
  });

  it('renders move information correctly', () => {
    const { getByText } = render(<MultiMovesMoveInfoList move={mockMoveRetirement} />);

    expect(getByText('Move Status')).toBeInTheDocument();
    expect(getByText('DRAFT')).toBeInTheDocument();

    expect(getByText('Orders Issue Date')).toBeInTheDocument();
    expect(getByText('2022-01-01')).toBeInTheDocument();

    expect(getByText('Orders Type')).toBeInTheDocument();
    expect(getByText('RETIREMENT')).toBeInTheDocument();

    expect(getByText('Retirement Date')).toBeInTheDocument();
    expect(getByText('2022-02-01')).toBeInTheDocument();

    expect(getByText('Current Duty Location')).toBeInTheDocument();
    expect(getByText('HOR, HOS, or PLEAD')).toBeInTheDocument();
  });

  it('renders move information correctly', () => {
    const { getByText } = render(<MultiMovesMoveInfoList move={mockMovePCS} />);

    expect(getByText('Move Status')).toBeInTheDocument();
    expect(getByText('DRAFT')).toBeInTheDocument();

    expect(getByText('Orders Issue Date')).toBeInTheDocument();
    expect(getByText('2022-01-01')).toBeInTheDocument();

    expect(getByText('Orders Type')).toBeInTheDocument();
    expect(getByText('PERMANENT_CHANGE_OF_DUTY_STATION')).toBeInTheDocument();

    expect(getByText('Report by Date')).toBeInTheDocument();
    expect(getByText('2022-02-01')).toBeInTheDocument();

    expect(getByText('Current Duty Location')).toBeInTheDocument();

    expect(getByText('Destination Duty Location')).toBeInTheDocument();
  });
});
