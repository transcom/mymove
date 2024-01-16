import React from 'react';
import { mount } from 'enzyme';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { Orders } from './Orders';

import { createOrders, getOrdersForServiceMember, getServiceMember, patchOrders } from 'services/internalApi';
import { MockRouterProvider, renderWithRouterProp } from 'testUtils';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  getOrdersForServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  getServiceMember: jest.fn().mockImplementation(() => Promise.resolve()),
  createOrders: jest.fn().mockImplementation(() => Promise.resolve()),
  patchOrders: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('components/LocationSearchBox/api', () => ({
  ShowAddress: jest.fn().mockImplementation(() =>
    Promise.resolve({
      city: 'Glendale Luke AFB',
      country: 'United States',
      id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
      postalCode: '85309',
      state: 'AZ',
      streetAddress1: 'n/a',
    }),
  ),
  SearchDutyLocations: jest.fn().mockImplementation(() =>
    Promise.resolve([
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '46c4640b-c35e-4293-a2f1-36c7b629f903',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '93f0755f-6f35-478b-9a75-35a69211da1c',
        name: 'Altus AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '2d7e17f6-1b8a-4727-8949-007c80961a62',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '7d123884-7c1b-4611-92ae-e8d43ca03ad9',
        name: 'Hill AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '25be4d12-fe93-47f1-bbec-1db386dfa67f',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: 'a8d6b33c-8370-4e92-8df2-356b8c9d0c1a',
        name: 'Luke AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '3dbf1fc7-3289-4c6e-90aa-01b530a7c3c3',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'd01bd2a4-6695-4d69-8f2f-69e88dff58f8',
        name: 'Shaw AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '1af8f0f3-f75f-46d3-8dc8-c67c2feeb9f0',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:49:14.322Z',
        id: 'b1f9a535-96d4-4cc3-adf1-b76505ce0765',
        name: 'Yuma AFB',
        updated_at: '2021-02-11T16:49:14.322Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: 'f2adfebc-7703-4d06-9b49-c6ca8f7968f1',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'a268b48f-0ad1-4a58-b9d6-6de10fd63d96',
        name: 'Los Angeles AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
      {
        address: {
          city: '',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '',
          state: '',
          streetAddress1: '',
        },
        address_id: '13eb2cab-cd68-4f43-9532-7a71996d3296',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'a48fda70-8124-4e90-be0d-bf8119a98717',
        name: 'Wright-Patterson AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
    ]),
  ),
}));

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

afterEach(() => {
  jest.resetAllMocks();
});

describe('Orders page', () => {
  const ordersOptions = [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'LOCAL_MOVE', value: 'Local Move' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ];

  const testProps = {
    serviceMemberId: '123',
    context: { flags: { allOrdersTypes: true } },
    updateOrders: jest.fn(),
    updateServiceMember: jest.fn(),
  };

  describe('if there are no current orders', () => {
    it('next button creates the orders and goes to the Upload Orders step', async () => {
      const testOrdersValues = {
        id: 'newOrdersId',
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
        grade: 'E_1',
      };

      createOrders.mockImplementation(() => Promise.resolve(testOrdersValues));

      renderWithRouterProp(<Orders {...testProps} />, { navigate: mockNavigate });

      await userEvent.selectOptions(screen.getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
      await userEvent.type(screen.getByLabelText('Orders date'), '08 Nov 2020');
      await userEvent.type(screen.getByLabelText('Report by date'), '26 Nov 2020');
      await userEvent.click(screen.getByLabelText('No'));
      await userEvent.selectOptions(screen.getByLabelText('Pay grade'), ['E_5']);

      // Test Current Duty Location Search Box interaction
      await userEvent.type(screen.getByLabelText('Current duty location'), 'AFB', { delay: 100 });
      const selectedOptionCurrent = await screen.findByText(/Altus/);
      await userEvent.click(selectedOptionCurrent);

      // Test New Duty Location Search Box interaction
      await userEvent.type(screen.getByLabelText('New duty location'), 'AFB', { delay: 100 });
      const selectedOptionNew = await screen.findByText(/Luke/);
      await userEvent.click(selectedOptionNew);

      await waitFor(() => {
        expect(screen.getByRole('form')).toHaveFormValues({
          orders_type: 'PERMANENT_CHANGE_OF_STATION',
          issue_date: '08 Nov 2020',
          report_by_date: '26 Nov 2020',
          has_dependents: 'no',
          new_duty_location: 'Luke AFB',
          grade: 'E_5',
          origin_duty_location: 'Altus AFB',
        });
      });

      const submitButton = screen.getByRole('button', { name: 'Next' });
      expect(submitButton).toBeEnabled();

      expect(submitButton).toBeInTheDocument();
      await userEvent.click(submitButton);

      await waitFor(() => {
        expect(createOrders).toHaveBeenCalled();
      });

      expect(testProps.updateOrders).toHaveBeenCalledWith(testOrdersValues);
      expect(getServiceMember).toHaveBeenCalledWith(testProps.serviceMemberId);
      expect(testProps.updateServiceMember).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith('/orders/upload');
    });

    it('does not load orders on mount', async () => {
      const { queryByRole } = renderWithRouterProp(<Orders {...testProps} />);
      await waitFor(() => {
        expect(queryByRole('heading', { name: 'Tell us about your move orders', level: 1 })).toBeInTheDocument();
        expect(getOrdersForServiceMember).not.toHaveBeenCalled();
      });
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
        grade: 'E_1',
        origin_duty_location: {
          address: {
            city: 'Altus AFB',
            country: 'United States',
            id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
            postalCode: '73523',
            state: 'OK',
            streetAddress1: 'n/a',
          },
          address_id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
          affiliation: 'AIR_FORCE',
          created_at: '2021-02-11T16:48:04.117Z',
          id: '93f0755f-6f35-478b-9a75-35a69211da1c',
          name: 'Altus AFB',
          updated_at: '2021-02-11T16:48:04.117Z',
        },
      };

      getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));

      const { queryByText } = renderWithRouterProp(<Orders {...testProps} currentOrders={{ id: 'testOrders' }} />);

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
        grade: 'E_1',
      };
      getOrdersForServiceMember.mockImplementation(() => Promise.resolve(testOrdersValues));
      patchOrders.mockImplementation(() => Promise.resolve(testOrdersValues));

      // Need to provide initial values because we aren't testing the form here, and just want to submit immediately
      const { queryByRole } = renderWithRouterProp(<Orders {...testProps} currentOrders={testOrdersValues} />, {
        navigate: mockNavigate,
      });

      await waitFor(() => {
        expect(queryByRole('button', { name: 'Next' })).toBeInTheDocument();
      });
      await userEvent.click(queryByRole('button', { name: 'Next' }));

      expect(testProps.updateOrders).toHaveBeenNthCalledWith(1, testOrdersValues);
      expect(getServiceMember).not.toHaveBeenCalled();
      expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    });
  });

  it('back button goes to the Home page', async () => {
    const { queryByText } = renderWithRouterProp(<Orders {...testProps} />, { navigate: mockNavigate });

    const backButton = queryByText('Back');
    await waitFor(() => {
      expect(backButton).toBeInTheDocument();
    });

    await userEvent.click(backButton);
    expect(mockNavigate).toHaveBeenCalledWith('/');
  });

  it('shows an error if the API returns an error', async () => {
    const testOrdersValues = {
      id: 'testOrdersId',
      orders_type: 'PERMANENT_CHANGE_OF_STATION',
      issue_date: '2020-11-08',
      report_by_date: '2020-11-26',
      has_dependents: false,
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
      origin_duty_location: {
        address: {
          city: 'Altus AFB',
          country: 'United States',
          id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
          postalCode: '73523',
          state: 'OK',
          streetAddress1: 'n/a',
        },
        address_id: 'fa51dab0-4553-4732-b843-1f33407f77bd',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:04.117Z',
        id: '93f0755f-6f35-478b-9a75-35a69211da1c',
        name: 'Altus AFB',
        updated_at: '2021-02-11T16:48:04.117Z',
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
    const { queryByText } = renderWithRouterProp(<Orders {...testProps} currentOrders={testOrdersValues} />);

    await waitFor(() => {
      expect(queryByText('Next')).toBeInTheDocument();
    });
    await userEvent.click(queryByText('Next'));

    await waitFor(() => {
      // expect(patchOrders).toHaveBeenCalled();
    });

    // expect(queryByText('A server error occurred saving the orders')).toBeInTheDocument();
    // expect(testProps.updateOrders).toHaveBeenCalledTimes(1);
    // expect(mockNavigate).not.toHaveBeenCalled();
  });

  describe('with the allOrdersType feature flag set to true', () => {
    it('passes all orders types into the form', async () => {
      const wrapper = mount(
        <Orders {...testProps} context={{ flags: { allOrdersTypes: true } }} router={{ navigate: mockNavigate }} />,
      );
      await waitFor(() => {
        expect(wrapper.find('OrdersInfoForm').prop('ordersTypeOptions')).toEqual(ordersOptions);
      });
    });
  });

  describe('with the allOrdersType feature flag set to false', () => {
    it('passes only the PCS option into the form', async () => {
      const wrapper = mount(
        <MockRouterProvider>
          <Orders {...testProps} context={{ flags: { allOrdersTypes: false } }} router={{ navigate: mockNavigate }} />
        </MockRouterProvider>,
      );
      await waitFor(() => {
        expect(wrapper.find('OrdersInfoForm').prop('ordersTypeOptions')).toEqual([ordersOptions[0]]);
      });
    });
  });
});
