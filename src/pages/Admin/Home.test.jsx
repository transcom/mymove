import React from 'react';
import { render, screen } from '@testing-library/react';
import fakeDataProvider from 'ra-data-fakerest';

import Home from './Home';

const dataProvider = fakeDataProvider({
  office_users: [],
});

describe('AdminHome tests', () => {
  describe('AdminHome component', () => {
    it('renders the default (first) route, without crashing', async () => {
      render(<Home basename="/" dataProvider={dataProvider} />);

      expect(await screen.findByText('Unclassified // For official use only')).toBeInTheDocument();
      const officeUsers = await screen.findByRole('menuitem', { name: 'Office Users' });
      expect(officeUsers.getAttribute('aria-current')).toEqual('page');
    });
  });
});
