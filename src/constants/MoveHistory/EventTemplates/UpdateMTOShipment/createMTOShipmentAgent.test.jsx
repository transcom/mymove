import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/createMTOShipmentAgent';

describe('when given a historyRecord that updates a receiving/releasing agent', () => {
  const historyRecord1 = {
    action: 'INSERT',
    eventName: 'updateMTOShipment',
    tableName: 'mto_agents',
    changedValues: {
      email: 'grace@email.com',
      first_name: 'Grace',
      last_name: 'Griffin',
      phone: '555-555-5555',
      agent_type: 'RELEASING_AGENT',
    },
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
  };
  const historyRecord2 = {
    action: 'INSERT',
    eventName: 'updateMTOShipment',
    tableName: 'mto_agents',
    changedValues: {
      email: 'catalina@email.com',
      first_name: 'Catalina',
      last_name: 'Washington',
      phone: '999-999-9999',
      agent_type: 'RECEIVING_AGENT',
    },
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
  };
  const historyRecord3 = {
    action: 'INSERT',
    eventName: 'createMTOShipment',
    tableName: 'mto_agents',
    changedValues: {
      email: 'grace@email.com',
      first_name: 'Grace',
      last_name: 'Griffin',
      phone: '555-555-5555',
      agent_type: 'RELEASING_AGENT',
    },
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
  };
  const historyRecord4 = {
    action: 'INSERT',
    eventName: 'createMTOShipment',
    tableName: 'mto_agents',
    changedValues: {
      email: 'catalina@email.com',
      first_name: 'Catalina',
      last_name: 'Washington',
      phone: '999-999-9999',
      agent_type: 'RECEIVING_AGENT',
    },
    context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
  };
  it.each([
    ['Releasing agent', ': Grace Griffin, 555-555-5555, grace@email.com', historyRecord1],
    ['Receiving agent', ': Catalina Washington, 999-999-9999, catalina@email.com', historyRecord2],
    ['Releasing agent', ': Grace Griffin, 555-555-5555, grace@email.com', historyRecord3],
    ['Receiving agent', ': Catalina Washington, 999-999-9999, catalina@email.com', historyRecord4],
  ])('should display the label `%s` for the value `%s`', async (label, value, historyRecord) => {
    const template = getTemplate(historyRecord);
    expect(template).toMatchObject(e);

    render(template.getDetails(historyRecord));
    expect(screen.getByText('HHG shipment #A1B2C')).toBeInTheDocument();
    expect(screen.getByText(label)).toBeInTheDocument();
    expect(screen.getByText(value)).toBeInTheDocument();
  });
});
