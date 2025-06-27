import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateShipmentDestinationAddress/updateShipmentAddressUpdate';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given a updated shipment address request, update move history record', () => {
  const historyRecord = {
    action: a.INSERT,
    eventName: o.updateShipmentDestinationAddress,
    tableName: t.shipment_address_updates,
    changedValues: { contractor_remarks: 'requesting an address change' },
  };

  it('correctly matches the record', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);
    expect(template.getEventNameDisplay()).toEqual('Shipment destination address request');
  });

  it('displays the proper value in the details field', () => {
    const template = getTemplate(historyRecord);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('Contractor remarks')).toBeInTheDocument();
    expect(screen.getByText(': requesting an address change')).toBeInTheDocument();
  });
});
