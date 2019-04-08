import React, { Component, Fragment } from 'react';
import helpIcon from 'shared/images/help-icon.png';
import Select, { components } from 'react-select';
import { reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import WizardHeader from '../WizardHeader';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import './PPMPaymentRequest.css';

const { Option } = components;
// make x-display-value values into an array with the img as label and x-display-value value
// the value is the key of the x-display-value
// refer to line 40 & 41 in JsonSchemaField.js
const vehicleTypes = [
  {
    value: 'Stuffs',
    label: (
      <Fragment>
        <img src={helpIcon} alt="" /> Stuffs
      </Fragment>
    ),
  },
];
export class WeightTicket extends Component {
  state = { vehicleOptions: '' };

  renderOption() {
    return <Option>Whatever</Option>;
  }

  render() {
    const { vehicleOptions } = this.state;
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
          <Select value={vehicleOptions} options={vehicleTypes} />
          {/*<SwaggerField fieldName="vehicle_options" swagger={schema} required />*/}
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
