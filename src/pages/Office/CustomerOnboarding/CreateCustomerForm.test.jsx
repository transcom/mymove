import React from 'react';
import { render, fireEvent, waitFor, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router';

import { CreateCustomerForm } from './CreateCustomerForm';

import { MockProviders } from 'testUtils';
import { createCustomerWithOktaOption, searchLocationByZipCityState } from 'services/ghcApi';
import { servicesCounselingRoutes } from 'constants/routes';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import departmentIndicators from 'constants/departmentIndicators';
import { roleTypes } from 'constants/userRoles';

const mockPickupLocation = [
  {
    city: 'BEVERLY HILLS',
    usPostRegionCitiesId: 'a13806fc-0e7d-4dc3-91ca-b802d9da50f1',
    postalCode: '90210',
    state: 'CA',
    county: 'LOS ANGELES',
  },
  {
    city: 'DRYDEN',
    usPostRegionCitiesId: 'a13806fc-0e7d-4dc3-91ca-b802d9da50f1',
    postalCode: '04225',
    state: 'ME',
    county: 'FRANKLIN',
  },
];

const mockSearchPickupLocation = () => Promise.resolve(mockPickupLocation);

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  createCustomerWithOktaOption: jest.fn(),
  searchLocationByZipCityState: jest.fn(),
}));

