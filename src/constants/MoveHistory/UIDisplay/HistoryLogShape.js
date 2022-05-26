import PropTypes from 'prop-types';

const HistoryLogValuesShape = PropTypes.object;
export const HistoryLogContextShape = PropTypes.arrayOf(PropTypes.object);

export const HistoryLogRecordShape = PropTypes.shape({
  action: PropTypes.string,
  changedValues: HistoryLogValuesShape,
  context: HistoryLogContextShape,
  eventName: PropTypes.string,
  oldValues: HistoryLogValuesShape,
  tableName: PropTypes.string,
});
