import { shape, string, func } from 'prop-types';

export const LocationShape = shape({
  key: string,
  pathname: string,
  search: string,
  hash: string,
  state: shape({}),
});

export const RouterShape = shape({
  navigate: func,
  location: LocationShape,
  params: shape({}),
});
