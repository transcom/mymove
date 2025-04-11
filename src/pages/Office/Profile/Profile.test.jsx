/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { render, screen } from '@testing-library/react';

import ConnectedProfile from './Profile';

import { MockProviders } from 'testUtils';

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useLocation: jest.fn(),
}));

describe('Profile component', () => {
  const testProps = {};

  it('renders the Profile Page', async () => {
    const mockState = {
      entities: {
        user: {
          userId123: {
            id: 'userId123',
            office_user: {
              id: '123',
              first_name: 'John',
              middle_name: 'M',
              last_name: 'Doe',
              telephone: '804-456-7890',
              email: 'john.doe@example.com',
            },
          },
        },
      },
    };

    render(
      <MockProviders initialState={mockState}>
        <ConnectedProfile {...testProps} />
      </MockProviders>,
    );
    expect(await screen.findByRole('heading', { name: 'Profile', level: 1 })).toBeInTheDocument();
    expect(await screen.findByRole('heading', { name: 'Contact info', level: 2 })).toBeInTheDocument();
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Return to Dashboard')).toBeInTheDocument();
  });
});
