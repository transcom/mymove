import React from 'react';
import * as reactRedux from 'react-redux';
import { push } from 'connected-react-router';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { MockProviders } from 'testUtils';
import ConnectedResidentialAddress, { ResidentialAddress } from 'pages/MyMove/Profile/ResidentialAddress';
import { customerRoutes } from 'constants/routes';
import { patchServiceMember } from 'services/internalApi';
import { ValidateZipRateData } from 'shared/api';

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchServiceMember: jest.fn(),
}));

jest.mock('shared/api', () => ({
  ...jest.requireActual('shared/api'),
  ValidateZipRateData: jest.fn(),
}));

describe('ResidentialAddress page', () => {
  const fakeAddress = {
    streetAddress1: '235 Prospect Valley Road SE',
    streetAddress2: '#125',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79912',
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

    render(<ResidentialAddress {...testProps} />);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Current mailing address', level: 1 })).toBeInTheDocument();
    });
  });

  it('validates zip code using api endpoint', async () => {
    const testProps = generateTestProps(blankAddress);

    ValidateZipRateData.mockImplementation(() => ({
      valid: true,
    }));

    render(<ResidentialAddress {...testProps} />);

    const postalCodeInput = await screen.findByLabelText('ZIP');

    const postalCode = '99999';

    userEvent.type(postalCodeInput, postalCode);
    userEvent.tab();

    await waitFor(() => {
      expect(ValidateZipRateData).toHaveBeenCalledWith(postalCode, 'origin');
    });
  });

  it('back button goes to the Current duty location step', async () => {
    const testProps = generateTestProps(blankAddress);

    render(<ResidentialAddress {...testProps} />);

    const backButton = await screen.findByRole('button', { name: 'Back' });
    expect(backButton).toBeInTheDocument();
    userEvent.click(backButton);

    expect(testProps.push).toHaveBeenCalledWith(customerRoutes.CURRENT_DUTY_LOCATION_PATH);
  });

  it('next button submits the form and goes to the Backup address step', async () => {
    const testProps = generateTestProps(blankAddress);

    const expectedServiceMemberPayload = { ...testProps.serviceMember, residential_address: fakeAddress };

    ValidateZipRateData.mockImplementation(() => ({
      valid: true,
    }));
    patchServiceMember.mockImplementation(() => Promise.resolve(expectedServiceMemberPayload));

    render(<ResidentialAddress {...testProps} />);

    userEvent.type(screen.getByLabelText('Address 1'), fakeAddress.streetAddress1);
    userEvent.type(screen.getByLabelText(/Address 2/), fakeAddress.streetAddress2);
    userEvent.type(screen.getByLabelText('City'), fakeAddress.city);
    userEvent.selectOptions(screen.getByLabelText('State'), [fakeAddress.state]);
    userEvent.type(screen.getByLabelText('ZIP'), fakeAddress.postalCode);

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(expectedServiceMemberPayload);
    expect(testProps.push).toHaveBeenCalledWith(customerRoutes.BACKUP_ADDRESS_PATH);
  });

  it('shows an error if the ValidateZipRateData API returns an error', async () => {
    const testProps = generateTestProps(fakeAddress);

    ValidateZipRateData.mockImplementation(() => ({
      valid: false,
    }));
    patchServiceMember.mockImplementation(() => Promise.resolve(testProps.serviceMember));

    render(<ResidentialAddress {...testProps} />);

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    const alert = await screen.findByRole('alert');

    expect(alert).toHaveTextContent(
      'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.',
    );
    expect(patchServiceMember).not.toHaveBeenCalled();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.push).not.toHaveBeenCalled();
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

    render(<ResidentialAddress {...testProps} />);

    const submitButton = screen.getByRole('button', { name: 'Next' });
    expect(submitButton).toBeInTheDocument();
    userEvent.click(submitButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(screen.getByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.push).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});

describe('requireCustomerState ResidentialAddress', () => {
  const useDispatchMock = jest.spyOn(reactRedux, 'useDispatch');
  const mockDispatch = jest.fn();

  beforeEach(() => {
    useDispatchMock.mockClear();
    mockDispatch.mockClear();
    useDispatchMock.mockReturnValue(mockDispatch);
  });

  const props = {
    updateServiceMember: jest.fn(),
    push: jest.fn(),
  };

  it('dispatches a redirect if the current state is earlier than the "DUTY STATION COMPLETE" state', async () => {
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
            rank: 'test rank',
            edipi: '1234567890',
            affiliation: 'ARMY',
            first_name: 'Tester',
            last_name: 'Testperson',
            telephone: '1234567890',
            personal_email: 'test@example.com',
            email_is_preferred: true,
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedResidentialAddress {...props} />
      </MockProviders>,
    );

    const h1 = screen.getByRole('heading', { name: 'Current mailing address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).toHaveBeenCalledWith(push(customerRoutes.CURRENT_DUTY_LOCATION_PATH));
    });
  });

  it('does not redirect if the current state equals the "DUTY STATION COMPLETE" state', async () => {
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
            rank: 'test rank',
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

    const h1 = screen.getByRole('heading', { name: 'Current mailing address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).not.toHaveBeenCalled();
    });
  });

  it('does not redirect if the current state is after the "DUTY STATION COMPLETE" state and profile is not complete', async () => {
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
            rank: 'test rank',
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

    const h1 = screen.getByRole('heading', { name: 'Current mailing address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).not.toHaveBeenCalled();
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
            rank: 'test rank',
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

    const h1 = screen.getByRole('heading', { name: 'Current mailing address', level: 1 });
    expect(h1).toBeInTheDocument();

    await waitFor(async () => {
      expect(mockDispatch).toHaveBeenCalledWith(push('/'));
    });
  });
});
