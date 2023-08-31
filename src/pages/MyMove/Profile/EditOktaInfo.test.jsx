import React from 'react';
import { MemoryRouter } from 'react-router';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import { EditOktaInfo } from './EditOktaInfo';

import { customerRoutes } from 'constants/routes';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

jest.mock('services/internalApi', () => ({
  ...jest.requireActual('services/internalApi'),
  patchBackupContact: jest.fn(),
  patchServiceMember: jest.fn(),
}));

beforeEach(() => {
  jest.resetAllMocks();
});

describe('EditOktaInfo page', () => {
  const testProps = {
    oktaInfo: {
      id: 'testServiceMemberID',
      oktaUsername: 'test@okta.mil',
      oktaEmail: 'test@okta.mil',
      oktaFirstName: 'Jim',
      oktaLastName: 'Dunk',
      oktaEdipi: '1234123412',
    },
    setFlashMessage: jest.fn(),
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

  it('goes back to the profile page when the cancel button is clicked', async () => {
    render(
      <MemoryRouter>
        <EditOktaInfo {...testProps} />
      </MemoryRouter>,
    );

    const cancelButton = await screen.findByRole('button', { name: 'Cancel' });

    await userEvent.click(cancelButton);

    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.PROFILE_PATH);
  });

  afterEach(jest.resetAllMocks);
});
