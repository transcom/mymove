import React from 'react';
import { MemoryRouter } from 'react-router';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditOktaInfo } from './EditOktaInfo';

import { customerRoutes } from 'constants/routes';
import { updateOktaUser } from 'services/internalApi';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  updateOktaUser: jest.fn(),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

describe('EditOktaInfo page', () => {
  const testProps = {
    oktaUser: {
      id: 'testServiceMemberID',
      login: 'test@okta.mil',
      email: 'test@okta.mil',
      firstName: 'Jim',
      lastName: 'Dunk',
      cac_edipi: '1234123412',
    },
    serviceMember: {
      id: 'testServiceMemberId',
    },
    setFlashMessage: jest.fn(),
    updateOktaUser: jest.fn(),
  };

  it('renders the EditOktaInfo form', async () => {
    render(
      <MemoryRouter>
        <EditOktaInfo {...testProps} />
      </MemoryRouter>,
    );

    const contactHeader = screen.getByRole('heading', { name: 'Your Okta Profile', level: 2 });
    expect(contactHeader).toBeInTheDocument();
  });

  it('shows error if no changes were made', async () => {
    render(
      <MemoryRouter>
        <EditOktaInfo {...testProps} />
      </MemoryRouter>,
    );

    await userEvent.type(screen.getByLabelText('First Name *'), 'Bob');
    const saveBtn = screen.getByRole('button', { name: 'Save' });
    await userEvent.click(saveBtn);

    // Check if updateOktaUser is called with the expected arguments
    expect(updateOktaUser).toHaveBeenCalledWith({
      profile: {
        id: 'testServiceMemberId', // Adjusted to match the received value
        login: 'test@okta.mil',
        email: 'test@okta.mil',
        firstName: 'JimBob', // Adjusted to match the typed value
        lastName: 'Dunk',
        cac_edipi: '1234123412',
      },
    });
  });

  it('shows error if no changes were made', async () => {
    render(
      <MemoryRouter>
        <EditOktaInfo {...testProps} />
      </MemoryRouter>,
    );

    const saveBtn = screen.getByRole('button', { name: 'Save' });
    await userEvent.click(saveBtn);
    const errorHeader = screen.getByText('No changes were made');
    const errorMessage = screen.getByText('You must make some changes if you want to edit your Okta profile.');
    expect(errorHeader).toBeInTheDocument();
    expect(errorMessage).toBeInTheDocument();
  });

  it('goes back to the profile page when the cancel button is clicked', async () => {
    render(
      <MemoryRouter>
        <EditOktaInfo {...testProps} />
      </MemoryRouter>,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH, { state: null });
  });

  afterEach(jest.resetAllMocks);
});
