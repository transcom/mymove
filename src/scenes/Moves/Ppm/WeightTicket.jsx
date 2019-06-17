import React, { Component, Fragment } from 'react';
import { reduxForm, getFormValues } from 'redux-form';
import { connect } from 'react-redux';
import { get, map } from 'lodash';
import PropTypes from 'prop-types';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import RadioButton from 'shared/RadioButton';
import Checkbox from 'shared/Checkbox';
import Uploader from 'shared/Uploader';
import Alert from 'shared/Alert';

import carTrailerImg from 'shared/images/car-trailer_mobile.png';
import carImg from 'shared/images/car_mobile.png';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';
import './PPMPaymentRequest.css';
import { createWeightTicketSetDocument } from 'shared/Entities/modules/weightTicketSetDocuments';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faQuestionCircle from '@fortawesome/fontawesome-free-solid/faQuestionCircle';

const vehicleTypes = {
  CarAndTrailer: 'CAR_TRAILER',
  Car: 'CAR',
  BoxTruck: 'BOX_TRUCK',
};

const nextBtnLabels = {
  SaveAndAddAnother: 'Save & Add Another',
  SaveAndContinue: 'Save & Continue',
};

const uploadEmptyTicketLabel =
  '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload empty weight ticket</span></span>';
const uploadFullTicketLabel =
  '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload full weight ticket</span></span>';
const uploadTrailerProofOfOwnership =
  '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload documentation</span></span>';

class WeightTicket extends Component {
  state = { ...this.initialState, weightTicketNumber: 1 };
  uploaders = {
    trailer: { uploaderRef: null, isMissingChecked: () => this.state.missingDocumentation },
    emptyWeight: { uploaderRef: null, isMissingChecked: () => this.state.missingEmptyWeightTicket },
    fullWeight: { uploaderRef: null, isMissingChecked: () => this.state.missingFullWeightTicket },
  };

  get initialState() {
    return {
      vehicleType: '',
      additionalWeightTickets: 'Yes',
      isValidTrailer: 'No',
      weightTicketSubmissionError: false,
      missingDocumentation: false,
      missingEmptyWeightTicket: false,
      missingFullWeightTicket: false,
    };
  }

  get isCarTrailer() {
    return this.state.vehicleType === vehicleTypes.CarAndTrailer;
  }

  hasWeightTicket = uploaderRef => {
    return !!(uploaderRef && !uploaderRef.isEmpty());
  };

  invalidState = uploader => {
    return this.hasWeightTicket(uploader.uploaderRef) === uploader.isMissingChecked();
  };

  uploaderWithInvalidState = () => {
    if (this.state.isValidTrailer === 'Yes' && (this.isCarTrailer && this.invalidState(this.uploaders.trailer))) {
      return true;
    }
    return this.invalidState(this.uploaders.emptyWeight) || this.invalidState(this.uploaders.fullWeight);
  };

  formIsIncomplete = () => {
    const { formValues } = this.props;
    return !(
      formValues &&
      formValues.vehicle_nickname &&
      formValues.vehicle_options &&
      formValues.empty_weight &&
      formValues.full_weight &&
      formValues.weight_ticket_date
    );
  };

  //  handleChange for vehicleType and additionalWeightTickets
  handleChange = (event, type) => {
    this.setState({ [type]: event.target.value });
  };

  handleCheckboxChange = event => {
    this.setState({
      [event.target.name]: event.target.checked,
    });
  };

  onAddFile = uploaderName => () => {
    this.setState({
      uploaderIsIdle: { ...this.state.uploaderIsIdle, [uploaderName]: false },
    });
  };

  onUploadChange = uploaderName => uploaderIsIdle => {
    this.setState({
      uploaderIsIdle: { ...this.state.uploaderIsIdle, [uploaderName]: uploaderIsIdle },
    });
  };

