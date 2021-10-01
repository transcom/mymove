/* eslint-disable import/first, import/order */
import React from 'react';
import { render, screen } from '@testing-library/react';
import fakeDataProvider from 'ra-data-fakerest';

import AdminWrapper from './index';

jest.mock('utils/api', () => ({
  ...jest.requireActual('utils/api'),
  GetLoggedInUser: jest.fn(),
}));

jest.mock('js-cookie', () => ({
  ...jest.requireActual('js-cookie'),
  get: jest.fn(() => 'mock-csrf-token'),
}));

import { GetLoggedInUser } from 'utils/api';

const dataProvider = fakeDataProvider({
  office_users: [],
});

describe('pages/Admin/AdminWrapper', () => {
  it('renders correctly if logged in', async () => {
    GetLoggedInUser.mockImplementationOnce(async () => {
      return null;
    });

    render(<AdminWrapper basename="/" dataProvider={dataProvider} />);

    // Displays logged in default (first) route
    const officeUsers = await screen.findByRole('menuitem', { name: 'Office Users' });
    expect(officeUsers.getAttribute('aria-current')).toEqual('page');
  });

  it('renders logged out if not logged in', async () => {
    GetLoggedInUser.mockImplementationOnce(async () => {
      throw new Error('hey it broke');
    });

    render(<AdminWrapper basename="/" dataProvider={dataProvider} />);

    expect(
      await screen.findByText('This is a new system from USTRANSCOM to support the relocation of families during PCS.'),
    ).toBeInTheDocument();
  });
});
