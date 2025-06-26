import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';

import OrdersTable from './OrdersTable';

import { ORDERS_PAY_GRADE_TYPE, ORDERS_RANK_OPTIONS, ORDERS_TYPE_OPTIONS } from 'constants/orders';

const mockProps = {
  hasDependents: true,
  issueDate: '2023-01-01',
  moveId: '123',
  newDutyLocationName: 'New Location',
  onEditClick: jest.fn(),
  orderType: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION,
  reportByDate: '2023-02-01',
  uploads: [{}, {}],
  originDutyLocationName: 'Current Location',
  payGrade: ORDERS_PAY_GRADE_TYPE.E_5,
  rank: { rankAbbv: ORDERS_RANK_OPTIONS.AIR_FORCE.SSgt },
  orderId: '456',
  counselingOfficeName: 'Counseling Office',
  accompaniedTour: true,
  dependentsUnderTwelve: 1,
  dependentsTwelveAndOver: 2,
};

describe('OrdersTable', () => {
  it('renders all fields correctly', () => {
    render(<OrdersTable {...mockProps} />);

    expect(screen.getByText('Orders type')).toBeInTheDocument();
    expect(screen.getByText('Orders date')).toBeInTheDocument();
    expect(screen.getByText('Current duty location')).toBeInTheDocument();
    expect(screen.getByText('Counseling office')).toBeInTheDocument();
    expect(screen.getByText('Counseling Office')).toBeInTheDocument();
    expect(screen.getByText('Dependents')).toBeInTheDocument();
    expect(screen.getByText('2 files')).toBeInTheDocument();
    expect(screen.getByText('Pay grade')).toBeInTheDocument();
    expect(screen.getByText('Rank')).toBeInTheDocument();
  });

  it('renders OCONUS fields when conditions are met', () => {
    render(<OrdersTable {...mockProps} />);

    expect(screen.getByText('Accompanied tour')).toBeInTheDocument();
    expect(screen.getByText('Dependents under twelve')).toBeInTheDocument();
    expect(screen.getByText('Dependents twelve and over')).toBeInTheDocument();
  });

  it('does not render OCONUS fields when conditions are not met', () => {
    const propsWithoutOCONUS = {
      ...mockProps,
      accompaniedTour: false,
      dependentsUnderTwelve: 0,
      dependentsTwelveAndOver: 0,
    };
    render(<OrdersTable {...propsWithoutOCONUS} />);

    expect(screen.queryByText('Accompanied tour')).not.toBeInTheDocument();
    expect(screen.queryByText('Dependents under twelve')).not.toBeInTheDocument();
    expect(screen.queryByText('Dependents twelve and over')).not.toBeInTheDocument();
  });

  it('renders correct label for retirement or separation orders', () => {
    const propsWithRetirement = { ...mockProps, orderType: 'RETIREMENT' };
    render(<OrdersTable {...propsWithRetirement} />);

    expect(screen.getByText('HOR, PLEAD or HOS')).toBeInTheDocument();
  });

  it('calls onEditClick with correct path when edit button is clicked', () => {
    render(<OrdersTable {...mockProps} />);
    const editButton = screen.getByText('Edit');
    fireEvent.click(editButton);

    expect(mockProps.onEditClick).toHaveBeenCalledWith('/move/123/review/edit-orders/456');
  });
});
