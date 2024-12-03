import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import { MockProviders } from 'testUtils';
import ConnectedResidentialAddress, { ResidentialAddress } from 'pages/MyMove/Profile/ResidentialAddress';
import { configureStore } from 'shared/store';
import { customerRoutes } from 'constants/routes';
import { patchServiceMember } from 'services/internalApi';
import { ValidateZipRateData } from 'shared/api';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

jest.mock('shared/api', () => ({
  ...jest.requireActual('shared/api'),
  ValidateZipRateData: jest.fn(),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

describe('ResidentialAddress page', () => {
  const fakeAddress = {
    streetAddress1: '235 Prospect Valley Road SE',
    streetAddress2: '#125',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79912',
    county: 'El Paso',
    usPostRegionCitiesID: '',
  };

  const blankAddress = Object.fromEntries(Object.keys(fakeAddress).map((k) => [k, '']));
  // TODO: We may want to change residential_address to residentialAddress
  const generateTestProps = (address) => ({
    updateServiceMember: jest.fn(),
    push: jest.fn(),
    serviceMember: {
      id: 'testServiceMemberId',
      residential_address: address,
    },
  });

  it('renders the ResidentialAddressForm', async () => {
    const testProps = generateTestProps(blankAddress);
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <ResidentialAddress {...testProps} />
      </Provider>,
    );

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Current address', level: 1 })).toBeInTheDocument();
    });
  });

  it('back button goes to the contact info step', async () => {
    const testProps = generateTestProps(blankAddress);
    const mockStore = configureStore({});

    render(
      <Provider store={mockStore.store}>
        <ResidentialAddress {...testProps} />
      </Provider>,
    );

    const backButton = await screen.findByRole('button', { name: 'Back' });
    expect(backButton).toBeInTheDocument();
    await userEvent.click(backButton);

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.CONTACT_INFO_PATH);
  });

  it('next button submits the form and goes to the Backup address step', async () => {
    const testProps = generateTestProps(fakeAddress);
    const mockStore = configureStore({});

    const expectedServiceMemberPayload = { ...testProps.serviceMember, residential_address: fakeAddress };

    ValidateZipRateData.mockImplementation(() => ({
      valid: true,
    }));
    patchServiceMember.mockImplementation(() => Promise.resolve(expectedServiceMemberPayload));

    render(
      <Provider store={mockStore.store}>
        <ResidentialAddress {...testProps} />
      </Provider>,
    );

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.BACKUP_ADDRESS_PATH);
  });

  it('shows an error if the ValidateZipRateData API returns an error', async () => {
    const testProps = generateTestProps(fakeAddress);
    const mockStore = configureStore({});

    ValidateZipRateData.mockImplementation(() => ({
      valid: false,
    }));
    patchServiceMember.mockImplementation(() => Promise.resolve(testProps.serviceMember));

    render(
      <Provider store={mockStore.store}>
        <ResidentialAddress {...testProps} />
      </Provider>,
    );

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    expect(patchServiceMember).not.toHaveBeenCalled();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it('shows an error if the patchServiceMember API returns an error', async () => {
    const testProps = generateTestProps(fakeAddress);

    ValidateZipRateData.mockImplementation(() => ({
      valid: true,
    }));
    patchServiceMember.mockImplementation(() =>
      // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
      // eslint-disable-next-line prefer-promise-reject-errors
      Promise.reject({
        message: 'A server error occurred saving the service member',
        response: {
          body: {
            detail: 'A server error occurred saving the service member',
          },
        },
      }),
    );

    render(
      <MockProviders>
        <ResidentialAddress {...testProps} />
      </MockProviders>,
    );

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(await screen.findByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState ResidentialAddress', () => {
  const props = {
    updateServiceMember: jest.fn(),
  };

  it('dispatches a redirect if the current state is earlier than the "CONTACT_INFO_PATH" state', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedResidentialAddress {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Current address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.CONTACT_INFO_PATH);
    });
  });

  it('does not redirect if the current state equals the "CONTACT_INFO_COMPLETE" state', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedResidentialAddress {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Current address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockNavigate).not.toHaveBeenCalled();
    });
  });

  it('does not redirect if the current state is after the "CONTACT_INFO_COMPLETE" state and profile is not complete', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
            },
            residential_address: {
              street: '123 Main St',
            },
            backup_mailing_address: {
              street: '456 Main St',
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedResidentialAddress {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Current address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockNavigate).not.toHaveBeenCalled();
    });
  });

  it('does redirect if the profile is complete', async () => {
    const mockState = {
      entities: {
        user: {
          testUserId: {
            id: 'testUserId',
            email: 'testuser@example.com',
            service_member: 'testServiceMemberId',
          },
        },
        serviceMembers: {
          testServiceMemberId: {
            id: 'testServiceMemberId',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
            current_location: {
              id: 'testDutyLocationId',
            },
            residential_address: {
              street: '123 Main St',
            },
            backup_mailing_address: {
              street: '456 Main St',
            },
            backup_contacts: [
              {
                id: 'testBackupContact',
              },
            ],
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedResidentialAddress {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Current address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
});
