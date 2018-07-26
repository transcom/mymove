import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import ShipmentDatePicker from 'scenes/Moves/Hhg/DatePicker';
import ShipmentAddress from 'scenes/Moves/Hhg/Address';

import './ShipmentForm.css';

const formName = 'shipment_form';
const ShipmentFormWizardForm = reduxifyWizardForm(formName);

export class ShipmentForm extends Component {
  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    const shipment = this.props.formValues;
    this.props
      .createOrUpdateShipment(moveId, shipment)
      .then(() => {
        console.log('You did it!');
      })
      .catch(err => {
        this.setState({
          shipmentError: err,
        });
      });
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      initialValues,
    } = this.props;

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
      >
        <div className="shipment-form">
          <div className="usa-grid">
            <h3 className="form-title">Shipment 1 (HHG)</h3>
          </div>
          <ShipmentDatePicker
            schema={this.props.schema}
            error={error}
            formValues={this.props.formValues}
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
  return bindActionCreators({}, dispatch);
}
function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swagger.spec.definitions.Shipment', {}),
    move: get(state, 'moves.currentMove', {}),
    initialValues: get(state, 'moves.currentMove.shipments[0]', {}),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentForm);
