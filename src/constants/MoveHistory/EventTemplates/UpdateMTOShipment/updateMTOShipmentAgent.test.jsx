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
      context: [{ shipment_type: 'HHG', shipment_locator: 'RQ38D4-01', shipment_id_abbr: 'a1b2c' }],
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
      context: [{ shipment_type: 'HHG', shipment_locator: 'RQ38D4-01', shipment_id_abbr: 'a1b2c' }],
    },
    DELETED_RECEIVING_AGENT: {
      action: 'UPDATE',
      eventName: 'updateMTOShipment',
      tableName: 'mto_agents',
      changedValues: {
        deleted_at: '2025-01-21T15:39:24.890356+00:00',
        email: null,
        first_name: null,
        last_name: null,
        phone: null,
      },
      oldValues: {
        agent_type: 'RECEIVING_AGENT',
        email: 'john.smith@email.com',
        first_name: 'John',
        last_name: 'Smith',
        phone: '555-765-4321',
      },
      context: [{ shipment_type: 'HHG', shipment_locator: 'RQ38D4-01', shipment_id_abbr: 'a1b2c' }],
    },
    DELETED_RELEASING_AGENT: {
      action: 'UPDATE',
      eventName: 'updateMTOShipment',
      tableName: 'mto_agents',
      changedValues: {
        deleted_at: '2025-01-21T16:39:24.890356+00:00',
        email: null,
        first_name: null,
        last_name: null,
        phone: null,
      },
      oldValues: {
        agent_type: 'RELEASING_AGENT',
        email: 'jane.smith@email.com',
        first_name: 'Jane',
        last_name: 'Smith',
        phone: '555-123-4567',
      },
      context: [{ shipment_type: 'NTS', shipment_locator: 'RQ38D4-01', shipment_id_abbr: 'a1b2c' }],
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
      expect(screen.getAllByText('HHG shipment #RQ38D4-01'));
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
      expect(screen.getAllByText('HHG shipment #RQ38D4-01'));
    });

    it('it displays the proper labeled details for the given receiving agent', () => {
      const template = getTemplate(historyRecord.RECEIVE);

      render(template.getDetails(historyRecord.RECEIVE));
      expect(screen.getByText('Receiving agent')).toBeInTheDocument();
      expect(screen.getByText(': Nancy Drew, 555-555-5555, nancy@email.com')).toBeInTheDocument();
    });
  });

  describe('when agent is deleted', () => {
    it('displays deleted receiving agent', () => {
      const template = getTemplate(historyRecord.DELETED_RECEIVING_AGENT);
      render(template.getDetails(historyRecord.DELETED_RECEIVING_AGENT));
      expect(screen.getByText('Deleted Receiving agent on HHG shipment #RQ38D4-01')).toBeInTheDocument();
      expect(screen.getByText('Receiving agent')).toBeInTheDocument();
      expect(screen.getByText(': John Smith, 555-765-4321, john.smith@email.com')).toBeInTheDocument();
    });

    it('displays deleted releasing agent', () => {
      const template = getTemplate(historyRecord.DELETED_RELEASING_AGENT);
      render(template.getDetails(historyRecord.DELETED_RELEASING_AGENT));
      expect(screen.getByText('Deleted Releasing agent on NTS shipment #RQ38D4-01')).toBeInTheDocument();
      expect(screen.getByText('Releasing agent')).toBeInTheDocument();
      expect(screen.getByText(': Jane Smith, 555-123-4567, jane.smith@email.com')).toBeInTheDocument();
    });
  });
});
