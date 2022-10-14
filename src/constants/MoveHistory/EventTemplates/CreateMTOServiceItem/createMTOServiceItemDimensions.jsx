import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { convertFromThousandthInchToInch } from 'utils/formatters';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const {
    type,
    height_thousandth_inches: heightThousandthInches,
    length_thousandth_inches: lengthThousandthInches,
    width_thousandth_inches: widthThousandthInches,
  } = historyRecord.changedValues;
  const height = convertFromThousandthInchToInch(heightThousandthInches);
  const length = convertFromThousandthInchToInch(lengthThousandthInches);
  const width = convertFromThousandthInchToInch(widthThousandthInches);

  const name = type === 'ITEM' ? 'item_size' : 'crate_size';

  const newChangedValues = {
    shipment_type: historyRecord.context[0]?.shipment_type,
    service_item_name: historyRecord.context[0]?.name,
    ...historyRecord.changedValues,
  };
  newChangedValues[name] = `${height}x${length}x${width} in`;

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: o.createMTOServiceItem,
  tableName: t.mto_service_item_dimensions,
  getEventNameDisplay: () => 'Requested service item',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
