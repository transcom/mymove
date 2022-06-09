import React from 'react';

import LabeledPaymentRequestDetails from './LabeledPaymentRequestDetails';
import LabeledDetails from './LabeledDetails';
import PaymentDetails from './PaymentDetails';

import { HistoryLogRecordShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';
import getTemplate from 'constants/MoveHistory/TemplateManager';
import detailsTypes from 'constants/MoveHistory/UIDisplay/DetailsTypes';

const MoveHistoryDetailsSelector = ({ historyRecord }) => {
  const eventTemplate = getTemplate(historyRecord);

  switch (eventTemplate.detailsType) {
    case detailsTypes.LABELED:
      return (
        <LabeledDetails
          historyRecord={historyRecord}
          getDetailsLabeledDetails={eventTemplate.getDetailsLabeledDetails}
        />
      );
    case detailsTypes.PAYMENT:
      return <PaymentDetails context={historyRecord.context} />;
    case detailsTypes.LABELED_PAYMENT_REQUEST:
      return (
        <LabeledPaymentRequestDetails
          context={historyRecord.context}
          getLabeledPaymentRequestDetails={eventTemplate.getLabeledPaymentRequestDetails}
        />
      );
    case detailsTypes.STATUS:
      return (
        <div>
          <b>Status</b>: {eventTemplate.getStatusDetails(historyRecord)}
        </div>
      );
    case detailsTypes.PLAIN_TEXT:
    default:
      return <div>{eventTemplate.getDetailsPlainText(historyRecord)}</div>;
  }
};

MoveHistoryDetailsSelector.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

MoveHistoryDetailsSelector.defaultProps = {
  historyRecord: [],
};

export default MoveHistoryDetailsSelector;
