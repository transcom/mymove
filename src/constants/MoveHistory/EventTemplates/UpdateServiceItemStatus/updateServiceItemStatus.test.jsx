import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateServiceItemStatus/updateServiceItemStatus';

describe('when given an approved service item history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'APPROVED' },
    context: [
      {
        name: 'Domestic origin price',
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
      },
    ],
    eventName: 'updateMTOServiceItemStatus',
    tableName: 'mto_service_items',
  };

  it('correctly matches the Approved service item event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('renders the correct values in the event and details column for a approved service item', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Approved service item')).toBeInTheDocument();
    expect(screen.getByText('HHG shipment #A1B2C, Domestic origin price')).toBeInTheDocument();
  });
});

describe('when given rejected service item history record', () => {
  const historyRecord = {
    action: 'UPDATE',
    changedValues: { status: 'REJECTED', rejection_reason: "I've chosen to reject this item" },
    context: [
      {
        name: 'Domestic origin price',
        shipment_type: 'HHG',
        shipment_id_abbr: 'a1b2c',
      },
    ],
    eventName: 'updateMTOServiceItemStatus',
    tableName: 'mto_service_items',
  };

  it('correctly matches the Approved service item event to the proper template', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
  });

  it('renders the correct values in the event and details column for a rejected service item', () => {
    const template = getTemplate(historyRecord);

    render(template.getEventNameDisplay(historyRecord));
    render(template.getDetails(historyRecord));
    expect(screen.getByText('Rejected service item')).toBeInTheDocument();
    expect(screen.getByText('HHG shipment #A1B2C, Domestic origin price')).toBeInTheDocument();
  });

  describe('When given a specific set of details for a rejected service item', () => {
    it.each([['Reason', ": I've chosen to reject this item"]])(
      'displays the proper details value for %s',
      async (label, value) => {
        const template = getTemplate(historyRecord);
        render(template.getDetails(historyRecord));
        expect(screen.getByText(label)).toBeInTheDocument();
        expect(screen.getByText(value)).toBeInTheDocument();
      },
    );
  });
});
