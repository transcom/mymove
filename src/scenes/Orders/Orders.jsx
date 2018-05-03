import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { Field } from 'redux-form';

import { createOrders, updateOrders, showCurrentOrders } from './ducks';
import { loadServiceMember } from 'scenes/ServiceMembers/ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './Orders.css';

const YesNoBoolean = props => {
  const {
    input: { value: rawValue, onChange },
  } = props;
  let value = rawValue;
  const yesChecked = value === true;
  const noChecked = value === false;

  const localOnChange = event => {
    if (event.target.id === 'yes') {
      onChange(true);
    } else {
      onChange(false);
    }
  };

  return (
    <Fragment>
      <input
        id="yes"
        type="radio"
        onChange={localOnChange}
        checked={yesChecked}
      />
      <label htmlFor="yes">Yes</label>
      <input
        id="no"
        type="radio"
        onChange={localOnChange}
        checked={noChecked}
      />
      <label htmlFor="no">No</label>
    </Fragment>
  );
};

const validateOrdersForm = (values, form) => {
  let errors = {};

  const required_fields = ['has_dependents', 'new_duty_station'];

  required_fields.forEach(fieldName => {
    if (values[fieldName] === undefined || values[fieldName] === '') {
      errors[fieldName] = 'Required.';
    }
  });

  return errors;
};

const formName = 'orders_info';
const OrdersWizardForm = reduxifyWizardForm(formName, validateOrdersForm);

export class Orders extends Component {
  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    // Update if orders object already extant
    if (pendingValues) {
      pendingValues['service_member_id'] = this.props.currentServiceMember.id;
      if (this.props.currentOrders) {
        this.props.updateOrders(this.props.currentOrders.id, pendingValues);
      } else {
        this.props.createOrders(pendingValues);
      }
    }
  };

  componentDidUpdate(prevProps, prevState) {
    // If we don't have a service member yet, fetch one when loggedInUser loads.
    if (
      !prevProps.user.loggedInUser &&
      this.props.user.loggedInUser &&
      !this.props.currentServiceMember
    ) {
      const serviceMemberID = this.props.user.loggedInUser.service_member.id;
      this.props.loadServiceMember(serviceMemberID);
      this.props.showCurrentOrders(serviceMemberID);
    }
  }

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
    } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentOrders ? currentOrders : null;

    return (
      <OrdersWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        hasSucceeded={hasSubmitSuccess}
        serverError={error}
        initialValues={initialValues}
      >
        <h1 className="sm-heading">Your Orders</h1>
        <SwaggerField
          fieldName="orders_type"
          swagger={this.props.schema}
          required
        />
        <SwaggerField
          fieldName="issue_date"
          swagger={this.props.schema}
          required
        />
        <SwaggerField
          fieldName="report_by_date"
          swagger={this.props.schema}
          required
        />
        <fieldset key="dependents">
          <legend htmlFor="dependents">
            Are dependents included in your orders?
          </legend>
          <Field name="has_dependents" component={YesNoBoolean} />
        </fieldset>
        <Field
          name="new_duty_station"
          component={DutyStationSearchBox}
          affiliation={this.props.affiliation}
        />
      </OrdersWizardForm>
    );
  }
}
Orders.propTypes = {
  schema: PropTypes.object.isRequired,
  updateOrders: PropTypes.func.isRequired,
  currentOrders: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateOrders, createOrders, showCurrentOrders, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  const currentServiceMember = state.serviceMember.currentServiceMember;
  const affiliation = currentServiceMember
    ? currentServiceMember.affiliation
    : null;
  const error = state.serviceMember.error || state.orders.error;
  const hasSubmitSuccess =
    state.serviceMember.hasSubmitSuccess || state.orders.hasSubmitSuccess;
  const props = {
    affiliation,
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
    user: state.loggedInUser,
    currentOrders: state.orders.currentOrders,
    error,
    hasSubmitSuccess,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateUpdateOrdersPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(Orders);
