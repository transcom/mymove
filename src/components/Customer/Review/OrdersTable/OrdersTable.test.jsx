import React from 'react';
import { render, screen, act, fireEvent } from '@testing-library/react';

import OrdersTable from './OrdersTable';

import { ORDERS_PAY_GRADE_TYPE, ORDERS_RANK_OPTIONS, ORDERS_TYPE_OPTIONS } from 'constants/orders';

const testOrders = {
  authorizedWeight: 5000,
  created_at: '2025-04-10T17:28:19.325Z',
  entitlement: {
    proGear: 2000,
    proGearSpouse: 500,
    ub_allowance: 0,
  },
  grade: 'E_1',
  has_dependents: false,
  id: '1c4f84f7-fdf7-4ecd-810b-84f5d6da1893',
  issue_date: '2025-04-03',
  new_duty_location: {
    address: {
      city: 'Scott',
      country: 'US',
      county: 'PULASKI',
      id: 'e70d2973-28fd-45c6-950b-223816962088',
      isOconus: false,
      postalCode: '72142',
      state: 'AR',
      streetAddress1: 'n/a',
      usPostRegionCitiesID: 'd229f2f3-c002-4e13-a21c-7829a1973a83',
    },
    address_id: 'e70d2973-28fd-45c6-950b-223816962088',
    affiliation: null,
    created_at: '2025-04-02T16:17:06.159Z',
    id: '5515a499-2800-454e-8b5f-ba74d4d164dc',
    name: 'Scott, AR 72142',
    provides_services_counseling: true,
    updated_at: '2025-04-02T16:17:06.159Z',
  },
  orders_type: ORDERS_TYPE_OPTIONS.PERMANENT_CHANGE_OF_STATION,
  originDutyLocationGbloc: 'HAFC',
  origin_duty_location: {
    address: {
      city: 'Scott',
      country: 'US',
      county: 'PULASKI',
      id: 'e70d2973-28fd-45c6-950b-223816962088',
      isOconus: false,
      postalCode: '72142',
      state: 'AR',
      streetAddress1: 'n/a',
      usPostRegionCitiesID: 'd229f2f3-c002-4e13-a21c-7829a1973a83',
    },
    address_id: 'e70d2973-28fd-45c6-950b-223816962088',
    affiliation: null,
    created_at: '2025-04-02T16:17:06.159Z',
    id: '5515a499-2800-454e-8b5f-ba74d4d164dc',
    name: 'Scott, AR 72142',
    provides_services_counseling: true,
    updated_at: '2025-04-02T16:17:06.159Z',
  },
  providesServicesCounseling: true,
  report_by_date: '2025-04-22',
  service_member_id: '71db5bbe-e319-429b-a7e1-b2c1c023f692',
  spouse_has_pro_gear: false,
  status: 'DRAFT',
  updated_at: '2025-04-10T17:28:19.325Z',
  uploaded_orders: {
    id: '88d68cbd-9966-47d5-bf45-130515242930',
    service_member_id: '71db5bbe-e319-429b-a7e1-b2c1c023f692',
    uploads: [
      {
        bytes: 787096,
        contentType: 'image/png',
        createdAt: '2025-04-10T17:28:23.555Z',
        filename: 'Screenshot 2025-01-17 at 12.10.19 PM (2).png-20250410132823',
        id: '961fdfea-71fc-4e27-9ea1-f4e9ce93ad67',
        status: 'CLEAN',
        updatedAt: '2025-04-10T17:28:23.555Z',
        uploadType: 'USER',
        url: '/storage/user/e88fd621-6fe7-4d48-8bd2-8c27b8d9cbfc/uploads/961fdfea-71fc-4e27-9ea1-f4e9ce93ad67?contentType=image%2Fpng\u0026filename=Screenshot+2025-01-17+at+12.10.19%E2%80%AFPM+%282%29.png-20250410132823',
      },
    ],
  },
};

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

describe('Orders table', () => {
  it('renders the Orders table with disabled edit buttons when move is locked by office user', async () => {
    await act(() => {
      render(
        <OrdersTable
          hasDependents={testOrders.has_dependents}
          issueDate={testOrders.issue_date}
          newDutyLocationName={testOrders.new_duty_location.name}
          orderType={testOrders.orders_type}
          reportByDate={testOrders.report_by_date}
          uploads={testOrders.uploaded_orders.uploads}
          payGrade={ORDERS_PAY_GRADE_TYPE[testOrders?.grade] || ''}
          originDutyLocationName={testOrders.origin_duty_location.name}
          orderId={testOrders.id}
          counselingOfficeName="Test Counseling Office"
          isMoveLocked
        />,
      );
    });

    expect(screen.getByTestId('edit-orders-table')).toBeDisabled();
  });

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

  it('calls onEditClick with correct path when edit button is clicked', () => {
    render(<OrdersTable {...mockProps} />);
    const editButton = screen.getByText('Edit');
    fireEvent.click(editButton);

    expect(mockProps.onEditClick).toHaveBeenCalledWith('/move/123/review/edit-orders/456');
  });
});
