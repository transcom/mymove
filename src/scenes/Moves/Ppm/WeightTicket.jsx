import React, { Component, Fragment } from 'react';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import Alert from 'shared/Alert';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import RadioButton from 'shared/RadioButton';
import Checkbox from 'shared/Checkbox';
import Uploader from 'shared/Uploader';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';
import './PPMPaymentRequest.css';

class WeightTicket extends Component {
  state = {
    vehicleType: 'CAR',
    additionalWeightTickets: 'Yes',
    missingEmptyWeightTicket: false,
    missingFullWeightTicket: false,
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

  render() {
    const { additionalWeightTickets, vehicleType, missingEmptyWeightTicket, missingFullWeightTicket } = this.state;
    const { schema } = this.props;
    const nextBtnLabel = additionalWeightTickets === 'Yes' ? 'Save & Add Another' : 'Save & Continue';
    const uploadEmptyTicketLabel =
      '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload upload empty weight ticket</span></span>';
    const uploadFullTicketLabel =
      '<span class="uploader-label">Drag & drop or <span class="filepond--label-action">click to upload upload full weight ticket</span></span>';
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
        <div className="usa-grid">
          <SwaggerField
            fieldName="vehicle_options"
            swagger={schema}
            onChange={event => this.handleChange(event, 'vehicleType')}
            value={vehicleType}
            required
          />
          <SwaggerField fieldName="vehicle_nickname" swagger={schema} required />
          {(vehicleType === 'CAR' || vehicleType === 'BOX_TRUCK') && (
            <>
              <div className="dashed-divider" />

              <div className="usa-grid-full" style={{ marginTop: '1em' }}>
                <div className="usa-width-one-third">
                  <strong className="input-header">Empty Weight</strong>
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
                  <Uploader options={{ labelIdle: uploadEmptyTicketLabel }} />
                  <Checkbox
                    label="I'm missing this weight ticket"
                    name="missingEmptyWeightTicket"
                    checked={missingEmptyWeightTicket}
                    onChange={this.handleCheckboxChange}
                    normalizeLabel
                  />
                  <Alert type="warning">
                    Contact your local Transportation Office (PPPO) to let them know you’re missing this weight ticket.
                    For now, keep going and enter the info you do have.
                  </Alert>
                </div>
              </div>
              <div className="usa-grid-full" style={{ marginTop: '1em' }}>
                <div className="usa-width-one-third">
                  <strong className="input-header">Full Weight</strong>
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
                  <Uploader options={{ labelIdle: uploadFullTicketLabel }} />
                  <Checkbox
                    label="I'm missing this weight ticket"
                    name="missingFullWeightTicket"
                    checked={missingFullWeightTicket}
                    onChange={this.handleCheckboxChange}
                    normalizeLabel
                  />
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
          <PPMPaymentRequestActionBtns nextBtnLabel={nextBtnLabel} />
        </div>
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

function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.WeightTicketPayload', {}),
  };
  return props;
}
export default connect(mapStateToProps)(WeightTicket);
