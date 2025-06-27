import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/ReviewShipmentAddressUpdate/reviewShipmentAddressUpdateShipment';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given a updated shipment address request, update move history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    eventName: o.reviewShipmentAddressUpdate,
    tableName: t.mto_shipments,
    context: [
      {
        shipment_id_abbr: '2fa5c',
        shipment_type: 'HHG',
        shipment_locator: 'ABC123-01',
      },
    ],
  };

  it('correctly matches the record', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Shipment destination address updated');
  });

  it('displays the proper value in the details field', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #ABC123-01', { exact: false })).toBeInTheDocument();
  });
});
