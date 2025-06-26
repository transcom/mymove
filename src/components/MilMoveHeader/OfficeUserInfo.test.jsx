import React from 'react';
import { getDefaultNormalizer, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import OfficeUserInfo from './OfficeUserInfo';

import MilMoveHeader from './index';

import { MockProviders } from 'testUtils';
import { formatOfficeProfileFirstAndLast } from 'utils/formatters';

describe('OfficeUserInfo', () => {
  it('renders the OfficeUser Info', async () => {
    const officeUser = {
      firstName: 'Sam',
      lastName: 'Samson',
    };

    const textContentToCheck = formatOfficeProfileFirstAndLast(officeUser);

    const signOutHandler = jest.fn();
    render(
      <MockProviders>
        <MilMoveHeader>
          <OfficeUserInfo
            lastName={officeUser.lastName}
            firstName={officeUser.firstName}
            handleLogout={signOutHandler}
          />
        </MilMoveHeader>
      </MockProviders>,
    );

    const profileLink = await screen.getByRole('link', {
      name: 'profile-link',
      normalizer: getDefaultNormalizer(),
    });

    expect(profileLink).toHaveTextContent(textContentToCheck);
    const signOutButton = screen.getByRole('button', { name: 'Sign out' });
    expect(signOutButton).toBeInstanceOf(HTMLButtonElement);
    await userEvent.click(signOutButton);
    await waitFor(() => {
      expect(signOutHandler).toHaveBeenCalled();
    });
  });
});
