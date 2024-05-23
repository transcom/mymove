import React from 'react';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ConusOrNot from './ConusOrNot';

import { renderWithProviders } from 'testUtils';
import { customerRoutes } from 'constants/routes';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

describe('ConusOrNot', () => {
  test('it should render all text for the component', async () => {
    renderWithProviders(<ConusOrNot />);

    expect(screen.getByText('Where are you moving?')).toBeInTheDocument();
    expect(screen.getByText('CONUS')).toBeInTheDocument();
    expect(screen.getByText('OCONUS')).toBeInTheDocument();
  });

  test('it selects an option and navigates the user', async () => {
    renderWithProviders(<ConusOrNot />);

    userEvent.click(screen.getByText('CONUS'));
    const nextBtn = await screen.findByRole('button', { name: 'Next' });
    await userEvent.click(nextBtn);
    expect(mockNavigate).toHaveBeenCalledWith(customerRoutes.DOD_INFO_PATH);
  });
});
