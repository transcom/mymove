import React, { Component } from 'react';
import { getFormValues, reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { get, map } from 'lodash';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { withLastLocation } from 'react-router-last-location';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import RadioButton from 'shared/RadioButton';
import Checkbox from 'shared/Checkbox';
import Uploader from 'shared/Uploader';
import Alert from 'shared/Alert';
import { formatDateForSwagger } from 'shared/dates';
import { documentSizeLimitMsg, WEIGHT_TICKET_SET_TYPE } from 'shared/constants';
import { selectCurrentPPM, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

import carTrailerImg from 'shared/images/car-trailer_mobile.png';
import carImg from 'shared/images/car_mobile.png';
import { createWeightTicketSetDocument } from 'shared/Entities/modules/weightTicketSetDocuments';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { withContext } from 'shared/AppContext';

import { getNextPage } from './utility';
import DocumentsUploaded from './PaymentReview/DocumentsUploaded';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';
import { formatToOrdinal } from 'shared/formatters';

import './PPMPaymentRequest.css';

const nextBtnLabels = {
  SaveAndAddAnother: 'Save & Add Another',
  SaveAndContinue: 'Save & Continue',
};

const reviewPagePath = '/ppm-payment-review';
const nextPagePath = '/ppm-expenses-intro';

const uploadEmptyTicketLabel =
  '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload empty weight ticket</span></span>';
const uploadFullTicketLabel =
  '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload full weight ticket</span></span>';
const uploadTrailerProofOfOwnership =
  '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload documentation</span></span>';

class WeightTicket extends Component {
  state = { ...this.initialState };
  uploaders = {
    trailer: { uploaderRef: null, isMissingChecked: () => this.state.missingDocumentation },
    emptyWeight: { uploaderRef: null, isMissingChecked: () => this.state.missingEmptyWeightTicket },
    fullWeight: { uploaderRef: null, isMissingChecked: () => this.state.missingFullWeightTicket },
  };

  get initialState() {
    return {
      weightTicketSetType: '',
      additionalWeightTickets: 'No',
      isValidTrailer: 'No',
      weightTicketSubmissionError: false,
      missingDocumentation: false,
      missingEmptyWeightTicket: false,
      missingFullWeightTicket: false,
    };
  }

  componentDidMount() {
    const { moveId } = this.props;
    this.props.getMoveDocumentsForMove(moveId);
  }

  get isCarTrailer() {
    return this.state.weightTicketSetType === WEIGHT_TICKET_SET_TYPE.CAR_TRAILER;
  }

  get isCar() {
    return this.state.weightTicketSetType === WEIGHT_TICKET_SET_TYPE.CAR;
  }

  get isProGear() {
    return this.state.weightTicketSetType === WEIGHT_TICKET_SET_TYPE.PRO_GEAR;
  }

  hasWeightTicket = (uploaderRef) => {
    return !!(uploaderRef && !uploaderRef.isEmpty());
  };

  invalidState = (uploader) => {
    if (uploader.isMissingChecked()) {
      return true;
    } else return !this.hasWeightTicket(uploader.uploaderRef);
  };

  carTrailerText = (isValidTrailer) => {
    if (this.isCarTrailer && isValidTrailer === 'Yes') {
      return (
        <div style={{ marginBottom: '1em' }}>
          You can claim this trailer's weight as part of the total weight of your trip.
        </div>
      );
    } else if (this.isCarTrailer) {
      return (
        <div style={{ marginBottom: '1em' }}>
          The weight of this trailer should be <strong>excluded</strong> from the total weight of this trip.
          <p>{documentSizeLimitMsg}</p>
        </div>
      );
    }
    // if there is no trailer, don't show text
    return undefined;
  };

  uploaderWithInvalidState = () => {
    // Validation for the vehicle type
    if (this.state.isValidTrailer === 'Yes' && this.isCarTrailer && this.invalidState(this.uploaders.trailer)) {
      return true;
    }
    // Full weight must be in a valid state to proceed.
    return this.invalidState(this.uploaders.fullWeight);
  };

  //  handleChange for weightTicketSetType and additionalWeightTickets
  handleChange = (event, type) => {
    this.setState({ [type]: event.target.value });
  };

  handleCheckboxChange = (event) => {
    this.setState({
      [event.target.name]: event.target.checked,
    });
  };

  onAddFile = (uploaderName) => () => {
    this.setState((prevState) => ({
      uploaderIsIdle: { ...prevState.uploaderIsIdle, [uploaderName]: false },
    }));
  };

  onUploadChange = (uploaderName) => (uploaderIsIdle) => {
    this.setState((prevState) => ({
      uploaderIsIdle: { ...prevState.uploaderIsIdle, [uploaderName]: uploaderIsIdle },
    }));
  };

  skipHandler = () => {
    const { moveId, history } = this.props;
    history.push(`/moves/${moveId}${nextPagePath}`);
  };

  nonEmptyUploaderKeys() {
    const uploadersKeys = Object.keys(this.uploaders);
    return uploadersKeys.filter((key) => this.uploaders[key].uploaderRef && !this.uploaders[key].uploaderRef.isEmpty());
  }

  saveAndAddHandler = (formValues) => {
    const { moveId, currentPpm, history } = this.props;
    const { additionalWeightTickets } = this.state;

    const uploaderKeys = this.nonEmptyUploaderKeys();
    const uploadIds = [];
    // eslint-disable-next-line no-unused-vars
    for (const key of uploaderKeys) {
      let files = this.uploaders[key].uploaderRef.getFiles();
      const documentUploadIds = map(files, 'id');
      uploadIds.push(...documentUploadIds);
    }
    const weightTicketSetDocument = {
      personally_procured_move_id: currentPpm.id,
      upload_ids: uploadIds,
      weight_ticket_set_type: formValues.weight_ticket_set_type,
      vehicle_nickname: formValues.vehicle_nickname,
      vehicle_make: formValues.vehicle_make,
      vehicle_model: formValues.vehicle_model,
      empty_weight_ticket_missing: this.state.missingEmptyWeightTicket,
      empty_weight: formValues.empty_weight,
      full_weight_ticket_missing: this.state.missingFullWeightTicket,
      full_weight: formValues.full_weight,
      weight_ticket_date: formatDateForSwagger(formValues.weight_ticket_date),
      trailer_ownership_missing: this.state.missingDocumentation,
      move_document_type: 'WEIGHT_TICKET_SET',
      notes: formValues.notes,
    };
    return this.props
      .createWeightTicketSetDocument(moveId, weightTicketSetDocument)
      .then(() => {
        this.cleanup();
        if (additionalWeightTickets === 'No') {
          const nextPage = getNextPage(`/moves/${moveId}${nextPagePath}`, this.props.lastLocation, reviewPagePath);
          history.push(nextPage);
        }
      })
      .catch((e) => {
        this.setState({ weightTicketSubmissionError: true });
      });
  };

  cleanup = () => {
    const { reset } = this.props;
    const uploaders = this.uploaders;
    const uploaderKeys = this.nonEmptyUploaderKeys();
    // eslint-disable-next-line no-unused-vars
    for (const key of uploaderKeys) {
      uploaders[key].uploaderRef.clearFiles();
    }
    reset();
    this.setState({ ...this.initialState });
  };

  render() {
    const {
      additionalWeightTickets,
      weightTicketSetType,
      missingEmptyWeightTicket,
      missingFullWeightTicket,
      missingDocumentation,
      isValidTrailer,
    } = this.state;
    const { handleSubmit, submitting, schema, weightTicketSets, invalid, moveId, transportationOffice } = this.props;
    const nextBtnLabel =
      additionalWeightTickets === 'Yes' ? nextBtnLabels.SaveAndAddAnother : nextBtnLabels.SaveAndContinue;
    const weightTicketSetOrdinal = formatToOrdinal(weightTicketSets.length + 1);
    const fullWeightTicketFieldsRequired = missingFullWeightTicket ? null : true;
    const emptyWeightTicketFieldsRequired = missingEmptyWeightTicket ? null : true;

    return (
      <div className="grid-container usa-prose">
        <WizardHeader
          title="Weight tickets"
          right={
            <ProgressTimeline>
              <ProgressTimelineStep name="Weight" current />
              <ProgressTimelineStep name="Expenses" />
              <ProgressTimelineStep name="Review" />
            </ProgressTimeline>
          }
        />
        <div className="grid-row">
          <div className="grid-col-12">
            <DocumentsUploaded moveId={moveId} />
          </div>
        </div>
        <div className="grid-row">
          <div className="grid-col-12">
            <form>
              {this.state.weightTicketSubmissionError && (
                <div className="grid-row">
                  <div className="grid-col-12 error-message">
                    <Alert type="error" heading="An error occurred">
                      Something went wrong contacting the server.
                    </Alert>
                  </div>
                </div>
              )}
              <div className="expenses-container">
                <h3 className="expenses-header">Weight Tickets - {weightTicketSetOrdinal} set</h3>
                Upload weight tickets for each vehicle trip and pro-gear weigh.{' '}
                <Link to="/weight-ticket-examples" className="usa-link">
                  <FontAwesomeIcon aria-hidden className="color_blue_link" icon="question-circle" />
                </Link>
                <SwaggerField
                  fieldName="weight_ticket_set_type"
                  swagger={schema}
                  onChange={(event) => this.handleChange(event, 'weightTicketSetType')}
                  value={weightTicketSetType}
                  required
                />
                {weightTicketSetType &&
                  (this.isCarTrailer || this.isCar ? (
                    <>
                      <SwaggerField
                        fieldName="vehicle_make"
                        data-testid="vehicle_make"
                        title="Vehicle make"
                        swagger={schema}
                        required={this.isCarTrailer || this.isCar}
                      />
                      <SwaggerField
                        fieldName="vehicle_model"
                        data-testid="vehicle_model"
                        title="Vehicle model"
                        swagger={schema}
                        required={this.isCarTrailer || this.isCar}
                      />
                    </>
                  ) : (
                    <SwaggerField
                      fieldName="vehicle_nickname"
                      data-testid="vehicle_nickname"
                      title={
                        this.isProGear
                          ? "Pro-gear type (ex. 'My pro-gear', 'Spouse pro-gear', 'Both')"
                          : 'Nickname (ex. "15-foot truck")'
                      }
                      swagger={schema}
                      required={!this.isCarTrailer && !this.isCar}
                    />
                  ))}
                {weightTicketSetType && this.isCarTrailer && (
                  <>
                    <div className="radio-group-wrapper normalize-margins">
                      <p className="radio-group-header">
                        Do you own this trailer, and does it meet all{' '}
                        <Link to="/trailer-criteria" className="usa-link">
                          trailer criteria
                        </Link>
                        ?
                      </p>
                      <RadioButton
                        inputClassName="usa-radio__input inline_radio"
                        labelClassName="usa-radio__label inline_radio"
                        label="Yes"
                        value="Yes"
                        name="isValidTrailer"
                        checked={isValidTrailer === 'Yes'}
                        onChange={(event) => this.handleChange(event, 'isValidTrailer')}
                      />

                      <RadioButton
                        inputClassName="usa-radio__input inline_radio"
                        labelClassName="usa-radio__label inline_radio"
                        label="No"
                        value="No"
                        name="isValidTrailer"
                        checked={isValidTrailer === 'No'}
                        onChange={(event) => this.handleChange(event, 'isValidTrailer')}
                      />
                    </div>
                    {isValidTrailer === 'Yes' && (
                      <>
                        <p className="normalize-margins" style={{ marginTop: '1em' }}>
                          Proof of ownership (ex. registration, bill of sale)
                        </p>
                        <p>{documentSizeLimitMsg}</p>
                        <span data-testid="trailer-upload">
                          <Uploader
                            options={{ labelIdle: uploadTrailerProofOfOwnership }}
                            onRef={(ref) => (this.uploaders.trailer.uploaderRef = ref)}
                            onChange={this.onUploadChange('trailer')}
                            onAddFile={this.onAddFile('trailer')}
                          />
                        </span>
                        <Checkbox
                          label="I don't have ownership documentation"
                          name="missingDocumentation"
                          checked={missingDocumentation}
                          onChange={this.handleCheckboxChange}
                        />
                        {missingDocumentation && (
                          <div className="grid-row">
                            <div className="grid-col-8" data-testid="trailer-warning">
                              <Alert type="warning">
                                If your state does not provide a registration or bill of sale for your trailer, you may
                                write and upload a signed and dated statement certifying that you or your spouse own the
                                trailer and meets the{' '}
                                <Link to="/trailer-criteria" className="usa-link">
                                  trailer criteria
                                </Link>
                                . Upload your statement using the proof of ownership field.
                              </Alert>
                            </div>
                          </div>
                        )}
                      </>
                    )}
                  </>
                )}
                {weightTicketSetType && (
                  <>
                    <div className="dashed-divider" />

                    <div className="grid-row">
                      <div className="grid-col-12" style={{ marginTop: '1em' }}>
                        {this.carTrailerText(isValidTrailer)}
                        <div className="grid-row grid-gap">
                          <div className="grid-col-4 input-group">
                            <strong className="input-header">
                              Empty Weight{' '}
                              {this.isCarTrailer &&
                                (isValidTrailer === 'Yes' ? (
                                  <>
                                    ( <img alt="car only" className="car-img" src={carImg} /> car only)
                                  </>
                                ) : (
                                  <>
                                    ( <img alt="car and trailer" className="car-img" src={carTrailerImg} /> car +
                                    trailer)
                                  </>
                                ))}
                            </strong>
                            <SwaggerField
                              className="short-field"
                              fieldName="empty_weight"
                              swagger={schema}
                              hideLabel
                              required={emptyWeightTicketFieldsRequired}
                            />{' '}
                            lbs
                          </div>
                          <div className="grid-col-8 uploader-wrapper">
                            <span data-testid="empty-weight-upload">
                              <Uploader
                                options={{ labelIdle: uploadEmptyTicketLabel }}
                                onRef={(ref) => (this.uploaders.emptyWeight.uploaderRef = ref)}
                                onChange={this.onUploadChange('emptyWeight')}
                                onAddFile={this.onAddFile('emptyWeight')}
                              />
                            </span>
                            <Checkbox
                              label="I'm missing this weight ticket"
                              name="missingEmptyWeightTicket"
                              checked={missingEmptyWeightTicket}
                              onChange={this.handleCheckboxChange}
                              normalizeLabel
                            />
                            {missingEmptyWeightTicket && (
                              <span data-testid="empty-warning">
                                <Alert type="warning">
                                  Contact your local Transportation Office (PPPO) to let them know you’re missing this
                                  weight ticket. For now, keep going and enter the info you do have.
                                </Alert>
                              </span>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                    <div className="grid-row grid-gap input-group" style={{ marginTop: '1em' }}>
                      <div className="grid-col-4">
                        <strong className="input-header">
                          Full Weight{' '}
                          {this.isCarTrailer && (
                            <>
                              ( <img alt="car and trailer" className="car-img" src={carTrailerImg} /> car + trailer)
                            </>
                          )}
                        </strong>
                        <label className="full-weight-label">Full weight at destination</label>
                        <SwaggerField
                          className="short-field"
                          fieldName="full_weight"
                          swagger={schema}
                          hideLabel
                          required={fullWeightTicketFieldsRequired}
                        />{' '}
                        lbs
                      </div>

                      <div className="grid-col-8 uploader-wrapper">
                        <div data-testid="full-weight-upload">
                          <Uploader
                            options={{ labelIdle: uploadFullTicketLabel }}
                            onRef={(ref) => (this.uploaders.fullWeight.uploaderRef = ref)}
                            onChange={this.onUploadChange('fullWeight')}
                            onAddFile={this.onAddFile('fullWeight')}
                          />
                        </div>
                        <Checkbox
                          label="I'm missing this weight ticket"
                          name="missingFullWeightTicket"
                          checked={missingFullWeightTicket}
                          onChange={this.handleCheckboxChange}
                          normalizeLabel
                        />
                        {missingFullWeightTicket && (
                          <div data-testid="full-warning">
                            <Alert type="warning">
                              <b>You can’t get paid without a full weight ticket.</b> See what you can do to find it,
                              because without certified documentation of the weight of your belongings, we can’t pay you
                              your incentive. Call the {transportationOffice.name} Transportation Office at{' '}
                              {transportationOffice.phone_lines[0]} if you have any questions.
                            </Alert>
                          </div>
                        )}
                      </div>
                    </div>

                    <SwaggerField
                      fieldName="weight_ticket_date"
                      swagger={schema}
                      required={fullWeightTicketFieldsRequired}
                    />
                    <div className="dashed-divider" />

                    <div className="radio-group-wrapper">
                      <p className="radio-group-header">Do you have more weight tickets for another vehicle or trip?</p>
                      <RadioButton
                        inputClassName="usa-radio__input inline_radio"
                        labelClassName="usa-radio__label inline_radio"
                        label="Yes"
                        value="Yes"
                        name="additional_weight_ticket"
                        checked={additionalWeightTickets === 'Yes'}
                        onChange={(event) => this.handleChange(event, 'additionalWeightTickets')}
                      />

                      <RadioButton
                        inputClassName="usa-radio__input inline_radio"
                        labelClassName="usa-radio__label inline_radio"
                        label="No"
                        value="No"
                        name="additional_weight_ticket"
                        checked={additionalWeightTickets === 'No'}
                        onChange={(event) => this.handleChange(event, 'additionalWeightTickets')}
                      />
                    </div>
                  </>
                )}
                <PPMPaymentRequestActionBtns
                  nextBtnLabel={nextBtnLabel}
                  hasConfirmation={true}
                  submitButtonsAreDisabled={this.uploaderWithInvalidState() || invalid}
                  submitting={submitting}
                  skipHandler={this.skipHandler}
                  saveAndAddHandler={handleSubmit(this.saveAndAddHandler)}
                  displaySkip={weightTicketSets.length >= 1}
                />
              </div>
            </form>
          </div>
        </div>
      </div>
    );
  }
}

const formName = 'weight_ticket_wizard';
WeightTicket = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(WeightTicket);

WeightTicket.propTypes = {
  schema: PropTypes.object.isRequired,
};

function mapStateToProps(state, ownProps) {
  const moveId = ownProps.match.params.moveId;
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const dutyLocationId = serviceMember?.current_location?.id;
  const transportationOffice = serviceMember?.current_location.transportation_office;

  return {
    moveId: moveId,
    formValues: getFormValues(formName)(state),
    genericMoveDocSchema: get(state, 'swaggerInternal.spec.definitions.CreateGenericMoveDocumentPayload', {}),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateWeightTicketDocumentsPayload', {}),
    currentPpm: selectCurrentPPM(state) || {},
    weightTicketSets: selectPPMCloseoutDocumentsForMove(state, moveId, ['WEIGHT_TICKET_SET']),
    transportationOffice: transportationOffice,
    dutyLocationId: dutyLocationId,
  };
}

const mapDispatchToProps = {
  getMoveDocumentsForMove,
  createWeightTicketSetDocument,
};

export default withContext(withLastLocation(connect(mapStateToProps, mapDispatchToProps)(WeightTicket)));
