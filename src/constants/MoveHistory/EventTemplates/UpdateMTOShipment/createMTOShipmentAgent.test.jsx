import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/createMTOShipmentAgent';

describe('when given a historyRecord that updates a receiving/releasing agent', () => {
  const historyRecord = {
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

  it('matches the events to the correct template', () => {
    const template = getTemplate(historyRecord);
    const template2 = getTemplate(historyRecord2);

    expect(template).toMatchObject(e);
    expect(template2).toMatchObject(e);
  });

  it('displays the proper shipment title in the details column', () => {
    const template = getTemplate(historyRecord);
    const template2 = getTemplate(historyRecord2);

    render(template.getDetails(historyRecord));
    render(template2.getDetails(historyRecord2));
    expect(screen.getAllByText('HHG shipment #A1B2C'));
  });

  describe('it displayes the proper labeled details for the given releasing agent', () => {
    it.each([
      ['Releasing agent', ': Grace Griffin, 555-555-5555, grace@email.com'],
      ['First name', ': Grace'],
      ['Last name', ': Griffin'],
    ])('should display the label `%s` for the value `%s`', async (label, value) => {
      const template = getTemplate(historyRecord);

      render(template.getDetails(historyRecord));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });

  describe('it displayes the proper labeled details for the given receiving agent', () => {
    it.each([
      ['Receiving agent', ': Catalina Washington, 999-999-9999, catalina@email.com'],
      ['First name', ': Catalina'],
      ['Last name', ': Washington'],
    ])('should display the label `%s` for the value `%s`', async (label, value) => {
      const template = getTemplate(historyRecord2);

      render(template.getDetails(historyRecord2));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
