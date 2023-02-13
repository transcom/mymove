import React from 'react';
import { render, screen } from '@testing-library/react';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

import PPMDocumentsStatus from 'constants/ppms';

const mockWeightTickets = [
  {
    status: PPMDocumentsStatus.APPROVED,
  },
  {
    status: PPMDocumentsStatus.REJECTED,
    reason: 'Rejection reason',
  },
];

describe('ReviewDocumentsSidePanel', () => {
  it('renders the component', async () => {
    render(<ReviewDocumentsSidePanel />);
    const h3 = await screen.getByRole('heading', { name: 'Send to customer?', level: 3 });
    expect(h3).toBeInTheDocument();
  });

  it('shows the appropriate statuses once weight tickets are reviewed', async () => {
    const { getAllByRole } = render(<ReviewDocumentsSidePanel weightTickets={mockWeightTickets} />);
    const listItems = await getAllByRole('listitem');
    expect(listItems).toHaveLength(2);
    expect(listItems[0]).toHaveTextContent(/Trip 1/);
    expect(listItems[0]).toHaveTextContent(/Accept/);
    expect(listItems[1]).toHaveTextContent(/Trip 2/);
    expect(listItems[1]).toHaveTextContent(/Reject/);
    expect(listItems[1]).toHaveTextContent(/Rejection reason/);
  });
});
