import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { debug } from 'jest-preview';

import ServicesCounselingAddOrders from './ServicesCounselingAddOrders';

import { MockProviders } from 'testUtils';
import { counselingCreateOrder } from 'services/ghcApi';
import { setCanAddOrders } from 'store/general/actions';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  counselingCreateOrder: jest.fn().mockImplementation(() => Promise.resolve()),
  showCounselingOffices: jest.fn().mockImplementation(() =>
    Promise.resolve({
      body: [
        {
          id: '3e937c1f-5539-4919-954d-017989130584',
          name: 'Albuquerque AFB',
        },
        {
          id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
          name: 'Glendale Luke AFB',
        },
      ],
    }),
  ),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
}));

jest.mock('store/general/actions', () => ({
  ...jest.requireActual('store/general/actions'),
  setCanAddOrders: jest.fn().mockImplementation(() => ({
    type: '',
    payload: '',
  })),
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
        provides_services_counseling: true,
      },
      {
        address: {
          city: 'Glendale Luke AFB',
          country: 'United States',
          id: 'fa51dab0-4553-4732-b843-1f33407f77bc',
          postalCode: '85309',
          state: 'AZ',
          streetAddress1: 'n/a',
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
      {
        address: {
          city: 'Cold',
          id: '00000000-0000-0000-0000-000000000000',
          postalCode: '12345',
          state: 'AK',
          streetAddress1: 'Worldly Avenue',
          isOconus: true,
        },
        address_id: '13eb2cab-cd68-4f43-9532-7a71996d3299',
        affiliation: 'AIR_FORCE',
        created_at: '2021-02-11T16:48:20.225Z',
        id: 'a48fda70-8124-4e90-be0d-bf8119a98717',
        name: 'Outta This World AFB',
        updated_at: '2021-02-11T16:48:20.225Z',
      },
    ]),
  ),
}));

const customer = {
  agency: 'NAVY',
  backupAddress: {
    city: 'Missoula',
    eTag: 'MjAyNC0wMy0yNlQyMjoyMjowNy45MzYwNzha',
    id: '9859fc20-3aa7-44b6-a0c4-0409817192e5',
    postalCode: '59801',
    state: 'MT',
    streetAddress1: '1909 Dearborn Ave',
    streetAddress2: '',
  },
  backup_contact: {
    email: 'paul.stonebraker@caci.com',
    name: 'Neighbor',
    phone: '324-321-1234',
  },
  current_address: {
    city: 'Missoula',
    eTag: 'MjAyNC0wMy0yNlQyMjoyMjowNy45Mjc5MDla',
    id: '9fe2a514-5497-4bda-8e4d-b560208a84b3',
    postalCode: '59801',
    state: 'MT',
    streetAddress1: '325 Dearborn Ave',
    streetAddress2: '',
  },
  dodID: '6565655555',
  eTag: 'MjAyNC0wMy0yNlQyMjoyMjowNy45NDAxNzJa',
  email: 'tio@3hhg.com',
  first_name: 'Tio',
  id: '31e86254-5b73-4822-9e0a-cfd0c4a028f2',
  last_name: 'Tester',
  phone: '324-321-1234',
  phoneIsPreferred: true,
  userID: '2527f5ad-b89b-49aa-822e-e4094a0059f1',
};

