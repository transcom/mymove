import React from 'react';
import { render, screen } from '@testing-library/react';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import PPMDocumentsStatus from 'constants/ppms';

const mockWeightTickets = [
  createCompleteWeightTicket(
    {},
    {
      status: PPMDocumentsStatus.APPROVED,
    },
  ),
  createCompleteWeightTicket(
    {},
    {
      status: PPMDocumentsStatus.REJECTED,
      reason: 'Rejection reason',
    },
  ),
];

describe('ReviewDocumentsSidePanel', () => {
  it('renders the component', async () => {
    render(<ReviewDocumentsSidePanel />);
    const h3 = await screen.getByRole('heading', { name: 'Send to customer?', level: 3 });
    expect(h3).toBeInTheDocument();
  });

  it('shows the appropriate statuses once weight tickets are reviewed', async () => {
    render(<ReviewDocumentsSidePanel weightTickets={mockWeightTickets} />);
    const listItems = await screen.getAllByRole('listitem');
    expect(listItems).toHaveLength(2);
    expect(listItems[0]).toHaveTextContent(/Trip 1/);
    expect(listItems[0]).toHaveTextContent(/Accept/);
    expect(listItems[1]).toHaveTextContent(/Trip 2/);
    expect(listItems[1]).toHaveTextContent(/Reject/);
    expect(listItems[1]).toHaveTextContent(/Rejection reason/);
  });
});