  saveForLaterHandler = formValues => {
    const { history } = this.props;
    return this.saveAndAddHandler(formValues).then(() => {
      if (this.state.weightTicketSubmissionError === false) {
        history.push('/');
      }
    });
  };

  nonEmptyUploaderKeys() {
    const uploadersKeys = Object.keys(this.uploaders);
    return uploadersKeys.filter(
      // eslint-disable-next-line security/detect-object-injection
      key => this.uploaders[key].uploaderRef && this.uploaders[key].uploaderRef.getFiles() > 0,
    );
  }

  saveAndAddHandler = formValues => {
    const { moveId, currentPpm, history } = this.props;
    const { additionalWeightTickets } = this.state;

    const uploaderKeys = this.nonEmptyUploaderKeys();
    const uploadIds = [];
    for (const key of uploaderKeys) {
      // eslint-disable-next-line security/detect-object-injection
      let files = this.uploaders[key].uploaderRef.getFiles();
      const documentUploadIds = map(files, 'id');
      uploadIds.push(...documentUploadIds);
    }
    const weightTicketSetDocument = {
      personally_procured_move_id: currentPpm.id,
      upload_ids: uploadIds,
      vehicle_options: formValues.vehicle_options,
      vehicle_nickname: formValues.vehicle_nickname,
      empty_weight_ticket_missing: this.state.missingEmptyWeightTicket,
      empty_weight: formValues.empty_weight,
      full_weight_ticket_missing: this.state.missingFullWeightTicket,
      full_weight: formValues.full_weight,
      weight_ticket_date: formValues.weight_ticket_date,
      trailer_ownership_missing: this.state.missingDocumentation,
      move_document_type: 'WEIGHT_TICKET_SET',
      notes: formValues.notes,
    };
    return this.props
      .createWeightTicketSetDocument(moveId, weightTicketSetDocument)
      .then(() => {
        this.setState({ weightTicketNumber: this.state.weightTicketNumber + 1 });
        this.cleanup();
        if (additionalWeightTickets === 'No') {
          history.push(`/moves/${moveId}/ppm-expenses-intro`);
        }
      })
      .catch(e => {
        this.setState({ weightTicketSubmissionError: true });
      });
  };

  cleanup = () => {
    const { reset } = this.props;
    const uploaders = this.uploaders;
    const uploaderKeys = this.nonEmptyUploaderKeys();
    for (const key of uploaderKeys) {
      // eslint-disable-next-line security/detect-object-injection
      uploaders[key].uploaderRef.clearFiles();
    }
    reset();
    this.setState({ ...this.initialState });
  };

  // maps int to int with ordinal 1 -> 1st, 2 -> 2nd, 3rd ...
  numberWithOrdinal = n => {
    const s = ['th', 'st', 'nd', 'rd'];
    const v = n % 100;
    // eslint-disable-next-line security/detect-object-injection
    return n + (s[(v - 20) % 10] || s[v] || s[0]);
  };

  submitButtonsAreDisabled = () => {
    return this.formIsIncomplete() || this.uploaderWithInvalidState();
  };

