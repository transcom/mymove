import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import { get } from 'lodash';

import { Field, reduxForm, getFormValues } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import SaveCancelButtons from './SaveCancelButtons';
import DatePicker from 'scenes/Moves/Hhg/DatePicker';
import { moveIsApproved } from 'scenes/Moves/ducks';
import { updateShipment, getShipment } from 'shared/Entities/modules/shipments';

import './Review.css';
import { selectShipment } from '../../shared/Entities/modules/shipments';

let EditShipmentForm = props => {
  const { onSubmit, submitting, valid, shipment } = props;
  const moveID = get(shipment, 'move_id');
  return (
    <form onSubmit={onSubmit}>
      {moveID && (
        <Fragment>
          <Field
            name="requested_pickup_date"
            component={DatePicker}
            availableMoveDates={props.availableMoveDates}
            currentShipment={shipment}
            moveID={moveID}
          />
          <div className="usa-grid">
            <div className="usa-width-one-whole">
              <SaveCancelButtons valid={valid} submitting={submitting} />
            </div>
          </div>
        </Fragment>
      )}
    </form>
  );
};

const editShipmentFormName = 'edit_shipment';
EditShipmentForm = reduxForm({ form: editShipmentFormName })(EditShipmentForm);

class EditShipment extends Component {
  componentDidMount() {
    if (this.props.onDidMount) {
      this.props.onDidMount();
    }
  }

  render() {
    const { error, shipment, schemaAffiliation, schemaRank, updateShipment, returnToReview } = this.props;

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
          onSubmit={updateShipment}
          onCancel={returnToReview}
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
    moveIsApproved: moveIsApproved(state),
    schemaRank: get(state, 'swaggerInternal.spec.definitions.ServiceMemberRank', {}),
    schemaAffiliation: get(state, 'swaggerInternal.spec.definitions.Affiliation', {}),
    formValues: getFormValues(editShipmentFormName)(state),
  };
}

const getShipmentLabel = 'EditShipment.getShipment';
function mapDispatchToProps(dispatch, ownProps) {
  const shipmentID = ownProps.match.params.shipmentId;
  return {
    onDidMount: function() {
      dispatch(getShipment(getShipmentLabel, shipmentID));
    },
    updateShipment: function(values, shipment) {
      dispatch(updateShipment(shipmentID, values)).then(function(action) {
        if (!action.error) {
          const moveID = Object.values(action.entities.shipments)[0].move_id;
          if (shipment.status !== 'DRAFT') {
            dispatch(push(`/moves/${moveID}/edit`));
          } else {
            dispatch(push(`/moves/${moveID}/review`));
          }
        }
      });
    },
    returnToReview: function(moveID) {
      dispatch(push(`/moves/${moveID}/review`));
    },
  };
}

function mergeProps(stateProps, dispatchProps, ownProps) {
  return Object.assign({}, stateProps, dispatchProps, ownProps, {
    updateShipment: function(event) {
      event.preventDefault();
      dispatchProps.updateShipment(stateProps.formValues, stateProps.shipment);
    },
    returnToReview: function(event) {
      event.preventDefault();
      const moveID = stateProps.shipment.move_id;
      dispatchProps.returnToReview(moveID);
    },
  });
}

export default connect(mapStateToProps, mapDispatchToProps, mergeProps)(EditShipment);
