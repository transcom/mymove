import { screen, render } from '@testing-library/react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import updateMTOShipmentDeprecatePaymentRequest from 'constants/MoveHistory/EventTemplates/UpdateMTOShipment/updateMTOShipmentDeprecatePaymentRequest';

describe('when updating a MTO shipment', () => {
  const historyRecord = {
    DEPRECATE: {
      action: 'UPDATE',
      eventName: 'updateMTOShipment',
      tableName: 'payment_requests',
      changedValues: {
        status: 'DEPRECATED',
      },
    },
    RECALCULATE: {
      action: 'UPDATE',
      eventName: 'updateMTOShipment',
      tableName: 'payment_requests',
      changedValues: {
        recalculation_of_payment_request_id: '1234-5789-1',
      },
    },
  };
  describe('and given a deprecated payment request history record', () => {
    it('matches the events to the correct template for update an MTO shipment template', () => {
      const template = getTemplate(historyRecord.DEPRECATE);
      expect(template).toMatchObject(updateMTOShipmentDeprecatePaymentRequest);
    });

    describe('it displays the proper labeled details for the component', () => {
      it.each([['Status', ': Deprecated']])(
        'should display the label `%s` for the value `%s`',
        async (label, value) => {
          const template = getTemplate(historyRecord.DEPRECATE);

          render(template.getDetails(historyRecord.DEPRECATE));
          expect(screen.getByText(label)).toBeInTheDocument();
          expect(screen.getByText(value)).toBeInTheDocument();
        },
      );
    });
  });

  describe('and given a recalculation_of_payment_request_id, displays recalculated payment request', () => {
    it('matches the events to the correct template for update an MTO shipment template', () => {
      const template = getTemplate(historyRecord.RECALCULATE);
      expect(template).toMatchObject(updateMTOShipmentDeprecatePaymentRequest);
    });

    describe('it displays the proper labeled details for the component', () => {
      const template = getTemplate(historyRecord.RECALCULATE);

      render(template.getDetails(historyRecord.RECALCULATE));
      expect(screen.getByText('Recalculated payment request')).toBeInTheDocument();
    });
  });
});
