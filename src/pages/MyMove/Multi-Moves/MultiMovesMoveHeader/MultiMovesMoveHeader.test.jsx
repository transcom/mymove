import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';

import MultiMovesMoveHeader from './MultiMovesMoveHeader';

describe('MultiMovesMoveHeader', () => {
  it('renders the move header with the correct title', () => {
    const title = 'Test Move';
    render(<MultiMovesMoveHeader title={title} />);

    const truckIcon = screen.getByTestId('truck-icon');
    const headerTitle = screen.getByText(title);

    expect(truckIcon).toBeInTheDocument();
    expect(headerTitle).toBeInTheDocument();
  });
});
