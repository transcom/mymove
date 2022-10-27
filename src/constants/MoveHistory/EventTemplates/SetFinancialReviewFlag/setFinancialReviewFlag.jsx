import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.setFinancialReviewFlag,
  tableName: t.moves,
  getEventNameDisplay: () => 'Flagged move',
  getDetails: (historyRecord) => (
    <>
      {historyRecord.changedValues?.financial_review_flag === 'true'
        ? 'Move flagged for financial review'
        : 'Move unflagged for financial review'}
      <LabeledDetails historyRecord={historyRecord} />
    </>
  ),
};
