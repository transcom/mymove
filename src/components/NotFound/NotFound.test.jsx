import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import NotFound from './NotFound';

const mockNavigate = jest.fn();
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate,
}));

afterEach(() => {
  jest.clearAllMocks();
});

describe('NotFound component', () => {
  const handleOnClick = jest.fn();
  it('Renders page not found', () => {
    render(<NotFound handleOnClick={handleOnClick} />);
    expect(screen.getByText('Error - 404')).toBeInTheDocument();
  });

  it('calls the handleOnClick when Go Back is clicked', async () => {
    render(<NotFound handleOnClick={handleOnClick} />);

    const goBackButton = screen.getByRole('button', { name: 'back home.' });
    await userEvent.click(goBackButton);

    expect(handleOnClick).toHaveBeenCalledTimes(1);
    expect(mockNavigate).toHaveBeenCalledTimes(0);
  });

  it('calls navigate(-1) when Go Back is clicked and no handleOnClick provided', async () => {
    render(<NotFound />);

    const goBackButton = screen.getByRole('button', { name: 'back home.' });
    await userEvent.click(goBackButton);

    expect(handleOnClick).toHaveBeenCalledTimes(0);
    expect(mockNavigate).toHaveBeenCalledTimes(1);
  });
});
