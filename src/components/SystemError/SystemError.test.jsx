import React from 'react';
import { render, screen } from '@testing-library/react';

import SystemError from './index';

describe('SystemError component', () => {
  it('renders without crashing', () => {
    render(<SystemError>Contact the Technical Help Desk and give them this code:</SystemError>);
    expect(screen.getByText(/Contact the Technical Help Desk and give them this code:/)).toBeInTheDocument();
  });
});
