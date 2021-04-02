import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { UploadOrders } from './UploadOrders';

import { getOrdersForServiceMember, deleteUpload } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  createUploadForDocument: jest.fn().mockImplementation(() => Promise.resolve()),
  deleteUpload: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('Orders Upload page', () => {
  const testProps = {
    serviceMemberId: '123',
    push: jest.fn(),
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
    const { queryByText, queryByRole } = render(<UploadOrders {...testProps} />);

    expect(queryByText('Loading, please wait...')).toBeInTheDocument();

    await waitFor(() => {
      expect(queryByText('Loading, please wait...')).not.toBeInTheDocument();
      expect(queryByRole('heading', { name: 'Upload your orders', level: 1 })).toBeInTheDocument();

      expect(getOrdersForServiceMember).toHaveBeenCalled();
      expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
    });
  });

  it('back button goes to the Orders Info page', async () => {
    const { queryByRole } = render(<UploadOrders {...testProps} />);

    await waitFor(() => {
      const backButton = queryByRole('button', { name: 'Back' });
      expect(backButton).toBeInTheDocument();
      userEvent.click(backButton);
    });

    expect(testProps.push).toHaveBeenCalledWith('/orders/info');
  });

  it('next button is disabled without any uploads', async () => {
    const { queryByRole } = render(<UploadOrders {...testProps} />);

    await waitFor(() => {
      const nextButton = queryByRole('button', { name: 'Next' });
      expect(nextButton).toBeInTheDocument();
      expect(nextButton).toBeDisabled();
    });
  });

  describe('when there are uploads', () => {
    const testUpload = {
      id: 'test upload',
      created_at: 'test date',
      bytes: 100,
      url: 'test url',
      filename: 'Test Upload',
    };

    it('renders the uploads table', async () => {
      const { queryByText } = render(<UploadOrders {...testProps} uploads={[testUpload]} />);

      await waitFor(() => {
        expect(queryByText(testUpload.filename)).toBeInTheDocument();
      });
    });

    it('implements the delete upload handler', async () => {
      deleteUpload.mockImplementation(() => Promise.resolve(testOrdersValues));

      const { queryByRole } = render(<UploadOrders {...testProps} uploads={[testUpload]} />);

      await waitFor(() => {
        const deleteButton = queryByRole('button', { name: 'Delete' });
        userEvent.click(deleteButton);
      });

      expect(deleteUpload).toHaveBeenCalledWith(testUpload.id);
      expect(getOrdersForServiceMember).toHaveBeenCalledTimes(2);
      expect(testProps.updateOrders).toHaveBeenNthCalledWith(1, testOrdersValues);
      expect(testProps.updateOrders).toHaveBeenNthCalledWith(2, testOrdersValues);
    });

    it('next button goes to the Home page if there are uploads', async () => {
      const { queryByRole } = render(<UploadOrders {...testProps} uploads={[testUpload]} />);

      await waitFor(() => {
        const nextButton = queryByRole('button', { name: 'Next' });
        expect(nextButton).toBeInTheDocument();
        expect(nextButton).not.toBeDisabled();

        userEvent.click(nextButton);
      });

      expect(testProps.push).toHaveBeenCalledWith('/');
    });
  });

  afterEach(jest.resetAllMocks);
});
