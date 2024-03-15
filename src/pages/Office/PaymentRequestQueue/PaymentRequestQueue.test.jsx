/* eslint-disable no-restricted-syntax */
import React from 'react';
import Select from 'react-select';
import { shallow, mount } from 'enzyme';
import { QueryClient } from '@tanstack/react-query';
import { NavTab } from 'react-router-tabs/cjs/react-router-tabs.min';
import { MemoryRouter, Route, Routes, useParams } from 'react-router';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import * as reactRouterDom from 'react-router-dom';

import PaymentRequestQueue from './PaymentRequestQueue';

import { MockProviders } from 'testUtils';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import PrivateRoute from 'containers/PrivateRoute';
import { roleTypes } from 'constants/userRoles';
import { initialState } from 'reducers/tacValidation';
import { MOVE_STATUS_OPTIONS, PAYMENT_REQUEST_STATUS_OPTIONS } from 'constants/queues';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'), // this line preserves the non-hook exports
  useParams: jest.fn(), // mock useParams
  useNavigate: jest.fn(), // mock useNavigate if needed
}));
jest.setTimeout(60000);
jest.mock('hooks/queries', () => ({
  useUserQueries: () => {
    return {
      isLoading: false,
      isError: false,
      data: {
        office_user: { transportation_office: { gbloc: 'TEST' } },
      },
    };
  },
  usePaymentRequestQueueQueries: () => {
    return {
      queueResult: {
        data: [
          {
            age: 0.8477863,
            customer: {
              agency: 'ARMY',
              dodID: '3305957632',
              eTag: 'MjAyMC0xMC0xNVQyMzo0ODozNC41ODQxOTZa',
              email: 'leo_spaceman_sm@example.com',
              first_name: 'Leo',
              id: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
              last_name: 'Spacemen',
              phone: '555-555-5555',
              userID: 'c4d59e2b-bff0-4fce-a31f-26a19b1ad34a',
            },
            departmentIndicator: 'AIR_AND_SPACE_FORCE',
            id: 'a2c34dba-015f-4f96-a38b-0c0b9272e208',
            locator: 'R993T7',
            moveID: '5d4b25bb-eb04-4c03-9a81-ee0398cb779e',
            originGBLOC: 'LKNQ',
            status: 'PENDING',
            submittedAt: '2020-10-15T23:48:35.420Z',
            originDutyLocation: {
              name: 'Scott AFB',
            },
          },
        ],
        totalCount: 1,
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
  useMoveSearchQueries: () => {
    return {
      searchResult: {
        data: [
          {
            age: 0.8477863,
            customer: {
              agency: 'ARMY',
              dodID: '3305957632',
              eTag: 'MjAyMC0xMC0xNVQyMzo0ODozNC41ODQxOTZa',
              email: 'leo_spaceman_sm@example.com',
              first_name: 'Leo',
              id: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
              last_name: 'Spacemen',
              phone: '555-555-5555',
              userID: 'c4d59e2b-bff0-4fce-a31f-26a19b1ad34a',
            },
            departmentIndicator: 'AIR_AND_SPACE_FORCE',
            id: 'a2c34dba-015f-4f96-a38b-0c0b9272e208',
            locator: 'R993T7',
            moveID: '5d4b25bb-eb04-4c03-9a81-ee0398cb779e',
            originGBLOC: 'LKNQ',
            status: 'PENDING',
            submittedAt: '2020-10-15T23:48:35.420Z',
            originDutyLocation: {
              name: 'Scott AFB',
            },
          },
        ],
        page: 0,
        perPage: 20,
        totalCount: 1,
      },
      isLoading: false,
      isError: false,
      isSuccess: true,
    };
  },
}));
const ExpectedPaymentRequestQueueColumns = [
  'Customer name',
  'DoD ID',
  'Status',
  'Age',
  'Submitted',
  'Move Code',
  'Branch',
  'Origin GBLOC',
  'Origin Duty Location',
];
const ExpectedOptions = ['Move Code', 'DOD ID', 'CustomerName'];
describe('PaymentRequestQueue', () => {
  const client = new QueryClient();

  it('renders the queue results text', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: '' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    screen.debug();
    expect(screen.queryByText('Payment requests (1)')).toBeInTheDocument();
  });

  it('renders the correct column headers', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: '' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    const columns = [
      'Customer name',
      'DoD ID',
      'Status',
      'Age',
      'Submitted',
      'Move Code',
      'Branch',
      'Origin GBLOC',
      'Origin Duty Location',
    ];

    // eslint-disable-next-line no-restricted-syntax, guard-for-in
    for (const col in ExpectedPaymentRequestQueueColumns) {
      expect(screen.findByText(columns[col], { selector: 'th' }));
    }
  });

  it('renders the table with data and expected values', () => {
    const wrapper = mount(
      <MockProviders client={client}>
        {/* eslint-disable-next-line react/jsx-props-no-spreading */}
        <PaymentRequestQueue />
      </MockProviders>,
    );
    expect(wrapper.find('Table').exists()).toBe(true);
    expect(wrapper.find('tbody tr').length).toBe(1);

    expect(wrapper.find('tbody tr td').at(0).text()).toBe('Spacemen, Leo');
    expect(wrapper.find('tbody tr td').at(1).text()).toBe('3305957632');
    expect(wrapper.find('tbody tr td').at(2).text()).toBe('Payment requested');
    expect(wrapper.find('tbody tr td').at(3).text()).toBe('Less than 1 day');
    expect(wrapper.find('tbody tr td').at(4).text()).toBe('15 Oct 2020');
    expect(wrapper.find('tbody tr td').at(5).text()).toBe('R993T7');
    expect(wrapper.find('tbody tr td').at(6).text()).toBe('Army');
    expect(wrapper.find('tbody tr td').at(7).text()).toBe('LKNQ');
    expect(wrapper.find('tbody tr td').at(8).text()).toBe('Scott AFB');
  });

  it('applies the sort to the age column in descending direction', () => {
    const wrapper = mount(
      <MockProviders client={client}>
        {/* eslint-disable-next-line react/jsx-props-no-spreading */}
        <PaymentRequestQueue />
      </MockProviders>,
    );
    expect(wrapper.find({ 'data-testid': 'age' }).at(0).hasClass('sortDescending')).toBe(true);
  });

  it('toggles the sort direction when clicked', () => {
    const wrapper = mount(
      <MockProviders client={client}>
        {/* eslint-disable-next-line react/jsx-props-no-spreading */}
        <PaymentRequestQueue />
      </MockProviders>,
    );
    const ageHeading = wrapper.find({ 'data-testid': 'age' }).at(0);

    ageHeading.simulate('click');
    wrapper.update();

    // no sort direction should be applied
    expect(wrapper.find({ 'data-testid': 'age' }).at(0).hasClass('sortAscending')).toBe(false);
    expect(wrapper.find({ 'data-testid': 'age' }).at(0).hasClass('sortDescending')).toBe(false);

    ageHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'age' }).at(0).hasClass('sortAscending')).toBe(true);

    const nameHeading = wrapper.find({ 'data-testid': 'lastName' }).at(0);
    nameHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'lastName' }).at(0).hasClass('sortAscending')).toBe(true);
  });

  it('filters the queue', () => {
    const wrapper = mount(
      <MockProviders client={client}>
        {/* eslint-disable-next-line react/jsx-props-no-spreading */}
        <PaymentRequestQueue />
      </MockProviders>,
    );
    const input = wrapper.find(Select).at(0).find('input');
    input.simulate('keyDown', { key: 'ArrowDown', keyCode: 40 });
    input.simulate('keyDown', { key: 'Enter', keyCode: 13 });

    wrapper.update();
    expect(wrapper.find('[data-testid="multi-value-container"]').text()).toEqual('Payment requested');
  });

  it('Displays the payment request ', async () => {
    // Setup initial state and mocks as necessary

    reactRouterDom.useParams.mockReturnValue({ queueType: 'Search' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    // Simulate user input and form submission
    const searchInput = screen.getByTestId('searchText'); // Adjust based on your actual labels/inputs
    await userEvent.type(searchInput, 'R993T7');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed
    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });

  it(' renders Search and Payment Request Queue tabs', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'Search' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.getByTestId('payment-request-queue-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('search-tab-link')).toBeInTheDocument();
  });
  it('renders SearchResultsTable when Search tab is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'Search' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.queryByTestId('payment-request-queue')).not.toBeInTheDocument();
    expect(screen.queryByTestId('move-search')).toBeInTheDocument();
  });
  it('renders TableQueue when Payment Request Queue tab is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'PaymentRequests' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    expect(screen.queryByTestId('payment-request-queue')).toBeInTheDocument();
    expect(screen.queryByTestId('move-search')).not.toBeInTheDocument();
  });
  it('submits search form and displays search results', async () => {
    // Setup initial state and mocks as necessary

    reactRouterDom.useParams.mockReturnValue({ queueType: 'Search' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    // Simulate user input and form submission
    const searchInput = screen.getByTestId('searchText'); // Adjust based on your actual labels/inputs
    await userEvent.type(searchInput, 'R993T7');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed
    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });
  it('submits search form and displays possible filters for status', async () => {
    // Setup initial state and mocks as necessary

    reactRouterDom.useParams.mockReturnValue({ queueType: 'Search' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    // Simulate user input and form submission
    const searchInput = screen.getByTestId('searchText'); // Adjust based on your actual labels/inputs
    await userEvent.type(searchInput, 'R993T7');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed

    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });
  it('Has 3 options for searches', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'Search' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    const options = ['Move Code', 'DoD ID', 'Customer Name'];

    // eslint-disable-next-line no-restricted-syntax, guard-for-in
    for (const col in options) {
      expect(screen.findByLabelText(options[col]));
    }
  });
  it('Has all status options for payment request search', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'Search' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    // eslint-disable-next-line no-restricted-syntax, guard-for-in
    for (const col in PAYMENT_REQUEST_STATUS_OPTIONS) {
      expect(screen.findByLabelText(PAYMENT_REQUEST_STATUS_OPTIONS[col]));
    }
  });

  it('Has all status options for payment request queue', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'PaymentRequest' });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    // eslint-disable-next-line no-restricted-syntax, guard-for-in
    for (const col in PAYMENT_REQUEST_STATUS_LABELS) {
      expect(screen.findByLabelText(PAYMENT_REQUEST_STATUS_OPTIONS[col]));
    }
  });
});
