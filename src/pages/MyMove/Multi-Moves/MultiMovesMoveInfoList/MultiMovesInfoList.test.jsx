import React from 'react';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';

import MultiMovesMoveInfoList from './MultiMovesMoveInfoList';

import { ORDERS_TYPE } from 'constants/orders';

describe('MultiMovesMoveInfoList', () => {
  const mockMoveSeparation = {
    status: 'DRAFT',
    orders: {
      issue_date: '2022-01-01',
      orders_type: ORDERS_TYPE.SEPARATION,
      report_by_date: '2022-02-01',
      origin_duty_location: {
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
      new_duty_location: {
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
      issue_date: '2022-01-01',
      orders_type: ORDERS_TYPE.RETIREMENT,
      report_by_date: '2022-02-01',
      origin_duty_location: {
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
      new_duty_location: {
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
      issue_date: '2022-01-01',
      orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
      report_by_date: '2022-02-01',
      origin_duty_location: {
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
      new_duty_location: {
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
    expect(getByText('01 Jan 2022')).toBeInTheDocument();

    expect(getByText('Orders Type')).toBeInTheDocument();
    expect(getByText('Separation')).toBeInTheDocument();

    expect(getByText('Separation Date')).toBeInTheDocument();
    expect(getByText('01 Feb 2022')).toBeInTheDocument();

    expect(getByText('Current Duty Location')).toBeInTheDocument();
    expect(getByText('HOR or PLEAD')).toBeInTheDocument();
  });

  it('renders move information correctly', () => {
    const { getByText } = render(<MultiMovesMoveInfoList move={mockMoveRetirement} />);

    expect(getByText('Move Status')).toBeInTheDocument();
    expect(getByText('DRAFT')).toBeInTheDocument();

    expect(getByText('Orders Issue Date')).toBeInTheDocument();
    expect(getByText('01 Jan 2022')).toBeInTheDocument();

    expect(getByText('Orders Type')).toBeInTheDocument();
    expect(getByText('Retirement')).toBeInTheDocument();

    expect(getByText('Retirement Date')).toBeInTheDocument();
    expect(getByText('01 Feb 2022')).toBeInTheDocument();

    expect(getByText('Current Duty Location')).toBeInTheDocument();
    expect(getByText('HOR, HOS, or PLEAD')).toBeInTheDocument();
  });

  it('renders move information correctly', () => {
    const { getByText } = render(<MultiMovesMoveInfoList move={mockMovePCS} />);

    expect(getByText('Move Status')).toBeInTheDocument();
    expect(getByText('DRAFT')).toBeInTheDocument();

    expect(getByText('Orders Issue Date')).toBeInTheDocument();
    expect(getByText('01 Jan 2022')).toBeInTheDocument();

    expect(getByText('Orders Type')).toBeInTheDocument();
    expect(getByText('Permanent Change of Station')).toBeInTheDocument();

    expect(getByText('Report by Date')).toBeInTheDocument();
    expect(getByText('01 Feb 2022')).toBeInTheDocument();

    expect(getByText('Current Duty Location')).toBeInTheDocument();

    expect(getByText('Destination Duty Location')).toBeInTheDocument();
  });
});
