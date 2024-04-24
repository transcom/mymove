import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import { formatDataForPPM } from 'utils/formatPPMData';

/** These keys are the ones that should be present when the user begins the PPM Document Upload process, as the
 *  first form in the process asks the user to put in their actual departure date, w2 address, etc.
 */
const objectKeysWhenUserStartsPPMUpload = [
  'actual_destination_postal_code',
  'actual_move_date',
  'actual_pickup_postal_code',
  'has_received_advance',
  'w2_address_id',
];

/**
 * Checks if the 'changed values' of the history record passed in contains all of the keys that are present when the user
 * starts the PPM Documentation process.
 * @param {object} historyRecord - Object containing history info like the action taken, values changed, etc.
 * @returns {boolean} True if the history record contains all of the keys in 'objectKeysWhenUserStartsPPMUpload'
 */
const hasAllKeys = (historyRecord) =>
  objectKeysWhenUserStartsPPMUpload.every((item) =>
    Object.prototype.hasOwnProperty.call(historyRecord.changedValues, item),
  );

const formatChangedValues = (historyRecord) => {
  const newChangedValues = {
    ...historyRecord.changedValues,
    ...getMtoShipmentLabel(historyRecord),
    ...formatDataForPPM(historyRecord),
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.ppm_shipments,
  getEventNameDisplay: (historyRecord) => {
    if (hasAllKeys(historyRecord)) return 'Customer Began PPM Document Process';
    return 'Updated shipment';
  },
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
