import React, { Component, Fragment } from 'react';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import './PPMPaymentRequest.css';
import Uploader from 'shared/Uploader';

class WeightTicket extends Component {
  state = {
    value: '',
  };
  labelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload upload empty weight ticket</span>';

  handleChange = event => {
    this.setState({ value: event.target.value });
  };

  render() {
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
          {this.state.value !== '' && (
            <div className="uploader-box">
              <Uploader options={{ allowMultiple: true, labelIdle: this.labelIdle }} />
            </div>
          )}
          <SwaggerField fieldName="vehicle_nickname" swagger={schema} required />

          {/* TODO: change onclick handler to go to next page in flow */}
          <PPMPaymentRequestActionBtns onClick={() => {}} nextBtnLabel="Save & Add Another" />
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
