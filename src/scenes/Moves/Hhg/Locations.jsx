import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { setCurrentShipmentID, getCurrentShipment } from 'shared/UI/ducks';
import { getLastError, getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import Alert from 'shared/Alert';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import Address from 'scenes/Moves/Hhg/Address';

import { createOrUpdateShipment, getShipment, getShipmentLabel } from 'shared/Entities/modules/shipments';

import './ShipmentWizard.css';

const formName = 'locations_form';

const LocationsWizardForm = reduxifyWizardForm(formName);

export class Locations extends Component {
  componentDidMount() {
    this.loadShipment();
  }

  componentDidUpdate(prevProps) {
    if (get(this.props, 'currentShipment.id') !== get(prevProps, 'currentShipment.id')) {
      this.loadShipment();
    }
  }

  loadShipment() {
    const shipmentID = get(this.props, 'currentShipment.id');
    if (shipmentID) {
      this.props.getShipment(shipmentID, this.props.currentShipment.move_id);
    }
  }

  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    const shipment = this.props.formValues;
    const currentShipmentId = get(this.props, 'currentShipment.id');

    return this.props
      .createOrUpdateShipment(moveId, shipment, currentShipmentId)
      .then(action => {
        return this.props.setCurrentShipmentID(Object.keys(action.entities.shipments)[0]);
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
      <LocationsWizardForm
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
          <Address schema={this.props.schema} error={error} formValues={this.props.formValues} />
        </div>
      </LocationsWizardForm>
    );
  }
}
Locations.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdateShipment, setCurrentShipmentID, getShipment }, dispatch);
}
function mapStateToProps(state) {
  const shipment = getCurrentShipment(state);
  const smAddress = get(state, 'serviceMember.currentServiceMember.residential_address', {});
  const props = {
    schema: getInternalSwaggerDefinition(state, 'Shipment'),
    move: get(state, 'moves.currentMove', {}),
    formValues: getFormValues(formName)(state),
    currentShipment: shipment,
    initialValues: { pickup_address: smAddress },
    error: getLastError(state, getShipmentLabel),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(Locations);
