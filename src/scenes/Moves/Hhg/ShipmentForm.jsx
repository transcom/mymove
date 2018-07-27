import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import { setCurrentShipment, currentShipment } from 'shared/UI/ducks';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import ShipmentDatePicker from 'scenes/Moves/Hhg/DatePicker';
import ShipmentAddress from 'scenes/Moves/Hhg/Address';

import { createOrUpdateShipment } from 'shared/Entities/modules/shipments';

import './ShipmentForm.css';

const formName = 'shipment_form';
const ShipmentFormWizardForm = reduxifyWizardForm(formName);

export class ShipmentForm extends Component {
  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    const shipment = this.props.formValues;
    const currentShipmentId = get(this, 'props.currentShipment.id');
    this.props
      .createOrUpdateShipment(moveId, shipment, currentShipmentId)
      .then(data => {
        this.props.setCurrentShipment(data.body);
      })
      .catch(err => {
        this.setState({
          shipmentError: err,
        });
      });
  };

  setDate = day => {
    this.setState({ requestedPickupDate: day });
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      initialValues,
    } = this.props;

    const requestedPickupDate = get(this.state, 'requestedPickupDate');

    // Shipment Wizard
    return (
      <ShipmentFormWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
        additionalValues={{ requested_pickup_date: requestedPickupDate }}
      >
        <div className="shipment-form">
          <div className="usa-grid">
            <h3 className="form-title">Shipment 1 (HHG)</h3>
          </div>
          <ShipmentDatePicker
            schema={this.props.schema}
            error={error}
            formValues={this.props.formValues}
            setDate={this.setDate}
          />
          <ShipmentAddress
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
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { createOrUpdateShipment, setCurrentShipment },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swagger.spec.definitions.Shipment', {}),
    move: get(state, 'moves.currentMove', {}),
    initialValues: get(state, 'moves.currentMove.shipments[0]', {}),
    formValues: getFormValues(formName)(state),
    currentShipment: currentShipment(state) || {},
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentForm);
