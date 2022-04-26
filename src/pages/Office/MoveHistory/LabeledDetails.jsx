import React from 'react';

import { HistoryLogValuesShape, dbFieldToDisplayName, dbWeightFields } from 'constants/historyLogUIDisplayName';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatMoveHistoryFullAddress } from 'utils/formatters';

const retrieveTextToDisplay = (fieldName, value) => {
  const displayName = dbFieldToDisplayName[fieldName];
  let displayValue = value;

  if (displayName === dbFieldToDisplayName.storage_in_transit) {
    displayValue = `${displayValue} days`;
  } else if (dbWeightFields.includes(fieldName)) {
    displayValue = `${displayValue} lbs`;
  }

  return {
    displayName,
    displayValue,
  };
};

const LabeledDetails = ({ changedValues, oldValues }) => {
  const backfilledChangedValues = {
    street_address_1: oldValues.street_address_1,
    street_address_2: oldValues.street_address_2,
    city: oldValues.city,
    state: oldValues.state,
    postal_code: oldValues.postal_code,
    ...changedValues,
  };

  const changedValuesWithFormattedAddress = {
    ...changedValues,
    address: formatMoveHistoryFullAddress(backfilledChangedValues),
  };

  const dbFieldsToDisplay = Object.keys(dbFieldToDisplayName).filter((dbField) => {
    return changedValuesWithFormattedAddress[dbField];
  });

  return (
    <div>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(
          modelField,
          changedValuesWithFormattedAddress[modelField],
        );

        return (
          <div key={modelField} className={descriptionListStyles.row}>
            <b>{displayName}</b>: {displayValue}
          </div>
        );
      })}
    </div>
  );
};

LabeledDetails.propTypes = {
  changedValues: HistoryLogValuesShape,
  oldValues: HistoryLogValuesShape,
};

LabeledDetails.defaultProps = {
  changedValues: {},
  oldValues: {},
};

export default LabeledDetails;
