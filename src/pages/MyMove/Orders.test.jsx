import React from 'react';
import { mount } from 'enzyme';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { Orders } from './Orders';

import { getServiceMember, getOrdersForServiceMember, createOrders, patchOrders } from 'services/internalApi';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  getServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  createOrders: jest.fn().mockImplementation(() => Promise.resolve()),
  patchOrders: jest.fn().mockImplementation(() => Promise.resolve()),
}));

describe('Orders page', () => {
  const ordersOptions = [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ];

  const testProps = {
    serviceMemberId: '123',
    push: jest.fn(),
    context: { flags: { allOrdersTypes: true } },
    updateOrders: jest.fn(),
    updateServiceMember: jest.fn(),
  };

  describe('if there are no current orders', () => {
    it('does not load orders on mount', async () => {
      const { queryByRole } = render(<Orders {...testProps} />);

      await waitFor(() => {
        expect(queryByRole('heading', { name: 'Tell us about your move orders', level: 1 })).toBeInTheDocument();
        expect(getOrdersForServiceMember).not.toHaveBeenCalled();
      });
    });

    it('next button creates the orders and goes to the Upload Orders step', async () => {
      const testOrdersValues = {
        // NO ID
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

      getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));

      createOrders.mockImplementation(() =>
        Promise.resolve({
          ...testOrdersValues,
          id: 'newOrdersId',
        }),
      );

      // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
      const { queryByRole } = render(<Orders {...testProps} currentOrders={testOrdersValues} />);

      await waitFor(() => {
        // TODO - fill out form instead of partial current orders
        const submitButton = queryByRole('button', { name: 'Next' });
        expect(submitButton).toBeInTheDocument();
        userEvent.click(submitButton);
      });

      await waitFor(() => {
        expect(createOrders).toHaveBeenCalled();
      });

      expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
      expect(getServiceMember).toHaveBeenCalledWith(testProps.serviceMemberId);
      expect(testProps.updateServiceMember).toHaveBeenCalled();
      expect(testProps.push).toHaveBeenCalledWith('/orders/upload');
    });
  });

  describe('if there are current orders', () => {
    it('loads orders on mount', async () => {
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

      getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));

      const { queryByText } = render(<Orders {...testProps} currentOrders={{ id: 'testOrders' }} />);

      expect(queryByText('Loading, please wait...')).toBeInTheDocument();

      await waitFor(() => {
        expect(queryByText('Loading, please wait...')).not.toBeInTheDocument();
        expect(getOrdersForServiceMember).toHaveBeenCalled();
        expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
      });
    });

    it('next button patches the orders and goes to the Upload Orders step', async () => {
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

      getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));
      patchOrders.mockImplementation(() => Promise.resolve(testOrdersValues));

      // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
      const { queryByRole } = render(<Orders {...testProps} currentOrders={testOrdersValues} />);

      await waitFor(() => {
        const submitButton = queryByRole('button', { name: 'Next' });
        expect(submitButton).toBeInTheDocument();
        userEvent.click(submitButton);
      });

      await waitFor(() => {
        expect(patchOrders).toHaveBeenCalled();
      });

      // updateOrders gets called twice: once on load, once on submit
      expect(testProps.updateOrders).toHaveBeenNthCalledWith(1, testOrdersValues);
      expect(testProps.updateOrders).toHaveBeenNthCalledWith(2, testOrdersValues);
      expect(getServiceMember).not.toHaveBeenCalled();
      expect(testProps.updateServiceMember).not.toHaveBeenCalled();
      expect(testProps.push).toHaveBeenCalledWith('/orders/upload');
    });
  });

  it('back button goes to the Home page', async () => {
    const { queryByText } = render(<Orders {...testProps} />);

    const backButton = queryByText('Back');
    await waitFor(() => {
      expect(backButton).toBeInTheDocument();
    });

    userEvent.click(backButton);
    expect(testProps.push).toHaveBeenCalledWith('/');
  });

  it('shows an error if the API returns an error', async () => {
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

    getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));

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

    // Need to provide complete & valid initial values because we aren't testing the form here, and just want to submit immediately
    const { queryByText } = render(<Orders {...testProps} currentOrders={testOrdersValues} />);

    await waitFor(() => {
      const submitButton = queryByText('Next');
      expect(submitButton).toBeInTheDocument();
      userEvent.click(submitButton);
    });

    await waitFor(() => {
      expect(patchOrders).toHaveBeenCalled();
    });

    expect(queryByText('A server error occurred saving the orders')).toBeInTheDocument();
    expect(testProps.updateOrders).toHaveBeenCalledTimes(1);
    expect(testProps.push).not.toHaveBeenCalled();
  });

  describe('with the allOrdersType feature flag set to true', () => {
    it('passes all orders types into the form', async () => {
      const wrapper = mount(<Orders {...testProps} context={{ flags: { allOrdersTypes: true } }} />);
      await waitFor(() => {
        expect(wrapper.find('OrdersInfoForm').prop('ordersTypeOptions')).toEqual(ordersOptions);
      });
    });
  });

  describe('with the allOrdersType feature flag set to false', () => {
    it('passes only the PCS option into the form', async () => {
      const wrapper = mount(<Orders {...testProps} context={{ flags: { allOrdersTypes: false } }} />);
      await waitFor(() => {
        expect(wrapper.find('OrdersInfoForm').prop('ordersTypeOptions')).toEqual([ordersOptions[0]]);
      });
    });
  });

  afterEach(jest.resetAllMocks);
});
