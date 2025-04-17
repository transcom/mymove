import React from 'react';
import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useLocation } from 'react-router-dom';

import { EditContactInfo } from './EditContactInfo';

import { patchBackupContact, patchServiceMember, getAllMoves } from 'services/internalApi';
import { customerRoutes } from 'constants/routes';
import { MockProviders } from 'testUtils';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
  useLocation: jest.fn(),
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchBackupContact: jest.fn(),
  patchServiceMember: jest.fn(),
  getAllMoves: jest.fn().mockImplementation(() => Promise.resolve()),
}));

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

const serviceMemberMoves = {
  currentMove: [
    {
      lockExpiresAt: '2099-04-07T17:21:30.450Z',
      createdAt: '2024-02-23T19:30:11.374Z',
      eTag: 'MjAyNC0wMi0yM1QxOTozMDoxMS4zNzQxN1o=',
      id: 'testMoveId',
      moveCode: '44649B',
      orders: {
        authorizedWeight: 11000,
        created_at: '2024-02-23T19:30:11.369Z',
        entitlement: {
          proGear: 2000,
          proGearSpouse: 500,
        },
        grade: 'E_7',
        has_dependents: false,
        id: 'testOrders1',
        issue_date: '2024-02-29',
        new_duty_location: {
          address: {
            city: 'Fort Irwin',
            country: 'United States',
            id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
            postalCode: '92310',
            state: 'CA',
            streetAddress1: 'n/a',
          },
          address_id: '77dca457-d0d6-4718-9ca4-a630b4614cf8',
          affiliation: 'ARMY',
          created_at: '2024-02-22T21:34:21.449Z',
          id: '12421bcb-2ded-4165-b0ac-05f76301082a',
          name: 'Fort Irwin, CA 92310',
          transportation_office: {
            address: {
              city: 'Fort Irwin',
              country: 'United States',
              id: '65a97b21-cf6a-47c1-a4b6-e3f885dacba5',
              postalCode: '92310',
              state: 'CA',
              streetAddress1: 'Langford Lake Rd',
              streetAddress2: 'Bldg 105',
            },
            created_at: '2018-05-28T14:27:37.312Z',
            gbloc: 'LKNQ',
            id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
            name: 'PPPO Fort Irwin - USA',
            phone_lines: [],
            updated_at: '2018-05-28T14:27:37.312Z',
          },
          transportation_office_id: 'd00e3ee8-baba-4991-8f3b-86c2e370d1be',
          updated_at: '2024-02-22T21:34:21.449Z',
        },
        originDutyLocationGbloc: 'BGAC',
        origin_duty_location: {
          address: {
            city: 'Fort Gregg-Adams',
            country: 'United States',
            id: '12270b68-01cf-4416-8b19-125d11bc8340',
            postalCode: '23801',
            state: 'VA',
            streetAddress1: 'n/a',
          },
          address_id: '12270b68-01cf-4416-8b19-125d11bc8340',
          affiliation: 'ARMY',
          created_at: '2024-02-22T21:34:26.430Z',
          id: '9cf15b8d-985b-4ca3-9f27-4ba32a263908',
          name: 'Fort Gregg-Adams, VA 23801',
          transportation_office: {
            address: {
              city: 'Fort Gregg-Adams',
              country: 'United States',
              id: '10dc88f5-d76a-427f-89a0-bf85587b0570',
              postalCode: '23801',
              state: 'VA',
              streetAddress1: '1401 B Ave',
              streetAddress2: 'Bldg 3400, Room 119',
            },
            created_at: '2018-05-28T14:27:42.125Z',
            gbloc: 'BGAC',
            id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
            name: 'PPPO Fort Gregg-Adams - USA',
            phone_lines: [],
            updated_at: '2018-05-28T14:27:42.125Z',
          },
          transportation_office_id: '4cc26e01-f0ea-4048-8081-1d179426a6d9',
          updated_at: '2024-02-22T21:34:26.430Z',
        },
        report_by_date: '2024-02-29',
        service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
        spouse_has_pro_gear: false,
        status: 'DRAFT',
        updated_at: '2024-02-23T19:30:11.369Z',
        uploaded_orders: {
          id: 'bd35c4c2-41c6-44a1-bf54-9098c68d87cc',
          service_member_id: '81aeac60-80f3-44d1-9b74-ba6d405ee2da',
          uploads: [
            {
              bytes: 92797,
              contentType: 'image/png',
              createdAt: '2024-02-26T18:43:58.515Z',
              filename: 'Screenshot 2024-02-08 at 12.57.43 PM.png',
              id: '786237dc-c240-449d-8859-3f37583b3406',
              status: 'PROCESSING',
              updatedAt: '2024-02-26T18:43:58.515Z',
              url: '/storage/user/5fe4d948-aa1c-4823-8967-b1fb40cf6679/uploads/786237dc-c240-449d-8859-3f37583b3406?contentType=image%2Fpng',
            },
          ],
        },
      },
      status: 'DRAFT',
      submittedAt: '0001-01-01T00:00:00.000Z',
      updatedAt: '0001-01-01T00:00:00.000Z',
    },
  ],
  previousMoves: [],
};

