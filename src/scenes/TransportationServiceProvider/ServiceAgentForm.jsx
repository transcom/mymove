import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { reduxForm, FormSection } from 'redux-form';

import { ServiceAgentEdit, OptionalServiceAgentEdit } from 'shared/TspPanel/ServiceAgentViews';

let ServiceAgentForm = props => {
  const { schema, onCancel, handleSubmit, submitting, valid } = props;
  const originValues = get(props, 'formValues.ORIGIN');
  const destinationValues = get(props, 'formValues.DESTINATION');

  return (
    <form className="infoPanel-wizard" onSubmit={handleSubmit}>
      <div className="infoPanel-wizard-header">Assign servicing agents</div>
      <FormSection name="origin_service_agent">
        <ServiceAgentEdit
          serviceAgentProps={{
            swagger: schema,
            values: originValues,
          }}
          saRole="Origin"
          columnSize="editable-panel-column"
        />
      </FormSection>
      <FormSection name="destination_service_agent">
        <OptionalServiceAgentEdit
          serviceAgentProps={{
            swagger: schema,
            values: destinationValues,
          }}
          saRole="Destination"
          columnSize="editable-panel-column"
        />
      </FormSection>

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
  const { shipment } = state.tsp;
  return {
    ...ownProps,
    initialValues: shipment,
  };
};

export default connect(mapStateToProps)(ServiceAgentForm);
