import { string, shape } from 'prop-types';

export const AlertStateShape = shape({
  alertType: string.isRequired,
  message: string.isRequired,
});

export default {
  AlertStateShape,
};
