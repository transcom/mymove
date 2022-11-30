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
      ['Status', ': SUBMITTED', historyRecord.RELOCATION],
      ['Report by date', ': 12 Oct 2022', historyRecord.RELOCATION],
      ['Orders date', ': 11 Oct 2022', historyRecord.RELOCATION],
      ['Orders type', ': Permanent Change Of Station (PCS)', historyRecord.RELOCATION],
      ['Origin duty location name', ': Los Angeles AFB', historyRecord.RELOCATION],
      ['New duty location name', ': Fairchild AFB', historyRecord.RELOCATION],
      ['Dependents included', ': No', historyRecord.RELOCATION],
      ['Rank', ': E-2', historyRecord.RELOCATION],
      ['Date of separation', ': 12 Oct 2022', historyRecord.SEPARATION],
      ['Orders type', ': Retirement', historyRecord.RETIREMENT],
      ['Date of retirement', ': 12 Oct 2022', historyRecord.RETIREMENT],
    ])('displays the proper details value for %s', async (label, value, record) => {
      const result = getTemplate(record);
      render(result.getDetails(record));
      expect(screen.getByText(label)).toBeInTheDocument();
      expect(screen.getByText(value)).toBeInTheDocument();
    });
  });
});