jest.mock('store/flash/actions', () => ({
  setFlashMessage: jest.fn(),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

jest.mock('store/general/actions', () => ({
  setCanAddOrders: jest.fn(),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const serviceCounselorState = {
  auth: {
    activeRole: roleTypes.SERVICES_COUNSELOR,
    isLoading: false,
    isLoggedIn: true,
  },
  entities: {
    user: {
      userId123: {
        id: 'userId123',
        roles: [{ roleType: roleTypes.SERVICES_COUNSELOR }],
        office_user: {
          id: 'officeId123',
          first_name: 'Amanda',
          last_name: 'Gorman',
          transportation_office: {
            gbloc: 'ABCD',
          },
        },
      },
    },
  },
};

const fakePayload = {
  affiliation: 'ARMY',
  edipi: '1234567890',
  first_name: 'Shish',
  middle_name: 'Ka',
  last_name: 'Bob',
  suffix: 'Mr.',
  telephone: '555-555-5555',
  secondary_telephone: '999-867-5309',
  personal_email: 'tastyAndDelicious@mail.mil',
  phone_is_preferred: true,
  email_is_preferred: '',
  residential_address: {
    streetAddress1: '8711 S Hungry Ave.',
    streetAddress2: '',
    streetAddress3: '',
    city: 'Starving',
    state: 'OK',
    postalCode: '74133',
  },
  backup_mailing_address: {
    streetAddress1: '420 S. Munchies Lane',
    streetAddress2: '',
    streetAddress3: '',
    city: 'Mustang',
    state: 'KS',
    postalCode: '73064',
  },
  backup_contact: {
    name: 'Silly String',
    telephone: '666-666-6666',
    email: 'allOverDaPlace@mail.com',
  },
  create_okta_account: 'true',
  cac_user: 'false',
  is_safety_move: false,
  is_bluebark: 'false',
};

const fakeResponse = {
  createdCustomer: {
    '7575b55a-0e14-4f11-8e42-10232d22b135': {
      affiliation: 'string',
      firstName: 'John',
      lastName: 'Doe',
      telephone: '216-421-1392',
      personalEmail: '73sGJ6jq7cS%6@PqElR.WUzkqFNvtduyyA',
      suffix: 'Jr.',
      middleName: 'David',
      residentialAddress: {
        id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
        streetAddress1: '123 Main Ave',
        streetAddress2: 'Apartment 9000',
        streetAddress3: 'Montmârtre',
        city: 'Anytown',
        eTag: 'string',
        state: 'AL',
        postalCode: '90210',
        country: 'USA',
      },
      backupContact: {
        name: 'string',
        email: 'backupContact@mail.com',
        phone: '381-100-5880',
      },
      id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
      edipi: 'string',
      userID: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
      oktaID: 'string',
      oktaEmail: 'string',
      phoneIsPreferred: true,
      emailIsPreferred: true,
      secondaryTelephone: '499-793-2722',
      backupAddress: {
        id: 'c56a4180-65aa-42ec-a945-5fd21dec0538',
        streetAddress1: '123 Main Ave',
        streetAddress2: 'Apartment 9000',
        streetAddress3: 'Montmârtre',
        city: 'Anytown',
        eTag: 'string',
        state: 'AL',
        postalCode: '90210',
        country: 'USA',
      },
    },
  },
};

const safetyPayload = {
  affiliation: 'ARMY',
  edipi: '1234567890',
  first_name: 'Shish',
  middle_name: 'Ka',
  last_name: 'Bob',
  suffix: 'Mr.',
  telephone: '555-555-5555',
  secondary_telephone: '999-867-5309',
  personal_email: 'tastyAndDelicious@mail.mil',
  phone_is_preferred: true,
  email_is_preferred: '',
  residential_address: {
    streetAddress1: '8711 S Hungry Ave.',
    streetAddress2: '',
    streetAddress3: '',
    city: 'Starving',
    state: 'OK',
    postalCode: '74133',
  },
  backup_mailing_address: {
    streetAddress1: '420 S. Munchies Lane',
    streetAddress2: '',
    streetAddress3: '',
    city: 'Mustang',
    state: 'KS',
    postalCode: '73064',
  },
  backup_contact: {
    name: 'Silly String',
    telephone: '666-666-6666',
    email: 'allOverDaPlace@mail.com',
  },
  create_okta_account: 'true',
  cac_user: 'false',
  is_safety_move: 'true',
  is_bluebark: 'false',
};

const bluebarkPayload = {
  is_safety_move: 'false',
  is_bluebark: 'true',
};

const mockUserPrivileges = [
  {
    createdAt: '0001-01-01T00:00:00.000Z',
    privilegeType: 'safety',
    updatedAt: '0001-01-01T00:00:00.000Z',
  },
];

const testProps = {
  setFlashMessage: jest.fn(),
  setCanAddOrders: jest.fn(),
  userPrivileges: mockUserPrivileges,
};

const ordersPath = generatePath(servicesCounselingRoutes.BASE_CUSTOMERS_ORDERS_ADD_PATH, {
  customerId: '7575b55a-0e14-4f11-8e42-10232d22b135',
});

describe('CreateCustomerForm', () => {
  it('renders without crashing', async () => {
    render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    // checking that all headers exist
    expect(screen.getByText('Create Customer Profile')).toBeInTheDocument();
    expect(screen.getByText('Customer Affiliation')).toBeInTheDocument();
    expect(screen.getByText('Customer Name')).toBeInTheDocument();
    expect(screen.getByText('Contact Info')).toBeInTheDocument();
    expect(screen.getByText('Pickup Address')).toBeInTheDocument();
    expect(screen.getByText('Backup Address')).toBeInTheDocument();
    expect(screen.getByText('Backup Contact')).toBeInTheDocument();
    expect(screen.getByText('Okta Account')).toBeInTheDocument();
    expect(screen.getByText('Non-CAC Users')).toBeInTheDocument();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();
    expect(saveBtn).toBeDisabled();
    const cancelBtn = await screen.findByRole('button', { name: 'Cancel' });
    expect(cancelBtn).toBeInTheDocument();
  });

  it('renders emplid input if branch is coast guard', async () => {
    const { getByLabelText } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    await user.selectOptions(getByLabelText('Branch of service'), [departmentIndicators.COAST_GUARD]);
    expect(screen.getByText('EMPLID')).toBeInTheDocument();
  });

  it('payload can have an empty secondary phone number', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));
    searchLocationByZipCityState.mockImplementation(mockSearchPickupLocation);

    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await user.selectOptions(getByLabelText('Branch of service'), [fakePayload.affiliation]);

    await user.type(getByLabelText('First name'), fakePayload.first_name);
    await user.type(getByLabelText('Last name'), fakePayload.last_name);

    await user.type(getByLabelText('Best contact phone'), fakePayload.telephone);
    await user.type(getByLabelText('Personal email'), fakePayload.personal_email);
    await userEvent.type(getByTestId('edipiInput'), fakePayload.edipi);

    await user.type(getByTestId('residential_address.streetAddress1'), fakePayload.residential_address.streetAddress1);

    const locationBox = screen.getAllByRole('combobox');
    await act(async () => {
      await userEvent.type(locationBox[1], 'BEVERLY HILLS');
      const selectedResidentialLocation = await screen.findByText(/90210/);
      await userEvent.click(selectedResidentialLocation);
    });

    await userEvent.type(
      getByTestId('backup_mailing_address.streetAddress1'),
      safetyPayload.backup_mailing_address.streetAddress1,
    );

    await act(async () => {
      await userEvent.type(locationBox[2], 'DRYDEN');
      const selectedBackupLocation = await screen.findByText(/04225/);
      await userEvent.click(selectedBackupLocation);
    });

    await user.type(getByLabelText('Name'), fakePayload.backup_contact.name);
    await user.type(getByRole('textbox', { name: 'Email' }), fakePayload.backup_contact.email);
    await user.type(getByRole('textbox', { name: 'Phone' }), fakePayload.backup_contact.telephone);

    await userEvent.type(getByTestId('create-okta-account-yes'), fakePayload.create_okta_account);

    await userEvent.type(getByTestId('cac-user-no'), fakePayload.cac_user);

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });

    const waiter = waitFor(() => {
      expect(createCustomerWithOktaOption).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith(ordersPath, {
        state: {
          isSafetyMoveSelected: false,
          isBluebarkMoveSelected: false,
        },
      });
    });

    await user.click(saveBtn);
    await waiter;
    expect(mockNavigate).toHaveBeenCalled();

    expect(createCustomerWithOktaOption.mock.calls[0][0]).not.toHaveProperty('secondary_number');
  }, 20000);

  it('navigates the user on cancel click', async () => {
    const { getByText } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );
    fireEvent.click(getByText('Cancel'));
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalled();
    });
  });

  it('submits the form and navigates the user once all required fields are filled out', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));
    searchLocationByZipCityState.mockImplementation(mockSearchPickupLocation);

    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await user.selectOptions(getByLabelText('Branch of service'), [fakePayload.affiliation]);

    await user.type(getByLabelText('First name'), fakePayload.first_name);
    await user.type(getByLabelText('Last name'), fakePayload.last_name);

    await user.type(getByLabelText('Best contact phone'), fakePayload.telephone);
    await user.type(getByLabelText('Personal email'), fakePayload.personal_email);

    await userEvent.type(getByTestId('edipiInput'), fakePayload.edipi);

    await userEvent.type(
      getByTestId('residential_address.streetAddress1'),
      fakePayload.residential_address.streetAddress1,
    );

    const locationBox = screen.getAllByRole('combobox');

    await act(async () => {
      await userEvent.type(locationBox[1], 'BEVERLY HILLS');
      const selectedResidentialLocation = await screen.findByText(/90210/);
      await userEvent.click(selectedResidentialLocation);
    });

    await userEvent.type(
      getByTestId('backup_mailing_address.streetAddress1'),
      safetyPayload.backup_mailing_address.streetAddress1,
    );

    await act(async () => {
      await userEvent.type(locationBox[2], 'DRYDEN');
      const selectedBackupLocation = await screen.findByText(/04225/);
      await userEvent.click(selectedBackupLocation);
    });

    await userEvent.type(getByLabelText('Name'), fakePayload.backup_contact.name);
    await userEvent.type(getByRole('textbox', { name: 'Email' }), fakePayload.backup_contact.email);
    await userEvent.type(getByRole('textbox', { name: 'Phone' }), fakePayload.backup_contact.telephone);

    await userEvent.type(getByTestId('create-okta-account-yes'), fakePayload.create_okta_account);

    await userEvent.type(getByTestId('cac-user-no'), fakePayload.cac_user);

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });
    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(createCustomerWithOktaOption).toHaveBeenCalled();
      expect(testProps.setCanAddOrders).toHaveBeenCalledWith(true);
      expect(mockNavigate).toHaveBeenCalledWith(ordersPath, {
        state: {
          isSafetyMoveSelected: false,
          isBluebarkMoveSelected: false,
        },
      });
    });
  }, 20000);

  it('validates emplid against a coast guard member', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));
    searchLocationByZipCityState.mockImplementation(mockSearchPickupLocation);

    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await user.selectOptions(getByLabelText('Branch of service'), 'COAST_GUARD');

    await user.type(getByLabelText('First name'), fakePayload.first_name);
    await user.type(getByLabelText('Last name'), fakePayload.last_name);

    await user.type(getByLabelText('Best contact phone'), fakePayload.telephone);
    await user.type(getByLabelText('Personal email'), fakePayload.personal_email);

    await userEvent.type(getByTestId('edipiInput'), fakePayload.edipi);

    await userEvent.type(
      getByTestId('residential_address.streetAddress1'),
      fakePayload.residential_address.streetAddress1,
    );

    const locationBox = screen.getAllByRole('combobox');

    await act(async () => {
      await userEvent.type(locationBox[1], 'BEVERLY HILLS');
      const selectedResidentialLocation = await screen.findByText(/90210/);
      await userEvent.click(selectedResidentialLocation);
    });

    await userEvent.type(
      getByTestId('backup_mailing_address.streetAddress1'),
      safetyPayload.backup_mailing_address.streetAddress1,
    );

    await act(async () => {
      await userEvent.type(locationBox[2], 'DRYDEN');
      const selectedBackupLocation = await screen.findByText(/04225/);
      await userEvent.click(selectedBackupLocation);
    });

    await userEvent.type(getByLabelText('Name'), fakePayload.backup_contact.name);
    await userEvent.type(getByRole('textbox', { name: 'Email' }), fakePayload.backup_contact.email);
    await userEvent.type(getByRole('textbox', { name: 'Phone' }), fakePayload.backup_contact.telephone);

    await userEvent.type(getByTestId('create-okta-account-yes'), fakePayload.create_okta_account);

    await userEvent.type(getByTestId('cac-user-no'), fakePayload.cac_user);

    await waitFor(() => {
      expect(saveBtn).toBeDisabled(); // EMPLID not set yet
    });
    await userEvent.type(getByTestId('emplidInput'), '1234567');
    await waitFor(() => {
      expect(saveBtn).toBeEnabled(); // EMPLID is set now, all validations true
    });
  }, 20000);

  it('allows safety privileged users to pass safety move status to orders screen', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    searchLocationByZipCityState.mockImplementation(mockSearchPickupLocation);

    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await userEvent.type(getByTestId('is-safety-move-yes'), safetyPayload.is_safety_move);
    await userEvent.type(getByTestId('is-bluebark-no'), safetyPayload.is_bluebark);
    await user.selectOptions(getByLabelText('Branch of service'), [safetyPayload.affiliation]);

    await user.type(getByLabelText('First name'), safetyPayload.first_name);
    await user.type(getByLabelText('Last name'), safetyPayload.last_name);

    await user.type(getByLabelText('Best contact phone'), safetyPayload.telephone);
    await user.type(getByLabelText('Personal email'), safetyPayload.personal_email);

    await userEvent.type(
      getByTestId('residential_address.streetAddress1'),
      safetyPayload.residential_address.streetAddress1,
    );

    const locationBox = screen.getAllByRole('combobox');

    await act(async () => {
      await userEvent.type(locationBox[1], 'BEVERLY HILLS');
      const selectedResidentialLocation = await screen.findByText(/90210/);
      await userEvent.click(selectedResidentialLocation);
    });

    await userEvent.type(
      getByTestId('backup_mailing_address.streetAddress1'),
      safetyPayload.backup_mailing_address.streetAddress1,
    );

    await act(async () => {
      await userEvent.type(locationBox[2], 'DRYDEN');
      const selectedBackupLocation = await screen.findByText(/04225/);
      await userEvent.click(selectedBackupLocation);
    });

    await userEvent.type(getByLabelText('Name'), safetyPayload.backup_contact.name);
    await userEvent.type(getByRole('textbox', { name: 'Email' }), safetyPayload.backup_contact.email);
    await userEvent.type(getByRole('textbox', { name: 'Phone' }), safetyPayload.backup_contact.telephone);

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });
    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(createCustomerWithOktaOption).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith(ordersPath, {
        state: {
          isSafetyMoveSelected: true,
          isBluebarkMoveSelected: false,
        },
      });
    });
  }, 20000);

  it('disables and populates DODID and EMPLID inputs when safety move is selected', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    searchLocationByZipCityState.mockImplementation(mockSearchPickupLocation);

    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const safetyMove = await screen.findByTestId('is-safety-move-no');
    expect(safetyMove).toBeChecked();

    // check the safety move box
    await userEvent.type(getByTestId('is-safety-move-yes'), safetyPayload.is_safety_move);
    await userEvent.type(getByTestId('is-bluebark-no'), safetyPayload.is_bluebark);
    expect(await screen.findByTestId('safetyMoveHint')).toBeInTheDocument();

    await user.selectOptions(getByLabelText('Branch of service'), ['COAST_GUARD']);

    // the input boxes should now be disabled
    expect(await screen.findByTestId('edipiInput')).toBeDisabled();
    expect(await screen.findByTestId('emplidInput')).toBeDisabled();

    // should be able to submit the form
    await user.type(getByLabelText('First name'), safetyPayload.first_name);
    await user.type(getByLabelText('Last name'), safetyPayload.last_name);

    await user.type(getByLabelText('Best contact phone'), safetyPayload.telephone);
    await user.type(getByLabelText('Personal email'), safetyPayload.personal_email);

    await userEvent.type(
      getByTestId('residential_address.streetAddress1'),
      safetyPayload.residential_address.streetAddress1,
    );

    const locationBox = screen.getAllByRole('combobox');

    await act(async () => {
      await userEvent.type(locationBox[1], 'BEVERLY HILLS');
      const selectedResidentialLocation = await screen.findByText(/90210/);
      await userEvent.click(selectedResidentialLocation);
    });

    await userEvent.type(
      getByTestId('backup_mailing_address.streetAddress1'),
      safetyPayload.backup_mailing_address.streetAddress1,
    );

    await act(async () => {
      await userEvent.type(locationBox[2], 'DRYDEN');
      const selectedBackupLocation = await screen.findByText(/04225/);
      await userEvent.click(selectedBackupLocation);
    });

    await userEvent.type(getByLabelText('Name'), safetyPayload.backup_contact.name);
    await userEvent.type(getByRole('textbox', { name: 'Email' }), safetyPayload.backup_contact.email);
    await userEvent.type(getByRole('textbox', { name: 'Phone' }), safetyPayload.backup_contact.telephone);

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });
    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(createCustomerWithOktaOption).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith(ordersPath, {
        state: {
          isSafetyMoveSelected: true,
          isBluebarkMoveSelected: false,
        },
      });
    });
  }, 10000);

  it('disables okta and non cac user inputs when bluebark move is selected', async () => {
    createCustomerWithOktaOption.mockImplementation(() => Promise.resolve(fakeResponse));
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    searchLocationByZipCityState.mockImplementation(mockSearchPickupLocation);

    const { getByLabelText, getByTestId, getByRole } = render(
      <MockProviders initialState={serviceCounselorState}>
        <CreateCustomerForm {...testProps} />
      </MockProviders>,
    );

    const user = userEvent.setup();

    const saveBtn = await screen.findByRole('button', { name: 'Save' });
    expect(saveBtn).toBeInTheDocument();

    // check the safety move box
    await userEvent.type(getByTestId('is-safety-move-no'), bluebarkPayload.is_safety_move);
    await userEvent.type(getByTestId('is-bluebark-yes'), bluebarkPayload.is_bluebark);

    await userEvent.selectOptions(getByLabelText('Branch of service'), ['ARMY']);
    await user.type(getByTestId('edipiInput'), safetyPayload.edipi);

    await user.type(getByLabelText('First name'), safetyPayload.first_name);
    await user.type(getByLabelText('Last name'), safetyPayload.last_name);

    await user.type(getByLabelText('Best contact phone'), safetyPayload.telephone);
    await user.type(getByLabelText('Personal email'), safetyPayload.personal_email);

    await userEvent.type(
      getByTestId('residential_address.streetAddress1'),
      safetyPayload.residential_address.streetAddress1,
    );

    const locationBox = screen.getAllByRole('combobox');

    await act(async () => {
      await userEvent.type(locationBox[1], 'BEVERLY HILLS');
      const selectedResidentialLocation = await screen.findByText(/90210/);
      await userEvent.click(selectedResidentialLocation);
    });

    await userEvent.type(
      getByTestId('backup_mailing_address.streetAddress1'),
      safetyPayload.backup_mailing_address.streetAddress1,
    );

    await act(async () => {
      await userEvent.type(locationBox[2], 'DRYDEN');
      const selectedBackupLocation = await screen.findByText(/04225/);
      await userEvent.click(selectedBackupLocation);
    });

    await userEvent.type(getByLabelText('Name'), safetyPayload.backup_contact.name);
    await userEvent.type(getByRole('textbox', { name: 'Email' }), safetyPayload.backup_contact.email);
    await userEvent.type(getByRole('textbox', { name: 'Phone' }), safetyPayload.backup_contact.telephone);

    expect(saveBtn).toBeInTheDocument();

    await waitFor(() => {
      expect(saveBtn).toBeEnabled();
    });
    await userEvent.click(saveBtn);

    await waitFor(() => {
      expect(createCustomerWithOktaOption).toHaveBeenCalled();
      expect(mockNavigate).toHaveBeenCalledWith(ordersPath, {
        state: {
          isSafetyMoveSelected: false,
          isBluebarkMoveSelected: true,
        },
      });
    });
  }, 10000);
});