describe('EditContactInfo page', () => {
  const testProps = {
    currentBackupContacts: [
      {
        id: 'backupContactID',
        name: 'Barbara St. Juste',
        email: 'bsj@example.com',
        telephone: '915-555-1234',
        permission: 'NONE',
      },
    ],
    serviceMember: {
      id: 'testServiceMemberID',
      telephone: '915-555-2945',
      secondary_telephone: '',
      personal_email: 'test@example.com',
      email_is_preferred: true,
      phone_is_preferred: false,
      residential_address: {
        streetAddress1: '148 S East St',
        streetAddress2: '',
        streetAddress3: '',
        city: 'Fake City',
        state: 'TX',
        postalCode: '79936',
        county: 'EL PASO',
        usPostRegionCitiesID: '',
      },
      backup_mailing_address: {
        streetAddress1: '10642 N Second Ave',
        streetAddress2: '',
        streetAddress3: '',
        city: 'Fake City',
        state: 'TX',
        postalCode: '79936',
        county: 'EL PASO',
        usPostRegionCitiesID: '',
      },
    },
    setFlashMessage: jest.fn(),
    updateBackupContact: jest.fn(),
    updateServiceMember: jest.fn(),
  };

  it('renders the EditContactInfo form', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const h1 = await screen.findByRole('heading', { name: 'Edit contact info', level: 1 });
    expect(h1).toBeInTheDocument();

    const contactHeader = screen.getByRole('heading', { name: 'Your contact info', level: 2 });
    expect(contactHeader).toBeInTheDocument();

    const addressHeader = screen.getByRole('heading', { name: 'Current address', level: 2 });
    expect(addressHeader).toBeInTheDocument();

    const backupAddressHeader = screen.getByRole('heading', { name: 'Backup address', level: 2 });
    expect(backupAddressHeader).toBeInTheDocument();

    const backupContactHeader = screen.getByRole('heading', { name: 'Backup contact', level: 2 });
    expect(backupContactHeader).toBeInTheDocument();
  });

  it('goes back to the profile page when the cancel button is clicked', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH, { state: { moveId: 'testMoveId' } });
  });

  it('saves backup contact info when it is updated and the save button is clicked', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    const newName = 'Rosalie Wexler';

    const expectedPayload = { ...testProps.currentBackupContacts[0], name: newName };

    const patchResponse = {
      ...expectedPayload,
      serviceMemberId: testProps.serviceMember.id,
      created_at: '2021-02-08T16:48:04.117Z',
      updated_at: '2021-02-11T16:48:04.117Z',
    };

    patchBackupContact.mockImplementation(() => Promise.resolve(patchResponse));
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const backupNameInput = await screen.findByLabelText(/Name/);

    await userEvent.clear(backupNameInput);

    await userEvent.type(backupNameInput, newName);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchBackupContact).toHaveBeenCalledWith(expectedPayload);
    });

    expect(testProps.updateBackupContact).toHaveBeenCalledWith(patchResponse);
  });

  it('shows an error if the patchBackupContact API returns an error', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    patchBackupContact.mockImplementation(() =>
      // Disable this rule because makeSwaggerRequest does not throw an error if the API call fails
      // eslint-disable-next-line prefer-promise-reject-errors
      Promise.reject({
        message: 'A server error occurred saving the backup contact',
        response: {
          body: {
            detail: 'A server error occurred saving the backup contact',
          },
        },
      }),
    );

    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const backupNameInput = await screen.findByLabelText(/Name/);

    await userEvent.clear(backupNameInput);

    await userEvent.type(backupNameInput, 'Rosalie Wexler');

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchBackupContact).toHaveBeenCalled();
    });

    expect(await screen.findByText('A server error occurred saving the backup contact')).toBeInTheDocument();
    expect(testProps.updateBackupContact).not.toHaveBeenCalled();
    expect(patchServiceMember).not.toHaveBeenCalled();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.setFlashMessage).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it('does not save backup contact info if it is not updated and the save button is clicked', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchBackupContact).not.toHaveBeenCalled();
    });

    expect(testProps.updateBackupContact).not.toHaveBeenCalled();
  });

  it('saves service member info when the save button is clicked', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    const secondaryPhone = '915-555-9753';

    const expectedPayload = { ...testProps.serviceMember, secondary_telephone: secondaryPhone };

    const patchResponse = {
      ...expectedPayload,
      created_at: '2021-02-08T16:48:04.117Z',
      updated_at: '2021-02-11T16:48:04.117Z',
    };

    patchServiceMember.mockImplementation(() => Promise.resolve(patchResponse));

    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const secondaryPhoneInput = await screen.findByLabelText(/Alt. phone/);

    await userEvent.clear(secondaryPhoneInput);

    await userEvent.type(secondaryPhoneInput, secondaryPhone);

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalledWith(expectedPayload);
    });

    expect(testProps.updateServiceMember).toHaveBeenCalledWith(patchResponse);
  });

  it('sets a flash message when the save button is clicked', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(testProps.setFlashMessage).toHaveBeenCalledWith(
        'EDIT_CONTACT_INFO_SUCCESS',
        'success',
        "You've updated your information.",
      );
    });
  });

  it('routes to the profile page when the save button is clicked', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH, { state: { moveId: 'testMoveId' } });
    });
  });

  it('converts the "Submit" button to the "Return to Home" button when the move has been locked by an office user', async () => {
    getAllMoves.mockResolvedValue(serviceMemberMoves);
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    await act(() => {
      render(
        <MockProviders>
          <EditContactInfo {...testProps} />
        </MockProviders>,
      );
    });
    expect(screen.getByRole('button', { name: 'Return home' })).toBeInTheDocument();
  });

  it('routes to the profile page when the cancel button is clicked', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
    patchServiceMember.mockImplementation(() => Promise.resolve());

    render(
      <MockProviders>
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const cancelButton = screen.getByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH, { state: { moveId: 'testMoveId' } });
    });
  });

  it('shows an error if the patchServiceMember API returns an error', async () => {
    useLocation.mockReturnValue({ state: { moveId: 'testMoveId' } });
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
        <EditContactInfo {...testProps} />
      </MockProviders>,
    );

    const saveButton = screen.getByRole('button', { name: 'Save' });

    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(patchServiceMember).toHaveBeenCalled();
    });

    expect(await screen.findByText('A server error occurred saving the service member')).toBeInTheDocument();
    expect(testProps.updateServiceMember).not.toHaveBeenCalled();
    expect(testProps.setFlashMessage).not.toHaveBeenCalled();
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  afterEach(jest.resetAllMocks);
});
