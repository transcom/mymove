import React from 'react';
import PropTypes from 'prop-types';

import PlainTextDetails from './PlainTextDetails';

import {
  eventNamesWithLabeledDetails,
  eventNamesWithServiceItemDetails,
  eventNamesWithPlainTextDetails,
  HistoryLogValuesShape,
} from 'constants/historyLogUIDisplayName';

const formatChangedValues = (values) => {
  return values
    ? Object.keys(values).map((key) => (
        <div key={`${key}-${values[key]}`}>
          {key}: {values[key]}
        </div>
      ))
    : '';
};

const MoveHistoryDetailsSelector = ({ eventName, oldValues, changedValues }) => {
  if (eventNamesWithLabeledDetails[eventName]) {
    return (
      <div>
        Labeled {eventName}
        <div>old Values {formatChangedValues(oldValues)}</div>
        <div>changed values {formatChangedValues(changedValues)}</div>
      </div>
    );
  }

  if (eventNamesWithServiceItemDetails[eventName]) {
    return (
      <div>
        Service Items {eventName}
        <div>old Values {formatChangedValues(oldValues)}</div>
        <div>changed values {formatChangedValues(changedValues)}</div>
      </div>
    );
  }

  if (eventNamesWithPlainTextDetails[eventName]) {
    return <PlainTextDetails eventName={eventName} changedValues={changedValues} />;
  }

  return (
    <div>
      - {eventName}
      <div>old Values {formatChangedValues(oldValues)}</div>
      <div>changed values {formatChangedValues(changedValues)}</div>
    </div>
  );
};

MoveHistoryDetailsSelector.propTypes = {
  eventName: PropTypes.string,
  oldValues: HistoryLogValuesShape,
  changedValues: HistoryLogValuesShape,
};

MoveHistoryDetailsSelector.defaultProps = {
  eventName: '',
  oldValues: {},
  changedValues: {},
};

export default MoveHistoryDetailsSelector;