const fakeResponse = {
  orders: {
    '80ac4b6b-96a9-40d0-a897-b6ae6891854a': {
      agency: 'NAVY',
      customer: {
        agency: 'NAVY',
        backupAddress: {
          city: 'Missoula',
          eTag: 'MjAyNC0wMy0yNlQyMjoyMjowNy45MzYwNzha',
          id: '9859fc20-3aa7-44b6-a0c4-0409817192e5',
          postalCode: '59801',
          state: 'MT',
          streetAddress1: '1909 Dearborn Ave',
          streetAddress2: '',
        },
        backup_contact: {
          email: 'paul.stonebraker@caci.com',
          name: 'Neighbor',
          phone: '324-321-1234',
        },
        current_address: {
          city: 'Missoula',
          eTag: 'MjAyNC0wMy0yNlQyMjoyMjowNy45Mjc5MDla',
          id: '9fe2a514-5497-4bda-8e4d-b560208a84b3',
          postalCode: '59801',
          state: 'MT',
          streetAddress1: '325 Dearborn Ave',
          streetAddress2: '',
        },
        dodID: '6565655555',
        eTag: 'MjAyNC0wMy0yNlQyMjoyMjowNy45NDAxNzJa',
        email: 'tio@3hhg.com',
        first_name: 'TioT',
        id: '31e86254-5b73-4822-9e0a-cfd0c4a028f2',
        last_name: 'Tester',
        phone: '324-321-1234',
        phoneIsPreferred: true,
        userID: '2527f5ad-b89b-49aa-822e-e4094a0059f1',
      },
      customerID: '31e86254-5b73-4822-9e0a-cfd0c4a028f2',
      date_issued: '2024-03-01',
      department_indicator: '',
      destinationDutyLocation: {
        address: {
          city: 'Silverton',
          country: 'United States',
          eTag: 'MjAyNC0wMy0yMVQxODo0Mjo1My4zMjM5OTha',
          id: '74198926-358c-4944-abb4-5c2f851f9dd6',
          postalCode: '83867',
          state: 'ID',
          streetAddress1: 'n/a',
        },
        address_id: '74198926-358c-4944-abb4-5c2f851f9dd6',
        eTag: 'MjAyNC0wMy0yMVQxODo0Mjo1My4zMjM5OTha',
        id: '1ffa181f-8de0-46e4-b513-f8a9353de8ee',
        name: 'Silverton, ID 83867',
      },
      eTag: 'MjAyNC0wMy0yN1QwNToyMjo1Ni45NjMwNDFa',
      entitlement: {
        authorizedWeight: 14000,
        dependentsAuthorized: true,
        eTag: 'MjAyNC0wMy0yN1QwNToyMjo1Ni45NTg3MjVa',
        id: '282a8bc4-9470-4748-b727-8e86341652d8',
        proGearWeight: 2000,
        proGearWeightSpouse: 500,
        storageInTransit: 90,
        totalWeight: 14000,
      },
      first_name: 'TioT',
      grade: 'E_8',
      id: '80ac4b6b-96a9-40d0-a897-b6ae6891854a',
      last_name: 'Tester',
      methodOfPayment: 'Payment will be made using the Third-Party Payment System (TPPS) Automated Payment System',
      moveCode: 'MM8CXJ',
      moveTaskOrderID: 'ab5c867f-6274-48af-a579-116e45642b01',
      naics: '488510 - FREIGHT TRANSPORTATION ARRANGEMENT',
      order_type: 'PERMANENT_CHANGE_OF_STATION',
      order_type_detail: '',
      originDutyLocation: {
        address: {
          city: 'Salineno',
          country: 'United States',
          eTag: 'MjAyNC0wMy0yMVQxODo0Mjo1My4zMjM5OTha',
          id: 'fdb0fe38-c9d0-4c17-b956-454f858a0021',
          postalCode: '78585',
          state: 'TX',
          streetAddress1: 'n/a',
        },
        address_id: 'fdb0fe38-c9d0-4c17-b956-454f858a0021',
        eTag: 'MjAyNC0wMy0yMVQxODo0Mjo1My4zMjM5OTha',
        id: '0202254d-556a-42da-ad86-158c42ff9fb4',
        name: 'Salineno, TX 78585',
      },
      originDutyLocationGBLOC: 'HAFC',
      packingAndShippingInstructions:
        'Packaging, packing, and shipping instructions as identified in the Conformed Copy of HTC111-11-1-1112 Attachment 1 Performance Work Statement',
      report_by_date: '2024-03-31',
      supplyAndServicesCostEstimate:
        'Prices for services under this task order will be in accordance with rates provided in GHC Attachment 2 - Pricing Rate Table. It is the responsibility of the contractor to provide the estimated weight quantity to apply to services on this task order, when applicable (See Attachment 1 - Performance Work Statement).',
      uploaded_order_id: '814ff664-849d-4dfc-b4b2-d346a67444e3',
    },
  },
};

const renderWithMocks = () => {
  const testProps = { customer, setCanAddOrders: jest.fn() };
  render(
    <MockProviders>
      <ServicesCounselingAddOrders {...testProps} />
    </MockProviders>,
  );
};

