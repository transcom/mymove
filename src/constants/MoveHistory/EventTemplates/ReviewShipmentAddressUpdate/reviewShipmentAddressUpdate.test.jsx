import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/EventTemplates/ReviewShipmentAddressUpdate/reviewShipmentAddressUpdate';

describe('when given a Review Shipment Address Update history record', () => {
  const context = [
    {
      name: 'Shipment',
    },
  ];

  it('displays the correct event name', () => {
    const historyRecord = {
      action: a.UPDATE,
      changedValues: {
        status: 'APPROVED',
      },
      context,
      eventName: o.reviewShipmentAddressUpdate,
      tableName: t.shipment_address_updates,
    };

    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay(historyRecord)).toEqual('Shipment destination address update');
  });

  it('displays the status as "Approved"', () => {
    const historyRecord = {
      action: a.UPDATE,
      changedValues: {
        status: 'APPROVED',
      },
      context,
      eventName: o.reviewShipmentAddressUpdate,
      tableName: t.shipment_address_updates,
    };

    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(/APPROVED/)).toBeInTheDocument();
  });

  it('displays the status as "Rejected"', () => {
    const historyRecord = {
      action: a.UPDATE,
      changedValues: {
        status: 'REJECTED',
      },
      context,
      eventName: o.reviewShipmentAddressUpdate,
      tableName: t.shipment_address_updates,
    };

    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(/REJECTED/)).toBeInTheDocument();
  });

  it('returns null if the status is not "APPROVED" or "REJECTED"', () => {
    const historyRecord = {
      action: a.UPDATE,
      changedValues: {
        status: 'PENDING',
      },
      context,
      eventName: o.reviewShipmentAddressUpdate,
      tableName: t.shipment_address_updates,
    };

    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText(/PENDING/)).toBeInTheDocument();
  });
});
