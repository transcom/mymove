import { screen, render } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateMTOShipmentAgent from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipmentAgent';

describe('when given an mto shipment agents update with mto agents table history record', () => {
  const historyRecord = {
    RELEASE: {
      action: 'UPDATE',
      eventName: 'updateMTOShipment',
      tableName: 'mto_agents',
      changedValues: {
        email: 'grace@email.com',
        first_name: 'Grace',
        phone: '555-555-5555',
      },
      oldValues: {
        agent_type: 'RELEASING_AGENT',
        email: 'gracie@email.com',
        first_name: 'Gracie',
        last_name: 'Griffin',
        phone: '555-555-5551',
      },
      context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
    },
    RECEIVE: {
      action: 'UPDATE',
      eventName: 'updateMTOShipment',
      tableName: 'mto_agents',
      changedValues: {
        email: 'nancy@email.com',
        first_name: 'Nancy',
        phone: '555-555-5555',
      },
      oldValues: {
        agent_type: 'RECEIVING_AGENT',
        email: 'nannye@email.com',
        first_name: 'Nanny',
        last_name: 'Drew',
        phone: '555-555-5551',
      },
      context: [{ shipment_type: 'HHG', shipment_id_abbr: 'a1b2c' }],
    },
  };

  describe('when agent type is Releasing', () => {
    it('matches the events to the correct template for update an MTO shipment template', () => {
      const template = getTemplate(historyRecord.RELEASE);
      expect(template).toMatchObject(updateMTOShipmentAgent);
    });

    it('displays the proper shipment title in the details column', () => {
      const template = getTemplate(historyRecord.RELEASE);
      render(template.getDetails(historyRecord.RELEASE));
      expect(screen.getAllByText('HHG shipment #A1B2C'));
    });

    it('it displays the proper labeled details for the given releasing agent', () => {
      const template = getTemplate(historyRecord.RELEASE);

      render(template.getDetails(historyRecord.RELEASE));
      expect(screen.getByText('Releasing agent')).toBeInTheDocument();
      expect(screen.getByText(': Grace Griffin, 555-555-5555, grace@email.com')).toBeInTheDocument();
    });
  });
  describe('when agent type is Receiving', () => {
    it('matches the events to the correct template for update an MTO shipment template', () => {
      const template = getTemplate(historyRecord.RECEIVE);
      expect(template).toMatchObject(updateMTOShipmentAgent);
    });

    it('displays the proper shipment title in the details column', () => {
      const template = getTemplate(historyRecord.RECEIVE);
      render(template.getDetails(historyRecord.RECEIVE));
      expect(screen.getAllByText('HHG shipment #A1B2C'));
    });

    it('it displays the proper labeled details for the given releasing agent', () => {
      const template = getTemplate(historyRecord.RECEIVE);

      render(template.getDetails(historyRecord.RECEIVE));
      expect(screen.getByText('Receiving agent')).toBeInTheDocument();
      expect(screen.getByText(': Nancy Drew, 555-555-5555, nancy@email.com')).toBeInTheDocument();
    });
  });
});
