import React from 'react';
import { mount } from 'enzyme';
import { QueryClient } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import * as reactRouterDom from 'react-router-dom';

import PaymentRequestQueue from './PaymentRequestQueue';

import { MockProviders } from 'testUtils';
import { generalRoutes, tioRoutes } from 'constants/routes';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'), // this line preserves the non-hook exports
  useParams: jest.fn(),
  useNavigate: jest.fn(),
}));
jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve()),
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
            assignedTo: {
              firstName: 'Alice',
              lastName: 'Bob',
              officeUserId: '404011d1-052a-4c34-b2e0-71dd5082718a',
            },
            customer: {
              agency: 'COAST_GUARD',
              dodID: '3305957632',
              emplid: '1253694',
              eTag: 'MjAyMC0xMC0xNVQyMzo0ODozNC41ODQxOTZa',
              email: 'leo_spaceman_sm@example.com',
              first_name: 'Leo',
              id: '6ac40a00-e762-4f5f-b08d-3ea72a8e4b63',
              last_name: 'Spacemen',
              phone: '555-555-5555',
              userID: 'c4d59e2b-bff0-4fce-a31f-26a19b1ad34a',
            },
            departmentIndicator: 'COAST_GUARD',
            id: 'a2c34dba-015f-4f96-a38b-0c0b9272e208',
            locator: 'R993T7',
            moveID: '5d4b25bb-eb04-4c03-9a81-ee0398cb779e',
            originGBLOC: 'LKNQ',
            counselingOffice: '67592323-fc7e-4b35-83a7-57faa53b7acf',
            status: 'PENDING',
            submittedAt: '2020-10-15T23:48:35.420Z',
            originDutyLocation: {
              name: 'Scott AFB',
            },
            lockExpiresAt: '2099-10-15T23:48:35.420Z',
            lockedByOfficeUserID: '2744435d-7ba8-4cc5-bae5-f302c72c966e',
          },
          {
            age: 0.8477863,
            availableOfficeUsers: [
              {
                firstName: 'Alice',
                lastName: 'Bob',
                officeUserId: '404011d1-052a-4c34-b2e0-71dd5082718a',
              },
            ],
            customer: {
              agency: 'NAVY',
              cacValidated: true,
              eTag: 'MjAyNC0xMC0xMFQyMjoyNDo1My4xNjYxNjNa',
              dodID: '1234567890',
              edipi: '1234567',
              email: '20241010222310-ae019978709c@example.com',
              first_name: 'Ooga',
              id: 'f23b5293-8ef0-453d-bd51-7c21d9730fcb',
              last_name: 'Booga',
              middle_name: '',
              phone: '211-111-1111',
              phoneIsPreferred: true,
              secondaryTelephone: '',
              suffix: '',
              userID: 'e9d421af-b598-4ff1-b102-d8b14d414129',
            },
            departmentIndicator: 'COAST_GUARD',
            id: '32b90458-2171-4451-bb0a-c5c0db897e34',
            locator: '0OOGAB',
            moveID: '8f29e53d-e816-4476-bfee-f38d07b94f2d',
            originGBLOC: 'LKNQ',
            counselingOffice: '67592323-fc7e-4b35-83a7-57faa53b7acf',
            status: 'PENDING',
            submittedAt: '2020-10-17T23:48:35.420Z',
            originDutyLocation: {
              name: 'Scott AFB',
            },
            lockExpiresAt: '2099-10-15T23:48:35.420Z',
            lockedByOfficeUserID: '2744435d-7ba8-4cc5-bae5-f302c72c966e',
          },
        ],
        totalCount: 2,
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
const SEARCH_OPTIONS = ['Move Code', 'DoD ID', 'Customer Name', 'Payment Request Number'];

