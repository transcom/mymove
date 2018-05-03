import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { Field } from 'redux-form';

import { createOrders, updateOrders, loadOrders } from './ducks';
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

const uiSchema = {
  title: 'Tell Us About Your Move Orders',
  order: [
    'orders_type',
    'issue_date',
    'report_by_date',
    'has_dependents',
    'new_duty_station',
  ],
  requiredFields: [
    'orders_type',
    'issue_date',
    'report_by_date',
    'has_dependents',
    'new_duty_station',
  ],
  custom_components: {
    has_dependents: YesNoBoolean,
    new_duty_station: DutyStationSearchBox,
  },
};
const subsetOfFields = [
  'orders_type',
  'issue_date',
  'report_by_date',
  'has_dependents',
  'new_duty_station',
];
const formName = 'orders_info';
const OrdersWizardForm = reduxifyWizardForm(formName);

export class Orders extends Component {
  // componentDidMount() {
  //   if (!this.props.currentOrders) {
  //     console.log('WHy are we reloading this?');
  //     this.props.loadOrders(this.props.match.params.orderId);
  //   }
  // }

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
    const initialValues = currentOrders
      ? pick(currentOrders, subsetOfFields)
      : null;
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
    { updateOrders, createOrders, loadOrders, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  const currentServiceMember = state.serviceMember.currentServiceMember;
  const affiliation = currentServiceMember
    ? currentServiceMember.affiliation
    : null;
  const props = {
    affiliation,
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
    user: state.loggedInUser,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateUpdateOrdersPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(Orders);
