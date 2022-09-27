import PropTypes from 'prop-types';

import { EvaluationReportShape } from './evaluationReport';

export const ReportViolationShape = PropTypes.shape({
  id: PropTypes.string,
  reportId: PropTypes.string,
  violationId: PropTypes.string,
  violation: EvaluationReportShape,
});

export default {
  ReportViolationShape,
};
