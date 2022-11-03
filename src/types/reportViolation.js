import PropTypes from 'prop-types';

import { PWSViolationShape } from './pwsViolation';

export const ReportViolationShape = PropTypes.shape({
  id: PropTypes.string,
  reportId: PropTypes.string,
  violationId: PropTypes.string,
  violation: PWSViolationShape,
});

export default {
  ReportViolationShape,
};
