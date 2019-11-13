import { debounce, get, bind, cloneDeep } from 'lodash';
import { push } from 'react-router-redux';
import PropTypes from 'prop-types';
import { getFormValues } from 'redux-form';
import SaveCancelButtons from './SaveCancelButtons';
import React, { Component, Fragment } from 'react';
import { reduxForm } from 'redux-form';
import Alert from 'shared/Alert'; // eslint-disable-line
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { createOrUpdatePpm, getPpmSitEstimate } from 'scenes/Moves/Ppm/ducks';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';
import { updatePPMEstimate } from 'shared/Entities/modules/ppms';
import 'scenes/Moves/Ppm/DateAndLocation.css';
import { editBegin, editSuccessful, entitlementChangeBegin } from './ducks';
import scrollToTop from 'shared/scrollToTop';

const sitEstimateDebounceTime = 300;

let EditDateAndLocationForm = props => {
  const { handleSubmit, currentOrders, getSitEstimate, schema, valid, sitReimbursement, submitting } = props;
  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <h1 className="sm-heading"> Edit PPM Dates & Locations </h1>
            <p>Changes could impact your move, including the estimated PPM incentive.</p>
            <h3 className="sm-heading-2"> Move Date </h3>
            <SwaggerField fieldName="original_move_date" onChange={getSitEstimate} swagger={schema} required />
            <hr className="spacer" />
            <h3 className="sm-heading-2">Pickup Location</h3>
            <SwaggerField fieldName="pickup_postal_code" onChange={getSitEstimate} swagger={schema} required />
            <SwaggerField fieldName="has_additional_postal_code" swagger={schema} component={YesNoBoolean} />
            {get(props, 'formValues.has_additional_postal_code', false) && (
              <Fragment>
                <SwaggerField fieldName="additional_pickup_postal_code" swagger={schema} required />
                <span className="grey">Making additional stops may decrease your PPM incentive.</span>
              </Fragment>
            )}
            <hr className="spacer" />
            <h3 className="sm-heading-2">Destination Location</h3>
            <p>
              Enter the ZIP for your new home if you know it, or for{' '}
              {currentOrders && currentOrders.new_duty_station.name} if you don't.
            </p>
            <SwaggerField fieldName="destination_postal_code" swagger={schema} onChange={getSitEstimate} required />
            <span className="grey">
              The ZIP code for {currentOrders && currentOrders.new_duty_station.name} is{' '}
              {currentOrders && currentOrders.new_duty_station.address.postal_code}{' '}
            </span>
            <SwaggerField fieldName="has_sit" swagger={schema} component={YesNoBoolean} />
            {get(props, 'formValues.has_sit', false) && (
              <Fragment>
                <SwaggerField
                  className="days-in-storage"
                  fieldName="days_in_storage"
                  swagger={schema}
                  onChange={getSitEstimate}
                  required
                />{' '}
                <span className="grey">You can choose up to 90 days.</span>
                {sitReimbursement && (
                  <div data-cy="storage-estimate" className="storage-estimate">
                    You can spend up to {sitReimbursement} on private storage. Save your receipts to submit with your
                    PPM paperwork.
                  </div>
                )}
              </Fragment>
            )}
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};

const editDateAndLocationFormName = 'edit_date_and_location';
EditDateAndLocationForm = reduxForm({ form: editDateAndLocationFormName })(EditDateAndLocationForm);

class EditDateAndLocation extends Component {
  handleSubmit = () => {
    const pendingValues = Object.assign({}, this.props.formValues);
    if (pendingValues) {
      pendingValues.has_additional_postal_code = pendingValues.has_additional_postal_code || false;
      pendingValues.has_sit = pendingValues.has_sit || false;
      if (!pendingValues.has_sit) {
        pendingValues.days_in_storage = null;
      }

      const moveId = this.props.match.params.moveId;
      return this.props.createOrUpdatePpm(moveId, pendingValues).then(({ payload }) => {
        this.props
          .updatePPMEstimate(moveId, payload.id)
          .then(() => {
            // This promise resolves regardless of error.
            if (!this.props.hasSubmitError) {
              this.props.editSuccessful();
              this.props.history.goBack();
            } else {
              scrollToTop();
            }
          })
          .catch(err => {
            // This promise resolves regardless of error.
            if (!this.props.hasSubmitError) {
              this.props.editSuccessful();
              this.props.history.goBack();
            } else {
              scrollToTop();
            }
            return err;
          });
      });
    }
  };

  getSitEstimate = (moveDate, sitDays, pickupZip, destZip, weight) => {
    if (sitDays <= 90 && pickupZip.length === 5 && destZip.length === 5) {
      this.props.getPpmSitEstimate(moveDate, sitDays, pickupZip, destZip, weight);
    }
  };

  debouncedSitEstimate = debounce(bind(this.getSitEstimate, this), sitEstimateDebounceTime);

  getDebouncedSitEstimate = (e, value, _, field) => {
    const { currentPpm, formValues } = this.props;
    const estimateValues = cloneDeep(formValues);
    // eslint-disable-next-line
    estimateValues[field] = value;
    this.debouncedSitEstimate(
      estimateValues.original_move_date,
      estimateValues.days_in_storage,
      estimateValues.pickup_postal_code,
      estimateValues.destination_postal_code,
      currentPpm.weight_estimate,
    );
  };

  componentDidMount() {
    this.props.editBegin();
    this.props.entitlementChangeBegin();
  }

  render() {
    const {
      initialValues,
      schema,
      formValues,
      sitReimbursement,
      currentOrders,
      error,
      entitiesSitReimbursement,
    } = this.props;
    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditDateAndLocationForm
            onSubmit={this.handleSubmit}
            getSitEstimate={this.getDebouncedSitEstimate}
            initialValues={initialValues}
            schema={schema}
            formValues={formValues}
            sitReimbursement={
              sitReimbursement !== entitiesSitReimbursement && sitReimbursement
                ? sitReimbursement
                : entitiesSitReimbursement
            }
            currentOrders={currentOrders}
            onCancel={this.returnToReview}
            createOrUpdatePpm={createOrUpdatePpm}
          />
        </div>
      </div>
    );
  }
}

EditDateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
  createOrUpdatePpm: PropTypes.func.isRequired,
  error: PropTypes.object,
};
function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.UpdatePersonallyProcuredMovePayload', {}),
    ...state.ppm,
    move: get(state, 'moves.currentMove'),
    currentOrders: get(state.orders, 'currentOrders'),
    currentPpm: get(state.ppm, 'currentPpm'),
    formValues: getFormValues(editDateAndLocationFormName)(state),
    entitlement: loadEntitlementsFromState(state),
    error: get(state, 'ppm.error'),
    hasSubmitError: get(state, 'ppm.hasSubmitError'),
    entitiesSitReimbursement: get(
      selectPPMForMove(state, get(state, 'moves.currentMove.id')),
      'estimated_storage_reimbursement',
      '',
    ),
  };
  const defaultPickupZip = get(state.serviceMember, 'currentServiceMember.residential_address.postal_code');
  props.initialValues = props.currentPpm
    ? props.currentPpm
    : defaultPickupZip
    ? {
        pickup_postal_code: defaultPickupZip,
      }
    : null;
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      push,
      createOrUpdatePpm,
      getPpmSitEstimate,
      editBegin,
      editSuccessful,
      entitlementChangeBegin,
      updatePPMEstimate,
    },
    dispatch,
  );
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(EditDateAndLocation);
