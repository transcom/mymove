import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import ShipmentDatePicker from 'scenes/Moves/Hhg/DatePicker';
import ShipmentAddress from 'scenes/Moves/Hhg/Address';

const formName = 'shipment_form';
const ShipmentFormWizardForm = reduxifyWizardForm(formName);
//hard-coded placeholder schema:
const schema = {
  properties: {
    planned_move_date: {
      type: 'string',
      format: 'date',
      example: '2018-04-26',
      title: 'Move Date',
      'x-nullable': true,
      'x-always-required': true,
    },
    has_secondary_pickup_address: {
      type: 'boolean',
      title: 'Do you have household goods at any other pickup location?',
      'x-nullable': true,
    },
    has_delivery_address: {
      type: 'boolean',
      title: 'Do you know your home address at your destination yet?',
      'x-nullable': true,
    },
  },
};

export class ShipmentForm extends Component {
  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
    } = this.props;
    const initialValues = get(currentServiceMember, 'residential_address');
    // Shipment Wizard
    return (
      <ShipmentFormWizardForm
        // handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
      >
        <div className="usa-grid">
          <h3 className="form-title">Shipment 1 (HHG)</h3>
          <ShipmentDatePicker
            schema={schema}
            error={error}
            formValues={this.props.formValues}
          />
          <ShipmentAddress
            schema={schema}
            error={error}
            formValues={this.props.formValues}
          />
        </div>
      </ShipmentFormWizardForm>
    );
  }
}
ShipmentForm.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
function mapStateToProps(state) {
  const props = {
    // schema: get(
    //   state,
    //   'swagger.spec.definitions.CreateHouseholdGoodsPayload',
    //   {},
    // ),
    schema,
    formValues: getFormValues(formName)(state),
    ...state.serviceMember,
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentForm);
