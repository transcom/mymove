import { React } from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { AmendOrders } from './AmendOrders';

import { getOrders, submitAmendedOrders } from 'services/internalApi';
import { customerRoutes } from 'constants/routes';
import { renderWithProviders } from 'testUtils';
import { selectOrdersForLoggedInUser } from 'store/entities/selectors';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('store/entities/selectors', () => ({
  ...jest.requireActual('store/entities/selectors'),
  selectOrdersForLoggedInUser: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  createUploadForDocument: jest.fn().mockImplementation(() => Promise.resolve()),
  deleteUpload: jest.fn().mockImplementation(() => Promise.resolve()),
  submitAmendedOrders: jest.fn(),
  getOrders: jest.fn().mockImplementation(() => Promise.resolve()),
}));

const testPropsWithUploads = {
  id: 'testOrderId',
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
  uploaded_orders: {
    id: 'testId',
    uploads: [
      {
        bytes: 1578588,
        contentType: 'image/png',
        createdAt: '2024-02-23T16:51:45.504Z',
        filename: 'Screenshot 2024-02-15 at 12.22.53 PM (2).png',
        id: 'fd88b0e6-ff6d-4a99-be6f-49458a244209',
        status: 'PROCESSING',
        updatedAt: '2024-02-23T16:51:45.504Z',
        url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/fd88b0e6-ff6d-4a99-be6f-49458a244209?contentType=image%2Fpng',
      },
    ],
  },
  moves: ['testMoveId'],
};

const testPropsNoUploads = {
  id: 'testOrderId2',
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
  uploaded_orders: {
    id: 'testId',
    service_member_id: 'testId',
    uploads: [],
  },
  moves: ['testMoveId'],
};

const testOrders = [testPropsWithUploads, testPropsNoUploads];

afterEach(() => {
  jest.resetAllMocks();
});

describe('Amended Orders Upload page', () => {
  const testProps = {
    serviceMemberId: '123',
    updateOrders: jest.fn(),
    orders: [testPropsWithUploads, testPropsNoUploads],
  };

  it('renders all content of AmendOrders', async () => {
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testPropsWithUploads);

    renderWithProviders(<AmendOrders {...testProps} />, {
      path: customerRoutes.ORDERS_AMEND_PATH,
      params: { orderId: 'testOrderId' },
    });

    await screen.findByRole('heading', { level: 1, name: 'Orders' });
    expect(screen.getByTestId('info-container')).toBeInTheDocument();
    expect(screen.getByTestId('upload-info-container')).toBeInTheDocument();
    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();
    const cancelBtn = await screen.findByRole('button', { name: 'Cancel' });
    expect(cancelBtn).toBeInTheDocument();
  });

  it('navigates user when cancel button is clicked', async () => {
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testPropsWithUploads);

    renderWithProviders(<AmendOrders {...testProps} />, {
      path: customerRoutes.ORDERS_AMEND_PATH,
      params: { orderId: 'testOrderId' },
    });

    const cancelBtn = await screen.findByRole('button', { name: 'Cancel' });
    expect(cancelBtn).toBeInTheDocument();

    await userEvent.click(cancelBtn);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(-1);
    });
  });

  it('navigates user when save button is clicked', async () => {
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testPropsWithUploads);
    submitAmendedOrders.mockImplementation(() => Promise.resolve());

    renderWithProviders(<AmendOrders {...testProps} />, {
      path: customerRoutes.ORDERS_AMEND_PATH,
      params: { orderId: 'testOrderId' },
    });

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalled();
    });
  });

  it('shows an error if the API returns an error', async () => {
    selectOrdersForLoggedInUser.mockImplementation(() => testOrders);
    getOrders.mockResolvedValue(testPropsWithUploads);
    submitAmendedOrders.mockImplementation(() =>
      // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
      // eslint-disable-next-line prefer-promise-reject-errors
      Promise.reject({
        message: 'A server error occurred saving the amended orders',
        response: {
          body: {
            detail: 'A server error occurred saving the amended orders',
          },
        },
      }),
    );

    // Need to provide complete & valid initial values because we aren't testing the form here, and just want to submit immediately
    renderWithProviders(<AmendOrders {...testProps} />, {
      path: customerRoutes.ORDERS_AMEND_PATH,
      params: { orderId: 'testOrderId' },
    });
    const saveButton = await screen.findByText('Save');
    expect(saveButton).toBeInTheDocument();
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(submitAmendedOrders).toHaveBeenCalled();
    });

    expect(await screen.queryByText('A server error occurred saving the amended orders')).toBeInTheDocument();
    expect(mockNavigate).not.toHaveBeenCalled();
  });
});
