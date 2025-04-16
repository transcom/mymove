import React from 'react';
import { render, screen } from '@testing-library/react';
import { v4 } from 'uuid';

import ReviewDocumentsSidePanel from './ReviewDocumentsSidePanel';

import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import PPMDocumentsStatus from 'constants/ppms';
import { MockProviders } from 'testUtils';
import { createCompleteProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createCompleteMovingExpense } from 'utils/test/factories/movingExpense';
import { expenseTypes } from 'constants/ppmExpenseTypes';

const serviceMemberId = v4();
const mockWeightTickets = [
  createCompleteWeightTicket(
    { serviceMemberId },
    {
      status: PPMDocumentsStatus.APPROVED,
    },
  ),
  createCompleteWeightTicket(
    { serviceMemberId },
    {
      status: PPMDocumentsStatus.REJECTED,
      reason: 'No weight ticket',
    },
  ),
];
const yesterday = new Date();
yesterday.setDate(yesterday.getDate() - 1);
const mockPPMShipment = {
  id: '2ecb311-edbe-4fd4-96ee-bd693113f3f3',
  expectedDepartureDate: new Date(yesterday.getDate() - 1),
  actualMoveDate: yesterday,
  miles: 300,
  estimatedWeight: 3000,
  actualWeight: 3500,
  showAllFields: false,
};

describe('ReviewDocumentsSidePanel', () => {
  it('renders the component', async () => {
    render(
      <MockProviders>
        <ReviewDocumentsSidePanel ppmShipmentInfo={mockPPMShipment} />
      </MockProviders>,
    );
    const h3 = await screen.getByRole('heading', { name: 'Send to customer?', level: 3 });
    expect(h3).toBeInTheDocument();
  });

  it('shows the appropriate statuses when multiple documents have been reviewed', async () => {
    const progearWeightTickets = [
      createCompleteProGearWeightTicket({ serviceMemberId }, { status: PPMDocumentsStatus.APPROVED }),
      createCompleteProGearWeightTicket(
        { serviceMemberId },
        { status: PPMDocumentsStatus.REJECTED, reason: 'Invalid weight ticket' },
      ),
    ];

    const movingExpenses = [
      createCompleteMovingExpense({ serviceMemberId }, { status: PPMDocumentsStatus.APPROVED }),
      createCompleteMovingExpense(
        { serviceMemberId },
        { status: PPMDocumentsStatus.REJECTED, reason: "We don't cover that expense." },
      ),
      createCompleteMovingExpense(
        { serviceMemberId },
        { movingExpenseType: expenseTypes.STORAGE, status: PPMDocumentsStatus.APPROVED },
      ),
      createCompleteMovingExpense(
        { serviceMemberId },
        { movingExpenseType: expenseTypes.STORAGE, status: PPMDocumentsStatus.EXCLUDED, reason: 'Invalid storage' },
      ),
    ];

    render(
      <MockProviders>
        <ReviewDocumentsSidePanel
          ppmShipmentInfo={mockPPMShipment}
          weightTickets={mockWeightTickets}
          proGearTickets={progearWeightTickets}
          expenseTickets={movingExpenses}
        />
      </MockProviders>,
    );

    const listItems = await screen.getAllByRole('listitem');
    expect(listItems).toHaveLength(9);

    // weight ticket 1
    expect(listItems[0]).toHaveTextContent(/Trip 1/);
    expect(listItems[0]).toHaveTextContent(/Accept/);

    // weight ticket 2
    expect(listItems[1]).toHaveTextContent(/Trip 2/);
    expect(listItems[1]).toHaveTextContent(/Reject/);
    expect(listItems[1]).toHaveTextContent(/No weight ticket/);

    // progear ticket 1
    expect(listItems[2]).toHaveTextContent(/Pro-gear 1/);
    expect(listItems[2]).toHaveTextContent(/Accept/);

    // progear ticket 2
    expect(listItems[3]).toHaveTextContent(/Pro-gear 2/);
    expect(listItems[1]).toHaveTextContent(/Reject/);
    expect(listItems[3]).toHaveTextContent(/Invalid weight ticket/);

    // moving expense 1 - non-storage 1
    expect(listItems[4]).toHaveTextContent(/Receipt 1/);
    expect(listItems[4]).toHaveTextContent(/Accept/);

    // moving expense 2 - non-storage 2
    expect(listItems[5]).toHaveTextContent(/Receipt 2/);
    expect(listItems[1]).toHaveTextContent(/Reject/);
    expect(listItems[5]).toHaveTextContent(/We don't cover that expense./);

    // moving expense 3 - storage 1
    expect(listItems[6]).toHaveTextContent(/Storage #1/);
    expect(listItems[6]).toHaveTextContent(/Accept/);

    // moving expense 4 - storage 2
    expect(listItems[7]).toHaveTextContent(/Storage #2/);
    expect(listItems[1]).toHaveTextContent(/Reject/);
    expect(listItems[7]).toHaveTextContent(/Invalid storage/);
  });

  it('renders read only view text', async () => {
    render(
      <MockProviders>
        <ReviewDocumentsSidePanel ppmShipmentInfo={mockPPMShipment} readOnly />
      </MockProviders>,
    );
    const h3 = await screen.getByRole('heading', { name: 'Sent to customer', level: 3 });
    expect(h3).toBeInTheDocument();
  });
});
