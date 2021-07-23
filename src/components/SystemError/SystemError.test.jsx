import React from 'react';
import { render, screen } from '@testing-library/react';

import SystemError from './index';

describe('SystemError component', () => {
  it('renders without crashing', () => {
    render(<SystemError />);
    expect(screen.getByText(/contact the Technical Help Desk and give them this code:/)).toBeInTheDocument();
  });
});
