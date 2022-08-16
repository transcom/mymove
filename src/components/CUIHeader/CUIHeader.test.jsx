import React from 'react';
import { render, screen } from '@testing-library/react';

import CUIHeader from './CUIHeader';

describe('CUIHeader', () => {
  it('Displays Controlled Unclassified Information', () => {
    render(<CUIHeader />);
    expect(screen.getByText('Controlled Unclassified Information')).toBeInTheDocument();
  });
});
