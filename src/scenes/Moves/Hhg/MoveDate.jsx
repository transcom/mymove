import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, Field } from 'redux-form';

import { setCurrentShipmentID, getCurrentShipment } from 'shared/UI/ducks';
import { getLastError } from 'shared/Swagger/selectors';
import Alert from 'shared/Alert';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import DatePicker from 'scenes/Moves/Hhg/DatePicker';
import { validateAdditionalFields } from 'shared/JsonSchemaForm';

import { createOrUpdateShipment, getShipment } from 'shared/Entities/modules/shipments';
import { getAvailableMoveDates, selectAvailableMoveDates } from 'shared/Entities/modules/calendar';

import './ShipmentWizard.css';

const validateMoveDateForm = validateAdditionalFields(['requested_pickup_date']);

const formName = 'move_date_form';
const getRequestLabel = 'MoveDate.getShipment';
const createOrUpdateRequestLabel = 'MoveDate.createOrUpdateShipment';
const availableMoveDatesLabel = 'MoveDate.getAvailableMoveDates';
const MoveDateWizardForm = reduxifyWizardForm(formName, validateMoveDateForm);

export class MoveDate extends Component {
  componentDidMount() {
    this.loadShipment();
    this.props.getAvailableMoveDates(availableMoveDatesLabel, this.props.today);
  }

  componentDidUpdate(prevProps) {
    if (get(this.props, 'currentShipment.id') !== get(prevProps, 'currentShipment.id')) {
      this.loadShipment();
    }
  }

  loadShipment() {
    const shipmentID = get(this.props, 'currentShipment.id');
    if (shipmentID) {
      this.props.getShipment(getRequestLabel, shipmentID, this.props.currentShipment.move_id);
    }
  }

  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    const shipment = this.props.formValues;
    const currentShipmentId = get(this.props, 'currentShipment.id');

    return this.props
      .createOrUpdateShipment(createOrUpdateRequestLabel, moveId, shipment, currentShipmentId)
      .then(action => {
        const id = Object.keys(action.entities.shipments)[0];
        return this.props.setCurrentShipmentID(id);
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
      <MoveDateWizardForm
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
          <Field
            name="requested_pickup_date"
            component={DatePicker}
            availableMoveDates={this.props.availableMoveDates}
            currentShipment={this.props.currentShipment}
          />
        </div>
      </MoveDateWizardForm>
    );
  }
}
MoveDate.propTypes = {
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      createOrUpdateShipment,
      setCurrentShipmentID,
      getShipment,
      getAvailableMoveDates,
    },
    dispatch,
  );
}
function mapStateToProps(state) {
  const shipment = getCurrentShipment(state);
  const today = new Date().toISOString().split('T')[0];
  const props = {
    move: get(state, 'moves.currentMove', {}),
    formValues: getFormValues(formName)(state),
    currentShipment: shipment,
    initialValues: shipment,
    today: today,
    availableMoveDates: selectAvailableMoveDates(state, today),
    error: getLastError(state, getRequestLabel),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(MoveDate);
