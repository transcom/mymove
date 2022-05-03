import React from 'react';

import LabeledDetails from './LabeledDetails';
import PaymentDetails from './PaymentDetails';

import { HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';
import getMoveHistoryEventTemplate, { detailsTypes } from 'constants/moveHistoryEventTemplate';

const MoveHistoryDetailsSelector = ({ historyRecord }) => {
  const eventTemplate = getMoveHistoryEventTemplate(historyRecord);

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
    case detailsTypes.STATUS:
      return (
        <div>
          <b>Status</b>: {eventTemplate.getStatusDetails()}
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
  historyRecord: {},
};

export default MoveHistoryDetailsSelector;
