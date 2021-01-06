import React, { Component, Fragment } from 'react';
import { debounce, get, bind, cloneDeep } from 'lodash';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import PropTypes from 'prop-types';
import { getFormValues, reduxForm } from 'redux-form';

import SaveCancelButtons from './SaveCancelButtons';
import { editBegin, editSuccessful, entitlementChangeBegin } from './ducks';

import Alert from 'shared/Alert';
import SectionWrapper from 'components/Customer/SectionWrapper';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { formatDateForSwagger } from 'shared/dates';
import { loadPPMs, updatePPM, selectActivePPMForMove } from 'shared/Entities/modules/ppms';
import scrollToTop from 'shared/scrollToTop';
import { formatCents } from 'shared/formatters';
import { persistPPMEstimate, calculatePPMSITEstimate } from 'services/internalApi';
import { updatePPM as updatePPMInRedux, updatePPMSitEstimate } from 'store/entities/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentMove,
  selectCurrentOrders,
  selectPPMSitEstimate,
} from 'store/entities/selectors';
import 'scenes/Moves/Ppm/DateAndLocation.css';

const sitEstimateDebounceTime = 300;

let EditDateAndLocationForm = (props) => {
  const {
    handleSubmit,
    currentOrders,
    getSitEstimate,
    sitEstimate,
    schema,
    valid,
    sitReimbursement,
    submitting,
  } = props;
  const displayedSitReimbursement = sitEstimate ? '$' + formatCents(sitEstimate) : sitReimbursement;

  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <h1>Edit PPM Dates & Locations</h1>
            <p>Changes could impact your move, including the estimated PPM incentive.</p>
            <SectionWrapper>
              <h2>Move Date</h2>
              <SwaggerField fieldName="original_move_date" onChange={getSitEstimate} swagger={schema} required />
            </SectionWrapper>
            <SectionWrapper>
              <h2>Pickup Location</h2>
              <SwaggerField fieldName="pickup_postal_code" onChange={getSitEstimate} swagger={schema} required />
              <SwaggerField fieldName="has_additional_postal_code" swagger={schema} component={YesNoBoolean} />
              {get(props, 'formValues.has_additional_postal_code', false) && (
                <Fragment>
                  <SwaggerField fieldName="additional_pickup_postal_code" swagger={schema} required />
                  <span className="grey">Making additional stops may decrease your PPM incentive.</span>
                </Fragment>
              )}
            </SectionWrapper>
            <SectionWrapper>
              <h2>Destination Location</h2>
              <p>
                Enter the ZIP for your new home if you know it, or for{' '}
                {currentOrders && currentOrders.new_duty_station.name} if you don't.
              </p>
              <SwaggerField fieldName="destination_postal_code" swagger={schema} required />
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
                  {displayedSitReimbursement && (
                    <div data-testid="storage-estimate" className="storage-estimate">
                      You can spend up to {displayedSitReimbursement} on private storage. Save your receipts to submit
                      with your PPM paperwork.
                    </div>
                  )}
                </Fragment>
              )}
            </SectionWrapper>
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};

const editDateAndLocationFormName = 'edit_date_and_location';
EditDateAndLocationForm = reduxForm({ form: editDateAndLocationFormName, enableReinitialize: true })(
  EditDateAndLocationForm,
);

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
      return this.props.updatePPM(moveId, this.props.currentPPM.id, pendingValues).then(({ response }) => {
        persistPPMEstimate(moveId, response.body.id)
          .then((response) => this.props.updatePPMInRedux(response))
          .then(() => {
            // This promise resolves regardless of error.
            if (!this.props.hasSubmitError) {
              this.props.editSuccessful();
              this.props.history.goBack();
            } else {
              scrollToTop();
            }
          })
          .catch((err) => {
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

  getSitEstimate = (ppmId, moveDate, sitDays, pickupZip, ordersID, weight) => {
    if (sitDays <= 90 && pickupZip.length === 5) {
      const formattedMoveDate = formatDateForSwagger(moveDate);
      calculatePPMSITEstimate(ppmId, formattedMoveDate, sitDays, pickupZip, ordersID, weight).then((response) =>
        this.props.updatePPMSitEstimate(response),
      );
    }
  };

  debouncedSitEstimate = debounce(bind(this.getSitEstimate, this), sitEstimateDebounceTime);

  getDebouncedSitEstimate = (e, value, _, field) => {
    const { currentPPM, formValues, currentOrders } = this.props;
    const estimateValues = cloneDeep(formValues);
    // eslint-disable-next-line
    estimateValues[field] = value;
    this.debouncedSitEstimate(
      currentPPM.id,
      estimateValues.original_move_date,
      estimateValues.days_in_storage,
      estimateValues.pickup_postal_code,
      currentOrders.id,
      currentPPM.weight_estimate,
    );
  };

  componentDidMount() {
    this.props.editBegin();
    this.props.entitlementChangeBegin();
    this.props.loadPPMs(this.props.match.params.moveId);
    scrollToTop();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.currentPPM !== this.props.currentPPM && prevProps.currentOrders !== this.props.currentOrders) {
      const currentPPM = this.props.currentPPM;
      calculatePPMSITEstimate(
        currentPPM.id,
        formatDateForSwagger(currentPPM.original_move_date),
        currentPPM.days_in_storage,
        currentPPM.pickup_postal_code,
        this.props.currentOrders.id,
        currentPPM.weight_estimate,
      ).then((response) => this.props.updatePPMSitEstimate(response));
    }
  }

  render() {
    const {
      initialValues,
      schema,
      formValues,
      sitReimbursement,
      currentOrders,
      error,
      sitEstimate,
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
            sitEstimate={sitEstimate}
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
          />
        </div>
      </div>
    );
  }
}

EditDateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
  updatePPM: PropTypes.func.isRequired,
  error: PropTypes.object,
};
function mapStateToProps(state) {
  const currentMove = selectCurrentMove(state) || {};
  const moveID = currentMove?.id;
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.UpdatePersonallyProcuredMovePayload', {}),
    move: currentMove,
    currentOrders: selectCurrentOrders(state) || {},
    currentPPM: selectActivePPMForMove(state, moveID),
    formValues: getFormValues(editDateAndLocationFormName)(state),
    entitlement: loadEntitlementsFromState(state),
    error: get(state, 'ppm.error'),
    hasSubmitError: get(state, 'ppm.hasSubmitError'),
    sitEstimate: selectPPMSitEstimate(state),
    entitiesSitReimbursement: get(selectActivePPMForMove(state, moveID), 'estimated_storage_reimbursement', ''),
  };

  const defaultPickupZip = serviceMember?.residential_address?.postal_code;

  props.initialValues = props.currentPPM
    ? props.currentPPM
    : defaultPickupZip
    ? {
        pickup_postal_code: defaultPickupZip,
      }
    : null;
  return props;
}

const mapDispatchToProps = {
  push,
  updatePPM,
  loadPPMs,
  editBegin,
  editSuccessful,
  entitlementChangeBegin,
  updatePPMInRedux,
  updatePPMSitEstimate,
};

export default connect(mapStateToProps, mapDispatchToProps)(EditDateAndLocation);
