import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { reduxForm } from 'redux-form';
import { selectShipment } from 'shared/Entities/modules/shipments';

let ServiceAgentForm = props => {
  const { onCancel, handleSubmit, submitting, valid } = props;

  return (
    <form className="infoPanel-wizard" onSubmit={handleSubmit}>
      <div className="infoPanel-wizard-header">Assign servicing agents</div>
      <div className="infoPanel-wizard-actions-container">
        <a className="infoPanel-wizard-cancel" onClick={onCancel}>
          Cancel
        </a>
        <button type="submit" disabled={submitting || !valid}>
          Done
        </button>
      </div>
    </form>
  );
};

ServiceAgentForm = reduxForm({ form: 'tsp_service_agents' })(ServiceAgentForm);

ServiceAgentForm.propTypes = {
  schema: PropTypes.object,
  onCancel: PropTypes.func,
  handleSubmit: PropTypes.func,
  submitting: PropTypes.bool,
  valid: PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const { shipmentId } = ownProps;
  const shipment = selectShipment(state, shipmentId);
  return {
    ...ownProps,
    initialValues: shipment,
  };
};

export default connect(mapStateToProps)(ServiceAgentForm);
