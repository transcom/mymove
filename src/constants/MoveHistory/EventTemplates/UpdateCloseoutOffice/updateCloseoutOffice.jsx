import React from 'react';

import styles from 'pages/Office/MoveHistory/LabeledDetails.module.scss';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateCloseoutOffice,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => {
    const result = historyRecord.context?.find((e) => e.closeout_office_name);
    const displayLineItem =
      (result?.closeout_office_name && (
        <>
          <b>Closeout office</b>: {result.closeout_office_name}
        </>
      )) ||
      '-';
    return (
      (historyRecord?.oldValues?.locator && (
        <>
          <span className={styles.shipmentType}>#{historyRecord?.oldValues?.locator}</span>
          {displayLineItem}
        </>
      )) ||
      displayLineItem
    );
  },
};
