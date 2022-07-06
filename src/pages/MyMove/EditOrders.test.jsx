import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditOrders } from './EditOrders';

import { patchOrders } from 'services/internalApi';

const mockPush = jest.fn();
const mockGoBack = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useHistory: () => ({
    push: mockPush,
    goBack: mockGoBack,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  patchOrders: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('EditOrders Page', () => {
  const testProps = {
    moveIsApproved: false,
    serviceMember: {
      id: 'id123',
      current_location: {
        address: {
          city: 'Fort Bragg',
          country: 'United States',
          id: 'f1ee4cea-6b23-4971-9947-efb51294ed32',
          postalCode: '29310',
          state: 'NC',
          streetAddress1: '',
        },
        address_id: 'f1ee4cea-6b23-4971-9947-efb51294ed32',
        affiliation: 'ARMY',
        created_at: '2020-10-19T17:01:16.114Z',
        id: 'dca78766-e76b-4c6d-ba82-81b50ca824b9"',
        name: 'Fort Bragg',
        updated_at: '2020-10-19T17:01:16.114Z',
      },
      rank: 'E_2',
      weight_allotment: {
        total_weight_self: 5000,
        total_weight_self_plus_dependents: 8000,
        pro_gear_weight: 2000,
        pro_gear_weight_spouse: 500,
      },
    },
    setFlashMessage: jest.fn(),
    updateOrders: jest.fn(),
    currentOrders: {
      id: 'testOrdersId',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: false,
      spouse_has_pro_gear: false,
      grade: 'E_2',
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
      moves: ['testMove'],
    },
    existingUploads: [
      {
        id: '123',
        created_at: '2020-11-08',
        bytes: 1,
        url: 'url',
        filename: 'Test Upload',
        content_type: 'application/pdf',
      },
    ],
    context: { flags: { allOrdersTypes: true } },
  };

  it('renders the edit orders form', async () => {
    render(<EditOrders {...testProps} />);

    const h1 = await screen.findByRole('heading', { name: 'Orders', level: 1 });
    expect(h1).toBeInTheDocument();

    const editOrdersHeader = await screen.findByRole('heading', { name: 'Edit Orders:', level: 2 });
    expect(editOrdersHeader).toBeInTheDocument();
  });

  it('goes back to the previous page when the cancel button is clicked', async () => {
    render(<EditOrders {...testProps} />);

    const cancel = await screen.findByText('Cancel');

    expect(cancel).toBeInTheDocument();

    userEvent.click(cancel);

    await waitFor(() => {
      expect(mockGoBack).toHaveBeenCalled();
    });
  });

  it('shows an error if the API returns an error', async () => {
    render(<EditOrders {...testProps} />);

    patchOrders.mockImplementation(() =>
      // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
      // eslint-disable-next-line prefer-promise-reject-errors
      Promise.reject({
        message: 'A server error occurred saving the orders',
        response: {
          body: {
            detail: 'A server error occurred saving the orders',
          },
        },
      }),
    );

    const submitButton = await screen.findByRole('button', { name: 'Save' });
    expect(submitButton).not.toBeDisabled();

    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalledTimes(1);
    });

    expect(screen.queryByText('A server error occurred saving the orders')).toBeInTheDocument();
    expect(mockPush).not.toHaveBeenCalled();
  });

  it('next button patches the orders and goes to the previous page', async () => {
    render(<EditOrders {...testProps} />);

    patchOrders.mockImplementation(() => Promise.resolve(testProps.currentOrders));

    const submitButton = await screen.findByRole('button', { name: 'Save' });
    expect(submitButton).not.toBeDisabled();

    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalledTimes(1);
    });

    expect(mockGoBack).toHaveBeenCalledTimes(1);
  });

  it('displays a flash message with the entitlement when the dependents value is updated', async () => {
    render(<EditOrders {...testProps} />);

    patchOrders.mockImplementation(() => Promise.resolve(testProps.currentOrders));

    const hasDependentsYes = await screen.findByLabelText('Yes');
    userEvent.click(hasDependentsYes);

    const submitButton = await screen.findByRole('button', { name: 'Save' });
    expect(submitButton).not.toBeDisabled();

    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalledTimes(1);
    });

    await waitFor(() => {
      expect(testProps.setFlashMessage).toHaveBeenCalledWith(
        'EDIT_ORDERS_SUCCESS',
        'info',
        'Your weight entitlement is now 8,000 lbs.',
        'Your changes have been saved. Note that the entitlement has also changed.',
      );
    });

    expect(mockGoBack).toHaveBeenCalledTimes(1);
  });

  afterEach(jest.clearAllMocks);
});
