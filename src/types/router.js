import { shape, bool, string, func } from 'prop-types';

export const MatchShape = shape({
  params: shape({}),
  isExact: bool,
  path: string,
  url: string,
});

export const LocationShape = shape({
  key: string,
  pathname: string,
  search: string,
  hash: string,
  state: shape({}),
});

export const HistoryShape = shape({
  push: func.isRequired,
});
