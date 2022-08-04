import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import LoggedOutHeader from './LoggedOutHeader';

describe('LoggedOutHeader', () => {
  it('renders the logged out header', () => {
    render(<LoggedOutHeader />);

    const homeLink = screen.getByTitle('Home');
    expect(homeLink).toBeInstanceOf(HTMLAnchorElement);

    const signInButton = screen.getByRole('button', { name: 'Sign In' });
    expect(signInButton).toBeInstanceOf(HTMLButtonElement);
  });

  it('shows the EULA modal when logging in', async () => {
    render(<LoggedOutHeader />);

    const signInButton = screen.getByRole('button', { name: 'Sign In' });
    expect(signInButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signInButton);

    const eulaModal = screen.getByText('Standard mandatory DoD Notice and consent Banner');
    expect(eulaModal).toBeInstanceOf(HTMLHeadingElement);
  });

  it('closes the EULA modal when cancel is clicked', async () => {
    render(<LoggedOutHeader />);

    const signInButton = screen.getByRole('button', { name: 'Sign In' });
    expect(signInButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(signInButton);

    const eulaModal = screen.getByText('Standard mandatory DoD Notice and consent Banner');
    expect(eulaModal).toBeInstanceOf(HTMLHeadingElement);

    const cancelButton = screen.getByRole('button', { name: 'Cancel' });
    expect(cancelButton).toBeInstanceOf(HTMLButtonElement);

    await userEvent.click(cancelButton);

    expect(screen.queryByText('Standard mandatory DoD Notice and consent Banner')).not.toBeInTheDocument();
  });
});
