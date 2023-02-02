import React from 'react';
import { render, screen } from '@testing-library/react';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

describe('ReviewDocumentsSidePanel', () => {
  it('renders the component', async () => {
    render(<ReviewDocumentsSidePanel />);
    const h3 = await screen.getByRole('heading', { name: 'Send to customer?', level: 3 });
    expect(h3).toBeInTheDocument();
  });
});
