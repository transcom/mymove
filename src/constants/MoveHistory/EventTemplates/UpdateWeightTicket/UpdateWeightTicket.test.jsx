import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When given a updated weight', () => {
  const netWeightRecord = {
    action: a.UPDATE,
    changedValues: {
      adjust_net_weight: '1000',
      net_weight_remarks: 'Adjusted based on uploaded doc',
      status: 'APPROVED',
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'f10be',
      },
    ],
    eventName: o.updateWeightTicket,
    tableName: t.weight_tickets,
  };

  const trailerRecord = {
    action: a.UPDATE,
    changedValues: {
      empty_weight: 1000,
      full_weight: 7999,
      owns_trailer: true,
      reason: 'They do not own it',
      status: 'REJECTED',
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'f10be',
      },
    ],
    eventName: o.updateWeightTicket,
    tableName: t.weight_tickets,
  };

  it('displays shipment type, shipment ID, and service item name properly', () => {
    const template = getTemplate(netWeightRecord);

    render(template.getDetails(netWeightRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01')).toBeInTheDocument();
  });

  it('displays an approved request', () => {
    const template = getTemplate(netWeightRecord);

    render(template.getDetails(netWeightRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': APPROVED')).toBeInTheDocument();
  });

  it('displays a rejected request with reason', () => {
    const template = getTemplate(trailerRecord);

    render(template.getDetails(trailerRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': REJECTED')).toBeInTheDocument();
    expect(screen.getByText('Reason')).toBeInTheDocument();
    expect(screen.getByText(': They do not own it')).toBeInTheDocument();
  });

  it('displays the contractor remarks', () => {
    const template = getTemplate(netWeightRecord);

    render(template.getDetails(netWeightRecord));
    expect(screen.getByText('Remarks')).toBeInTheDocument();
    expect(screen.getByText(': Adjusted based on uploaded doc')).toBeInTheDocument();
  });

  it('displays trailer data', () => {
    const template = getTemplate(trailerRecord);

    render(template.getDetails(trailerRecord));
    expect(screen.getByText('Trailer')).toBeInTheDocument();
    expect(screen.getByText(': Yes')).toBeInTheDocument();
  });

  it('displays trailer weight data', () => {
    const template = getTemplate(trailerRecord);

    render(template.getDetails(trailerRecord));
    expect(screen.getByText('Empty weight')).toBeInTheDocument();
    expect(screen.getByText(': 1,000 lbs')).toBeInTheDocument();
    expect(screen.getByText('Full weight')).toBeInTheDocument();
    expect(screen.getByText(': 7,999 lbs')).toBeInTheDocument();
  });
});
