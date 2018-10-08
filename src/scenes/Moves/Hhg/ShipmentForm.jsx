import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { setCurrentShipment, currentShipment } from 'shared/UI/ducks';
import { getLastError, getSwaggerDefinition } from 'shared/Swagger/selectors';
import Alert from 'shared/Alert';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import Address from 'scenes/Moves/Hhg/Address';
import WeightEstimates from 'scenes/Moves/Hhg/WeightEstimates';

import {
  createOrUpdateShipment,
  getShipment,
} from 'shared/Entities/modules/shipments';

import './ShipmentWizard.css';

const formName = 'shipment_form';
const getRequestLabel = 'ShipmentForm.getShipment';
const createOrUpdateRequestLabel = 'ShipmentForm.createOrUpdateShipment';

const ShipmentFormWizardForm = reduxifyWizardForm(formName);

export class ShipmentForm extends Component {
  componentDidMount() {
    this.loadShipment();
  }

  componentDidUpdate(prevProps) {
    if (
      get(this.props, 'currentShipment.id') !==
      get(prevProps, 'currentShipment.id')
    ) {
      this.loadShipment();
    }
  }

  loadShipment() {
    const shipmentID = get(this.props, 'currentShipment.id');
    if (shipmentID) {
      this.props.getShipment(
        getRequestLabel,
        shipmentID,
        this.props.currentShipment.move_id,
      );
    }
  }

  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    const shipment = this.props.formValues;
    const currentShipmentId = get(this.props, 'currentShipment.id');

    return this.props
      .createOrUpdateShipment(
        createOrUpdateRequestLabel,
        moveId,
        shipment,
        currentShipmentId,
      )
      .then(data => {
        return this.props.setCurrentShipment(data.body);
      })
      .catch(err => {
        this.setState({
          shipmentError: err,
        });
        return { error: err };
      });
  };

  render() {
    const { pages, pageKey, error, initialValues } = this.props;

    // Shipment Wizard
    return (
      <ShipmentFormWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
        initialValues={initialValues}
      >
        <Fragment>
          {this.props.error && (
            <div className="usa-grid">
              <div className="usa-width-one-whole error-message">
                <Alert type="error" heading="An error occurred">
                  Something went wrong contacting the server.
                </Alert>
              </div>
            </div>
          )}
        </Fragment>
        <div className="shipment-wizard">
          <div className="usa-grid">
            <h3 className="form-title">Shipment 1 (HHG)</h3>
          </div>
          <Address
            schema={this.props.schema}
            error={error}
            formValues={this.props.formValues}
          />
          <WeightEstimates
            schema={this.props.schema}
            error={error}
            formValues={this.props.formValues}
          />
        </div>
      </ShipmentFormWizardForm>
    );
  }
}
ShipmentForm.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { createOrUpdateShipment, setCurrentShipment, getShipment },
    dispatch,
  );
}
function mapStateToProps(state) {
  const shipment = currentShipment(state);
  const props = {
    schema: getSwaggerDefinition(state, 'Shipment'),
    move: get(state, 'moves.currentMove', {}),
    formValues: getFormValues(formName)(state),
    currentShipment: shipment,
    initialValues: shipment,
    error: getLastError(state, getRequestLabel),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentForm);
