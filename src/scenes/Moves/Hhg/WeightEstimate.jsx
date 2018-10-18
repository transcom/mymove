import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { setCurrentShipmentID, getCurrentShipment } from 'shared/UI/ducks';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { getLastError, getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import Alert from 'shared/Alert';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { loadEntitlementsFromState } from 'shared/entitlements';

import { createOrUpdateShipment, getShipment } from 'shared/Entities/modules/shipments';

import './ShipmentWizard.css';

const formName = 'weight_form';
const getRequestLabel = 'WeightForm.getShipment';
const createOrUpdateRequestLabel = 'WeightForm.createOrUpdateShipment';

const ShipmentFormWizardForm = reduxifyWizardForm(formName);

export class WeightEstimate extends Component {
  constructor(props) {
    super(props);
    this.state = { showInfo: false };
  }

  openInfo = e => {
    e.preventDefault();
    this.setState({ showInfo: true });
  };

  closeInfo = e => {
    e.preventDefault();
    this.setState({ showInfo: false });
  };

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
    const { pages, pageKey, error, initialValues, entitlement } = this.props;
    const weight_estimate = get(this.props, 'formValues.weight_estimate');

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
          <div className="form-section">
            <div className="usa-grid">
              <h3 className="instruction-heading">Estimate Weight</h3>
              <div className="usa-width-one-whole">
                {entitlement ? (
                  <div className="weight-info-box">
                    <b>How much are you entitled to move?</b>
                    <br />
                    {entitlement.weight.toLocaleString()} lbs. + {entitlement.pro_gear.toLocaleString()} lbs. of
                    pro-gear + {entitlement.pro_gear_spouse.toLocaleString()} lbs. of spouse's pro-gear.{' '}
                    <a href="" onClick={this.openInfo}>
                      What's this?
                    </a>
                  </div>
                ) : (
                  <LoadingPlaceholder />
                )}
                {this.state.showInfo && (
                  <Alert type="info" heading="">
                    Your entitlement represents the weight the military is willing to move for you. If you exceed this
                    weight, you'll have to pay for the excess out of pocket. Pro-gear is any gear you need to perform
                    your official duties at your next or later destination, such as reference materials, tools for a
                    technician or mechanic or specialized clothing that's not a typical uniform (such as diving or
                    astronaut suits.{' '}
                    <a href="" onClick={this.closeInfo}>
                      Close
                    </a>
                  </Alert>
                )}
                <p className="review-todo">TODO</p>
                <hr className="weight-estimate-hr" />
                <SwaggerField
                  title="Your estimated weight (in pounds):"
                  className="weight-estimate"
                  fieldName="weight_estimate"
                  swagger={this.props.schema}
                  required
                />
                <div className="weight-estimate-help">
                  If you already know the weight of your stuff, type it in the box.
                </div>
                {entitlement &&
                  (weight_estimate && weight_estimate > entitlement.weight) && (
                    <Alert type="warning" heading="Entitlement exceeded">
                      You have exceeded your entitlement weight of {entitlement.weight.toLocaleString()} lbs.
                    </Alert>
                  )}
              </div>
            </div>
          </div>
        </div>
      </ShipmentFormWizardForm>
    );
  }
}
WeightEstimate.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdateShipment, setCurrentShipmentID, getShipment }, dispatch);
}
function mapStateToProps(state) {
  const shipment = getCurrentShipment(state);
  const props = {
    schema: getInternalSwaggerDefinition(state, 'Shipment'),
    move: get(state, 'moves.currentMove', {}),
    formValues: getFormValues(formName)(state),
    currentShipment: shipment,
    initialValues: shipment,
    error: getLastError(state, getRequestLabel),
    entitlement: loadEntitlementsFromState(state),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(WeightEstimate);