describe('PaymentRequestQueue', () => {
  const client = new QueryClient();

  it('renders the queue results text', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    expect(screen.queryByText('Payment requests (2)')).toBeInTheDocument();
  });

  it('renders the table with data and expected values without queue management ff', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });
    const wrapper = mount(
      <MockProviders client={client}>
        {/* eslint-disable-next-line react/jsx-props-no-spreading */}
        <PaymentRequestQueue />
      </MockProviders>,
    );
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );

    const paymentRequests = wrapper.find('tbody tr');
    const firstPaymentRequest = paymentRequests.at(0);

    expect(firstPaymentRequest.find('td.customerName').text()).toBe('Spacemen, Leo');
    expect(firstPaymentRequest.find('td.dodID').text()).toBe('3305957632');
    expect(firstPaymentRequest.find('td.emplid').text()).toBe('1253694');
    expect(firstPaymentRequest.find('td.status').text()).toBe('Payment requested');
    expect(firstPaymentRequest.find('td.age').text()).toBe('Less than 1 day');
    expect(firstPaymentRequest.find('td.submittedAt').text()).toBe('15 Oct 2020');
    expect(firstPaymentRequest.find('td.locator').text()).toBe('R993T7');
    expect(firstPaymentRequest.find('td.branch').text()).toBe('Coast Guard');
    expect(firstPaymentRequest.find('td.originGBLOC').text()).toBe('LKNQ');
    expect(firstPaymentRequest.find('td.originDutyLocation').text()).toBe('Scott AFB');
    expect(firstPaymentRequest.find('td.counselingOffice').text()).toBe('67592323-fc7e-4b35-83a7-57faa53b7acf');

    const secondPaymentRequest = paymentRequests.at(1);
    expect(secondPaymentRequest.find('td.customerName').text()).toBe('Booga, Ooga');
    expect(secondPaymentRequest.find('td.dodID').text()).toBe('1234567890');
    expect(secondPaymentRequest.find('td.emplid').text()).toBe('');
    expect(secondPaymentRequest.find('td.status').text()).toBe('Payment requested');
    expect(secondPaymentRequest.find('td.age').text()).toBe('Less than 1 day');
    expect(secondPaymentRequest.find('td.submittedAt').text()).toBe('17 Oct 2020');
    expect(secondPaymentRequest.find('td.locator').text()).toBe('0OOGAB');
    expect(secondPaymentRequest.find('td.branch').text()).toBe('Navy');
    expect(secondPaymentRequest.find('td.originGBLOC').text()).toBe('LKNQ');
    expect(secondPaymentRequest.find('td.originDutyLocation').text()).toBe('Scott AFB');
    expect(secondPaymentRequest.find('td.counselingOffice').text()).toBe('67592323-fc7e-4b35-83a7-57faa53b7acf');
  });

  it('renders the table with data and expected values with queue management ff', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });
    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue isQueueManagementFFEnabled />
      </reactRouterDom.BrowserRouter>,
    );

    await waitFor(() => {
      const assignedSelect = screen.queryAllByTestId('assigned-col')[0];
      expect(assignedSelect).toBeInTheDocument();
    });
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

    const nameHeading = wrapper.find({ 'data-testid': 'customerName' }).at(0);
    nameHeading.simulate('click');
    wrapper.update();

    expect(wrapper.find({ 'data-testid': 'customerName' }).at(0).hasClass('sortAscending')).toBe(true);
  });

  it('displays the payment request ', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // Simulate user input and form submission
    const searchInput = screen.getByTestId('searchText');
    await userEvent.type(searchInput, 'R993T7');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed
    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });

  it(' renders Search and Payment Request Queue tabs', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    expect(screen.getByTestId('payment-request-queue-tab-link')).toBeInTheDocument();
    expect(screen.getByTestId('search-tab-link')).toBeInTheDocument();
  });
  it('renders SearchResultsTable when Search tab is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    expect(screen.queryByTestId('payment-request-queue')).not.toBeInTheDocument();
    expect(screen.queryByTestId('move-search')).toBeInTheDocument();
  });
  it('renders TableQueue when Payment Request Queue tab is selected', () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    expect(screen.queryByTestId('payment-request-queue')).toBeInTheDocument();
    expect(screen.queryByTestId('move-search')).not.toBeInTheDocument();
  });
  it('submits search form and displays search results', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // Simulate user input and form submission
    const searchInput = screen.getByTestId('searchText');
    await userEvent.type(searchInput, 'R993T7');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed
    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });
  it('searches by Move Code and displays possible filters for status', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // Simulate user input and form submission
    const searchSelection = screen.getByLabelText('Move Code');
    await userEvent.click(searchSelection);

    const searchInput = screen.getByTestId('searchText');
    await userEvent.type(searchInput, 'R993T7');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed

    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });
  it('searches by Customer Name', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // Simulate user input and form submission
    const searchSelection = screen.getByLabelText('Move Code');
    await userEvent.click(searchSelection);

    const searchInput = screen.getByTestId('searchText');
    await userEvent.type(searchInput, 'R993T7');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed

    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });
  it('searches by DOD ID', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // Simulate user input and form submission
    const searchSelection = screen.getByLabelText('Customer Name');
    await userEvent.click(searchSelection);

    const searchInput = screen.getByTestId('searchText');
    await userEvent.type(searchInput, '3305957632');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed

    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });
  it('searches by Move Code and displays possible filters for status', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // Simulate user input and form submission
    const searchSelection = screen.getByLabelText('Payment Request Number');
    await userEvent.click(searchSelection);

    const searchInput = screen.getByTestId('searchText');
    await userEvent.type(searchInput, '1234-5678-9');
    await userEvent.click(screen.getByTestId('searchTextSubmit'));
    // Assert search results are displayed

    expect(screen.queryByText('Results (1)')).toBeInTheDocument();
    expect(screen.queryByTestId('table-queue')).toBeInTheDocument();
  });
  it('has 4 options for searches', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: generalRoutes.QUEUE_SEARCH_PATH });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // This pattern allows minimal test changes if the search options were ever to change.
    SEARCH_OPTIONS.forEach((option) => expect(screen.findByLabelText(option)));
  });
  it('only displays payment requests with a status of Payment Requested', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    // expect Payment requested status to appear in the TIO queue
    expect(screen.getAllByText('Payment requested')).toHaveLength(2);
    // expect other statuses NOT to appear in the TIO queue
    expect(screen.queryByText('Deprecated')).not.toBeInTheDocument();
    expect(screen.queryByText('Error')).not.toBeInTheDocument();
    expect(screen.queryByText('Rejected')).not.toBeInTheDocument();
    expect(screen.queryByText('Reviewed')).not.toBeInTheDocument();
  });
  it('renders a 404 if a bad route is provided', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: 'BadRoute' });
    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    await expect(screen.getByText('Error - 404')).toBeInTheDocument();
    await expect(screen.getByText("We can't find the page you're looking for")).toBeInTheDocument();
  });
  it('renders a lock icon when move lock flag is on', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });
    isBooleanFlagEnabled.mockResolvedValue(true);

    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    await waitFor(() => {
      const lockIcon = screen.queryAllByTestId('lock-icon')[0];
      expect(lockIcon).toBeInTheDocument();
    });
  });
  it('does NOT render a lock icon when move lock flag is off', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });
    isBooleanFlagEnabled.mockResolvedValue(false);

    render(
      <MockProviders>
        <PaymentRequestQueue />
      </MockProviders>,
    );
    await waitFor(() => {
      const lockIcon = screen.queryByTestId('lock-icon');
      expect(lockIcon).not.toBeInTheDocument();
    });
  });
  it('renders assignedTo columns when the queue management flag is on', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });

    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue isQueueManagementFFEnabled />
      </reactRouterDom.BrowserRouter>,
    );
    await waitFor(() => {
      const assignedToColumn = screen.queryByText('Assigned');
      expect(assignedToColumn).toBeInTheDocument();
    });
  });
  it('does NOT render assignedTo columns when the queue management flag is on', async () => {
    reactRouterDom.useParams.mockReturnValue({ queueType: tioRoutes.PAYMENT_REQUEST_QUEUE });

    render(
      <reactRouterDom.BrowserRouter>
        <PaymentRequestQueue />
      </reactRouterDom.BrowserRouter>,
    );
    await waitFor(() => {
      const assignedToColumn = screen.queryByText('Assigned');
      expect(assignedToColumn).not.toBeInTheDocument();
    });
  });
});
