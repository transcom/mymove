import PropTypes from 'prop-types';

export const MatchShape = PropTypes.shape({
  // eslint-disable-next-line react/forbid-prop-types
  params: PropTypes.object,
  isExact: PropTypes.bool,
  path: PropTypes.string,
  url: PropTypes.string,
});

export const LocationShape = PropTypes.shape({
  key: PropTypes.string,
  pathname: PropTypes.string,
  search: PropTypes.string,
  hash: PropTypes.string,
  // eslint-disable-next-line react/forbid-prop-types
  state: PropTypes.object,
});

export const HistoryShape = PropTypes.shape({
  push: PropTypes.func.isRequired,
});
