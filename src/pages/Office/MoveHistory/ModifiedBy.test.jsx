import React from 'react';
import { render, screen } from '@testing-library/react';

import ModifiedBy from './ModifiedBy';

describe('Move History ModifiedBy view', () => {
  it('formats name and contact info correctly', () => {
    render(<ModifiedBy firstName="Leo" lastName="Spaceman" email="leos@example.com" phone="555-555-5555" />);
    expect(screen.getByText('Spaceman, Leo')).toBeInTheDocument();
    expect(screen.getByText('leos@example.com | 555-555-5555')).toBeInTheDocument();
  });
});
