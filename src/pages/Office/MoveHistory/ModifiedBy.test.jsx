import React from 'react';
import { render, screen } from '@testing-library/react';

import ModifiedBy from './ModifiedBy';

describe('Move History ModifiedBy view', () => {
  it('formats name and contact info correctly', () => {
    render(<ModifiedBy firstName="Leo" lastName="Spaceman" email="leos@example.com" phone="555-555-5555" />);
    expect(screen.getByText('Spaceman, Leo')).toBeInTheDocument();
    expect(screen.getByText('leos@example.com | 555-555-5555')).toBeInTheDocument();
  });

  it('supports system modifications when all attributes are blank', () => {
    render(<ModifiedBy firstName="" lastName="" email="" phone="" />);
    expect(screen.getByText('MilMove')).toBeInTheDocument();
  });

  it('supports system modifications made by the Prime', () => {
    render(<ModifiedBy firstName="Prime" lastName="" email="" phone="" />);
    expect(screen.getByText('Prime')).toBeInTheDocument();
  });
});
