import React from 'react';

import styles from './LabeledDetails.module.scss';

import dateFields from 'constants/MoveHistory/Database/DateFields';
import fieldMappings from 'constants/MoveHistory/Database/FieldMappings';
import weightFields from 'constants/MoveHistory/Database/WeightFields';
import { shipmentTypes } from 'constants/shipments';
import { HistoryLogRecordShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';
import optionFields from 'constants/MoveHistory/Database/OptionFields';
import { formatCustomerDate, formatWeight } from 'utils/formatters';

const retrieveTextToDisplay = (fieldName, value) => {
  const emptyValue = '—';
  const displayName = fieldMappings[fieldName];
  let displayValue = value;

  if (displayName === fieldMappings.storage_in_transit) {
    displayValue = `${displayValue} days`;
  } else if (weightFields[fieldName]) {
    // turn string value into number so it can be formatted correctly
    displayValue = formatWeight(Number(displayValue));
  } else if (optionFields[displayValue]) {
    displayValue = optionFields[displayValue];
  } else if (dateFields[fieldName]) {
    displayValue = formatCustomerDate(displayValue);
  }

  if (!displayValue) {
    displayValue = emptyValue;
  }

  return {
    displayName,
    displayValue,
  };
};

const LabeledDetails = ({ historyRecord }) => {
  const changedValuesToUse = historyRecord.changedValues;
  let shipmentDisplay = '';

  // Check for shipment_type to use it as a header for the row
  if ('shipment_type' in changedValuesToUse) {
    shipmentDisplay = shipmentTypes[changedValuesToUse.shipment_type];
    shipmentDisplay += ` shipment #${changedValuesToUse.shipment_id_display}`;
    delete changedValuesToUse.shipment_type;
  }

  if ('service_item_name' in changedValuesToUse) {
    shipmentDisplay += `, ${changedValuesToUse.service_item_name}`;
    delete changedValuesToUse.service_item_name;
  }

  /* Filter out empty values unless they used to be non-empty
     These values may be non-nullish in oldValues and nullish in changedValues */
  const dbFieldsToDisplay = Object.keys(fieldMappings).filter((dbField) => {
    return (
      changedValuesToUse[dbField] ||
      (dbField in changedValuesToUse && historyRecord.oldValues && historyRecord.oldValues[dbField])
    );
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
};

LabeledDetails.defaultProps = {
  historyRecord: {},
};

export default LabeledDetails;
