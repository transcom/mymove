import React from 'react';
import PropTypes from 'prop-types';

import styles from './LabeledDetails.module.scss';

import dateFields from 'constants/MoveHistory/Database/DateFields';
import fieldMappings from 'constants/MoveHistory/Database/FieldMappings';
import weightFields from 'constants/MoveHistory/Database/WeightFields';
import { shipmentTypes } from 'constants/shipments';
import { HistoryLogRecordShape } from 'constants/MoveHistory/UIDisplay/HistoryLogShape';
import optionFields from 'constants/MoveHistory/Database/OptionFields';
import { formatCustomerDate } from 'utils/formatters';

const retrieveTextToDisplay = (fieldName, value) => {
  const emptyValue = '—';
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

  if (!displayValue) {
    displayValue = emptyValue;
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
  // TODO: [ MB-12182 ] This will include a shipment ID label in the future
  if ('shipment_type' in changedValuesToUse) {
    shipmentDisplay = shipmentTypes[changedValuesToUse.shipment_type];
    shipmentDisplay += ' shipment';
    delete changedValuesToUse.shipment_type;
  }

  if ('service_item_name' in changedValuesToUse) {
    shipmentDisplay += `, ${changedValuesToUse.service_item_name}`;
    delete changedValuesToUse.service_item_name;
  }

  // Filter out empty values unless they used to be non-empty
  const dbFieldsToDisplay = Object.keys(fieldMappings).filter((dbField) => {
    return (
      changedValuesToUse[dbField] ||
      (dbField in changedValuesToUse && historyRecord.oldValues && dbField in historyRecord.oldValues)
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
  getDetailsLabeledDetails: PropTypes.func,
};

LabeledDetails.defaultProps = {
  historyRecord: {},
  getDetailsLabeledDetails: null,
};

export default LabeledDetails;
