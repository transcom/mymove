import React from 'react';
import PropTypes from 'prop-types';

import styles from './LabeledDetails.module.scss';

import dateFields from 'constants/MoveHistory/Database/DateFields';
import fieldMappings from 'constants/MoveHistory/Database/FieldMappings';
import weightFields from 'constants/MoveHistory/Database/WeightFields';
import { shipmentTypes } from 'constants/shipments';
import { HistoryLogRecordShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';
import optionFields from 'constants/MoveHistory/Database/Orders';
import { formatCustomerDate } from 'utils/formatters';

const retrieveTextToDisplay = (fieldName, value) => {
  const displayName = fieldMappings[fieldName];
  let displayValue = value;

  if (displayName === fieldMappings.storage_in_transit) {
    displayValue = `${displayValue} days`;
  } else if (weightFields[fieldName]) {
    displayValue = `${displayValue} lbs`;
  } else if (optionFields[displayValue]) {
    displayValue = optionFields[displayValue];
  } else if (dateFields[fieldName]) {
    displayValue = formatCustomerDate(displayValue);
  }

  return {
    displayName,
    displayValue,
  };
};

const LabeledDetails = ({ historyRecord, getDetailsLabeledDetails }) => {
  let changedValuesToUse = historyRecord.changedValues;
  let shipmentDisplay = '';
  // run custom function to mutate changedValues to display if not null
  if (getDetailsLabeledDetails) {
    changedValuesToUse = getDetailsLabeledDetails(historyRecord);
  }

  // Check for shipment_type to use it as a header for the row
  if ('shipment_type' in changedValuesToUse) {
    shipmentDisplay = shipmentTypes[changedValuesToUse.shipment_type];
    shipmentDisplay += ' shipment';
    delete changedValuesToUse.shipment_type;
  }

  if ('service_item_name' in changedValuesToUse) {
    shipmentDisplay += `, ${changedValuesToUse.service_item_name}`;
    delete changedValuesToUse.service_item_name;
  }

  const dbFieldsToDisplay = Object.keys(fieldMappings).filter((dbField) => {
    return changedValuesToUse[dbField];
  });

  return (
    <>
      <span className={styles.shipmentType}>{shipmentDisplay}</span>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(modelField, changedValuesToUse[modelField]);

        return (
          <div key={modelField}>
            <b>{displayName}</b>: {displayValue}
          </div>
        );
      })}
    </>
  );
};

LabeledDetails.propTypes = {
  historyRecord: HistoryLogRecordShape,
  getDetailsLabeledDetails: PropTypes.func,
};

LabeledDetails.defaultProps = {
  historyRecord: {},
  getDetailsLabeledDetails: null,
};

export default LabeledDetails;
