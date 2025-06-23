import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateShipmentAfterSitExtCreation from 'constants/MoveHistory/EventTemplates/UpdateShipmentAfterSitExtCreation/updateShipmentAfterSitExtCreation';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

describe('when given a updateShipmentAfterSitExtCreation history record', () => {
  const historyRecord = {
    action: a.UPDATE,
    changedValues: {
      status: 'APPROVALS_REQUESTED',
    },
    context: [
      {
        shipment_id_abbr: 'bb22c',
        shipment_locator: 'J3XBDR-01',
        shipment_type: 'HHG',
      },
    ],
    eventName: o.createSITExtension,
    tableName: t.mto_shipments,
  };

  it('matches the template from getTemplate', () => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(updateShipmentAfterSitExtCreation);
  });

  it('returns the correct event display name', () => {
    expect(updateShipmentAfterSitExtCreation.getEventNameDisplay()).toEqual('Updated shipment');
  });

  it('renders the details via LabeledDetails with merged changed values', () => {
    const template = getTemplate(historyRecord);
    render(template.getDetails(historyRecord));

    // Check for the presence of the changed values.
    // The actual keys and values displayed depend on your LabeledDetails implementation.
    // Here we expect the values from changedValues to be rendered.
    expect(screen.getByText(/Status/i)).toBeInTheDocument();
    expect(screen.getByText(/APPROVALS REQUESTED/i)).toBeInTheDocument();
  });
});