describe('ServicesCounselingAddOrders component', () => {
  it('renders the Services Counseling Add Orders Form', async () => {
    renderWithMocks();

    const h1 = await screen.getByRole('heading', { name: 'Tell us about the orders', level: 1 });
    await waitFor(() => {
      expect(h1).toBeInTheDocument();
    });
  });

  it('routes to the move details page when the next button is clicked and we receive 200 res ', async () => {
    renderWithMocks();

    counselingCreateOrder.mockImplementation(() => Promise.resolve(fakeResponse));

    const user = userEvent.setup();

    await user.selectOptions(screen.getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    await user.type(screen.getByLabelText('Orders date'), '08 Nov 2020');
    await user.type(screen.getByLabelText('Report by date'), '26 Nov 2020');
    await user.click(screen.getByLabelText('No'));
    await user.selectOptions(screen.getByLabelText('Pay grade'), ['E-5']);

    // Test Current Duty Location Search Box interaction
    await user.type(screen.getByLabelText('Current duty location'), 'AFB', { delay: 500 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await user.click(selectedOptionCurrent);

    await user.type(screen.getByLabelText('New duty location'), 'AFB', { delay: 500 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await user.click(selectedOptionNew);

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    await waitFor(() => {
      expect(nextBtn).toBeEnabled();
    });

    // debug();
    await userEvent.click(nextBtn);

    await waitFor(() => {
      expect(setCanAddOrders).toHaveBeenCalledWith(false);
      expect(mockNavigate).toHaveBeenCalledWith('/counseling/moves/MM8CXJ/details');
    });
  });

  it('Displays the counseling office dropdown', async () => {
    renderWithMocks();

    counselingCreateOrder.mockImplementation(() => Promise.resolve(fakeResponse));

    const user = userEvent.setup();

    await user.selectOptions(screen.getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    await user.type(screen.getByLabelText('Orders date'), '08 Nov 2020');
    await user.type(screen.getByLabelText('Report by date'), '29 Nov 2020');
    await user.click(screen.getByLabelText('No'));
    await user.selectOptions(screen.getByLabelText('Pay grade'), ['E-5']);

    // Test Current Duty Location Search Box interaction
    await user.type(screen.getByLabelText('Current duty location'), 'AFB', { delay: 500 });
    const selectedOptionCurrent = await screen.findByText(/Hill/);
    await user.click(selectedOptionCurrent);

    await waitFor(async () => {
      expect(screen.getByLabelText(/Counseling office/));
    });

    await user.type(screen.getByLabelText('New duty location'), 'AFB', { delay: 500 });
    const selectedOptionNew = await screen.findByText(/Luke/);
    await user.click(selectedOptionNew);
    debug();
  });

  it('routes to the move details page when the next button is clicked for OCONUS orders', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    renderWithMocks();

    counselingCreateOrder.mockImplementation(() => Promise.resolve(fakeResponse));

    const user = userEvent.setup();

    await user.selectOptions(screen.getByLabelText('Orders type'), 'PERMANENT_CHANGE_OF_STATION');
    await user.type(screen.getByLabelText('Orders date'), '08 Nov 2020');
    await user.type(screen.getByLabelText('Report by date'), '26 Nov 2020');
    await user.click(screen.getByLabelText('No'));
    await user.selectOptions(screen.getByLabelText('Pay grade'), ['E-5']);

    await user.type(screen.getByLabelText('Current duty location'), 'AFB', { delay: 500 });
    const selectedOptionCurrent = await screen.findByText(/Altus/);
    await user.click(selectedOptionCurrent);

    await user.type(screen.getByLabelText('New duty location'), 'AFB', { delay: 500 });
    const selectedOptionNew = await screen.findByText(/Outta This World/);
    await user.click(selectedOptionNew);

    await user.click(screen.getByTestId('hasDependentsYes'));
    await user.click(screen.getByTestId('isAnAccompaniedTourYes'));
    await user.type(screen.getByTestId('dependentsUnderTwelve'), '2');
    await user.type(screen.getByTestId('dependentsTwelveAndOver'), '1');

    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    await waitFor(() => {
      expect(nextBtn).toBeEnabled();
    });

    await userEvent.click(nextBtn);

    await waitFor(() => {
      expect(setCanAddOrders).toHaveBeenCalledWith(false);
      expect(mockNavigate).toHaveBeenCalledWith('/counseling/moves/MM8CXJ/details');
    });
  });

  it('navigates back to Customer Info on back click', async () => {
    renderWithMocks();
    const backBtn = screen.getByRole('button', { name: 'Back' });
    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(-1);
    });
  });
});
