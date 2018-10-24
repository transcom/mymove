import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { Field, reduxForm } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import SaveCancelButtons from './SaveCancelButtons';
import DatePicker from 'scenes/Moves/Hhg/DatePicker';
import { moveIsApproved } from 'scenes/Moves/ducks';
import { createOrUpdateShipment, getShipment } from 'shared/Entities/modules/shipments';

import './Review.css';
import profileImage from './images/profile.png';
import { selectShipment } from '../../shared/Entities/modules/shipments';

const editShipmentFormName = 'edit_shipment';

let EditShipmentForm = props => {
  const { handleSubmit, submitting, valid, shipment } = props;
  const moveID = get(shipment, 'move_id');
  return (
    <form onSubmit={handleSubmit}>
      {moveID && (
        <Fragment>
          <Field
            name="requested_pickup_date"
            component={DatePicker}
            availableMoveDates={props.availableMoveDates}
            currentShipment={shipment}
            moveID={moveID}
          />
          <div class="usa-grid">
            <div class="usa-width-one-whole">
              <SaveCancelButtons valid={valid} submitting={submitting} />
            </div>
          </div>
        </Fragment>
      )}
    </form>
  );
};
EditShipmentForm = reduxForm({ form: editShipmentFormName })(EditShipmentForm);

class EditShipment extends Component {
  componentDidMount() {
    if (this.props.onDidMount) {
      this.props.onDidMount();
    }
  }

  render() {
    const { error, schema, shipment, schemaAffiliation, schemaRank } = this.props;

    return (
      <Fragment>
        {error && (
          <div className="usa-grid">
            <div className="usa-width-one-whole error-message">
              <Alert type="error" heading="An error occurred">
                {error.message}
              </Alert>
            </div>
          </div>
        )}
        <EditShipmentForm
          initialValues={shipment}
          shipment={shipment}
          onSubmit={this.updateShipment}
          onCancel={this.returnToReview}
          schema={schema}
          schemaRank={schemaRank}
          schemaAffiliation={schemaAffiliation}
        />
      </Fragment>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const shipment = selectShipment(state, ownProps.match.params.shipmentId);
  return {
    shipment,
    move: get(state, 'moves.currentMove'),
    error: get(state, 'serviceMember.error'),
    hasSubmitError: get(state, 'serviceMember.hasSubmitError'),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    moveIsApproved: moveIsApproved(state),
    schemaRank: get(state, 'swaggerInternal.spec.definitions.ServiceMemberRank', {}),
    schemaAffiliation: get(state, 'swaggerInternal.spec.definitions.Affiliation', {}),
  };
}

const getShipmentLabel = 'EditShipment.getShipment';
function mapDispatchToProps(dispatch, ownProps) {
  const shipmentID = ownProps.match.params.shipmentId;
  return {
    onDidMount: function() {
      dispatch(getShipment(getShipmentLabel, shipmentID));
    },
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(EditShipment);