  render() {
    const {
      additionalWeightTickets,
      vehicleType,
      missingEmptyWeightTicket,
      missingFullWeightTicket,
      missingDocumentation,
      isValidTrailer,
    } = this.state;
    const { handleSubmit, submitting, schema } = this.props;
    const nextBtnLabel =
      additionalWeightTickets === 'Yes' ? nextBtnLabels.SaveAndAddAnother : nextBtnLabels.SaveAndContinue;

    return (
      <Fragment>
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
        <form>
          {this.state.weightTicketSubmissionError && (
            <div className="usa-grid">
              <div className="usa-width-one-whole error-message">
                <Alert type="error" heading="An error occurred">
                  Something went wrong contacting the server.
                </Alert>
              </div>
            </div>
          )}
          <div className="usa-grid expenses-container">
            <h3 className="expenses-header">
              Weight Tickets - {this.numberWithOrdinal(this.state.weightTicketNumber)} set
            </h3>
            Upload an <strong>empty</strong> & <strong>full</strong> weight ticket below for <em>only</em>{' '}
            <strong>one</strong> vehicle or trip at a time until they're all uploaded.{' '}
            <Link to="/weight-ticket-examples">
              <FontAwesomeIcon aria-hidden className="color_blue_link" icon={faQuestionCircle} />
            </Link>
            <SwaggerField
              fieldName="vehicle_options"
              swagger={schema}
              onChange={event => this.handleChange(event, 'vehicleType')}
              value={vehicleType}
              required
            />
            <SwaggerField fieldName="vehicle_nickname" swagger={schema} required />
            {vehicleType &&
              this.isCarTrailer && (
                <>
                  <div className="radio-group-wrapper normalize-margins">
                    <p className="radio-group-header">
                      Is this a <strong>different</strong> trailer you own and does it meet the{' '}
                      <Link to="/trailer-criteria">trailer criteria</Link>?
                    </p>
                    <RadioButton
                      inputClassName="inline_radio"
                      labelClassName="inline_radio"
                      label="Yes"
                      value="Yes"
                      name="isValidTrailer"
                      checked={isValidTrailer === 'Yes'}
                      onChange={event => this.handleChange(event, 'isValidTrailer')}
                    />

                    <RadioButton
                      inputClassName="inline_radio"
                      labelClassName="inline_radio"
                      label="No"
                      value="No"
                      name="isValidTrailer"
                      checked={isValidTrailer === 'No'}
                      onChange={event => this.handleChange(event, 'isValidTrailer')}
                    />
                  </div>
                  {isValidTrailer === 'Yes' && (
                    <>
                      <p className="normalize-margins" style={{ marginTop: '1em' }}>
                        Proof of ownership (ex. registration, bill of sale)
                      </p>
                      <span data-cy="trailer-upload">
                        <Uploader
                          options={{ labelIdle: uploadTrailerProofOfOwnership }}
                          onRef={ref => (this.uploaders.trailer.uploaderRef = ref)}
                          onChange={this.onUploadChange('trailer')}
                          onAddFile={this.onAddFile('trailer')}
                        />
                      </span>
                      <Checkbox
                        label="I don't have ownership documentation"
                        name="missingDocumentation"
                        checked={missingDocumentation}
                        onChange={this.handleCheckboxChange}
                        normalizeLabel
                      />
                      {missingDocumentation && (
                        <div className="one-half" data-cy="trailer-warning">
                          <Alert type="warning">
                            If your state does not provide a registration or bill of sale for your trailer, you may
                            write and upload a signed and dated statement certifying that you or your spouse own the
                            trailer and meets the <Link to="/trailer-criteria">trailer criteria</Link>. Upload your
                            statement using the proof of ownership field.
                          </Alert>
                        </div>
                      )}
                    </>
                  )}
                </>
              )}
            {vehicleType && (
              <>
                <div className="dashed-divider" />

                <div className="usa-grid-full" style={{ marginTop: '1em' }}>
                  {this.isCarTrailer && isValidTrailer === 'Yes' ? (
                    <div style={{ marginBottom: '1em' }}>
                      You can claim this trailer's weight as part of the total weight of your trip.
                    </div>
                  ) : (
                    <div style={{ marginBottom: '1em' }}>
                      The weight of this trailer should be <strong>excluded</strong> from the total weight of this trip.
                    </div>
                  )}
                  <div className="usa-width-one-third input-group">
                    <strong className="input-header">
                      Empty Weight{' '}
                      {this.isCarTrailer &&
                        (isValidTrailer === 'Yes' ? (
                          <>
                            ( <img alt="car only" className="car-img" src={carImg} /> car only)
                          </>
                        ) : (
                          <>
                            ( <img alt="car and trailer" className="car-img" src={carTrailerImg} /> car + trailer)
                          </>
                        ))}
                    </strong>
                    <SwaggerField
                      className="short-field"
                      fieldName="empty_weight"
                      swagger={schema}
                      hideLabel
                      required
                    />{' '}
                    lbs
                  </div>
                  <div className="usa-width-two-thirds uploader-wrapper">
                    <span data-cy="empty-weight-upload">
                      <Uploader
                        options={{ labelIdle: uploadEmptyTicketLabel }}
                        onRef={ref => (this.uploaders.emptyWeight.uploaderRef = ref)}
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
                      <span data-cy="empty-warning">
                        <Alert type="warning">
                          Contact your local Transportation Office (PPPO) to let them know you’re missing this weight
                          ticket. For now, keep going and enter the info you do have.
                        </Alert>
                      </span>
                    )}
                  </div>
                </div>
                <div className="usa-grid-full input-group" style={{ marginTop: '1em' }}>
                  <div className="usa-width-one-third">
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
                      required
                    />{' '}
                    lbs
                  </div>

                  <div className="usa-width-two-thirds uploader-wrapper">
                    <span data-cy="full-weight-upload">
                      <Uploader
                        options={{ labelIdle: uploadFullTicketLabel }}
                        onRef={ref => (this.uploaders.fullWeight.uploaderRef = ref)}
                        onChange={this.onUploadChange('fullWeight')}
                        onAddFile={this.onAddFile('fullWeight')}
                      />
                    </span>
                    <Checkbox
                      label="I'm missing this weight ticket"
                      name="missingFullWeightTicket"
                      checked={missingFullWeightTicket}
                      onChange={this.handleCheckboxChange}
                      normalizeLabel
                    />
                    {missingFullWeightTicket && (
                      <span data-cy="full-warning">
                        <Alert type="warning">
                          Contact your local Transportation Office (PPPO) to let them know you’re missing this weight
                          ticket. For now, keep going and enter the info you do have.
                        </Alert>
                      </span>
                    )}
                  </div>
                </div>

                <SwaggerField fieldName="weight_ticket_date" swagger={schema} required />
                <div className="dashed-divider" />

                <div className="radio-group-wrapper">
                  <p className="radio-group-header">Do you have more weight tickets for another vehicle or trip?</p>
                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="Yes"
                    value="Yes"
                    name="additional_weight_ticket"
                    checked={additionalWeightTickets === 'Yes'}
                    onChange={event => this.handleChange(event, 'additionalWeightTickets')}
                  />

                  <RadioButton
                    inputClassName="inline_radio"
                    labelClassName="inline_radio"
                    label="No"
                    value="No"
                    name="additional_weight_ticket"
                    checked={additionalWeightTickets === 'No'}
                    onChange={event => this.handleChange(event, 'additionalWeightTickets')}
                  />
                </div>
              </>
            )}
            {/* TODO: change onclick handler to go to next page in flow */}
            <PPMPaymentRequestActionBtns
              nextBtnLabel={nextBtnLabel}
              submitButtonsAreDisabled={this.submitButtonsAreDisabled()}
              submitting={submitting}
              saveForLaterHandler={handleSubmit(this.saveForLaterHandler)}
              saveAndAddHandler={handleSubmit(this.saveAndAddHandler)}
              displaySaveForLater={true}
            />
          </div>
        </form>
      </Fragment>
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
  return {
    moveId: moveId,
    formValues: getFormValues(formName)(state),
    genericMoveDocSchema: get(state, 'swaggerInternal.spec.definitions.CreateGenericMoveDocumentPayload', {}),
    moveDocSchema: get(state, 'swaggerInternal.spec.definitions.MoveDocumentPayload', {}),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateWeightTicketDocumentsPayload', {}),
    currentPpm: get(state, 'ppm.currentPpm'),
  };
}
const mapDispatchToProps = {
  createWeightTicketSetDocument,
};

export default connect(mapStateToProps, mapDispatchToProps)(WeightTicket);
