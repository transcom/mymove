import React, { Component, Fragment } from 'react';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import RadioButton from 'shared/RadioButton';
import './PPMPaymentRequest.css';
import Uploader from 'shared/Uploader';

export const SaveAddAnotherButton = () => (
  <PPMPaymentRequestActionBtns onClick={() => {}} nextBtnLabel="Save & Add Another" />
);

export const SaveContinueButton = () => (
  <PPMPaymentRequestActionBtns onClick={() => {}} nextBtnLabel="Save & Continue" />
);

export class WeightTicket extends Component {
  state = {
    value: 'Yes',
  };

  labelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload upload empty weight ticket</span>';

  handleRadioChange = event => {
    this.setState({ value: event.target.value });
  };

  render() {
    const { value } = this.state;
    const { schema } = this.props;
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
            onChange={this.handleChange}
            value={this.state.value}
            required
          />
          <SwaggerField fieldName="vehicle_nickname" swagger={schema} required />
          <div className="dashed-divider" />

          <div className="usa-grid-full">
            <div className="usa-width-one-third">
              <SwaggerField className="short-field" fieldName="empty_weight" swagger={schema} title=" " required /> lbs
            </div>
            <div className="usa-width-two-third">
              <Uploader options={{ labelIdle: this.labelIdle }} />
            </div>
          </div>

          <div className="usa-grid-full">
            <div className="usa-width-one-third">
              <SwaggerField
                className="short-field"
                fieldName="full_weight"
                swagger={schema}
                title="Full weight at destination"
                required
              />{' '}
              lbs
            </div>

            <div className="usa-width-two-third uploader-container">
              <Uploader options={{ labelIdle: this.labelIdle }} />
            </div>
          </div>

          <SwaggerField fieldName="weight_ticket_date" swagger={schema} required />
          <div className="dashed-divider" />

          <p className="radio-group-header">Do you have more weight tickets for another vehicle or trip?</p>
          <RadioButton
            inputClassName="inline_radio"
            labelClassName="inline_radio"
            label="Yes"
            value="Yes"
            name="additional_weight_ticket"
            checked={value === 'Yes'}
            onChange={this.handleRadioChange}
          />

          <RadioButton
            inputClassName="inline_radio"
            labelClassName="inline_radio"
            label="No"
            value="No"
            name="additional_weight_ticket"
            checked={value === 'No'}
            onChange={this.handleRadioChange}
          />

          {/* TODO: change onclick handler to go to next page in flow */}
          {value === 'Yes' ? (
            <SaveAddAnotherButton moreWeightTickets={value} />
          ) : (
            <SaveContinueButton moreWeightTickets={value} />
          )}
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
