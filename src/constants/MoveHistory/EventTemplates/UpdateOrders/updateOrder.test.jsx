import { render, screen } from '@testing-library/react';

import getTemplate from 'constants/MoveHistory/TemplateManager';
import e from 'constants/MoveHistory/EventTemplates/UpdateOrders/updateOrder';

describe('when given an Order update history record', () => {
  const historyRecord = {
    RELOCATION: {
      action: 'UPDATE',
      eventName: 'updateOrders',
      tableName: 'orders',
      eventNameDisplay: 'Updated orders',
      changedValues: {
        status: 'SUBMITTED',
        report_by_date: '2022-10-12',
        issue_date: '2022-10-11',
        orders_type: 'PERMANENT_CHANGE_OF_STATION',
        origin_duty_location_id: 'ID2',
        new_duty_location_id: 'ID2',
        has_dependents: 'false',
        grade: 'E_2',
      },
      context: [
        {
          new_duty_location_name: 'Fairchild AFB',
          origin_duty_location_name: 'Los Angeles AFB',
        },
      ],
    },
    SEPARATION: {
      action: 'UPDATE',
      eventName: 'updateOrders',
      tableName: 'orders',
      eventNameDisplay: 'Updated orders',
      changedValues: {
        report_by_date: '2022-10-12',
      },
      oldValues: {
        orders_type: 'SEPARATION',
      },
    },
    RETIREMENT: {
      action: 'UPDATE',
      eventName: 'updateOrders',
      tableName: 'orders',
      eventNameDisplay: 'Updated orders',
      changedValues: {
        report_by_date: '2022-10-12',
        orders_type: 'RETIREMENT',
      },
    },
  };
  it('correctly matches to the proper template', () => {
    const template = getTemplate(historyRecord.RELOCATION);
    expect(template).toMatchObject(e);
  });
  describe('When given a specific set of details for updated orders', () => {
    it.each([
      ['Status', ': SUBMITTED'],
      ['Report by date', ': 12 Oct 2022'],
      ['Orders date', ': 11 Oct 2022'],
      ['Orders type', ': Permanent Change Of Station (PCS)'],
      ['Origin duty location name', ': Los Angeles AFB'],
      ['New duty location name', ': Fairchild AFB'],
      ['Dependents included', ': false'],
      ['Rank', ': E-2'],
    ])('displays the proper details value for %s', async (label, value) => {
      const result = getTemplate(historyRecord.RELOCATION);
      render(result.getDetails(historyRecord.RELOCATION));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
  describe('When given orders type is SEPARATION, it correctly displays label', () => {
    it.each([['Date of separation', ': 12 Oct 2022']])(
      'displays the proper details value for %s',
      async (label, value) => {
        const template = getTemplate(historyRecord.SEPARATION);
        render(template.getDetails(historyRecord.SEPARATION));
        expect(screen.getByText(label)).toBeInTheDocument();
        expect(screen.getByText(value)).toBeInTheDocument();
      },
    );
  });
  describe('When given orders type is RETIREMENT, it correctly displays label', () => {
    it.each([
      ['Orders type', ': Retirement'],
      ['Date of retirement', ': 12 Oct 2022'],
    ])('displays the proper details value for %s', async (label, value) => {
      const template = getTemplate(historyRecord.RETIREMENT);
      render(template.getDetails(historyRecord.RETIREMENT));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
