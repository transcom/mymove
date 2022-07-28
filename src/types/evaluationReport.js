import PropTypes from 'prop-types';

export const EvaluationReportShape = PropTypes.shape({
  id: PropTypes.string,
  location: PropTypes.string,
  violations: PropTypes.bool,
  submitted_at: PropTypes.string,
  shipment_id: PropTypes.string,
});

export default {
  EvaluationReportShape,
};
