/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { screen } from '@testing-library/react';

import { Dashboard } from './index';

import { renderWithRouterProp } from 'testUtils';

describe('Dashboard component', () => {
  describe('with default props', () => {
    renderWithRouterProp(<Dashboard />);

    it('should successfully render multiple moves dashboard', async () => {
      expect(await screen.findByText('MULTIPLE MOVES DASHBOARD')).toBeInTheDocument();
    });
  });
});
