import { bool, string, shape } from 'prop-types';

export const MatchShape = shape({
  isExact: bool.isRequired,
  params: shape({
    moveCode: string.isRequired,
  }),
  path: string.isRequired,
  url: string.isRequired,
});

export default {
  MatchShape,
};
