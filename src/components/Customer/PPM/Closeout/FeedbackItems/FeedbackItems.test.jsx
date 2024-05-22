import React from 'react';
import { render, screen } from '@testing-library/react';

import FeedbackItems from './FeedbackItems';

import { FEEDBACK_DOCUMENT_TYPES } from 'constants/ppmFeedback';

const formattedWeightTickets = [
  [
    {
      key: 'vehicleDescription',
      label: 'Vehicle description: ',
      value: '2023 F-150',
    },
    {
      key: 'emptyWeight',
      label: 'Empty: ',
      secondaryKey: 'submittedEmptyWeight',
      value: '1,000 lbs',
    },
    {
      key: 'fullWeight',
      label: 'Full: ',
      secondaryKey: 'submittedFullWeight',
      value: '8,000 lbs',
    },
    {
      key: 'tripWeight',
      label: 'Trip weight: ',
      value: '7,000 lbs',
    },
    {
      key: 'ownsTrailer',
      label: 'Trailer: ',
      value: 'yes',
    },
    {
      key: 'status',
      value: 'APPROVED',
    },
  ],
];

const formattedProGearWeightTickets = [
  [
    {
      key: 'belongsToSelf',
      label: '',
      value: 'Pro-Gear',
    },
    {
      key: 'description',
      label: 'Description: ',
      value: 'my gear',
    },
    {
      key: 'weight',
      label: 'Weight: ',
      secondaryKey: 'submittedWeight',
      value: '1,222 lbs',
    },
    {
      key: 'status',
      value: "This doesn't make sense, this receipt is from 1999",
      label: 'REJECTED: ',
    },
  ],
];

const formattedMovingExpenses = [
  [
    {
      key: 'movingExpenseType',
      label: 'Type: ',
      value: 'Storage',
    },
    {
      key: 'description',
      label: 'Description: ',
      value: 'A spot to store things',
    },
    {
      key: 'amount',
      label: 'Amount: ',
      secondaryKey: 'submittedAmount',
      value: '$2,000.00',
      secondaryValue: '$2,200.00',
    },
    {
      key: 'sitStartDate',
      label: 'SIT start date: ',
      secondaryKey: 'submittedSitStartDate',
      value: '02 Apr 2024',
      secondaryValue: '01 Apr 2024',
    },
    {
      key: 'sitEndDate',
      label: 'SIT end date: ',
      secondaryKey: 'submittedSitEndDate',
      value: '01 Jun 2024',
      secondaryValue: '02 Jun 2024',
    },
    {
      key: 'status',
      value: 'EDITED',
    },
  ],
];

const formattedProGearWeightTicketsExcluded = [
  [
    {
      key: 'belongsToSelf',
      label: '',
      value: 'Pro-Gear',
    },
    {
      key: 'description',
      label: 'Description: ',
      value: 'my gear',
    },
    {
      key: 'weight',
      label: 'Weight: ',
      secondaryKey: 'submittedWeight',
      value: '1,222 lbs',
    },
    {
      key: 'status',
      value: 'Claim on taxes',
      label: 'EXCLUDED: ',
    },
  ],
];

describe('FeedbackItems component', () => {
  it('displays feedback for weight tickets', () => {
    render(<FeedbackItems documents={formattedWeightTickets} doctype={FEEDBACK_DOCUMENT_TYPES.WEIGHT} />);

    expect(screen.getByText('Vehicle description:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('2023 F-150')).toBeInTheDocument();
    expect(screen.getByText('Empty:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('1,000 lbs')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Full:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('8,000 lbs')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Trip weight:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('7,000 lbs')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Trailer:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('yes')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('APPROVED')).toBeInstanceOf(HTMLSpanElement);
  });

  it('displays feedback for pro-gear sets', () => {
    render(<FeedbackItems documents={formattedProGearWeightTickets} doctype={FEEDBACK_DOCUMENT_TYPES.PRO_GEAR} />);

    expect(screen.getByText('Pro-Gear')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Description:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('my gear')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Weight:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('1,222 lbs')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('REJECTED:')).toBeInstanceOf(HTMLSpanElement);
  });

  it('displays feedback for moving expenses', () => {
    render(<FeedbackItems documents={formattedMovingExpenses} doctype={FEEDBACK_DOCUMENT_TYPES.MOVING_EXPENSE} />);

    expect(screen.getByText('Type:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Storage')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Description:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Amount:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('$2,000.00')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('SIT start date:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('02 Apr 2024')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('SIT end date:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('01 Jun 2024')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('EDITED')).toBeInstanceOf(HTMLSpanElement);
  });

  it('displays the edited values in parentheses', () => {
    render(<FeedbackItems documents={formattedMovingExpenses} doctype={FEEDBACK_DOCUMENT_TYPES.MOVING_EXPENSE} />);

    expect(screen.getByText('SIT start date:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('02 Apr 2024')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('(01 Apr 2024)')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('SIT end date:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('01 Jun 2024')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('(02 Jun 2024)')).toBeInstanceOf(HTMLSpanElement);
  });

  it('displays the edited status when values were edited', () => {
    render(<FeedbackItems documents={formattedMovingExpenses} doctype={FEEDBACK_DOCUMENT_TYPES.MOVING_EXPENSE} />);

    expect(screen.getByText('EDITED')).toBeInstanceOf(HTMLSpanElement);
  });

  it('displays the approved status when status is approved', () => {
    render(<FeedbackItems documents={formattedWeightTickets} doctype={FEEDBACK_DOCUMENT_TYPES.WEIGHT} />);

    expect(screen.getByText('APPROVED')).toBeInstanceOf(HTMLSpanElement);
  });

  it('displays the excluded status with remarks when status is excluded', () => {
    render(
      <FeedbackItems documents={formattedProGearWeightTicketsExcluded} doctype={FEEDBACK_DOCUMENT_TYPES.PRO_GEAR} />,
    );

    expect(screen.getByText('EXCLUDED:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText('Claim on taxes')).toBeInstanceOf(HTMLSpanElement);
  });

  it('displays the rejected status with remarks when status is rejected', () => {
    render(<FeedbackItems documents={formattedProGearWeightTickets} doctype={FEEDBACK_DOCUMENT_TYPES.PRO_GEAR} />);

    expect(screen.getByText('REJECTED:')).toBeInstanceOf(HTMLSpanElement);
    expect(screen.getByText(`This doesn't make sense, this receipt is from 1999`)).toBeInstanceOf(HTMLSpanElement);
  });
});
