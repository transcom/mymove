import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { generatePath } from 'react-router-dom';

import { generalRoutes } from 'constants/routes';
import ProGear from 'pages/MyMove/PPM/Closeout/ProGear/ProGear';

const mockPush = jest.fn();
const mockReplace = jest.fn();
const homePath = generatePath(generalRoutes.HOME_PATH);
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockPush,
    replace: mockReplace,
  }),
  useLocation: () => ({}),
}));

describe('Pro-gear page', () => {
  it('displays the page', () => {
    render(<ProGear />);
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Pro-gear');
  });
  it('displays reminder to include pro-gear weight in total', () => {
    render(<ProGear />);
    expect(screen.getByText(/This pro-gear should be included in your total weight moved./)).toBeInTheDocument();
  });
  it('routes back to home when finish later is clicked', async () => {
    render(<ProGear />);

    expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));
    expect(mockPush).toHaveBeenCalledWith(homePath);
  });
});
