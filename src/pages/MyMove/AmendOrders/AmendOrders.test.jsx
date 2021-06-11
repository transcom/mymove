import { React } from 'react';
import { render, waitFor } from '@testing-library/react';

import { AmendOrders } from './AmendOrders';

import { getOrdersForServiceMember } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  createUploadForDocument: jest.fn().mockImplementation(() => Promise.resolve()),
  deleteUpload: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('Amended Orders Upload page', () => {
  const testProps = {
    serviceMemberId: '123',
    updateOrders: jest.fn(),
  };

  const testOrdersValues = {
    id: 'testOrdersId',
    orders_type: 'PERMANENT_CHANGE_OF_STATION',
    issue_date: '2020-11-08',
    report_by_date: '2020-11-26',
    has_dependents: false,
    new_duty_station: {
      address: {
        city: 'Des Moines',
        country: 'US',
        id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
        postal_code: '50309',
        state: 'IA',
        street_address_1: '987 Other Avenue',
        street_address_2: 'P.O. Box 1234',
        street_address_3: 'c/o Another Person',
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
    const { queryByText, queryByRole } = render(<AmendOrders {...testProps} />);

    expect(queryByText('Loading, please wait...')).toBeInTheDocument();

    await waitFor(() => {
      expect(queryByText('Loading, please wait...')).not.toBeInTheDocument();
      expect(queryByRole('heading', { name: 'Upload orders', level: 5 })).toBeInTheDocument();

      expect(getOrdersForServiceMember).toHaveBeenCalled();
      expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
    });
  });

  it('renders the save button', async () => {
    const { queryByText } = render(<AmendOrders {...testProps} uploads={[]} />);

    await waitFor(() => {
      expect(queryByText('Save')).toBeInTheDocument();
    });
  });
  it('renders the cancel button', async () => {
    const { queryByText } = render(<AmendOrders {...testProps} uploads={[]} />);

    await waitFor(() => {
      expect(queryByText('Cancel')).toBeInTheDocument();
    });
  });
});
