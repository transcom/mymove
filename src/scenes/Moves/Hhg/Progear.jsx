import { get, isEmpty } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues, Field } from 'redux-form';

import { setCurrentShipmentID, getCurrentShipment } from 'shared/UI/ducks';
import { getLastError, getInternalSwaggerDefinition } from 'shared/Swagger/selectors';
import Alert from 'shared/Alert';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';

import { createOrUpdateShipment, getShipment, getShipmentLabel } from 'shared/Entities/modules/shipments';

import './ShipmentWizard.css';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { loadEntitlementsFromState } from 'shared/entitlements';

const formName = 'progear_form';
const ProgearWizardForm = reduxifyWizardForm(formName);

export class Progear extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showInfo: false,
    };
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
      this.props.getShipment(shipmentID, this.props.currentShipment.move_id);
    }
  }

  handleSubmit = () => {
    const moveId = this.props.match.params.moveId;
    let shipment = this.props.formValues;
    const currentShipmentId = get(this.props, 'currentShipment.id');

    if (!shipment.has_pro_gear) {
      shipment = Object.assign({}, shipment, { progear_weight_estimate: 0, spouse_progear_weight_estimate: 0 });
    } else {
      let overrides = {};
      if (!shipment.progear_weight_estimate) {
        overrides.progear_weight_estimate = 0;
      }
      if (!shipment.spouse_progear_weight_estimate) {
        overrides.spouse_progear_weight_estimate = 0;
      }
      if (!isEmpty(overrides)) {
        shipment = Object.assign({}, shipment, overrides);
      }
    }

    return this.props
      .createOrUpdateShipment(moveId, shipment, currentShipmentId)
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
    const { pages, pageKey, error, initialValues, entitlement, currentShipment, schema } = this.props;
    const { showInfo } = this.state;

    const hasProgear = get(this.props, 'formValues.has_pro_gear');

    let progearExceeded = false;
    let spouseProgearExceeded = false;
    if (entitlement) {
      const progearWeightEstimate = get(this.props, 'formValues.progear_weight_estimate');
      if (progearWeightEstimate && progearWeightEstimate > entitlement.pro_gear) {
        progearExceeded = true;
      }
      const spouseProgearWeightEstimate = get(this.props, 'formValues.spouse_progear_weight_estimate');
      if (spouseProgearWeightEstimate && spouseProgearWeightEstimate > entitlement.pro_gear_spouse) {
        spouseProgearExceeded = true;
      }
    }

    // Shipment Wizard
    return (
      <ProgearWizardForm
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
              <h3 className="instruction-heading">Estimate Pro-Gear</h3>
              <div className="usa-width-one-whole">
                {entitlement ? (
                  <div className="weight-info-box">
                    <b>How much Pro-Gear can you have?</b>{' '}
                    <a href="" onClick={this.openInfo}>
                      What qualifies as Pro-gear?
                    </a>
                    <br />
                    You are entitled to move up to {entitlement.pro_gear.toLocaleString()} lbs. of pro-gear and{' '}
                    {entitlement.pro_gear_spouse.toLocaleString()} lbs of spouse pro-gear. Pro-gear includes any gear
                    you or your spouse need to perform your jobs.
                  </div>
                ) : (
                  <LoadingPlaceholder />
                )}
                {showInfo && (
                  <Alert type="info" heading="">
                    Pro-gear includes reference materials, instruments, tools and equipment for technicians, mechanics
                    and similar professions, specialized clothing (diving, astronaut and flying suits and helmets, band
                    uniforms, chaplain vestments and other specialized apparel that’s not usual uniform or clothing),
                    specially-issued field clothing and equipment, and communications equipment used by a member in
                    association with the Military Affiliated Radio System. You'll need to provide a Pro-gear declaration
                    form with a detailed list of items.{' '}
                    <a href="" onClick={this.closeInfo}>
                      Close
                    </a>
                  </Alert>
                )}
                {currentShipment && (
                  <div>
                    <div className="usa-input radio-title">
                      <label className="usa-input-label">Do you or your spouse have any Pro-Gear?</label>
                      <Field name="has_pro_gear" component={YesNoBoolean} />
                    </div>
                    {hasProgear && (
                      <Fragment>
                        <SwaggerField title="Your Pro-Gear:" fieldName="progear_weight_estimate" swagger={schema} />
                        {progearExceeded && (
                          <Alert type="warning" heading="Entitlement exceeded">
                            You have exceeded your entitlement pro-gear weight of{' '}
                            {entitlement.pro_gear.toLocaleString()} lbs.
                          </Alert>
                        )}
                        <SwaggerField
                          title="Spouse's Pro-Gear:"
                          fieldName="spouse_progear_weight_estimate"
                          swagger={schema}
                        />
                        {spouseProgearExceeded && (
                          <Alert type="warning" heading="Entitlement exceeded">
                            You have exceeded your entitlement spouse pro-gear weight of{' '}
                            {entitlement.pro_gear_spouse.toLocaleString()} lbs.
                          </Alert>
                        )}
                      </Fragment>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </ProgearWizardForm>
    );
  }
}
Progear.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createOrUpdateShipment, setCurrentShipmentID, getShipment }, dispatch);
}
function mapStateToProps(state) {
  const shipment = getCurrentShipment(state);

  let initialHasProgear = false;
  if (shipment) {
    if (
      (shipment.progear_weight_estimate && shipment.progear_weight_estimate > 0) ||
      (shipment.spouse_progear_weight_estimate && shipment.spouse_progear_weight_estimate > 0)
    ) {
      initialHasProgear = true;
    }
  }

  const props = {
    schema: getInternalSwaggerDefinition(state, 'Shipment'),
    move: get(state, 'moves.currentMove', {}),
    formValues: getFormValues(formName)(state),
    currentShipment: shipment,
    initialValues: Object.assign({}, shipment, { has_pro_gear: initialHasProgear }),
    error: getLastError(state, getShipmentLabel),
    entitlement: loadEntitlementsFromState(state),
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(Progear);
