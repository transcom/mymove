import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('When given a updated weight', () => {
  const progearRecord = {
    action: a.UPDATE,
    changedValues: {
      weight: '1000',
      description: 'Pro-gear Weight',
      status: 'APPROVED',
      reason: 'Test reason',
    },
    oldValues: {
      belongs_to_self: true,
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'f10be',
      },
    ],
    eventName: o.updateProGearWeightTicket,
    tableName: t.progear_weight_tickets,
  };

  const spouseProgearRecord = {
    action: a.UPDATE,
    changedValues: {
      weight: '277',
      description: 'Spouse Pro-gear Weight',
      status: 'REJECTED',
    },
    oldValues: {
      belongs_to_self: false,
    },
    context: [
      {
        shipment_type: 'PPM',
        shipment_locator: 'RQ38D4-01',
        shipment_id_abbr: 'f10be',
      },
    ],
    eventName: o.updateProGearWeightTicket,
    tableName: t.progear_weight_tickets,
  };

  it('displays shipment type, shipment ID, and service member properly', () => {
    const template = getTemplate(progearRecord);

    render(template.getDetails(progearRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01, Service member pro-gear')).toBeInTheDocument();
  });

  it('displays shipment type, shipment ID, and spouse properly', () => {
    const template = getTemplate(spouseProgearRecord);

    render(template.getDetails(spouseProgearRecord));
    expect(screen.getByText('PPM shipment #RQ38D4-01, Spouse pro-gear')).toBeInTheDocument();
  });

  it('displays weight updates', () => {
    const template = getTemplate(progearRecord);

    render(template.getDetails(progearRecord));
    expect(screen.getByText('Weight')).toBeInTheDocument();
    expect(screen.getByText(': 1,000 lbs')).toBeInTheDocument();
  });

  it('displays an approved request', () => {
    const template = getTemplate(progearRecord);

    render(template.getDetails(progearRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': APPROVED')).toBeInTheDocument();
    expect(screen.getByText('Reason')).toBeInTheDocument();
    expect(screen.getByText(': Test reason')).toBeInTheDocument();
  });

  it('displays a rejected request with reason', () => {
    const template = getTemplate(spouseProgearRecord);

    render(template.getDetails(spouseProgearRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(': REJECTED')).toBeInTheDocument();
    expect(screen.getByText('Description')).toBeInTheDocument();
    expect(screen.getByText(': Spouse Pro-gear Weight')).toBeInTheDocument();
  });
});
