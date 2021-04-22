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

// Components that are rendered at a route (<Route />) will receive these props automatically
export const RouteProps = {
  match: MatchShape.isRequired,
  location: LocationShape.isRequired,
  history: HistoryShape.isRequired,
};
