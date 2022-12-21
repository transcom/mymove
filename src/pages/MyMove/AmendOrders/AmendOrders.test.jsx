import { React } from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { AmendOrders } from './AmendOrders';

import { getOrdersForServiceMember, submitAmendedOrders } from 'services/internalApi';
import { generalRoutes } from 'constants/routes';

const mockPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: () => ({
    pathname: 'localhost:3000/',
  }),
  useHistory: () => ({
    push: mockPush,
  }),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  createUploadForDocument: jest.fn().mockImplementation(() => Promise.resolve()),
  deleteUpload: jest.fn().mockImplementation(() => Promise.resolve()),
  submitAmendedOrders: jest.fn(),
}));

describe('Amended Orders Upload page', () => {
  const testProps = {
    serviceMemberId: '123',
    updateOrders: jest.fn(),
    currentOrders: {
      moves: ['testMove'],
    },
  };

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

  beforeEach(() => {
    getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));
  });

  it('loads orders on mount', async () => {
    const { queryByText, findByRole, getByRole } = render(<AmendOrders {...testProps} />);

    expect(getByRole('heading', { name: 'Loading, please wait...', level: 2 })).toBeInTheDocument();

    expect(await findByRole('heading', { name: 'Upload orders', level: 5 })).toBeInTheDocument();
    expect(queryByText('Loading, please wait...')).not.toBeInTheDocument();

    expect(getOrdersForServiceMember).toHaveBeenCalled();
    expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
  });

  it('renders the save button', async () => {
    const { findByText } = render(<AmendOrders {...testProps} uploads={[]} />);

    expect(await findByText('Save')).toBeInTheDocument();
  });

  it('renders the cancel button', async () => {
    const { findByText } = render(<AmendOrders {...testProps} uploads={[]} />);

    expect(await findByText('Cancel')).toBeInTheDocument();
  });

  describe('when the user clicks cancel', () => {
    it('redirects to the home page', async () => {
      render(<AmendOrders {...testProps} moveIsInDraft={false} />);

      const cancelButton = await screen.findByText('Cancel');
      expect(cancelButton).toBeInTheDocument();
      userEvent.click(cancelButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith(generalRoutes.HOME_PATH);
      });
    });
  });

  describe('when the user saves', () => {
    it('submits the form and redirects to the home page', async () => {
      submitAmendedOrders.mockImplementation(() => Promise.resolve());
      render(<AmendOrders {...testProps} moveIsInDraft={false} />);

      const saveButton = await screen.findByText('Save');
      expect(saveButton).toBeInTheDocument();
      userEvent.click(saveButton);

      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith(generalRoutes.HOME_PATH);
      });
    });

    it('shows an error if the API returns an error', async () => {
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
      render(<AmendOrders {...testProps} moveIsInDraft={false} />);

      const saveButton = await screen.findByText('Save');
      expect(saveButton).toBeInTheDocument();
      userEvent.click(saveButton);

      await waitFor(() => {
        expect(submitAmendedOrders).toHaveBeenCalled();
      });

      expect(await screen.queryByText('A server error occurred saving the amended orders')).toBeInTheDocument();
      expect(mockPush).not.toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
