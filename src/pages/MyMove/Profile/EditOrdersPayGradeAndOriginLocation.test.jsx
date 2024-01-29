import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';

import { EditServiceInfo } from './EditServiceInfo';

import { createOrders, getOrdersForServiceMember, patchOrders } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  createOrders: jest.fn().mockImplementation(() => Promise.resolve()),
  patchOrders: jest.fn().mockImplementation(() => Promise.resolve()),
}));
const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

afterEach(() => {
  jest.resetAllMocks();
});
const testProps = {
  serviceMemberId: '123',
  context: { flags: { allOrdersTypes: true } },
  updateOrders: jest.fn(),
  updateServiceMember: jest.fn(),
  setFlashMessage: jest.fn(),
};
describe('EditServiceInfo page updates orders table information', () => {
  it('save button on profile page patches the orders', async () => {
    const testOrdersValues = {
      id: 'testOrdersId',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: false,
      new_duty_location: {
        address: {
          city: 'Des Moines',
          country: 'US',
          id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
          postalCode: '50309',
          state: 'IA',
          streetAddress1: '987 Other Avenue',
          streetAddress2: 'P.O. Box 1234',
          streetAddress3: 'c/o Another Person',
        },
        address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        affiliation: 'AIR_FORCE',
        created_at: '2020-10-19T17:01:16.114Z',
        id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
        name: 'Yuma AFB',
        updated_at: '2020-10-19T17:01:16.114Z',
      },
    };
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      rank: 'E_5',
      current_location: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postalCode: '12345',
          state: 'NY',
          streetAddress1: '123 Main St',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      weight_allotment: {
        total_weight_self: 7000,
        total_weight_self_plus_dependents: 9000,
        pro_gear_weight: 2000,
        pro_gear_weight_spouse: 500,
      },
    };
    createOrders.mockImplementation(() => Promise.resolve(testOrdersValues));

    render(
      <MemoryRouter>
        <EditServiceInfo {...testProps} serviceMember={testServiceMemberValues} currentOrders={testOrdersValues} />
      </MemoryRouter>,
    );

    getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));
    patchOrders.mockImplementation(() => Promise.resolve(testOrdersValues));

    const payGradeInput = await screen.findByLabelText('Pay grade');
    await userEvent.selectOptions(payGradeInput, ['E_2']);

    const submitButton = await screen.findByText('Save');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalled();
    });

    expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
    expect(testProps.setFlashMessage).toHaveBeenCalledWith(
      'EDIT_SERVICE_INFO_SUCCESS',
      'info',
      `Your weight entitlement is now 7,000 lbs.`,
      'Your changes have been saved. Note that the entitlement has also changed.',
    );
  });

  it('if pay grade does not change, entitlement does not change', async () => {
    const testOrdersValues = {
      grade: 'E_5',
      id: 'testOrdersId',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: false,
      new_duty_location: {
        address: {
          city: 'Des Moines',
          country: 'US',
          id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
          postalCode: '50309',
          state: 'IA',
          streetAddress1: '987 Other Avenue',
          streetAddress2: 'P.O. Box 1234',
          streetAddress3: 'c/o Another Person',
        },
        address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        affiliation: 'AIR_FORCE',
        created_at: '2020-10-19T17:01:16.114Z',
        id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
        name: 'Yuma AFB',
        updated_at: '2020-10-19T17:01:16.114Z',
      },
    };
    const testServiceMemberValues = {
      id: 'testServiceMemberId',
      first_name: 'Leo',
      last_name: 'Spaceman',
      affiliation: 'NAVY',
      edipi: '1234567890',
      rank: 'E_5',
      current_location: {
        address: {
          city: 'Test City',
          id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
          postalCode: '12345',
          state: 'NY',
          streetAddress1: '123 Main St',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      weight_allotment: {
        total_weight_self: 7000,
        total_weight_self_plus_dependents: 9000,
        pro_gear_weight: 2000,
        pro_gear_weight_spouse: 500,
      },
    };
    createOrders.mockImplementation(() => Promise.resolve(testOrdersValues));

    render(
      <MemoryRouter>
        <EditServiceInfo {...testProps} serviceMember={testServiceMemberValues} currentOrders={testOrdersValues} />
      </MemoryRouter>,
    );

    getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));
    patchOrders.mockImplementation(() => Promise.resolve(testOrdersValues));

    const payGradeInput = await screen.findByLabelText('Pay grade');
    await userEvent.selectOptions(payGradeInput, ['E_5']);

    const submitButton = await screen.findByText('Save');
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalled();
    });

    expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
    expect(testProps.setFlashMessage).toHaveBeenCalledWith(
      'EDIT_SERVICE_INFO_SUCCESS',
      'success',
      '',
      'Your changes have been saved.',
    );
  });
});
